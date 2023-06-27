package traefik_query_append_url_test

import (
	"context"
	traefikqueryappendurl "github.com/axyi/traefik-query-append-url"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddQueryParam_NoPrevious(t *testing.T) {
	cfg := traefikqueryappendurl.CreateConfig()
	cfg.QueryParamName = "url"
	expected := "url=http%3A%2F%2Flocalhost%2Ftest"

	assertQueryModification(t, cfg, "", expected)
}

func TestAddQueryParam_OtherPrevious(t *testing.T) {
	cfg := traefikqueryappendurl.CreateConfig()
	cfg.QueryParamName = "url"
	expected := "a=b&url=http%3A%2F%2Flocalhost%2Ftest"
	previous := "a=b"

	assertQueryModification(t, cfg, previous, expected)
}

func TestAddQueryParam_AddPrevious(t *testing.T) {
	cfg := traefikqueryappendurl.CreateConfig()
	cfg.QueryParamName = "newparam"
	expected := "newparam=oldvalue&newparam=http%3A%2F%2Flocalhost%2Ftest"
	previous := "newparam=oldvalue"

	assertQueryModification(t, cfg, previous, expected)
}

func TestAddQueryParam_Previous(t *testing.T) {
	cfg := traefikqueryappendurl.CreateConfig()
	cfg.QueryParamName = "newparam"
	expected := "a=b&newparam=http%3A%2F%2Flocalhost%2Ftest"
	previous := "a=b"

	assertQueryModification(t, cfg, previous, expected)
}

func createReqAndRecorder(cfg *traefikqueryappendurl.Config) (http.Handler, error, *httptest.ResponseRecorder, *http.Request) {
	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	handler, err := traefikqueryappendurl.New(ctx, next, cfg, "query-modification-plugin")
	if err != nil {
		return nil, err, nil, nil
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/test", nil)
	return handler, err, recorder, req
}

func assertQueryModification(t *testing.T, cfg *traefikqueryappendurl.Config, previous, expected string) {
	handler, err, recorder, req := createReqAndRecorder(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	req.URL.RawQuery = previous
	handler.ServeHTTP(recorder, req)

	if req.URL.Query().Encode() != expected {
		t.Errorf("Expected %s, got %s", expected, req.URL.Query().Encode())
	}
}
