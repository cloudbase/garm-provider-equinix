// Copyright 2024 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-equinix/internal/spec"
	"github.com/google/uuid"
	"github.com/juju/clock"
	"github.com/juju/retry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	metal "github.com/equinix/equinix-sdk-go/services/metalv1"
)

var (
	errStopRetry = errors.New("stop retry")
)

var statusMap = map[metal.DeviceState]params.InstanceStatus{
	metal.DEVICESTATE_QUEUED:       params.InstanceRunning,
	metal.DEVICESTATE_PROVISIONING: params.InstanceRunning,
	metal.DEVICESTATE_ACTIVE:       params.InstanceRunning,
	metal.DEVICESTATE_INACTIVE:     params.InstanceStopped,
	metal.DEVICESTATE_FAILED:       params.InstanceError,
	metal.DEVICESTATE_DELETED:      params.InstanceStopped,
	metal.DEVICESTATE_POWERING_ON:  params.InstanceRunning,
	metal.DEVICESTATE_POWERING_OFF: params.InstanceStopped,
	metal.DeviceState(""):          params.InstanceStatusUnknown,
}

func equinixToGarmInstance(device metal.Device) (params.ProviderInstance, error) {
	tags := extractTagsAsMap(device)
	name, ok := tags["Name"]
	if !ok {
		return params.ProviderInstance{}, fmt.Errorf("missing Name property")
	}

	instance := params.ProviderInstance{
		ProviderID: device.GetId(),
		Name:       name,
		Status:     statusMap[device.GetState()],
	}

	if instance.ProviderID == "" {
		return params.ProviderInstance{}, fmt.Errorf("device ID is empty")
	}

	for key, val := range tags {
		switch key {
		case "OSType":
			instance.OSType = params.OSType(val)
		case "OSArch":
			instance.OSArch = params.OSArch(val)
		case "Name":
			instance.Name = val
		}
	}

	for _, address := range device.GetIpAddresses() {
		addrType := params.PrivateAddress
		if address.GetPublic() {
			addrType = params.PublicAddress
		}
		addr := params.Address{
			Address: address.GetAddress(),
			Type:    addrType,
		}
		instance.Addresses = append(instance.Addresses, addr)
	}
	return instance, nil
}

func (a *equinixProvider) waitDeviceActive(ctx context.Context, deviceID string) (params.ProviderInstance, error) {
	var p params.ProviderInstance
	err := retry.Call(retry.CallArgs{
		IsFatalError: func(err error) bool {
			return errors.Is(err, errStopRetry)
		},
		Func: func() error {
			var err error
			device, _, err := a.cli.DevicesApi.FindDeviceById(ctx, deviceID).Execute()
			if err != nil {
				return fmt.Errorf("failed to find device: %w", errStopRetry)
			}
			if device == nil {
				return fmt.Errorf("device not found: %w", errStopRetry)
			}

			state := device.GetState()
			switch state {
			case metal.DEVICESTATE_FAILED:
				return fmt.Errorf("device failed: %w", errStopRetry)
			case metal.DEVICESTATE_DELETED:
				return fmt.Errorf("device deleted: %w", errStopRetry)
			case metal.DEVICESTATE_POWERING_OFF, metal.DEVICESTATE_INACTIVE:
				return fmt.Errorf("invalid state change: %w", errStopRetry)
			case metal.DEVICESTATE_ACTIVE:
				p, err = equinixToGarmInstance(*device)
				if err != nil {
					return fmt.Errorf("failed to convert device: %w", errStopRetry)
				}
				return nil
			}

			// Return potentially a lot earlier than the instance is active. Equinix seems to wait for
			// user-data scripts to finish running before it phones home. This might cause status updates
			// to arive while the instance is still in "creating" state.
			// This is currently allowed, although it might change in the future.
			if device.GetProvisioningPercentage() >= 90 {
				p, err = equinixToGarmInstance(*device)
				if err != nil {
					return fmt.Errorf("failed to convert device: %w", errStopRetry)
				}
				return nil
			}
			return fmt.Errorf("instance not active yet")
		},
		// Roughly 20 minutes. Might be longer if API calls take longer.
		Attempts: 240,
		Delay:    5 * time.Second,
		Clock:    clock.WallClock,
	})

	if err != nil {
		return params.ProviderInstance{}, fmt.Errorf("failed to wait for instance to become active: %w", err)
	}
	return p, nil
}

func (a *equinixProvider) findInstancesByName(ctx context.Context, instance string) ([]metal.Device, error) {
	ret := []metal.Device{}

	devices, _, err := a.cli.DevicesApi.FindProjectDevices(ctx, a.cfg.ProjectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	for _, dev := range devices.Devices {
		tags := extractTagsAsMap(dev)
		name, ok := tags["Name"]
		if !ok {
			continue
		}
		controllerID, ok := tags[spec.ControllerIDTagName]
		if !ok {
			continue
		}
		if controllerID == a.controllerID && name == instance {
			ret = append(ret, dev)
		}
	}

	return ret, nil
}

func extractTagsAsMap(device metal.Device) map[string]string {
	ret := map[string]string{}
	for _, tag := range device.GetTags() {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 {
			continue
		}
		ret[parts[0]] = parts[1]
	}
	return ret
}

func (a *equinixProvider) waitForErrorGroupOrContextCancelled(ctx context.Context, g *errgroup.Group) error {
	if g == nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		waitErr := g.Wait()
		done <- waitErr
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *equinixProvider) deleteOneInstance(ctx context.Context, instanceID string) error {
	_, err := uuid.Parse(instanceID)
	if err != nil {
		return fmt.Errorf("invalid instance ID %s: %w", instanceID, err)
	}
	device, resp, err := a.cli.DevicesApi.FindDeviceById(ctx, instanceID).Execute()
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
			return nil
		}
		return fmt.Errorf("failed to find device: %w", err)
	}
	state := device.GetState()
	if state == metal.DEVICESTATE_PROVISIONING || state == metal.DEVICESTATE_QUEUED {
		if _, err := a.waitDeviceActive(ctx, instanceID); err != nil {
			return fmt.Errorf("failed to wait for device: %w", err)
		}
	}
	resp, err = a.cli.DevicesApi.DeleteDevice(ctx, instanceID).Execute()
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
			return nil
		}
		if state == metal.DEVICESTATE_DELETED || state == metal.DEVICESTATE_FAILED {
			return nil
		}
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return nil
}
