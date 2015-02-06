package tracker

import (
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"strings"
	"time"
)

type BuildJobId struct {
	Id     string `json:"id"`
	Uri    string `json:"uri"`
	Status string `json:"status"`
}

func NewBuildJobId(id, uri string) *BuildJobId {
	return &BuildJobId{Id: id, Uri: uri}
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
   Holds a single packge request job for a given project. i.e. 1 per
   project regardless of the no. of distros.  Each distro container
   id is stored along with the uri for the host container.

   This is used to retreive details about the builds byt the user.
*/
type BuildJob struct {
	Id        string       `json:"_id,omitempty"`
	Timestamp float64      `json:"timestamp"`
	Username  string       `json:"username"`
	Project   string       `json:"project"`
	URL       string       `json:"url"`
	TagBranch string       `json:"tagbranch"`
	Version   string       `json:"version"`
	Jobs      []BuildJobId `json:"jobs"`
}

func NewBuildJob(pkgReq *specer.PackageRequest, buildIds []string, uri string) *BuildJob {

	bj := BuildJob{
		Timestamp: float64(time.Now().UnixNano()) / 1000000000,
		Username:  pkgReq.Package.Packager,
		URL:       pkgReq.Package.URL,
		Project:   pkgReq.Package.Name,
		TagBranch: pkgReq.Package.TagBranch,
		Version:   pkgReq.Package.Version,
		Jobs:      make([]BuildJobId, len(buildIds)),
	}

	for i, v := range buildIds {
		jid, _ := NewBuildJobIdFromString(v + "@" + uri)
		bj.Jobs[i] = *jid
	}

	return &bj
}

/*
	Add BuildJob to datastore
*/
func (b *BuildJob) Record(ds IJobstore) error {
	return ds.Add(*b)
}

func (b *BuildJob) GetSubJob(id string) (BuildJobId, error) {
	for _, v := range b.Jobs {
		if v.Id == id {
			return v, nil
		}
	}
	return BuildJobId{}, fmt.Errorf("Not found")
}
