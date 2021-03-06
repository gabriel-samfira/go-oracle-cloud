// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package api

import (
	"fmt"
	"net/http"
)

// Authenticate this request returns an authentication token in the Set-Cookie response header.
// The token expires after 30 minutes. A valid (that is, unexpired) authentication
// token must be included in every request to the service, in the Cookie: request header.
// The client making the API call must examine the cookie expiry time and discard it if the cookie has expired.
// Requests sent with expired cookies will result in an Unauthorized error in the response.
func (c *Client) Authenticate() (err error) {
	if c.isAuth() {
		return ErrAlreadyAuth
	}

	// build the json authentication
	auth := map[string]string{
		"user":     fmt.Sprintf("/Compute-%s/%s", c.identify, c.username),
		"password": c.password,
	}

	return request(paramsRequest{
		client: &c.http,
		cookie: c.cookie,
		url:    fmt.Sprintf("%s/%s/", c.endpoint, "authenticate"),
		verb:   "POST",
		body:   auth,
		treat: func(resp *http.Response) (err error) {
			// if the operation is successful then we will recive 204 http status
			// if this is not the case then we should stop and return a friendly error
			if resp.StatusCode != http.StatusNoContent {
				return fmt.Errorf(
					"go-oracle-cloud: Error api response %d %s",
					resp.StatusCode, dumpApiError(resp.Body),
				)
			}

			// the orcale api uses cookies to manage sessions
			// once a cookie is taken then we can make
			// more connections to other api resources
			cookies := resp.Cookies()
			if len(cookies) != 1 {
				return fmt.Errorf("go-oracle-cloud: Invalid number of session cookies: %q", cookies)
			}
			// take the cookie
			c.cookie = cookies[0]
			return nil
		},
		resp: nil,
	})
}
