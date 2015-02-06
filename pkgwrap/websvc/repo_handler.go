package websvc

import (
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"net/http"
)

type RepoHandler struct {
	DefaultMethodHandler

	Repository repository.BuildRepository

	logger *logging.Logger
}

func NewRepoHandler(repo repository.BuildRepository, logger *logging.Logger) *RepoHandler {
	rh := RepoHandler{Repository: repo}
	if logger == nil {
		rh.logger = logging.NewStdLogger()
	} else {
		rh.logger = logger
	}
	return &rh
}

/*
   Params:
       args : name, version, distroLabel, package (rpm or deb)
*/
func (rh *RepoHandler) GET(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	pkgr := args[0]
	project := args[1]
	rh.logger.Trace.Printf("%s/%s (%d)", pkgr, project, len(args))

	switch len(args) {
	case 2:
		list, err := rh.Repository.ListPackageVersions(pkgr, project)
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}
		return nil, list, 200
	case 3:
		list, err := rh.Repository.ListPackageDistros(pkgr, project, args[2])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}
		return nil, list, 200
	case 4:
		list, err := rh.Repository.ListPackages(pkgr, project, args[2], args[3])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}
		return nil, list, 200
	case 5:
		// Send package to client
		pkgPath, err := rh.Repository.GetPackagePathForDistro(pkgr, project, args[2], args[3], args[4])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}
		http.ServeFile(w, r, pkgPath)
		// Tells parent not to send default response as we will be sending it here.
		return nil, nil, -1
	default:
		//return header, data, code
		return nil, map[string]string{"error": "Bad request"}, 400
	}
}
