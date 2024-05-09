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
	spec "github.com/cloudbase/garm-provider-equinix/internal/spec"
	metal "github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEquinixToGarmInstance(t *testing.T) {
	deviceID := "mock-id"
	tests := []struct {
		name           string
		device         metal.Device
		expectedOutput params.ProviderInstance
		errString      string
	}{
		{
			name: "valid device",
			device: metal.Device{
				Id: spec.Ptr("mock-id"),
				Tags: []string{
					"Name=mock-name",
					"OSType=ubuntu",
					"OSArch=amd64",
				},
				IpAddresses: []metal.IPAssignment{
					{
						Address: spec.Ptr("100.10.0.4"),
						Network: spec.Ptr("public"),
						Public:  spec.Ptr(true),
					},
					{
						Address: spec.Ptr("10.10.0.4"),
						Network: spec.Ptr("private"),
					},
				},
				State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
			},
			expectedOutput: params.ProviderInstance{
				ProviderID: deviceID,
				Name:       "mock-name",
				Addresses: []params.Address{
					{
						Address: "100.10.0.4",
						Type:    "public",
					},
					{
						Address: "10.10.0.4",
						Type:    "private",
					},
				},
				Status: params.InstanceRunning,
				OSType: params.OSType("ubuntu"),
				OSArch: params.OSArch("amd64"),
			},
			errString: "",
		},
		{
			name: "missing Name tag",
			device: metal.Device{
				Id: spec.Ptr("mock-id"),
			},
			expectedOutput: params.ProviderInstance{},
			errString:      "missing Name property",
		},
		{
			name: "empty device ID",
			device: metal.Device{
				Id: spec.Ptr(""),
			},
			expectedOutput: params.ProviderInstance{},
			errString:      "device ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, err := equinixToGarmInstance(tt.device)
			if tt.errString != "" {
				require.Error(t, err, tt.errString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, instance)
		})
	}
}

func TestWaitDeviceActive(t *testing.T) {
	ctx := context.Background()
	deviceID := "mock-id"
	device := metal.Device{
		Id: spec.Ptr(deviceID),
		Tags: []string{
			"Name=mock-name",
		},
		IpAddresses: []metal.IPAssignment{
			{
				Address: spec.Ptr("10.10.0.4"),
			},
		},
		State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
	}
	DefaultExecuteFindDeviceByID = func(r metal.ApiFindDeviceByIdRequest) (*metal.Device, *http.Response, error) {
		return &device, &http.Response{StatusCode: http.StatusOK}, nil
	}
	cli := new(MockClient)
	a := &equinixProvider{
		cli:          cli,
		cfg:          &config.Config{},
		controllerID: "mock-controller-id",
	}
	expectedOutput := params.ProviderInstance{
		ProviderID: "mock-id",
		Name:       "mock-name",
		Addresses: []params.Address{
			{
				Address: "10.10.0.4",
				Type:    "private",
			},
		},
		Status: params.InstanceRunning,
	}

	cli.On("FindDeviceById", ctx, deviceID).Return(metal.ApiFindDeviceByIdRequest{
		ApiService: &metal.DevicesApiService{},
	}, nil)
	output, err := a.waitDeviceActive(ctx, deviceID)
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestWaitDeviceActiveFails(t *testing.T) {
	ctx := context.Background()
	deviceID := "mock-id"
	cli := new(MockClient)
	a := &equinixProvider{
		cli:          cli,
		cfg:          &config.Config{},
		controllerID: "mock-controller-id",
	}
	tests := []struct {
		name      string
		device    *metal.Device
		errString string
	}{
		{
			name:      "empty device ID",
			device:    nil,
			errString: "device not found",
		},
		{
			name: "failed device",
			device: &metal.Device{
				State: spec.Ptr(metal.DEVICESTATE_FAILED),
			},
			errString: "device failed",
		},
		{
			name: "device deleted",
			device: &metal.Device{
				State: spec.Ptr(metal.DEVICESTATE_DELETED),
			},
			errString: "device deleted",
		},
		{
			name: "invalid state change",
			device: &metal.Device{
				State: spec.Ptr(metal.DEVICESTATE_POWERING_OFF),
			},
			errString: "invalid state change",
		},
		{
			name: "failed to convert device",
			device: &metal.Device{
				State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
			},
			errString: "failed to convert device",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DefaultExecuteFindDeviceByID = func(r metal.ApiFindDeviceByIdRequest) (*metal.Device, *http.Response, error) {
				return tt.device, &http.Response{StatusCode: http.StatusOK}, nil
			}
			cli.On("FindDeviceById", ctx, deviceID).Return(metal.ApiFindDeviceByIdRequest{
				ApiService: &metal.DevicesApiService{},
			}, nil)
			_, err := a.waitDeviceActive(ctx, deviceID)
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.errString)
		})
	}
}

