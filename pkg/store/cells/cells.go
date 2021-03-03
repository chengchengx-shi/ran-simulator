// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package cells

import (
	"context"
	"math/rand"
	"reflect"
	"sync"

	"github.com/google/uuid"

	"github.com/onosproject/ran-simulator/pkg/store/event"

	"github.com/onosproject/ran-simulator/pkg/store/watcher"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	liblog "github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
)

var log = liblog.GetLogger("store", "cells")

// Store tracks inventory of simulated cells.
type Store interface {
	// Add adds the specified cell to the registry
	Add(ctx context.Context, cell *model.Cell) error

	// Get retrieves the cell with the specified ECGI
	Get(ctx context.Context, ecgi types.ECGI) (*model.Cell, error)

	// Update updates the cell
	Update(ctx context.Context, Cell *model.Cell) error

	// Delete deletes the cell with the specified ECGI
	Delete(ctx context.Context, ecgi types.ECGI) (*model.Cell, error)

	// Watch watches the cell inventory events using the supplied channel
	Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error

	// List list all of the cells
	List(ctx context.Context) ([]*model.Cell, error)

	// GetRandomCell retrieves a random cell from the registry
	GetRandomCell() (*model.Cell, error)
}

// WatchOptions allows tailoring the WatchCells behaviour
type WatchOptions struct {
	Replay  bool
	Monitor bool
}

type store struct {
	mu        sync.RWMutex
	cells     map[types.ECGI]*model.Cell
	nodeStore nodes.Store
	watchers  *watcher.Watchers
}

// NewCellRegistry creates a new store abstraction from the specified fixed cell map.
func NewCellRegistry(cells map[string]model.Cell, nodeStore nodes.Store) Store {
	log.Infof("Creating registry from model with %d cells", len(cells))
	watchers := watcher.NewWatchers()
	reg := &store{
		mu:        sync.RWMutex{},
		cells:     make(map[types.ECGI]*model.Cell),
		nodeStore: nodeStore,
		watchers:  watchers,
	}

	// Copy the Cells into our own map
	for _, c := range cells {
		cell := c // avoids scopelint issue
		reg.cells[cell.ECGI] = &cell
	}

	log.Infof("Created registry primed with %d cells", len(reg.cells))
	return reg
}

// Add adds a cell
func (s *store) Add(ctx context.Context, cell *model.Cell) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cells[cell.ECGI]; ok {
		return errors.New(errors.NotFound, "cell with EnbID already exists")
	}

	s.cells[cell.ECGI] = cell
	cellEvent := event.Event{
		Key:   cell.ECGI,
		Value: cell,
		Type:  Created,
	}
	s.watchers.Send(cellEvent)
	return nil

}

// Get gets a cell
func (s *store) Get(ctx context.Context, ecgi types.ECGI) (*model.Cell, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if cell, ok := s.cells[ecgi]; ok {
		return cell, nil
	}

	return nil, errors.New(errors.NotFound, "cell not found")
}

// Update updates a cell
func (s *store) Update(ctx context.Context, cell *model.Cell) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if prevCell, ok := s.cells[cell.ECGI]; ok {
		s.cells[cell.ECGI] = cell
		prevNeighbors := prevCell.Neighbors
		equalNeighborsResult := equalNeighbors(prevNeighbors, cell.Neighbors)
		if !equalNeighborsResult {
			cellEvent := event.Event{
				Key:   cell.ECGI,
				Value: cell,
				Type:  UpdatedNeighbors,
			}
			s.watchers.Send(cellEvent)
		}

		cellEvent := event.Event{
			Key:   cell.ECGI,
			Value: cell,
			Type:  Updated,
		}
		s.watchers.Send(cellEvent)
		return nil
	}

	return errors.New(errors.NotFound, "cell not found")
}

// Delete deletes a cell
func (s *store) Delete(ctx context.Context, ecgi types.ECGI) (*model.Cell, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cell, ok := s.cells[ecgi]; ok {
		delete(s.cells, ecgi)
		deleteEvent := event.Event{
			Key:   cell.ECGI,
			Value: cell,
			Type:  Deleted,
		}
		s.watchers.Send(deleteEvent)
		err := s.nodeStore.PruneCell(ctx, ecgi)
		if err != nil {
			return nil, err
		}
		return cell, nil
	}
	return nil, errors.New(errors.NotFound, "cell not found")
}

// Watch watch cell events
func (s *store) Watch(ctx context.Context, ch chan<- event.Event, options ...WatchOptions) error {
	log.Debug("Watching cell changes")
	replay := len(options) > 0 && options[0].Replay
	id := uuid.New()
	err := s.watchers.AddWatcher(id, ch)
	if err != nil {
		log.Error(err)
	}
	go func() {
		<-ctx.Done()
		err = s.watchers.RemoveWatcher(id)
		if err != nil {
			log.Error(err)
		}
		close(ch)

	}()

	if replay {
		go func() {
			for _, cell := range s.cells {
				ch <- event.Event{
					Key:   cell.ECGI,
					Value: cell,
					Type:  None,
				}
			}
		}()

	}
	return nil
}

// List returns list of cells
func (s *store) List(ctx context.Context) ([]*model.Cell, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Cell, 0, len(s.cells))
	for _, cell := range s.cells {
		list = append(list, cell)
	}
	return list, nil
}

func (s *store) GetRandomCell() (*model.Cell, error) {
	keys := reflect.ValueOf(s.cells).MapKeys()
	ecgi := types.ECGI(keys[rand.Intn(len(keys))].Uint())
	return s.cells[ecgi], nil
}