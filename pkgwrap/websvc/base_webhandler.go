package websvc

import (
	"encoding/json"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"net/http"
	"strings"
)

type IMethodHandler interface {
	GET(r *http.Request, args ...string) (map[string]string, interface{}, int)
	POST(r *http.Request, args ...string) (map[string]string, interface{}, int)
	PUT(r *http.Request, args ...string) (map[string]string, interface{}, int)
	DELETE(r *http.Request, args ...string) (map[string]string, interface{}, int)
	OPTIONS(r *http.Request, args ...string) (map[string]string, interface{}, int)
}

type DefaultMethodHandler struct{}

func (d *DefaultMethodHandler) GET(r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) PATCH(r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) DELETE(r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) OPTIONS(r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) POST(r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) PUT(r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}

type PkgBuilderHandler struct {
	Prefix        string
	MethodHandler IMethodHandler
	logger        *logging.Logger
}

func NewPkgBuilderHandler(prefix string, methodHandler IMethodHandler, logger *logging.Logger) *PkgBuilderHandler {
	if logger == nil {
		logger = logging.NewStdLogger()
	}
	p := PkgBuilderHandler{prefix, methodHandler, logger}

	http.Handle(prefix, &p)
	if strings.HasSuffix(prefix, "/") {
		http.Handle(prefix[:len(prefix)-1], &p)
	} else {
		http.Handle(prefix+"/", &p)
	}
	return &p
}

func (p *PkgBuilderHandler) parsePath(path string) []string {
	sPath := strings.TrimPrefix(path, p.Prefix)
	parts := make([]string, 0)
	for _, part := range strings.Split(sPath, "/") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func (p *PkgBuilderHandler) callMethod(r *http.Request, pkgName, pkgVersion string) (map[string]string, interface{}, int) {

	switch r.Method {
	case "GET":
		return p.MethodHandler.GET(r, pkgName, pkgVersion)
	case "POST":
		return p.MethodHandler.POST(r, pkgName, pkgVersion)
	case "DELETE":
		return p.MethodHandler.DELETE(r, pkgName, pkgVersion)
	case "PUT":
		return p.MethodHandler.PUT(r, pkgName, pkgVersion)
	case "OPTIONS":
		return p.MethodHandler.OPTIONS(r, pkgName, pkgVersion)
	default:
		return nil, fmt.Sprintf(`{"error": "Method not supported: %s"}`, r.Method), 405
	}
}

func (p *PkgBuilderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		pkgName    string
		pkgVersion = ""
		pathParts  = p.parsePath(r.URL.Path)
	)
	//p.logger.Trace.Printf("%v\n", pathParts)

	switch len(pathParts) {
	case 1:
		pkgName = pathParts[0]
		break
	case 2:
		pkgName = pathParts[0]
		pkgVersion = pathParts[1]
		break
	default:
		p.writeJsonResponse(w, r, nil, []byte(fmt.Sprintf(`{"error": "Bad request: %s", "code": 404}`, r.URL.Path)), 404)
		return
	}

	headers, data, code := p.callMethod(r, pkgName, pkgVersion)
	p.writeJsonResponse(w, r, headers, data, code)
}

func (p *PkgBuilderHandler) writeResponse(w http.ResponseWriter, r *http.Request, headers map[string]string, data []byte, respCode int) {
	if headers != nil {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(respCode)
	w.Write(data)
	p.logger.Info.Printf("%s %d %s\n", r.Method, respCode, r.URL.RequestURI())
}

func (p *PkgBuilderHandler) writeJsonResponse(w http.ResponseWriter, r *http.Request, headers map[string]string, data interface{}, respCode int) {
	var b []byte
	switch data.(type) {
	case string:
		s, _ := data.(string)
		b = []byte(s)
		break
	case []byte:
		b, _ = data.([]byte)
		break
	default:
		b, _ = json.Marshal(&data)
		break
	}

	w.Header().Set("Content-Type", "application/json")
	p.writeResponse(w, r, nil, b, respCode)
}
