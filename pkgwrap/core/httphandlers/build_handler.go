package httphandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/core/request"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"github.com/naveabe/pkgwrap/pkgwrap/tracker"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

/* Request params mapping */
var (
	HTTP_REQ_PARAMS = map[string]string{
		"PACKAGE":    "package",
		"BUILD_CONF": "conf",
	}
)

/*
 * endpoint:  /:name/:version=optional
 */
type PkgBuilderMethodHandler struct {
	DefaultMethodHandler

	Config     *config.AppConfig
	Repository repository.BuildRepository

	Logger *logging.Logger

	// This channel will be read to get PackageRequests
	RequestChan chan request.PackageRequest

	Datastore *tracker.TrackerStore
}

/*
	Params:
		binary : Build package from pre-compiled data
		dryrun : Runs through accepting the request but does not actually
				 submit for building
*/
func (m *PkgBuilderMethodHandler) POST(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	var (
		err     error = nil
		pkgReq  request.PackageRequest
		pkgFile *multipart.FileHeader
	)
	if len(args) != 3 {
		return nil, map[string]string{"error": "Invalid request"}, 400
	}
	// Determines build type: 'source' or 'binary'
	if _, ok := r.URL.Query()["binary"]; ok {
		m.Logger.Debug.Printf("Binary build request!\n")
		if pkgReq, pkgFile, err = m.assembleBuiltPkgReq(r, args...); err == nil {
			err = m.downloadUserPackage(pkgReq.Package, pkgFile)
		}
	} else {
		m.Logger.Debug.Printf("Source build request!\n")
		pkgReq, err = m.assembleBuildPkgReq(r, args...)
	}
	// Final check
	if err != nil {
		m.Logger.Warning.Printf("%s\n", err)
		return nil, map[string]string{"error": err.Error()}, 400
	}
	// Do not send over channel
	if _, ok := r.URL.Query()["dryrun"]; !ok {
		m.RequestChan <- pkgReq
	}

	return nil, pkgReq, 200
}

func (m *PkgBuilderMethodHandler) GET(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	rslts, err := m.Datastore.GetBuildInfo(args...)
	/*
		Throws an error if index is empty while sorting on '_timestamp'
		as it's an internal field and will only be available after at least 1
		item has been indexed.
	*/
	if err != nil {
		m.Logger.Error.Printf("Could not get history: %s\n", err)
		return nil, fmt.Sprintf(`{"error": "%s"}`, err), 400
	}
	return nil, rslts, 200
}

/*
 * Parse FORM params into a PackageRequest struct.
 * Contains the pre-built tarball
 *
 * Initiate user binary package upload. This is a tarball containing
 * the filesystem overlay, binaries and all. This uses http FORM params
 *
 * e.g. curl -XPOST http://.... -F package=@path/to/file.tgz -F conf=/path/to/conf.json
 */
func (m *PkgBuilderMethodHandler) assembleBuiltPkgReq(req *http.Request, args ...string) (request.PackageRequest, *multipart.FileHeader, error) {
	//pkgr, pkgname, pkgversion
	var (
		pkgReq  = request.NewPackageRequest(args[1])
		pkgFile *multipart.FileHeader
		err     error
	)

	if args[2] != "" {
		pkgReq.Version = args[2]
	}
	pkgReq.Package.Packager = args[0]

	if err = req.ParseMultipartForm(128); err != nil {
		return *pkgReq, pkgFile, err
	}
	// Binary tarball //
	if uFiles, ok := req.MultipartForm.File[HTTP_REQ_PARAMS["PACKAGE"]]; ok {
		pkgFile = uFiles[0]
	} else {
		return *pkgReq, pkgFile, fmt.Errorf("'package' not specified!")
	}
	// Check if build config provided //
	confFile, ok := req.MultipartForm.File[HTTP_REQ_PARAMS["BUILD_CONF"]]
	if !ok {
		return *pkgReq, pkgFile, fmt.Errorf("Build config not provided (conf)!")
	}
	// Build config to PackageRequest //
	if err = m.getPackageRequestFromConf(confFile[0], pkgReq); err != nil {
		return *pkgReq, pkgFile, err
	}

	pkgReq.Package.BuildType = request.BUILDTYPE_BIN
	// TODO: Fix to account for repository
	pkgReq.Package.Path = fmt.Sprintf("%s/%s/%s",
		pkgReq.Package.Name, pkgReq.Package.Version, pkgFile.Filename)

	return *pkgReq, pkgFile, nil
}

/*
	Returns:
		PackageRequest assembled from user supplied data
		error
*/
func (m *PkgBuilderMethodHandler) assembleBuildPkgReq(r *http.Request, args ...string) (request.PackageRequest, error) {
	//pkgr, pkgname, pkgversion
	pkgReq := request.NewPackageRequest(args[1])
	pkgReq.Version = args[2]
	pkgReq.Package.Version = args[2]
	pkgReq.Package.Packager = args[0]

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return *pkgReq, err
	}
	if err = json.Unmarshal(b, pkgReq); err != nil {
		return *pkgReq, err
	}
	if err = pkgReq.Validate(false); err != nil {
		return *pkgReq, err
	}
	//pkgReq.Package.InitScript, _ = initscript.NewBasicInitScript(pkgReq.Name)
	pkgReq.Package.BuildType = request.BUILDTYPE_SOURCE
	// TODO: Fix to account for repository
	pkgReq.Package.Path = fmt.Sprintf("%s/%s/%s",
		pkgReq.Package.Name, pkgReq.Package.Version, pkgReq.Package.Name)

	return *pkgReq, nil
}

func (m *PkgBuilderMethodHandler) getPackageRequestFromConf(usrFile *multipart.FileHeader, pkgReq *request.PackageRequest) error {
	fh, err := usrFile.Open()
	if err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	if _, err = io.Copy(buff, fh); err != nil {
		return err
	}
	if err = json.Unmarshal(buff.Bytes(), pkgReq); err != nil {
		return err
	}

	return pkgReq.Validate(true)
}

func (m *PkgBuilderMethodHandler) downloadUserPackage(pkg *request.UserPackage, pkgfile *multipart.FileHeader) error {
	pkgFilepath := m.Repository.BuildDir(pkg)
	os.MkdirAll(pkgFilepath, 0755)

	pkgFilepath += "/" + pkgfile.Filename
	fh, err := os.Create(pkgFilepath)
	if err != nil {
		return err
	}
	defer fh.Close()

	usrFh, err := pkgfile.Open()
	if err != nil {
		return err
	}
	defer usrFh.Close()

	if _, err = io.Copy(fh, usrFh); err != nil {
		return err
	}
	return nil
}

func SetupBuildHandler(cfg *config.AppConfig, datastore *tracker.TrackerStore,
	repo repository.BuildRepository, reqChan chan request.PackageRequest, logger *logging.Logger) {

	methodHandler := PkgBuilderMethodHandler{
		Config:      cfg,
		Repository:  repo,
		Logger:      logger,
		RequestChan: reqChan, // testing
		Datastore:   datastore,
	}

	NewRestHandler(cfg.Endpoints.Builder, &methodHandler, logger)
	logger.Warning.Printf("Builder API: %s\n", cfg.Endpoints.Builder)
}
