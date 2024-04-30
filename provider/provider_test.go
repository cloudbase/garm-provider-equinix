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
	"net/http"
	"testing"

	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-equinix/config"
	"github.com/cloudbase/garm-provider-equinix/internal/spec"
	metal "github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateInstance(t *testing.T) {
	ctx := context.Background()
	cli := new(MockClient)
	a := &equinixProvider{
		cli: cli,
		cfg: &config.Config{
			AuthToken:             "token",
			MetroCode:             "AM",
			HardwareReservationID: nil,
			ProjectID:             "project",
		},
		controllerID: "mock-controller-id",
	}

	tests := []struct {
		name            string
		bootstrapParams params.BootstrapInstance
		device          metal.Device
		expectedOutput  params.ProviderInstance
		errString       string
		err             error
	}{
		{
			name: "valid linux instance",
			bootstrapParams: params.BootstrapInstance{
				Name:          "test-instance",
				InstanceToken: "test-token",
				OSArch:        params.Amd64,
				OSType:        params.Linux,
				Image:         "ubuntu_22_04",
				Flavor:        "c3.small.x86",
				Tools: []params.RunnerApplicationDownload{
					{
						OS:                spec.Ptr("linux"),
						Architecture:      spec.Ptr("x64"),
						DownloadURL:       spec.Ptr("http://test.com"),
						Filename:          spec.Ptr("runner.tar.gz"),
						SHA256Checksum:    spec.Ptr("sha256:1123"),
						TempDownloadToken: spec.Ptr("test-token"),
					},
				},
				ExtraSpecs: []byte(`{"metro_code": "AM"}`),
				PoolID:     "test-pool",
			},
			device: metal.Device{
				Id: spec.Ptr("mock-id"),
				Tags: []string{
					"Name=mock-name",
					"garm-pool-id=test-pool",
					"garm-controller-id=mock-controller-id",
					"OSType=linux",
					"OSArch=amd64",
				},
				IpAddresses: []metal.IPAssignment{
					{
						Address: spec.Ptr("10.10.0.4"),
					},
				},
				State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
			},
			expectedOutput: params.ProviderInstance{
				ProviderID: "mock-id",
				Name:       "mock-name",
				OSType:     params.Linux,
				OSArch:     params.Amd64,
				Addresses: []params.Address{
					{
						Address: "10.10.0.4",
						Type:    "private",
					},
				},
				Status: params.InstanceRunning,
			},
			errString: "",
			err:       nil,
		},
		{
			name: "valid windows instance",
			bootstrapParams: params.BootstrapInstance{
				Name:          "test-instance-windows",
				InstanceToken: "test-token",
				OSArch:        params.Amd64,
				OSType:        params.Windows,
				Image:         "windows-server-2022",
				Flavor:        "c3.small.x86",
				Tools: []params.RunnerApplicationDownload{
					{
						OS:                spec.Ptr("windows"),
						Architecture:      spec.Ptr("x64"),
						DownloadURL:       spec.Ptr("http://test.com"),
						Filename:          spec.Ptr("runner.tar.gz"),
						SHA256Checksum:    spec.Ptr("sha256:1123"),
						TempDownloadToken: spec.Ptr("test-token"),
					},
				},
				ExtraSpecs: []byte(`{"metro_code": "AM"}`),
				PoolID:     "test-pool",
			},
			device: metal.Device{
				Id: spec.Ptr("mock-id"),
				Tags: []string{
					"Name=mock-name",
					"garm-pool-id=test-pool",
					"garm-controller-id=mock-controller-id",
					"OSType=windows",
					"OSArch=amd64",
				},
				IpAddresses: []metal.IPAssignment{
					{
						Address: spec.Ptr("10.10.0.4"),
					},
				},
				State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
			},
			expectedOutput: params.ProviderInstance{
				ProviderID: "mock-id",
				Name:       "mock-name",
				OSType:     params.Windows,
				OSArch:     params.Amd64,
				Addresses: []params.Address{
					{
						Address: "10.10.0.4",
						Type:    "private",
					},
				},
				Status: params.InstanceRunning,
			},
			errString: "",
			err:       nil,
		},
		{
			name: "failed to create device",
			bootstrapParams: params.BootstrapInstance{
				Name:          "test-instance",
				InstanceToken: "test-token",
				OSArch:        params.Amd64,
				OSType:        params.Linux,
				Image:         "ubuntu_22_04",
				Flavor:        "c3.small.x86",
				Tools: []params.RunnerApplicationDownload{
					{
						OS:                spec.Ptr("linux"),
						Architecture:      spec.Ptr("x64"),
						DownloadURL:       spec.Ptr("http://test.com"),
						Filename:          spec.Ptr("runner.tar.gz"),
						SHA256Checksum:    spec.Ptr("sha256:1123"),
						TempDownloadToken: spec.Ptr("test-token"),
					},
				},
				ExtraSpecs: []byte(`{"metro_code": "AM"}`),
				PoolID:     "test-pool",
			},
			device: metal.Device{
				Id: spec.Ptr("mock-id"),
				Tags: []string{
					"Name=mock-name",
					"garm-pool-id=test-pool",
					"garm-controller-id=mock-controller-id",
					"OSType=linux",
					"OSArch=amd64",
				},
				IpAddresses: []metal.IPAssignment{
					{
						Address: spec.Ptr("10.10.0.4"),
					},
				},
			},
			expectedOutput: params.ProviderInstance{},
			errString:      "failed to create device",
			err:            fmt.Errorf("failed to create device"),
		},
		{
			name: "invalid instance",
			bootstrapParams: params.BootstrapInstance{
				OSArch: params.Amd64,
				OSType: params.Linux,
				Tools: []params.RunnerApplicationDownload{
					{},
				},
			},
			device:         metal.Device{},
			expectedOutput: params.ProviderInstance{},
			errString:      "failed to get runner spec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec.DefaultToolFetch = func(osType params.OSType, osArch params.OSArch, tools []params.RunnerApplicationDownload) (params.RunnerApplicationDownload, error) {
				return tt.bootstrapParams.Tools[0], nil
			}
			spec.DefaultGetCloudconfig = func(bootstrapParams params.BootstrapInstance, tools params.RunnerApplicationDownload, runnerName string) (string, error) {
				return "cloudconfig", nil
			}
			cli.On("CreateDevice", ctx, a.cfg.ProjectID).Return(metal.ApiCreateDeviceRequest{
				ApiService: &metal.DevicesApiService{},
			}, nil)
			DefaultExecuteCreateDevice = func(r metal.ApiCreateDeviceRequest) (*metal.Device, *http.Response, error) {
				return &tt.device, &http.Response{StatusCode: http.StatusOK}, tt.err
			}
			cli.On("FindDeviceById", ctx, "mock-id").Return(metal.ApiFindDeviceByIdRequest{
				ApiService: &metal.DevicesApiService{},
			}, nil)
			DefaultExecuteFindDeviceByID = func(r metal.ApiFindDeviceByIdRequest) (*metal.Device, *http.Response, error) {
				return &tt.device, &http.Response{StatusCode: http.StatusOK}, tt.err
			}
			output, err := a.CreateInstance(ctx, tt.bootstrapParams)
			if tt.errString != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestGetInstance(t *testing.T) {
	ctx := context.Background()
	cli := new(MockClient)
	a := &equinixProvider{
		cli: cli,
		cfg: &config.Config{
			AuthToken:             "token",
			MetroCode:             "AM",
			HardwareReservationID: nil,
			ProjectID:             "project",
		},
		controllerID: "mock-controller-id",
	}
	device := metal.Device{
		Id: spec.Ptr("mock-id"),
		Tags: []string{
			"Name=mock-name",
			"garm-pool-id=test-pool",
			"garm-controller-id=mock-controller-id",
			"OSType=linux",
			"OSArch=amd64",
		},
		IpAddresses: []metal.IPAssignment{
			{
				Address: spec.Ptr("10.10.0.4"),
			},
		},
		State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
	}
	cli.On("FindDeviceById", ctx, "mock-id").Return(metal.ApiFindDeviceByIdRequest{
		ApiService: &metal.DevicesApiService{},
	}, nil)
	DefaultExecuteFindDeviceByID = func(r metal.ApiFindDeviceByIdRequest) (*metal.Device, *http.Response, error) {
		return &device, &http.Response{StatusCode: http.StatusOK}, nil
	}
	expectedOutput := params.ProviderInstance{
		ProviderID: "mock-id",
		Name:       "mock-name",
		OSType:     params.Linux,
		OSArch:     params.Amd64,
		Addresses: []params.Address{
			{
				Address: "10.10.0.4",
				Type:    "private",
			},
		},
		Status: params.InstanceRunning,
	}

	output, err := a.GetInstance(ctx, "mock-id")
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestListInstances(t *testing.T) {
	ctx := context.Background()
	cli := new(MockClient)
	poolID := "test-pool"
	a := &equinixProvider{
		cli: cli,
		cfg: &config.Config{
			AuthToken:             "token",
			MetroCode:             "AM",
			HardwareReservationID: nil,
			ProjectID:             "test-pool",
		},
		controllerID: "mock-controller-id",
	}
	devicesList := metal.DeviceList{
		Devices: []metal.Device{
			{
				Id: spec.Ptr("mock-id"),
				Tags: []string{
					"Name=mock-name",
					"garm-pool-id=test-pool",
					"garm-controller-id=mock-controller-id",
					"OSType=linux",
					"OSArch=amd64",
				},
				IpAddresses: []metal.IPAssignment{
					{
						Address: spec.Ptr("10.10.0.4"),
					},
				},
				State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
			},
		},
	}
	cli.On("FindProjectDevices", ctx, a.cfg.ProjectID).Return(
		metal.ApiFindProjectDevicesRequest{
			ApiService: &metal.DevicesApiService{},
		}, nil)
	DefaultExecuteFindProjectDevices = func(r metal.ApiFindProjectDevicesRequest) (*metal.DeviceList, *http.Response, error) {
		return &devicesList, &http.Response{StatusCode: http.StatusOK}, nil
	}
	expectedOutput := params.ProviderInstance{
		ProviderID: "mock-id",
		Name:       "mock-name",
		OSType:     params.Linux,
		OSArch:     params.Amd64,
		Addresses: []params.Address{
			{
				Address: "10.10.0.4",
				Type:    "private",
			},
		},
		Status: params.InstanceRunning,
	}

	output, err := a.ListInstances(ctx, poolID)
	require.NoError(t, err)
	assert.Equal(t, []params.ProviderInstance{expectedOutput}, output)
}

func TestDeleteInstance(t *testing.T) {
	ctx := context.Background()
	cli := new(MockClient)
	instanceID := "mock-id-1"
	a := &equinixProvider{
		cli: cli,
		cfg: &config.Config{
			AuthToken:             "token",
			MetroCode:             "AM",
			HardwareReservationID: nil,
			ProjectID:             "test-pool",
		},
		controllerID: "mock-controller-id",
	}
	devicesList := metal.DeviceList{
		Devices: []metal.Device{
			{
				Id: spec.Ptr("mock-id"),
				Tags: []string{
					"Name=mock-name",
					"garm-pool-id=test-pool",
					"garm-controller-id=mock-controller-id",
					"OSType=linux",
					"OSArch=amd64",
				},
				IpAddresses: []metal.IPAssignment{
					{
						Address: spec.Ptr("10.10.0.4"),
					},
				},
				State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
			},
		},
	}
	cli.On("FindProjectDevices", ctx, "test-pool").Return(
		metal.ApiFindProjectDevicesRequest{
			ApiService: &metal.DevicesApiService{},
		}, nil)
	DefaultExecuteFindProjectDevices = func(r metal.ApiFindProjectDevicesRequest) (*metal.DeviceList, *http.Response, error) {
		return &devicesList, &http.Response{StatusCode: http.StatusOK}, nil
	}
	DefaultExecuteFindDeviceByID = func(r metal.ApiFindDeviceByIdRequest) (*metal.Device, *http.Response, error) {
		return &devicesList.Devices[0], &http.Response{StatusCode: http.StatusOK}, nil
	}
	cli.On("FindDeviceById", ctx, instanceID).Return(metal.ApiFindDeviceByIdRequest{
		ApiService: &metal.DevicesApiService{},
	}, nil)
	cli.On("DeleteDevice", ctx, instanceID).Return(metal.ApiDeleteDeviceRequest{
		ApiService: &metal.DevicesApiService{},
	}, nil)
	DefaultExecuteDeleteDevice = func(r metal.ApiDeleteDeviceRequest) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	err := a.DeleteInstance(ctx, instanceID)
	require.NoError(t, err)
}

func TestStop(t *testing.T) {
	ctx := context.Background()
	cli := new(MockClient)
	instanceID := "mock-id-1"
	a := &equinixProvider{
		cli: cli,
		cfg: &config.Config{
			AuthToken:             "token",
			MetroCode:             "AM",
			HardwareReservationID: nil,
			ProjectID:             "test-pool",
		},
		controllerID: "mock-controller-id",
	}
	cli.On("PerformAction", ctx, instanceID).Return(metal.ApiPerformActionRequest{
		ApiService: &metal.DevicesApiService{},
	}, nil)
	DefaultExecutePerformAction = func(r metal.ApiPerformActionRequest) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	err := a.Stop(ctx, instanceID, true)
	require.NoError(t, err)
}
