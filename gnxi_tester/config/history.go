/* Copyright 2020 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Device stores connection details of a target.
type Device struct {
	Address string `mapstructure:"address"`
	Ca      string `mapstructure:"ca"`
	CaKey   string `mapstructure:"cakey"`
}

// SetTarget adds any new target to the list of known targets.
func SetTarget(targetName, targetAddress, ca, caKey string) error {
	if err := prepareTarget(targetName, targetAddress, ca, caKey); err != nil {
		return err
	}
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

// prepareTarget parses provided details and creates or modifies target entry.
func prepareTarget(targetName, targetAddress, ca, caKey string) error {
	devices := GetDevices()
	if devices == nil {
		devices = map[string]Device{}
	}
	if targetName == "" {
		if len(devices) > 0 {
			return nil
		}
		return errors.New("No targets in history and no target specified")
	}
	if _, exists := devices[targetName]; !exists {
		if targetAddress == "" || ca == "" || caKey == "" {
			return errors.New("Device not found")
		}
		devices[targetName] = Device{
			Address: targetAddress,
			Ca:      ca,
			CaKey:   caKey,
		}
	} else {
		device := devices[targetName]
		if targetAddress != "" {
			device.Address = targetAddress
		}
		if ca != "" {
			device.Ca = ca
		}
		if caKey != "" {
			device.CaKey = caKey
		}
		devices[targetName] = device
	}
	viper.Set("targets.last_target", targetName)
	viper.Set("targets.devices", devices)
	return nil
}
