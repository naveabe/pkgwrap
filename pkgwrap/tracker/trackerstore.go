package tracker

import (
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	elastigo "github.com/mattbaird/elastigo/lib"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"strings"
)

/*
	Elasticsearch doc types
*/
const (
	DTYPE_CONTAINER string = "container"
	DTYPE_REQUEST   string = "pkgreq"
)

/*
   Complete build information.  Contains:
    - package request
    - array of container info per distro.
    - timestamp
*/

type BuildInfo struct {
	Request *request.PackageRequest
	// Container info index
	Containers map[string]*docker.Container
	//Containers map[string]BuildContainerInfo
	Timestamp float64
}

func NewBuildInfo() *BuildInfo {
	return &BuildInfo{Containers: map[string]*docker.Container{}}
}

/*
	Stores and retrieves package build information. Two data structures are recorded:

		1. docker container config
		2. package request
*/
type TrackerStore struct {
	EssDatastore
}

func NewTrackerStore(cfg *config.DatastoreConfig, logger *logging.Logger) (*TrackerStore, error) {
	eds, err := NewEssDatastore(cfg, logger)
	if err != nil {
		return nil, err
	}
	return &TrackerStore{EssDatastore: *eds}, nil
}

/*
	Add container details
*/
func (t *TrackerStore) AddContainer(id string, contInfo interface{}) (string, error) {
	return t.AddWithId(DTYPE_CONTAINER, id, contInfo)
}

/*
	Add build request (packaging request)
*/
func (t *TrackerStore) AddRequest(pkgReq request.PackageRequest) (string, error) {
	return t.Add(DTYPE_REQUEST, pkgReq)
}

/*
	Update container details
*/
func (t *TrackerStore) UpdateContainer(id string, contInfo interface{}) error {
	return t.Update(DTYPE_CONTAINER, id, contInfo)
}

/*
	Update packaging request. Usually after the build has started
*/
func (t *TrackerStore) UpdateRequest(id string, pkgReq request.PackageRequest) error {
	return t.Update(DTYPE_REQUEST, id, pkgReq)
}

/*
	IN PROGRESS
	Get the request containing the container id
*/
func (t *TrackerStore) GetRequestByContainerId(id string) (*request.PackageRequest, error) {
	preq := &request.PackageRequest{}

	terms := map[string]interface{}{
		"term": map[string]string{"Distributions.id": id},
	}

	eRslt, err := t.performQuery("pkgreq", terms, nil)
	if err != nil {
		return preq, err
	}

	if len(eRslt.Hits.Hits) < 1 {
		return preq, fmt.Errorf("Not found: %s", id)
	}

	t.logger.Trace.Printf("%s\n", *eRslt.Hits.Hits[0].Source)

	if err = json.Unmarshal(*eRslt.Hits.Hits[0].Source, preq); err != nil {
		return preq, err
	}
	return preq, nil
}

/*
	Get build container.  This contains all the container information.

	Args:
		id : container id

	Return:
		docker container configuration
*/
func (t *TrackerStore) GetContainer(id string) (*docker.Container, error) {
	var (
		err   error
		dcont docker.Container
	)

	result, err := t.conn.Get(t.index, DTYPE_CONTAINER, id, nil)
	if err != nil {
		return &dcont, err
	}

	return &dcont, json.Unmarshal(*result.Source, &dcont)
}

/*
	Gets the combined information containing the request, container/s info
	and posted timestamp.

	Args:
		0 : repo
		1 : user
		2 : project
		3 : version

	Return:
		array of the latest build/s
*/
func (t *TrackerStore) GetBuildInfo(args ...string) ([]*BuildInfo, error) {
	var (
		err       error
		terms     map[string]interface{}
		out       []*BuildInfo
		queryOpts = map[string]interface{}{
			"fields": []string{"_source", "_timestamp"},
			"sort":   map[string]string{"_timestamp": "desc"},
		}
	)

	if terms, err = t.makeTermsQuery(args); err != nil {
		return out, err
	}

	resp, err := t.performQuery(DTYPE_REQUEST, terms, queryOpts)
	if err != nil {
		return out, err
	}
	/* convert responses to []BuildInfo */
	return t.assembleBuildInfo(resp)
}

/*
	Assemble BuildInfo - add timestamp, package request  and containers from ess result
*/
func (t *TrackerStore) assembleBuildInfo(resp elastigo.SearchResult) ([]*BuildInfo, error) {
	var (
		err error
		out = make([]*BuildInfo, len(resp.Hits.Hits))
	)

	for i, hit := range resp.Hits.Hits {
		// Get timestamp
		flds := map[string]interface{}{}
		// maybe check ??
		json.Unmarshal(*hit.Fields, &flds)
		timestamp, _ := flds["_timestamp"].(float64)

		out[i] = NewBuildInfo()
		out[i].Timestamp = timestamp
		if err = json.Unmarshal(*hit.Source, &out[i].Request); err != nil {
			return out, err
		}
		for _, val := range out[i].Request.Distributions {
			if out[i].Containers[val.Id], err = t.GetContainer(val.Id); err != nil {
				return out, err
			}
		}
	}
	return out, nil
}

/*
	Assemble terms query for elastic search

		filtered
			filter
				bool
					must : []
*/
func (t *TrackerStore) makeTermsQuery(args []string) (map[string]interface{}, error) {
	var mustTerms []interface{}
	query := make(map[string]interface{})

	switch len(args) {
	case 0:
		//latest jobs
		break
	case 2: //packager
		mustTerms = []interface{}{
			map[string]interface{}{
				"term": map[string]string{
					"Package.packager": strings.ToLower(args[1]),
				},
			},
		}
		break
	case 3: //project
		mustTerms = []interface{}{
			map[string]interface{}{
				"term": map[string]interface{}{
					"Package.packager": strings.ToLower(args[1]),
				},
			},
			map[string]interface{}{
				"term": map[string]interface{}{
					"Package.name": strings.ToLower(args[2]),
				},
			},
		}
		break
	case 4: //version
		mustTerms = []interface{}{
			map[string]interface{}{
				"term": map[string]interface{}{
					"Package.packager": strings.ToLower(args[1]),
				},
			},
			map[string]interface{}{
				"term": map[string]interface{}{
					"Package.name": strings.ToLower(args[2]),
				},
			},
			map[string]interface{}{
				"term": map[string]interface{}{
					"Package.version": strings.ToLower(args[3]),
				},
			},
		}
		break
	default:
		return query, fmt.Errorf("Invalid request: %v", args)
		break
	}

	if mustTerms != nil {
		query = map[string]interface{}{
			"filtered": map[string]interface{}{
				"filter": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": mustTerms,
					},
				},
			},
		}
	}

	return query, nil
}

/*
	Generic function to perform arbitrary queries

	Args:
		docType : container | pkgreq
		terms   : terms to query
		opts    : extra options to pass to elasticsearch

	Return:
		elasticsearch search result
*/
func (t *TrackerStore) performQuery(docType string,
	qterms map[string]interface{}, opts map[string]interface{}) (elastigo.SearchResult, error) {

	query := map[string]interface{}{}
	if qterms != nil && len(qterms) > 0 {
		query["query"] = qterms
	}
	if opts != nil {
		for k, v := range opts {
			query[k] = v
		}
	}

	t.logger.Trace.Printf("Tracker query: %s %s %s\n", t.index, docType, query)
	return t.conn.Search(t.index, docType, nil, query)
}
