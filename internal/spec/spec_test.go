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
	"testing"

	"github.com/cloudbase/garm-provider-common/cloudconfig"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExtraSpecsFromBootstrapData(t *testing.T) {
	tests := []struct {
		name           string
		specs          params.BootstrapInstance
		expectedOutput extraSpecs
		errString      string
	}{
		{
			name: "full specs",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"metro_code": "AM",
				"hardware_reservation_id": "hw-res-id",
				"disable_updates": true,
				"enable_boot_debug": false,
				"extra_packages": ["package1", "package2"],
				"runner_install_template": "IyEvYmluL2Jhc2gKZWNobyBJbnN0YWxsaW5nIHJ1bm5lci4uLg==",
				"pre_install_scripts": {"setup.sh": "IyEvYmluL2Jhc2gKZWNobyBTZXR1cCBzY3JpcHQuLi4="},
				"extra_context": {"key": "value"}
				}`),
			},
			expectedOutput: extraSpecs{
				MetroCode:             "AM",
				HardwareReservationID: Ptr("hw-res-id"),
				DisableUpdates:        Ptr(true),
				EnableBootDebug:       Ptr(false),
				ExtraPackages:         []string{"package1", "package2"},
				CloudConfigSpec: cloudconfig.CloudConfigSpec{
					RunnerInstallTemplate: []byte("#!/bin/bash\necho Installing runner..."),
					PreInstallScripts: map[string][]byte{
						"setup.sh": []byte("#!/bin/bash\necho Setup script..."),
					},
					ExtraContext: map[string]string{"key": "value"},
				},
			},
			errString: "",
		},
		{
			name: "specs just with metro code",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"metro_code": "AM"}`),
			},
			expectedOutput: extraSpecs{
				MetroCode: "AM",
			},
			errString: "",
		},
		{
			name: "specs just with HardwareReservationID",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"hardware_reservation_id": "hw-res-id"}`),
			},
			expectedOutput: extraSpecs{
				HardwareReservationID: Ptr("hw-res-id"),
			},
			errString: "",
		},
		{
			name: "specs just with DisableUpdates",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"disable_updates": true}`),
			},
			expectedOutput: extraSpecs{
				DisableUpdates: Ptr(true),
			},
			errString: "",
		},
		{
			name: "specs just with EnableBootDebug",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"enable_boot_debug": false}`),
			},
			expectedOutput: extraSpecs{
				EnableBootDebug: Ptr(false),
			},
			errString: "",
		},
		{
			name: "specs just with ExtraPackages",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"extra_packages": ["package1", "package2"]}`),
			},
			expectedOutput: extraSpecs{
				ExtraPackages: []string{"package1", "package2"},
			},
			errString: "",
		},
		{
			name: "specs just with RunnerInstallTemplate",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"runner_install_template": "IyEvYmluL2Jhc2gKZWNobyBJbnN0YWxsaW5nIHJ1bm5lci4uLg=="}`),
			},
			expectedOutput: extraSpecs{
				CloudConfigSpec: cloudconfig.CloudConfigSpec{
					RunnerInstallTemplate: []byte("#!/bin/bash\necho Installing runner..."),
				},
			},
			errString: "",
		},
		{
			name: "specs just with PreInstallScripts",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"pre_install_scripts": {"setup.sh": "IyEvYmluL2Jhc2gKZWNobyBTZXR1cCBzY3JpcHQuLi4="}}`),
			},
			expectedOutput: extraSpecs{
				CloudConfigSpec: cloudconfig.CloudConfigSpec{
					PreInstallScripts: map[string][]byte{
						"setup.sh": []byte("#!/bin/bash\necho Setup script..."),
					},
				},
			},
			errString: "",
		},
		{
			name: "specs just with ExtraContext",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"extra_context": {"key": "value"}}`),
			},
			expectedOutput: extraSpecs{
				CloudConfigSpec: cloudconfig.CloudConfigSpec{
					ExtraContext: map[string]string{"key": "value"},
				},
			},
			errString: "",
		},
		{
			name: "Empty specs",
			specs: params.BootstrapInstance{
				ExtraSpecs: nil,
			},
			expectedOutput: extraSpecs{},
			errString:      "",
		},
		{
			name: "invalid json specs",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"extra_context": }`),
			},
			expectedOutput: extraSpecs{},
			errString:      "failed to validate extra specs",
		},
		{
			name: "invalid input for metro code - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"metro_code": 1}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "metro_code: Invalid type. Expected: string, given: integer",
		},
		{
			name: "invalid input for hardware reservation id - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"hardware_reservation_id": 1}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "hardware_reservation_id: Invalid type. Expected: string, given: integer",
		},
		{
			name: "invalid input for disable updates - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"disable_updates": "true"}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "disable_updates: Invalid type. Expected: boolean, given: string",
		},
		{
			name: "invalid input for enable boot debug - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"enable_boot_debug": "false"}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "enable_boot_debug: Invalid type. Expected: boolean, given: string",
		},
		{
			name: "invalid input for extra packages - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"extra_packages": "package1"}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "extra_packages: Invalid type. Expected: array, given: string",
		},
		{
			name: "invalid input for runner install template - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"runner_install_template": 1}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "runner_install_template: Invalid type. Expected: string, given: integer",
		},
		{
			name: "invalid input for pre install scripts - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"pre_install_scripts": "setup.sh"}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "pre_install_scripts: Invalid type. Expected: object, given: string",
		},
		{
			name: "invalid input for extra context - wrong data type",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"extra_context": "key"}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "extra_context: Invalid type. Expected: object, given: string",
		},
		{
			name: "invalid input - additional property",
			specs: params.BootstrapInstance{
				ExtraSpecs: []byte(`{"additional_property": "key"}`),
			},
			expectedOutput: extraSpecs{},
			errString:      "Additional property additional_property is not allowed",
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
