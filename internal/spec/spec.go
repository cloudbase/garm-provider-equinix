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
	"fmt"

	"github.com/cloudbase/garm-provider-common/cloudconfig"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm-provider-common/util"
	"github.com/xeipuuv/gojsonschema"
)

const (
	ControllerIDTagName = "garm-controller-id"
	PoolIDTagName       = "garm-pool-id"
)

const jsonSchema string = `
	{
		"$schema": "http://cloudbase.it/garm-provider-equinix/schemas/extra_specs#",
		"type": "object",
		"description": "Schema defining supported extra specs for the Garm Equinix Metal Provider",
		"properties": {
			"metro_code": {
				"type": "string",
				"description": "The metro in which this pool will create runners."
			},
			"hardware_reservation_id": {
				"type": "string",
				"description": "The hardware reservation ID to use."
			}
		},
		"additionalProperties": false
	}
`

func jsonSchemaValidation(schema json.RawMessage) error {
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
	extraSpecsLoader := gojsonschema.NewBytesLoader(schema)
	result, err := gojsonschema.Validate(schemaLoader, extraSpecsLoader)
	if err != nil {
		return fmt.Errorf("failed to validate schema: %w", err)
	}
	if !result.Valid() {
		return fmt.Errorf("schema validation failed: %s", result.Errors())
	}
	return nil
}

func newExtraSpecsFromBootstrapData(data params.BootstrapInstance) (extraSpecs, error) {
	spec := extraSpecs{}

	if len(data.ExtraSpecs) > 0 {
		if err := jsonSchemaValidation(data.ExtraSpecs); err != nil {
			return extraSpecs{}, fmt.Errorf("failed to validate extra specs: %w", err)
		}

		if err := json.Unmarshal(data.ExtraSpecs, &spec); err != nil {
			return extraSpecs{}, fmt.Errorf("failed to unmarshal extra specs: %w", err)
		}
	}

	return spec, nil
}

type extraSpecs struct {
	// MetroCode is the metro (usually a two letter code) to use for the instance.
	// See: https://deploy.equinix.com/developers/docs/metal/locations/metros/
	MetroCode string `toml:"metro_code"`
	// HardwareReservationID is the UUID representing the hardware reservation to use.
	HardwareReservationID *string `toml:"hardware_reservation_id,omitempty"`
}

func GetRunnerSpecFromBootstrapParams(data params.BootstrapInstance, controllerID string) (*RunnerSpec, error) {
	tools, err := util.GetTools(data.OSType, data.OSArch, data.Tools)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools: %s", err)
	}

	extraSpecs, err := newExtraSpecsFromBootstrapData(data)
	if err != nil {
		return nil, fmt.Errorf("error loading extra specs: %w", err)
	}

	tags := []string{
		fmt.Sprintf("%s=%s", PoolIDTagName, data.PoolID),
		fmt.Sprintf("%s=%s", ControllerIDTagName, controllerID),
		fmt.Sprintf("OSArch=%s", data.OSArch),
		fmt.Sprintf("Name=%s", data.Name),
	}

	spec := &RunnerSpec{
		BootstrapParams: data,
		Tools:           tools,
		Tags:            tags,
	}
	spec.MergeExtraSpecs(extraSpecs)

	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("error validating spec: %w", err)
	}

	return spec, nil
}

type RunnerSpec struct {
	ProjectID             string
	MetroCode             string
	HardwareReservationID *string
	Tools                 params.RunnerApplicationDownload
	Tags                  []string
	BootstrapParams       params.BootstrapInstance
}

func (r RunnerSpec) Validate() error {
	if r.Tools.DownloadURL == nil {
		return fmt.Errorf("missing tools")
	}

	if r.BootstrapParams.Name == "" || r.BootstrapParams.OSType == "" || r.BootstrapParams.InstanceToken == "" {
		return fmt.Errorf("invalid bootstrap params")
	}

	return nil
}

func (r *RunnerSpec) MergeExtraSpecs(spec extraSpecs) {
	if spec.HardwareReservationID != nil {
		r.HardwareReservationID = spec.HardwareReservationID
	}

	if spec.MetroCode != "" {
		r.MetroCode = spec.MetroCode
	}
}

func (r *RunnerSpec) ComposeUserData() (string, error) {
	switch r.BootstrapParams.OSType {
	case params.Linux, params.Windows:
	default:
		return "", fmt.Errorf("unsupported OS type for cloud config: %s", r.BootstrapParams.OSType)
	}

	udata, err := cloudconfig.GetCloudConfig(r.BootstrapParams, r.Tools, r.BootstrapParams.Name)
	if err != nil {
		return "", fmt.Errorf("failed to generate userdata: %w", err)
	}
	if r.BootstrapParams.OSType == params.Windows {
		udata = fmt.Sprintf("#ps1_sysnative\n%s", udata)
	}
	return udata, nil
}
