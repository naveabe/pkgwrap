package tracker

import (
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"io/ioutil"
	"os"
)

const ESS_DEFAULT_RESULT_SIZE int = 10000000

type EssMapping struct {
	Meta             map[string]interface{} `json:"_meta"`
	DynamicTemplates []interface{}          `json:"dynamic_templates"`
}

type EssDatastore struct {
	conn   *elastigo.Conn
	index  string
	logger *logging.Logger
}

func NewEssDatastore(cfg *config.DatastoreConfig, logger *logging.Logger) (*EssDatastore, error) {
	ed := EssDatastore{}

	if logger == nil {
		ed.logger = logging.NewStdLogger()
	} else {
		ed.logger = logger
	}

	c := elastigo.NewConn()
	c.Domain = cfg.Host
	c.Port = fmt.Sprintf("%d", cfg.Port)

	ed.conn = c
	ed.index = cfg.Index

	exists, err := c.ExistsIndex(cfg.Index, "", nil)
	if err != nil {
		if err.Error() == "record not found" {
			exists = false
		} else {
			return &ed, err
		}
	}

	if !exists {
		return &ed, ed.initializeIndex(cfg.MappingFile)
	}
	return &ed, nil
}

func (e *EssDatastore) initializeIndex(mappingFile string) error {
	resp, err := e.conn.CreateIndex(e.index)
	if err != nil {
		return err
	}
	e.logger.Warning.Printf("Index created: %s %s\n", e.index, resp)

	if _, err := os.Stat(mappingFile); err != nil {
		return fmt.Errorf("Mapping file not found %s: %s", mappingFile, err)
	}

	mappingDataBytes, err := ioutil.ReadFile(mappingFile)
	if err != nil {
		return err
	}
	b, err := e.conn.DoCommand("PUT", fmt.Sprintf("/%s/_mapping/_default_", e.index), nil, mappingDataBytes)
	if err != nil {
		return err
	}
	e.logger.Warning.Printf("Updated _default_ mapping for %s: %s\n", e.index, b)
	return nil
}

func (e *EssDatastore) Add(docType string, data interface{}) (string, error) {
	resp, err := e.conn.Index(e.index, docType, "", nil, data)
	if err != nil {
		e.logger.Trace.Printf("%s\n", err)
		return "", err
	}
	if !resp.Created {
		return "", fmt.Errorf("Failed to record job: %s", resp)
	}

	return resp.Id, nil
}

func (e *EssDatastore) Update(docType, id string, data interface{}) error {
	resp, err := e.conn.Index(e.index, docType, id, nil, data)
	if err != nil {
		e.logger.Trace.Printf("%s\n", err)
		return err
	}
	e.logger.Trace.Printf("%s\n", resp)
	/*
		if !resp.Created {
			return fmt.Errorf("Failed to record job: %s", resp)
		}
	*/
	return nil
}
