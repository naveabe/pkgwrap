package tracker

import (
	"encoding/json"
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

func (e *EssDatastore) Add(job BuildJob) error {
	resp, err := e.conn.Index(e.index, "job", "", nil, job)
	//e.conn.Flush()
	if err != nil {
		e.logger.Trace.Printf("%s\n", err)
		return err
	}

	if !resp.Created {
		return fmt.Errorf("Failed to record job: %s", resp)
	}

	return nil
}

func (e *EssDatastore) performQuery(terms map[string]interface{}) ([]BuildJob, error) {
	var (
		out []BuildJob
		//filters = map[string]map[string]string{"term": terms}
	)

	q := map[string]interface{}{
		"sort":   map[string]string{"timestamp": "desc"},
		"filter": map[string]map[string]interface{}{"term": terms},
	}
	e.logger.Trace.Printf("Query: %v\n", q)

	resp, err := e.conn.Search(e.index, "job", nil, q)
	if err != nil {
		return out, err
	}

	out = make([]BuildJob, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(*hit.Source, &out[i]); err != nil {
			return out, err
		}
		out[i].Id = hit.Id
	}

	return out, nil
}

func (e *EssDatastore) GetBuildsForPackageVersion(pkgr, name, version string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": pkgr,
		"project":  name,
		"version":  version,
	}

	return e.performQuery(terms)
}

func (e *EssDatastore) GetBuildsForPackage(pkgr, name string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": pkgr,
		"project":  name,
	}
	return e.performQuery(terms)
}

func (e *EssDatastore) GetBuildsForUser(pkgr string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": pkgr,
	}
	return e.performQuery(terms)
}

func (e *EssDatastore) GetBuild(id string) (*BuildJob, error) {
	terms := map[string]interface{}{
		"jobs.id": []string{id},
	}

	bJobs, err := e.performQuery(terms)
	if err != nil {
		return nil, err
	}
	if len(bJobs) > 1 {
		e.logger.Trace.Printf("** BIG BIG PROBLEM! MORE THAN ONE FOUND: %d **\n", len(bJobs))
	} else if len(bJobs) == 0 {
		return nil, fmt.Errorf("Not found")
	}

	//bJobs[0].Id
	//data := map[string]map[string]string{"doc":{}}
	//e.conn.Update(e.index, "job", bJobs[0].Id, nil, data)

	return &bJobs[0], nil
}
