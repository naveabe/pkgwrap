package tracker

import (
	"encoding/json"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
)

type EssJobstore struct {
	EssDatastore
}

func NewEssJobstore(cfg *config.DatastoreConfig, logger *logging.Logger) (*EssJobstore, error) {
	eds, err := NewEssDatastore(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &EssJobstore{EssDatastore: *eds}, nil
}

func (e *EssJobstore) AddRequest(pkgReq specer.PackageRequest) (string, error) {
	return e.Add("pkgreq", pkgReq)
}

func (e *EssJobstore) UpdateRequest(id string, pkgReq specer.PackageRequest) error {
	return e.Update("pkgreq", id, pkgReq)
}

func (e *EssJobstore) AddJob(job BuildJob) (string, error) {
	return e.Add("job", job)
}

func (e *EssJobstore) performQuery(terms map[string]interface{}) ([]BuildJob, error) {
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

func (e *EssJobstore) GetBuildsForPackageVersion(pkgr, name, version string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": pkgr,
		"project":  name,
		"version":  version,
	}

	return e.performQuery(terms)
}

func (e *EssJobstore) GetBuildsForPackage(pkgr, name string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": pkgr,
		"project":  name,
	}
	return e.performQuery(terms)
}

func (e *EssJobstore) GetBuildsForUser(pkgr string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": pkgr,
	}
	return e.performQuery(terms)
}

func (e *EssJobstore) GetBuild(id string) (*BuildJob, error) {
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

	return &bJobs[0], nil
}
