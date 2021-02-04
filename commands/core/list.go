// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package core

import (
	"github.com/arduino/arduino-cli/arduino/cores"
	"github.com/arduino/arduino-cli/commands"
	"github.com/pkg/errors"
)

// GetPlatforms returns a list of installed platforms, optionally filtered by
// those requiring an update.
func GetPlatforms(instanceID int32, updatableOnly, all bool) ([]*cores.PlatformRelease, error) {
	i := commands.GetInstance(instanceID)
	if i == nil {
		return nil, errors.Errorf("unable to find an instance with ID: %d", instanceID)
	}

	packageManager := i.PackageManager
	if packageManager == nil {
		return nil, errors.New("invalid instance")
	}

	if updatableOnly && all {
		return nil, errors.New("can't use both updatableOnly and all flags at the same time")
	}

	res := []*cores.PlatformRelease{}
	for _, targetPackage := range packageManager.Packages {
		for _, platform := range targetPackage.Platforms {
			platformRelease := packageManager.GetInstalledPlatformRelease(platform)
			if all {
				res = append(res, platform.GetLatestRelease())
				continue
			}
			if platformRelease != nil {
				if updatableOnly {
					if latest := platform.GetLatestRelease(); latest == nil || latest == platformRelease {
						continue
					}
				}
				res = append(res, platformRelease)
			}
		}
	}

	return res, nil
}
