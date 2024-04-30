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

package spec

import (
	"encoding/json"
	"testing"

	"github.com/cloudbase/garm-provider-common/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonSchemaValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     json.RawMessage
		errString string
	}{
		{
			name: "Valid input",
			input: json.RawMessage(`{
				"metro_code": "AM",
				"hardware_reservation_id": "id"
			}`),
			errString: "",
		},
		{
			name: "Invalid input - wrong data type",
			input: json.RawMessage(`{
				"metro_code": true,
				"hardware_reservation_id": "id"
			}`),
			errString: "schema validation failed: [metro_code: Invalid type. Expected: string, given: boolean]",
		},
		{
			name: "Invalid input - additional property",
			input: json.RawMessage(`{
				"additional_property": true
			}`),
			errString: "Additional property additional_property is not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonSchemaValidation(tt.input)
			if tt.errString == "" {
				assert.NoError(t, err, "Expected no error, got %v", err)
			} else {
				assert.Error(t, err, "Expected an error")
				if err != nil {
					assert.Contains(t, err.Error(), tt.errString, "Error message does not match")
				}
			}
		})
	}
}

func TestNewExtraSpecsFromBootstrapData(t *testing.T) {
	tests := []struct {
		name           string
		specs          params.BootstrapInstance
		expectedOutput extraSpecs
		errString      string
	}{
		{
			name: "Empty specs",
			specs: params.BootstrapInstance{
				ExtraSpecs: nil,
			},
			expectedOutput: extraSpecs{},
			errString:      "",
		},
		{
			name: "Valid specs",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"metro_code": "AM", "hardware_reservation_id": "id"}`),
			},
			expectedOutput: extraSpecs{
				MetroCode:             "AM",
				HardwareReservationID: Ptr("id"),
			},
			errString: "",
		},
		{
			name: "Invalid specs",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"metro_code": true}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "metro_code: Invalid type. Expected: string, given: boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := newExtraSpecsFromBootstrapData(tt.specs)
			if tt.errString != "" {
				assert.ErrorContains(t, err, tt.errString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		spec      RunnerSpec
		errString string
	}{
		{
			name: "Valid spec",
			spec: RunnerSpec{
				BootstrapParams: params.BootstrapInstance{
					Name:          "name",
					OSType:        "os",
					InstanceToken: "token",
				},
				Tools: params.RunnerApplicationDownload{
					DownloadURL: Ptr("url"),
				},
			},
			errString: "",
		},
		{
			name: "Missing tools",
			spec: RunnerSpec{
				BootstrapParams: params.BootstrapInstance{
					Name:          "name",
					OSType:        "os",
					InstanceToken: "token",
				},
				Tools: params.RunnerApplicationDownload{},
			},
			errString: "missing tools",
		},
		{
			name: "Missing bootstrap params",
			spec: RunnerSpec{
				BootstrapParams: params.BootstrapInstance{},
				Tools: params.RunnerApplicationDownload{
					DownloadURL: Ptr("url"),
				},
			},
			errString: "invalid bootstrap params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate()
			if tt.errString != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetRunnerSpecFromBootstrapParams(t *testing.T) {
	bootstrapParams := params.BootstrapInstance{
		Name:          "test-instance",
		InstanceToken: "test-token",
		OSArch:        params.Amd64,
		OSType:        params.Linux,
		Image:         "ubuntu_22_04",
		Flavor:        "c3.small.x86",
		Tools: []params.RunnerApplicationDownload{
			{
				OS:                Ptr("linux"),
				Architecture:      Ptr("x64"),
				DownloadURL:       Ptr("http://test.com"),
				Filename:          Ptr("runner.tar.gz"),
				SHA256Checksum:    Ptr("sha256:1123"),
				TempDownloadToken: Ptr("test-token"),
			},
		},
		ExtraSpecs: []byte(`{"metro_code": "AM"}`),
		PoolID:     "test-pool",
	}
	controllerID := "test-controller"
	DefaultToolFetch = func(osType params.OSType, osArch params.OSArch, tools []params.RunnerApplicationDownload) (params.RunnerApplicationDownload, error) {
		return bootstrapParams.Tools[0], nil
	}

	expectedOutput := RunnerSpec{
		MetroCode:       "AM",
		BootstrapParams: bootstrapParams,
		Tools:           bootstrapParams.Tools[0],
		Tags: []string{
			"garm-pool-id=test-pool",
			"garm-controller-id=test-controller",
			"OSArch=amd64",
			"Name=test-instance",
		},
	}

	output, err := GetRunnerSpecFromBootstrapParams(bootstrapParams, controllerID)
	require.NoError(t, err)
	assert.Equal(t, expectedOutput, *output)
}

func TestComposeUserData(t *testing.T) {
	spec := RunnerSpec{
		BootstrapParams: params.BootstrapInstance{
			Name:          "test-instance",
			InstanceToken: "test-token",
			OSArch:        params.Amd64,
			Flavor:        "c3.small.x86",
			Tools: []params.RunnerApplicationDownload{
				{
					OS:                Ptr("linux"),
					Architecture:      Ptr("x64"),
					DownloadURL:       Ptr("http://test.com"),
					Filename:          Ptr("runner.tar.gz"),
					SHA256Checksum:    Ptr("sha256:1123"),
					TempDownloadToken: Ptr("test-token"),
				},
			},
			ExtraSpecs: []byte(`{"metro_code": "AM"}`),
			PoolID:     "test-pool",
		},
		MetroCode: "AM",
		Tags: []string{
			"garm-pool-id=test-pool",
			"garm-controller-id=test-controller",
			"OSArch=amd64",
			"Name=test-instance",
		},
	}
	DefaultGetCloudconfig = func(bootstrapParams params.BootstrapInstance, tools params.RunnerApplicationDownload, runnerName string) (string, error) {
		return "cloudconfig", nil
	}
	tests := []struct {
		name           string
		imageType      params.OSType
		expectedOutput string
		errString      string
	}{
		{
			name:           "Valid spec linux",
			imageType:      params.Linux,
			expectedOutput: "cloudconfig",
			errString:      "",
		},
		{
			name:           "Valid spec windows",
			imageType:      params.Windows,
			expectedOutput: "#ps1_sysnative\ncloudconfig",
			errString:      "",
		},
		{
			name:           "Invalid image type",
			imageType:      params.Unknown,
			expectedOutput: "",
			errString:      "unsupported OS type for cloud config: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec.BootstrapParams.OSType = tt.imageType
			out, err := spec.ComposeUserData()
			if tt.errString != "" {
				assert.ErrorContains(t, err, tt.errString)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, out)
			}
		})
	}
}
