package api_test

import (
	enc "encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/juju/go-oracle-cloud/api"
	gc "gopkg.in/check.v1"
)

type clientTest struct{}

var _ = gc.Suite(&clientTest{})

const (
	// content-type values for json requests and directory request
	json = "application/oracle-compute-v3+json"
	dir  = "application/oracle-compute-v3+directory+json"
)

const (
	username = "gooraclecloudadmin@oracle.com"
	password = "some1337eleetpasswd"
	identify = "qraffd"
)

// findHostAndPortFromURL extracts the
// host and port part of a full url
// (like http://address:port/path#fragment)
func findHostAndPortFromURL(rawurl string) (string, int, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return "", 0, err
	}
	host, port, err := net.SplitHostPort(url.Host)
	if err != nil {
		return "", 0, err
	}
	iport, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return host, iport, nil
}

type httpParams struct {
	// body for providing the test function
	// to send content back to the api client
	body []byte

	// check is the assert library
	check *gc.C

	// handler for manual testing the
	// http.ResponseWriter and http.Reqeust
	handler http.HandlerFunc

	// handler manual the header
	// and header status
	manualHeaderStatus bool
}

func hdr(value string, w http.ResponseWriter, r *http.Request) (err error) {
	switch value {
	case json:
		w.Header().Set("Accept", json)
		w.Header().Set("Content-Type", json)
		return nil
	case dir:
		w.Header().Set("Accept", dir)
		w.Header().Set("Content-Type", dir)
		return nil
	default:
		return fmt.Errorf("oracle api does not support this header value %s", value)
	}
}

func handlerHeaderStatus(w http.ResponseWriter, r *http.Request) (err error) {
	value := r.Header.Get("Content-Type")
	switch r.Method {
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		return hdr(value, w, r)
	case http.MethodGet:
		return hdr(r.Header.Get("Accept"), w, r)
	case http.MethodPut:
		w.WriteHeader(http.StatusOK)
		return hdr(value, w, r)
	case http.MethodDelete:
		w.WriteHeader(http.StatusNoContent)
		return hdr(value, w, r)
	default:
		return fmt.Errorf("oracle api does not support this method %s", r.Method)
	}
}

// StartTestServer will start an httptest server on a random port
// with the given httpParams and then return the oracle client implementation
func (cli clientTest) StartTestServer(
	params httpParams,
) (*httptest.Server, *api.Client) {
	// create a new http server for testing
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			if !params.manualHeaderStatus {
				// handler the header if it's somthing wrong then
				// handler the status based on the method
				err := handlerHeaderStatus(w, r)
				params.check.Assert(err, gc.IsNil)
			}

			// treat the w and *r here in a custom way
			if params.handler != nil {
				params.handler(w, r)
			}

			// if the caller is expecting the
			// server to send some response,
			// we should include that response from body
			// and send it back
			m := len(params.body)
			if m > 0 {
				// write the response
				n, err := w.Write(params.body)
				params.check.Assert(err, gc.IsNil)
				params.check.Assert(m, gc.Equals, n)
			}

			defer func() {
				err := r.Body.Close()
				params.check.Assert(err, gc.IsNil)
			}()
		}))

	// find the host and port from the http server url
	host, port, err := findHostAndPortFromURL(ts.URL)
	// always test if the finding is valid
	params.check.Assert(err, gc.IsNil)

	// create a new config
	cfg := api.Config{
		Username: username,
		Password: password,
		Identify: identify,
		Endpoint: fmt.Sprintf("http://%s:%d", host, port),
	}

	// create a new client based on the config
	client, err := api.NewClient(cfg)
	params.check.Assert(err, gc.IsNil)
	params.check.Assert(client, gc.NotNil)

	// return the new started server and the oracle client
	return ts, client
}

func makeClientAuth(
	c *gc.C,
	w http.ResponseWriter,
	r *http.Request,
) {
	c.Assert(r.Method, gc.Equals, http.MethodPost)
	c.Assert(r.Header.Get("Content-Type"), gc.DeepEquals, json)
	c.Assert(r.Header.Get("Accept"), gc.DeepEquals, json)

	auth := struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}{}

	err := enc.NewDecoder(r.Body).Decode(&auth)
	c.Assert(err, gc.IsNil)
	c.Assert(auth.User, gc.DeepEquals, fmt.Sprintf("/Compute-%s/%s",
		identify, username))
	c.Assert(auth.Password, gc.DeepEquals, password)

	// give the client a new cookie
	w.Header().Set("Set-Cookie", cookie)
	w.Header().Set("Content-Type", json)
	w.WriteHeader(http.StatusNoContent)
}

// StartTestServer will start an httptest server on a random port
// with the given httpParams and then return the oracle client implementation
// that has been already authenticated
func (cl clientTest) StartTestServerAuth(
	params httpParams,
) (*httptest.Server, *api.Client) {
	i := 0

	// create a new http server for testing
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			// we are dealing with the first auth request
			if i == 0 {
				makeClientAuth(params.check, w, r)
				i++
				return
			}

			if !params.manualHeaderStatus {
				// handler the header if it's somthing wrong then
				// handler the status based on the method
				err := handlerHeaderStatus(w, r)
				params.check.Assert(err, gc.IsNil)
			}

			// treat the w and *r here in a custom way
			if params.handler != nil {
				params.handler(w, r)
			}

			// if the caller is expecting the
			// server to send some response,
			// we should include that response from body
			// and send it back
			m := len(params.body)
			if m > 0 {
				// write the response
				n, err := w.Write(params.body)
				params.check.Assert(err, gc.IsNil)
				params.check.Assert(m, gc.Equals, n)
			}

			defer func() {
				err := r.Body.Close()
				params.check.Assert(err, gc.IsNil)
			}()
		}))

	// find the host and port from the http server url
	host, port, err := findHostAndPortFromURL(ts.URL)
	// always test if the finding is valid
	params.check.Assert(err, gc.IsNil)

	// create a new config
	cfg := api.Config{
		Username: username,
		Password: password,
		Identify: identify,
		Endpoint: fmt.Sprintf("http://%s:%d", host, port),
	}

	// create a new client based on the config
	client, err := api.NewClient(cfg)
	params.check.Assert(err, gc.IsNil)
	params.check.Assert(client, gc.NotNil)

	// make the client authenticate
	err = client.Authenticate()
	params.check.Assert(err, gc.IsNil)

	// return the new started server and the oracle client
	// that has been already authenticated
	return ts, client

}
