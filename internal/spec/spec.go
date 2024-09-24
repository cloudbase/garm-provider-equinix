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
	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

const (
	ControllerIDTagName = "garm-controller-id"
	PoolIDTagName       = "garm-pool-id"
)

type ToolFetchFunc func(osType params.OSType, osArch params.OSArch, tools []params.RunnerApplicationDownload) (params.RunnerApplicationDownload, error)

type GetCloudConfigFunc func(bootstrapParams params.BootstrapInstance, tools params.RunnerApplicationDownload, runnerName string) (string, error)

var (
	DefaultToolFetch      ToolFetchFunc      = util.GetTools
	DefaultGetCloudconfig GetCloudConfigFunc = cloudconfig.GetCloudConfig
)

func generateJSONSchema() *jsonschema.Schema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
	}
	// Reflect the extraSpecs struct
	schema := reflector.Reflect(extraSpecs{})

	return schema
}

func jsonSchemaValidation(schema json.RawMessage) error {
	jsonSchema := generateJSONSchema()
	schemaLoader := gojsonschema.NewGoLoader(jsonSchema)
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
	MetroCode string `json:"metro_code,omitempty" jsonschema:"description=The metro in which this pool will create runners."`
	// HardwareReservationID is the UUID representing the hardware reservation to use.
	HardwareReservationID *string  `json:"hardware_reservation_id,omitempty" jsonschema:"description=The hardware reservation ID to use for the runner."`
	DisableUpdates        *bool    `json:"disable_updates,omitempty" jsonschema:"description=Disable automatic updates on the VM."`
	EnableBootDebug       *bool    `json:"enable_boot_debug,omitempty" jsonschema:"description=Enable boot debug on the VM."`
	ExtraPackages         []string `json:"extra_packages,omitempty" jsonschema:"description=Extra packages to install on the VM."`
	// The Cloudconfig struct from common package
	cloudconfig.CloudConfigSpec
}

func GetRunnerSpecFromBootstrapParams(data params.BootstrapInstance, controllerID string) (*RunnerSpec, error) {
	tools, err := DefaultToolFetch(data.OSType, data.OSArch, data.Tools)
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
		ExtraPackages:   extraSpecs.ExtraPackages,
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
	DisableUpdates        bool
	ExtraPackages         []string
	EnableBootDebug       bool
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

	if spec.DisableUpdates != nil {
		r.DisableUpdates = *spec.DisableUpdates
	}

	if spec.EnableBootDebug != nil {
		r.EnableBootDebug = *spec.EnableBootDebug
	}
}

func (r *RunnerSpec) ComposeUserData() (string, error) {
	bootstrapParams := r.BootstrapParams
	bootstrapParams.UserDataOptions.DisableUpdatesOnBoot = r.DisableUpdates
	bootstrapParams.UserDataOptions.ExtraPackages = r.ExtraPackages
	bootstrapParams.UserDataOptions.EnableBootDebug = r.EnableBootDebug
	switch bootstrapParams.OSType {
	case params.Linux, params.Windows:
	default:
		return "", fmt.Errorf("unsupported OS type for cloud config: %s", bootstrapParams.OSType)
	}

	udata, err := DefaultGetCloudconfig(bootstrapParams, r.Tools, bootstrapParams.Name)
	if err != nil {
		return "", fmt.Errorf("failed to generate userdata: %w", err)
	}
	if bootstrapParams.OSType == params.Windows {
		udata = fmt.Sprintf("#ps1_sysnative\n%s", udata)
	}
	return udata, nil
}

func Ptr[T any](v T) *T {
	return &v
}