func TestFindInstancesByName(t *testing.T) {
	ctx := context.Background()
	instanceName := "test-instance"
	cli := new(MockClient)
	a := &equinixProvider{
		cli: cli,
		cfg: &config.Config{
			ProjectID: "mock-project-id",
		},
		controllerID: "mock-controller-id",
	}
	devicesList := metal.DeviceList{
		Devices: []metal.Device{
			{
				Id: spec.Ptr("mock-id-1"),
				Tags: []string{
					"Name=test-instance",
					"garm-controller-id=mock-controller-id",
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
	cli.On("FindProjectDevices", ctx, "mock-project-id").Return(
		metal.ApiFindProjectDevicesRequest{
			ApiService: &metal.DevicesApiService{},
		}, nil)
	DefaultExecuteFindProjectDevices = func(r metal.ApiFindProjectDevicesRequest) (*metal.DeviceList, *http.Response, error) {
		return &devicesList, &http.Response{StatusCode: http.StatusOK}, nil
	}

	output, err := a.findInstancesByName(ctx, instanceName)
	require.NoError(t, err)
	assert.Equal(t, devicesList.Devices, output)
}

func TestDeleteOneInstance(t *testing.T) {
	ctx := context.Background()
	instanceID := "76e33e9e-6155-472e-ae76-37b5401f888f"
	device := metal.Device{
		Id: spec.Ptr(instanceID),
		Tags: []string{
			"Name=mock-name",
		},
		IpAddresses: []metal.IPAssignment{
			{
				Address: spec.Ptr("10.10.0.4"),
			},
		},
		State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
	}
	DefaultExecuteFindDeviceByID = func(r metal.ApiFindDeviceByIdRequest) (*metal.Device, *http.Response, error) {
		return &device, &http.Response{StatusCode: http.StatusOK}, nil
	}
	cli := new(MockClient)
	a := &equinixProvider{
		cli:          cli,
		cfg:          &config.Config{},
		controllerID: "mock-controller-id",
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

	err := a.deleteOneInstance(ctx, instanceID)
	require.NoError(t, err)
}

func TestDeleteOneInstanceErrors(t *testing.T) {
	ctx := context.Background()
	cli := new(MockClient)
	a := &equinixProvider{
		cli:          cli,
		cfg:          &config.Config{},
		controllerID: "mock-controller-id",
	}
	tests := []struct {
		name          string
		deviceID      string
		httpsResponse http.Response
		errString     string
		err           error
	}{
		{
			name:      "empty device ID",
			deviceID:  "",
			errString: "invalid instance ID",
		},
		{
			name:     "failed to get device",
			deviceID: "76e33e9e-6155-472e-ae76-37b5401f888f",
			httpsResponse: http.Response{
				StatusCode: 500,
			},
			errString: "failed to find device",
			err:       fmt.Errorf("failed to find device"),
		},
		{
			name:     "failed to delete device",
			deviceID: "76e33e9e-6155-472e-ae76-37b5401f888f",
			httpsResponse: http.Response{
				StatusCode: 200,
			},
			errString: "failed to delete device",
			err:       fmt.Errorf("failed to delete device"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := metal.Device{
				Id: spec.Ptr(tt.deviceID),
				Tags: []string{
					"Name=mock-name",
				},
				IpAddresses: []metal.IPAssignment{
					{
						Address: spec.Ptr("10.10.0.4"),
					},
				},
				State: spec.Ptr(metal.DEVICESTATE_ACTIVE),
			}
			DefaultExecuteFindDeviceByID = func(r metal.ApiFindDeviceByIdRequest) (*metal.Device, *http.Response, error) {
				return &device, &tt.httpsResponse, tt.err
			}
			DefaultExecuteDeleteDevice = func(r metal.ApiDeleteDeviceRequest) (*http.Response, error) {
				return &tt.httpsResponse, tt.err
			}
			cli.On("FindDeviceById", ctx, tt.deviceID).Return(metal.ApiFindDeviceByIdRequest{
				ApiService: &metal.DevicesApiService{},
			}, nil)
			cli.On("DeleteDevice", ctx, tt.deviceID).Return(metal.ApiDeleteDeviceRequest{
				ApiService: &metal.DevicesApiService{},
			}, fmt.Errorf("failed to delete device"))
			err := a.deleteOneInstance(ctx, tt.deviceID)
			if tt.errString != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
