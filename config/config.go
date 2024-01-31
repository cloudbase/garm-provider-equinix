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

package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func NewConfig(cfgFile string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(cfgFile, &config); err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}
	return &config, nil

}

type Config struct {
	// AuthToken is the authentication token for the Equinix Metal API.
	AuthToken string `toml:"auth_token"`
	// MetroCode is the metro (usually a two letter code) to use for the instance.
	// See: https://deploy.equinix.com/developers/docs/metal/locations/metros/
	MetroCode string `toml:"metro_code"`
	// HardwareReservationID is the UUID representing the hardware reservation to use.
	HardwareReservationID *string `toml:"hardware_reservation_id,omitempty"`
	// ProjectID is the UUID representing the project to use.
	ProjectID string `toml:"project_id"`
}

func (c *Config) Validate() error {
	if c.AuthToken == "" {
		return fmt.Errorf("auth_token is required")
	}

	if c.MetroCode == "" {
		return fmt.Errorf("metro_code is required")
	}

	if c.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}
	return nil
}
