package provider

import (
	"context"
	"fmt"

	"github.com/cloudbase/garm-provider-common/execution"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-equinix/config"
	"github.com/cloudbase/garm-provider-equinix/internal/spec"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	metal "github.com/equinix/equinix-sdk-go/services/metalv1"
)

var _ execution.ExternalProvider = &equinixProvider{}

func NewEquinixProvider(configPath, controllerID string) (execution.ExternalProvider, error) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	configuration := metal.NewConfiguration()
	configuration.AddDefaultHeader("X-Auth-Token", conf.AuthToken)

	api_client := metal.NewAPIClient(configuration)

	return &equinixProvider{
		cfg:          conf,
		cli:          api_client,
		controllerID: controllerID,
	}, nil
}

type equinixProvider struct {
	cli          *metal.APIClient
	cfg          *config.Config
	controllerID string
}

func (a *equinixProvider) CreateInstance(ctx context.Context, bootstrapParams params.BootstrapInstance) (instance params.ProviderInstance, err error) {
	spec, err := spec.GetRunnerSpecFromBootstrapParams(bootstrapParams, a.controllerID)
	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("failed to get runner spec: %w", err)
	}
	userdata, err := spec.ComposeUserData()
	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("failed to compose userdata: %w", err)
	}

	metro := a.cfg.MetroCode
	if spec.MetroCode != "" {
		metro = spec.MetroCode
	}
	hostname := bootstrapParams.Name
	if bootstrapParams.OSType == params.Windows {
		// Equnix has a maximum of 15 characters for Windows
		hostname = hostname[:15]
	}
	deviceRequest := metal.CreateDeviceRequest{
		DeviceCreateInMetroInput: &metal.DeviceCreateInMetroInput{
			Metro:                 metro,
			Plan:                  bootstrapParams.Flavor,
			OperatingSystem:       bootstrapParams.Image,
			Tags:                  spec.Tags,
			Userdata:              &userdata,
			HardwareReservationId: spec.HardwareReservationID,
			Hostname:              &hostname,
		},
	}

	var device *metal.Device

	defer func() {
		if err != nil {
			nameOrID := bootstrapParams.Name
			if device != nil && device.GetId() != "" {
				nameOrID = device.GetId()
			}
			if err := a.DeleteInstance(ctx, nameOrID); err != nil {
				fmt.Printf("failed to delete instance: %q", err)
			}
		}
	}()

	device, _, err = a.cli.DevicesApi.CreateDevice(context.Background(), a.cfg.ProjectID).CreateDeviceRequest(deviceRequest).Execute()
	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("failed to create device: %w", err)
	}
	if device == nil || device.GetId() == "" {
		return params.ProviderInstance{}, fmt.Errorf("device ID is empty")
	}
	return a.waitDeviceActive(ctx, device.GetId())
}

// GetInstance will return details about one instance.
func (a *equinixProvider) GetInstance(ctx context.Context, instance string) (params.ProviderInstance, error) {
	device, _, err := a.cli.DevicesApi.FindDeviceById(ctx, instance).Execute()
	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("failed to find device: %w", err)
	}
	if device == nil {
		return params.ProviderInstance{}, fmt.Errorf("device not found")
	}
	ret, err := equinixToGarmInstance(*device)
	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("failed to convert device to garm instance: %w", err)
	}
	return ret, nil
}

// ListInstances will list all instances for a provider.
func (a *equinixProvider) ListInstances(ctx context.Context, poolID string) ([]params.ProviderInstance, error) {
	devices, _, err := a.cli.DevicesApi.FindProjectDevices(ctx, a.cfg.ProjectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	ret := []params.ProviderInstance{}
	for _, device := range devices.GetDevices() {
		tags := extractTagsAsMap(device)
		devicePoolID, ok := tags[spec.PoolIDTagName]
		if !ok || devicePoolID != poolID {
			continue
		}

		deviceControllerID, ok := tags[spec.ControllerIDTagName]
		if !ok || deviceControllerID != a.controllerID {
			continue
		}

		instance, err := equinixToGarmInstance(device)
		if err != nil {
			return nil, fmt.Errorf("failed to convert device to garm instance: %w", err)
		}
		ret = append(ret, instance)

	}
	return ret, nil
}

// Delete instance will delete the instance in a provider.
func (a *equinixProvider) DeleteInstance(ctx context.Context, instance string) error {
	if instance == "" {
		return fmt.Errorf("instance ID is empty")
	}

	_, err := uuid.Parse(instance)
	if err != nil {
		instances, err := a.findInstancesByName(ctx, instance)
		if err != nil {
			return fmt.Errorf("failed to find instances by name: %w", err)
		}
		if len(instances) == 0 {
			return nil
		}

		if len(instances) > 1 {
			g, ctx := errgroup.WithContext(ctx)
			for _, inst := range instances {
				inst := inst
				g.Go(func() error {
					return a.deleteOneInstance(ctx, inst.GetId())
				})
			}
			return a.waitForErrorGroupOrContextCancelled(ctx, g)
		}
		instance = instances[0].GetId()
	}
	return a.deleteOneInstance(ctx, instance)
}

// RemoveAllInstances will remove all instances created by this provider.
func (a *equinixProvider) RemoveAllInstances(ctx context.Context) error {
	return nil
}

// Stop shuts down the instance.
func (a *equinixProvider) Stop(ctx context.Context, instance string, force bool) error {
	_, err := a.cli.DevicesApi.PerformAction(ctx, instance).DeviceActionInput(metal.DeviceActionInput{
		Type: metal.DEVICEACTIONINPUTTYPE_POWER_OFF,
	}).Execute()
	if err != nil {
		return fmt.Errorf("failed to stop device: %w", err)
	}
	return nil
}

// Start boots up an instance.
func (a *equinixProvider) Start(ctx context.Context, instance string) error {
	_, err := a.cli.DevicesApi.PerformAction(ctx, instance).DeviceActionInput(metal.DeviceActionInput{
		Type: metal.DEVICEACTIONINPUTTYPE_POWER_ON,
	}).Execute()
	if err != nil {
		return fmt.Errorf("failed to start device: %w", err)
	}
	return nil
}
