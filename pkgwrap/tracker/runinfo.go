package tracker

import (
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"strings"
	"time"
)

type BuildJobId struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
}

func NewBuildJobId(id, host string) *BuildJobId {
	return &BuildJobId{id, host}
}

func NewBuildJobIdFromString(jobId string) (*BuildJobId, error) {
	b := BuildJobId{}
	parts := strings.Split(jobId, "@")
	if len(parts) != 2 {
		return &b, fmt.Errorf("Invalid job ID: %s", jobId)
	}
	b.Id = parts[0]
	b.Uri = parts[1]

	return &b, nil
}

/*
type JobId string

func (j *JobId) HostURI() string {
	return strings.Split(string(*j), "@")[1]
}
func (j *JobId) ContainerId() string {
	return strings.Split(string(*j), "@")[0]
}
*/

/*
   Holds a single packge request job for a given project.
   i.e. 1 per project regardless of the no. of distros
*/
type BuildJob struct {
	Timestamp float64      `json:"timestamp"`
	Username  string       `json:"username"`
	Project   string       `json:"project"`
	URL       string       `json:"url"`
	TagBranch string       `json:"tagbranch"`
	Version   string       `json:"version"`
	Jobs      []BuildJobId `json:"jobs"`
}

func NewBuildJob(pkgReq *specer.PackageRequest, buildIds []string, uri string) *BuildJob {
	jBldIds := make([]BuildJobId, len(buildIds))
	for i, v := range buildIds {
		jid, _ := NewBuildJobIdFromString(v + "@" + uri)
		jBldIds[i] = *jid
	}
	return &BuildJob{
		Timestamp: float64(time.Now().UnixNano()) / 1000000000,
		Username:  pkgReq.Package.Packager,
		URL:       pkgReq.Package.URL,
		Project:   pkgReq.Package.Name,
		TagBranch: pkgReq.Package.TagBranch,
		Version:   pkgReq.Package.Version,
		Jobs:      jBldIds,
	}
}

/*
func NewBuildJobFromURL(url, tag, version string, ids []string) *BuildJob {
	b := BuildJob{
		URL:       url,
		TagBranch: tag,
		Version:   version,
	}

    return &b
}
*/
