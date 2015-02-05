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

	if len(args) == 1 || (len(args) == 2 && args[1] == "") {
		list, err := rh.Repository.ListPackageVersions(args[0])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}
		return nil, list, 200
	} else if len(args) == 2 || (len(args) == 3 && args[2] == "") {
		list, err := rh.Repository.ListPackageDistros(args[0], args[1])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}
		return nil, list, 200
	} else if len(args) == 3 || (len(args) == 4 && args[3] == "") {
		list, err := rh.Repository.ListPackages(args[0], args[1], args[2])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}
		return nil, list, 200
	} else if len(args) == 4 {
		// Send package to client
		pkgPath, err := rh.Repository.GetPackagePathForDistro(args[0], args[1], args[2], args[3])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 500
		}

		http.ServeFile(w, r, pkgPath)
		// Tells parent not to send default response as we will be sending it here.
		return nil, nil, -1
	}
	//return header, data, code
	return nil, map[string]string{"error": "Bad request"}, 400
}
