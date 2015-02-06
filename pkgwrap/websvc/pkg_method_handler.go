package websvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"github.com/naveabe/pkgwrap/pkgwrap/repository"
	"github.com/naveabe/pkgwrap/pkgwrap/specer"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	RequestChan chan specer.PackageRequest
}

func (m *PkgBuilderMethodHandler) getPackageRequestFromConf(usrFile *multipart.FileHeader, pkgReq *specer.PackageRequest) error {
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

/*
 * Parse FORM params into a PackageRequest struct.
 * Contains the pre-built tarball
 *
 * Initiate user binary package upload. This is a tarball containing
 * the filesystem overlay, binaries and all. This uses http FORM params
 *
 * e.g. curl -XPOST http://.... -F package=@path/to/file.tgz -F conf=/path/to/conf.json
 */
func (m *PkgBuilderMethodHandler) assembleBuiltPkgReq(req *http.Request, pkgr, pkgName, pkgVersion string) (specer.PackageRequest, *multipart.FileHeader, error) {
	var (
		pkgReq  = specer.NewPackageRequest(pkgName)
		pkgFile *multipart.FileHeader
		err     error
	)

	if pkgVersion != "" {
		pkgReq.Version = pkgVersion
	}
	pkgReq.Package.Packager = pkgr

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

	pkgReq.Package.BuildType = specer.BUILDTYPE_BIN
	pkgReq.Package.Path = fmt.Sprintf("%s/%s/%s",
		pkgReq.Package.Name, pkgReq.Package.Version, pkgFile.Filename)

	return *pkgReq, pkgFile, nil
}

/*
	Returns:
		PackageRequest assembled from user supplied data
		error
*/
func (m PkgBuilderMethodHandler) assembleBuildPkgReq(r *http.Request, pkgr, pkgname, pkgversion string) (specer.PackageRequest, error) {
	pkgReq := specer.NewPackageRequest(pkgname)
	pkgReq.Version = pkgversion
	pkgReq.Package.Version = pkgversion
	pkgReq.Package.Packager = pkgr

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

	pkgReq.Package.BuildType = specer.BUILDTYPE_SOURCE

	pkgReq.Package.Path = fmt.Sprintf("%s/%s/%s",
		pkgReq.Package.Name, pkgReq.Package.Version, pkgReq.Package.Name)

	return *pkgReq, nil
}

func (m PkgBuilderMethodHandler) downloadUserPackage(pkgr, pkgname, pkgversion string, pkgfile *multipart.FileHeader) error {
	pkgFilepath := filepath.Join(m.Config.Repository, pkgr, pkgname, pkgversion)
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

func (m *PkgBuilderMethodHandler) POST(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {

	var (
		contentType       = r.Header.Get("Content-Type")
		err         error = nil

		pkgReq  specer.PackageRequest
		pkgFile *multipart.FileHeader
	)

	if len(args) < 3 {
		return nil, map[string]string{"error": "Invalid request"}, 400
	}

	// Determines build type: 'source' or 'binary'
	if strings.HasPrefix(contentType, "application/json") {
		pkgReq, err = m.assembleBuildPkgReq(r, args[0], args[1], args[2])
	} else {
		pkgReq, pkgFile, err = m.assembleBuiltPkgReq(r, args[0], args[1], args[2])
		if err == nil {
			err = m.downloadUserPackage(args[0], args[1], args[2], pkgFile)
		}
	}
	// Final check
	if err != nil {
		m.Logger.Warning.Printf("%s\n", err)
		return nil, map[string]string{"error": err.Error()}, 400
	}

	params := r.URL.Query()
	if _, ok := params["dryrun"]; !ok {

		m.RequestChan <- pkgReq
	}

	return nil, pkgReq, 200
}
