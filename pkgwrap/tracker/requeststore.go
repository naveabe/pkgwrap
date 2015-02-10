package tracker

import (
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
)

type EssRequeststore struct {
	EssDatastore
}

func NewEssRequeststore(cfg *config.DatastoreConfig, logger *logging.Logger) (*EssRequeststore, error) {
	eds, err := NewEssDatastore(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &EssRequeststore{EssDatastore: *eds}, nil
}
