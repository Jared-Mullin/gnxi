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
	"os"
	"path"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/spf13/viper"
)

type Want struct {
	targets map[string]Target
	err     error
}

func TestPrepareTarget(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %v", err)
	}
	certPath := path.Join(dir, "ca.crt")
	certKeyPath := path.Join(dir, "ca.key")
	tests := []struct {
		name          string
		targetName    string
		targetAddress string
		targetCA      string
		targetCAKey   string
		config        map[string]Target
		lastTarget    string
		want          Want
	}{
		{
			name: "No targets in history, no target specified",
			want: Want{
				err:     errors.New("No targets in history and no target specified"),
				targets: map[string]Target{},
			},
		},
		{
			name: "No target specified",
			config: map[string]Target{"myhost.com": {
				Address: "localhost:9339",
				Ca:      certPath,
				CaKey:   certKeyPath,
			}},
			lastTarget: "myhost.com",
			want: Want{
				err: nil,
				targets: map[string]Target{"myhost.com": {
					Address: "localhost:9339",
					Ca:      certPath,
					CaKey:   certKeyPath,
				}},
			},
		},
		{
			name: "Non-existent target",
			config: map[string]Target{"myhost.com": {
				Address: "localhost:9339",
				Ca:      certPath,
				CaKey:   certKeyPath,
			}},
			targetName: "nonexistenttarget",
			lastTarget: "myhost.com",
			want: Want{
				err: errors.New("Target not found"),
				targets: map[string]Target{"myhost.com": {
					Address: "localhost:9339",
					Ca:      certPath,
					CaKey:   certKeyPath,
				}},
			},
		},
		{
			name:          "Add new target",
			targetName:    "myhost.com",
			targetAddress: "localhost:9339",
			targetCA:      certPath,
			targetCAKey:   certKeyPath,
			want: Want{
				targets: map[string]Target{
					"myhost.com": {
						Address: "localhost:9339",
						Ca:      certPath,
						CaKey:   certKeyPath,
					},
				},
			},
		},
		{
			name:          "Update existing target's address",
			targetName:    "myhost.com",
			targetAddress: "newhost:9340",
			config: map[string]Target{"myhost.com": {
				Address: "localhost:9339",
				Ca:      certPath,
				CaKey:   certKeyPath,
			}},
			want: Want{
				targets: map[string]Target{
					"myhost.com": {
						Address: "newhost:9340",
						Ca:      certPath,
						CaKey:   certKeyPath,
					},
				},
			},
		},
		{
			name:       "Update existing target's ca",
			targetName: "myhost.com",
			targetCA:   "newca.crt",
			config: map[string]Target{"myhost.com": {
				Address: "localhost:9339",
				Ca:      certPath,
				CaKey:   certKeyPath,
			}},
			want: Want{
				targets: map[string]Target{
					"myhost.com": {
						Address: "localhost:9339",
						Ca:      path.Join(dir, "newca.crt"),
						CaKey:   certKeyPath,
					},
				},
			},
		},
		{
			name:        "Update existing target's ca_key",
			targetName:  "myhost.com",
			targetCAKey: "newca.key",
			config: map[string]Target{"myhost.com": {
				Address: "localhost:9339",
				Ca:      certPath,
				CaKey:   certKeyPath,
			}},
			want: Want{
				targets: map[string]Target{
					"myhost.com": {
						Address: "localhost:9339",
						Ca:      certPath,
						CaKey:   path.Join(dir, "newca.key"),
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			viper.Reset()
			var targets map[string]Target
			viper.Set("targets.last_target", test.lastTarget)
			viper.Set("targets.devices", test.config)
			err := prepareTarget(test.targetName, test.targetAddress, test.targetCA, test.targetCAKey, true)
			viper.UnmarshalKey("targets.devices", &targets)
			got := Want{
				targets: targets,
				err:     err,
			}
			if diff := pretty.Compare(test.want, got); diff != "" {
				t.Errorf("prepareTarget(%s, %s, %s, %s): (-want +got)\n%s", test.targetName, test.targetAddress, test.targetCA, test.targetCAKey, diff)
			}
		})
	}
}
