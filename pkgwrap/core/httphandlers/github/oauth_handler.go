package github

import (
	"encoding/json"
	//"fmt"
	"github.com/naveabe/pkgwrap/pkgwrap/config"
	"github.com/naveabe/pkgwrap/pkgwrap/core/httphandlers"
	"github.com/naveabe/pkgwrap/pkgwrap/logging"
	"io/ioutil"
	"net/http"
)

type GithubAccessToken struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	TokenType string `json:"token_type"`
}

type GithubOauthHandler struct {
	httphandlers.DefaultMethodHandler

	logger *logging.Logger
	cfg    config.CodeRepoConfig
}

func NewGithubOauthHandler(cfg config.CodeRepoConfig, logger *logging.Logger) *GithubOauthHandler {
	g := GithubOauthHandler{cfg: cfg}
	if logger == nil {
		g.logger = logging.NewStdLogger()
	} else {
		g.logger = logger
	}
	return &g
}

func (g *GithubOauthHandler) getAccessToken(code string) (*GithubAccessToken, error) {
	var (
		body  []byte
		err   error
		token GithubAccessToken
	)

	urlStr := g.cfg.OAuth.TokenURL +
		"?code=" + code +
		"&client_id=" + g.cfg.OAuth.ClientId +
		"&client_secret=" + g.cfg.OAuth.ClientSecret

	client := http.Client{}
	req, err := http.NewRequest("POST", urlStr, nil)
	req.Header["Accept"] = []string{"application/json"}

	resp, err := client.Do(req)
	if err != nil {
		return &token, err
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return &token, err
	}
	g.logger.Trace.Printf("%s\n", body)
	if err = json.Unmarshal(body, &token); err != nil {
		return &token, err
	}

	return &token, nil
}

func (g *GithubOauthHandler) GET(w http.ResponseWriter, r *http.Request,
	args ...string) (map[string]string, interface{}, int) {

	var (
		code       []string
		redirectTo = g.cfg.OAuth.ClientRedirect
		ok         bool
		err        error
		token      *GithubAccessToken
	)

	if code, ok = r.URL.Query()["code"]; ok {
		if token, err = g.getAccessToken(code[0]); err == nil {
			redirectTo += "?access_token=" + token.Token +
				"&token_type=" + token.TokenType +
				"&scope=" + token.Scope
		} else {
			redirectTo += "?error=unauthorized_client" +
				"&error_description=" + err.Error()
		}
	} else {
		redirectTo += "?error=invalid_request&error_description=code missing"
	}

	g.logger.Trace.Printf("%s\n", token)

	http.Redirect(w, r, redirectTo, 302)

	return nil, nil, -1
}

func SetupGithubOauthHandler(cfgs []config.CodeRepoConfig, logger *logging.Logger) {
	for _, v := range cfgs {
		if v.Type == "github" {
			ghOauthHandler := NewGithubOauthHandler(v, logger)
			httphandlers.NewRestHandler(v.LocalEndpoint, ghOauthHandler, logger)
			logger.Warning.Printf("Github oauth API: %s\n", v.LocalEndpoint)
			return
		}
	}
	logger.Warning.Printf("Github config not found. Disabling...!\n")
}
