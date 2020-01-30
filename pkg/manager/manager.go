// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package manager is is the main coordinator for the ONOS RAN subsystem.
package manager

import (
	"github.com/OpenNetworkingFoundation/gmap-ran/api/types"
	log "k8s.io/klog"
)

var mgr Manager

// NewManager initializes the RAN subsystem.
func NewManager() (*Manager, error) {
	log.Info("Creating Manager")
	mgr = Manager{}
	return &mgr, nil
}

// Manager single point of entry for the trafficsim system.
type Manager struct {
	Towers    map[string]*types.Tower
	Locations map[string]*Location
	Routes    map[string]*types.Route
}

// Run starts a synchronizer based on the devices and the northbound services.
func (m *Manager) Run(towerparams TowersParams, locParams LocationsParams, routesParams RoutesParams) {
	log.Infof("Starting Manager with %v %v %v", towerparams, locParams, routesParams)
	m.Towers = newTowers(towerparams)
	m.Locations = newLocations(locParams, towerparams)

	var err error
	m.Routes, err = m.newRoutes(routesParams)
	if err != nil {
		log.Fatalf("Error calculating routes", err)
	}
}

//Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}

// GetManager returns the initialized and running instance of manager.
// Should be called only after NewManager and Run are done.
func GetManager() *Manager {
	return &mgr
}