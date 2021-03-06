// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package api_test

import (
	"github.com/hoenirvili/go-oracle-cloud/api"
	gc "gopkg.in/check.v1"
)

type clientTest struct{}

var _ = gc.Suite(&clientTest{})

// NewConfig returns a predefined config for testing
// the client interactions
func (cli clientTest) NewConfig(c *gc.C) api.Config {
	return api.Config{
		Username: "oracleusername@oracle.com",
		Password: "Password123",
		Identify: "myIdentify",
		Endpoint: "http://localhost",
	}
}

func (cl clientTest) TestNewClient(c *gc.C) {
	cli, err := api.NewClient(api.Config{})
	c.Assert(err, gc.NotNil)
	c.Assert(cli, gc.IsNil)
	c.Assert(err.Error(), gc.DeepEquals,
		"go-oracle-cloud: Empty identify endpoint name")

	cli, err = api.NewClient(api.Config{
		Identify: "someIdentify",
	})
	c.Assert(err, gc.NotNil)
	c.Assert(cli, gc.IsNil)
	c.Assert(err.Error(), gc.DeepEquals,
		"go-oracle-cloud: Empty client username")

	cli, err = api.NewClient(api.Config{
		Identify: "someIdentify",
		Username: "sometestuser",
	})
	c.Assert(err, gc.NotNil)
	c.Assert(cli, gc.IsNil)
	c.Assert(err.Error(), gc.DeepEquals,
		"go-oracle-cloud: Empty client password")

	cli, err = api.NewClient(api.Config{
		Identify: "someIdentify",
		Username: "sometestuser",
		Password: "providesomepasswrd",
	})
	c.Assert(err, gc.NotNil)
	c.Assert(cli, gc.IsNil)
	c.Assert(err.Error(), gc.DeepEquals,
		"go-oracle-cloud: Empty endpoint url basepath")

	cli, err = api.NewClient(api.Config{
		Identify: "someIdentify",
		Username: "sometestuser",
		Password: "providesomepasswrd",
		Endpoint: "s",
	})

	c.Assert(err, gc.NotNil)
	c.Assert(cli, gc.IsNil)
	c.Assert(err.Error(), gc.DeepEquals,
		"go-oracle-cloud: The endpoint provided is invalid")

	// provide some valid configuration
	cli, err = api.NewClient(cl.NewConfig(c))
	c.Assert(err, gc.IsNil)
	c.Assert(cli, gc.NotNil)
}
