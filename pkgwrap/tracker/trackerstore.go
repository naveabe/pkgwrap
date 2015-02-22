package tracker

import (
	"encoding/json"
	"fmt"
	elastigo "github.com/mattbaird/elastigo/lib"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"strings"
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

/*
	Performs generic query
*/
func (e *EssJobstore) performQuery(idxType string, terms map[string]interface{}, sortTimestamp bool) (elastigo.SearchResult, error) {
	var q map[string]interface{}
	if sortTimestamp {
		q = map[string]interface{}{
			"sort":   map[string]string{"timestamp": "desc"},
			"filter": map[string]map[string]interface{}{"term": terms},
		}
	} else {
		q = map[string]interface{}{
			"filter": map[string]map[string]interface{}{"term": terms},
		}
	}
	e.logger.Trace.Printf("Query: %v\n", q)

	return e.conn.Search(e.index, idxType, nil, q)
}

func (e *EssJobstore) AddRequest(pkgReq specer.PackageRequest) (string, error) {
	return e.Add("pkgreq", pkgReq)
}

func (e *EssJobstore) GetRequests(args ...string) ([]specer.PackageRequest, error) {
	var (
		out []specer.PackageRequest
	)
	terms := map[string]interface{}{}
	switch len(args) {
	case 2:
		terms["Package.packager"] = strings.ToLower(args[1])
		//packager
		break
	case 3:
		terms["Package.packager"] = strings.ToLower(args[1])
		terms["Package.name"] = strings.ToLower(args[2])
		//project
		break
	case 4:
		terms["Package.packager"] = strings.ToLower(args[1])
		terms["Package.name"] = strings.ToLower(args[2])
		terms["Package.version"] = strings.ToLower(args[3])
		//version
		break
	default:
		return out, fmt.Errorf("Invalid request: %v", args)
		break
	}
	e.logger.Trace.Printf("%s\n", terms)

	resp, err := e.performQuery("pkgreq", terms, false)
	if err != nil {
		return out, err
	}

	out = make([]specer.PackageRequest, len(resp.Hits.Hits))
	for i, hit := range resp.Hits.Hits {
		if err := json.Unmarshal(*hit.Source, &out[i]); err != nil {
			return out, err
		}
		out[i].Id = hit.Id
	}

	return out, nil
}

func (e *EssJobstore) UpdateRequest(id string, pkgReq specer.PackageRequest) error {
	return e.Update("pkgreq", id, pkgReq)
}

func (e *EssJobstore) AddJob(job BuildJob) (string, error) {
	return e.Add("job", job)
}

func (e *EssJobstore) performBuildJobQuery(terms map[string]interface{}) ([]BuildJob, error) {
	var (
		out []BuildJob
		//filters = map[string]map[string]string{"term": terms}
	)
	/*
			q := map[string]interface{}{
				"sort":   map[string]string{"timestamp": "desc"},
				"filter": map[string]map[string]interface{}{"term": terms},
			}
			e.logger.Trace.Printf("Query: %v\n", q)

			resp, err := e.conn.Search(e.index, "job", nil, q)

		resp, err := e.conn.Search("job", terms)
	*/
	resp, err := e.performQuery("job", terms, true)
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
		"username": strings.ToLower(pkgr),
		"project":  strings.ToLower(name),
		"version":  version,
	}

	return e.performBuildJobQuery(terms)
}

func (e *EssJobstore) GetBuildsForPackage(pkgr, name string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": strings.ToLower(pkgr),
		"project":  strings.ToLower(name),
	}
	return e.performBuildJobQuery(terms)
}

func (e *EssJobstore) GetBuildsForUser(pkgr string) ([]BuildJob, error) {
	terms := map[string]interface{}{
		"username": strings.ToLower(pkgr),
	}
	return e.performBuildJobQuery(terms)
}

func (e *EssJobstore) GetBuild(id string) (*BuildJob, error) {
	terms := map[string]interface{}{
		"jobs.id": []string{id},
	}

	bJobs, err := e.performBuildJobQuery(terms)
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
