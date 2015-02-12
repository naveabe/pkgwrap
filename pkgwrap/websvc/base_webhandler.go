package websvc

import (
	"encoding/json"
	"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"net/http"
	"strings"
)

type IMethodHandler interface {
	GET(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int)
	POST(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int)
	PUT(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int)
	DELETE(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int)
	OPTIONS(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int)
}

type DefaultMethodHandler struct{}

func (d *DefaultMethodHandler) GET(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) PATCH(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) DELETE(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) OPTIONS(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) POST(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}
func (d *DefaultMethodHandler) PUT(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {
	return map[string]string{}, map[string]string{"error": "Method not allowed"}, 405
}

type RestHandler struct {
	Prefix        string
	MethodHandler IMethodHandler
	logger        *logging.Logger
}

func NewRestHandler(prefix string, methodHandler IMethodHandler, logger *logging.Logger) *RestHandler {
	if logger == nil {
		logger = logging.NewStdLogger()
	}
	p := RestHandler{prefix, methodHandler, logger}

	http.Handle(prefix, &p)
	if strings.HasSuffix(prefix, "/") {
		http.Handle(prefix[:len(prefix)-1], &p)
	} else {
		http.Handle(prefix+"/", &p)
	}
	return &p
}

func (p *RestHandler) parsePath(path string) []string {
	sPath := strings.TrimPrefix(path, p.Prefix)
	parts := make([]string, 0)
	for _, part := range strings.Split(sPath, "/") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func (p *RestHandler) callMethod(w http.ResponseWriter, r *http.Request, args ...string) (map[string]string, interface{}, int) {

	switch r.Method {
	case "GET":
		return p.MethodHandler.GET(w, r, args...)
	case "POST":
		return p.MethodHandler.POST(w, r, args...)
	case "DELETE":
		return p.MethodHandler.DELETE(w, r, args...)
	case "PUT":
		return p.MethodHandler.PUT(w, r, args...)
	case "OPTIONS":
		return p.MethodHandler.OPTIONS(w, r, args...)
	default:
		return nil, fmt.Sprintf(`{"error": "Method not supported: %s"}`, r.Method), 405
	}
}

func (p *RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathParts := p.parsePath(r.URL.Path)

	if len(pathParts) < 1 {
		p.writeJsonResponse(w, r, nil, []byte(fmt.Sprintf(`{"error": "Bad request: %s", "code": 404}`, r.URL.Path)), 404)
		return
	}

	headers, data, code := p.callMethod(w, r, pathParts...)
	// Don't write response.  The method must write the response
	// and return -1 for the status code.
	if code != -1 {
		p.writeJsonResponse(w, r, headers, data, code)
	} else {
		p.logger.Trace.Printf("Not writing HTTP response!\n")
	}
}

func (p *RestHandler) writeResponse(w http.ResponseWriter, r *http.Request, headers map[string]string, data []byte, respCode int) {
	if headers != nil {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	w.WriteHeader(respCode)
	w.Write(data)
	p.logger.Info.Printf("%s %d %s\n", r.Method, respCode, r.URL.RequestURI())
}

func (p *RestHandler) writeJsonResponse(w http.ResponseWriter, r *http.Request, headers map[string]string, data interface{}, respCode int) {
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
		if _, ok := r.URL.Query()["pretty"]; ok {
			b, _ = json.MarshalIndent(&data, "", "  ")
		} else {
			b, _ = json.Marshal(&data)
		}
		break
	}

	w.Header().Set("Content-Type", "application/json")
	p.writeResponse(w, r, headers, b, respCode)
}
