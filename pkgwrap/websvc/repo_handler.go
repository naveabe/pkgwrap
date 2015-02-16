package websvc

import (
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"net/http"
)

var (
	ALL_ORIGIN_ACL = map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
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
	//repo := args[0]
	//pkgr := args[1]
	//project := args[1]
	rh.logger.Trace.Printf("Path - %v\n", args)

	switch len(args) {
	case 2:
		list, err := rh.Repository.ListUserProjects(args[0], args[1])
		if err != nil {
			return ALL_ORIGIN_ACL, map[string]string{"error": err.Error()}, 404
		}
		return ALL_ORIGIN_ACL, list, 200
	case 3:
		list, err := rh.Repository.ListPackageVersions(args[0], args[1], args[2])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 400
		}
		return ALL_ORIGIN_ACL, list, 200
	case 4:
		list, err := rh.Repository.ListPackageDistros(args[0], args[1], args[2], args[3])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 400
		}
		return ALL_ORIGIN_ACL, list, 200
	case 5:
		//rh.logger.Trace.Printf("Path - %v\n", args)
		list, err := rh.Repository.ListPackages(args[0], args[1], args[2], args[3], args[4])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 400
		}
		return ALL_ORIGIN_ACL, list, 200
	case 6:
		// Send package to client
		pkgPath, err := rh.Repository.GetPackagePathForDistro(args[0], args[1], args[2], args[3], args[4], args[5])
		if err != nil {
			return nil, map[string]string{"error": err.Error()}, 400
		}
		http.ServeFile(w, r, pkgPath)
		// Tells parent not to send default response as we will be sending it here.
		return nil, nil, -1
	default:
		//return header, data, code
		//break
		return nil, map[string]string{"error": "Invalid path"}, 400
	}
}
