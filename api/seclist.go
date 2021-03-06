// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hoenirvili/go-oracle-cloud/response"
)

// CreatesSecList a security list. After creating security
// lists, you can add instances to them by using the HTTP request,
// CreateSecAssociation (Create a Security Association).
func (c Client) CreateSecList(
	description string,
	name string,
	outbound_cidr_policy string,
	policy string,
) (resp response.SecList, err error) {
	if !c.isAuth() {
		return resp, ErrNotAuth
	}

	if name == "" {
		return resp, errors.New("go-oracle-cloud: Empty secure list name")
	}

	params := struct {
		Description          string `json:"description,omitempty"`
		Name                 string `json:"name"`
		Outbound_cidr_policy string `json:"outbound_cidr_policy"`
		Policy               string `json:"policy"`
	}{
		Description:          description,
		Name:                 name,
		Outbound_cidr_policy: outbound_cidr_policy,
		Policy:               policy,
	}

	url := fmt.Sprintf("%s/seclist/", c.endpoint)

	if err = request(paramsRequest{
		client: &c.http,
		cookie: c.cookie,
		url:    url,
		body:   &params,
		verb:   "POST",
		treat:  defaultPostTreat,
		resp:   &resp,
	}); err != nil {
		return resp, err
	}

	strip(&resp.Account)
	strip(&resp.Name)

	return resp, nil
}

// DeleteSecList the specified security list. No response is returned.<Paste>
func (c Client) DeleteSecList(name string) (err error) {
	if !c.isAuth() {
		return ErrNotAuth
	}

	if name == "" {
		return errors.New("go-oracle-cloud: Empty secure list")
	}

	url := fmt.Sprintf("%s/seclist/Compute-%s/%s/%s",
		c.endpoint, c.identify, c.username, name)

	if err = request(paramsRequest{
		client: &c.http,
		cookie: c.cookie,
		url:    url,
		verb:   "DELETE",
		treat:  defaultDeleteTreat,
	}); err != nil {
		return err
	}

	return nil
}

// AllSecList retrieves details of the security lists that are in the specified
// container and match the specified query criteria.
func (c Client) AllSecList() (resp response.AllSecList, err error) {
	if !c.isAuth() {
		return resp, ErrNotAuth
	}

	url := fmt.Sprintf("%s/seclist/Compute-%s/%s/",
		c.endpoint, c.identify, c.username)

	if err = request(paramsRequest{
		client: &c.http,
		cookie: c.cookie,
		url:    url,
		verb:   "GET",
		treat:  defaultTreat,
		resp:   &resp,
	}); err != nil {

		return resp, err
	}

	for key, _ := range resp.Result {
		strip(&resp.Result[key].Account)
		strip(&resp.Result[key].Name)
	}

	return resp, nil
}

// SecListDetails retrieves information about the specified security list.
func (c Client) SecListDetails(name string) (resp response.SecList, err error) {
	if !c.isAuth() {
		return resp, ErrNotAuth
	}

	url := fmt.Sprintf("%s/seclist/Compute-%s/%s/%s",
		c.endpoint, c.identify, c.username, name)

	if err = request(paramsRequest{
		client: &c.http,
		cookie: c.cookie,
		url:    url,
		verb:   "GET",
		treat:  defaultTreat,
		resp:   &resp,
	}); err != nil {
		return resp, err
	}

	strip(&resp.Account)
	strip(&resp.Name)

	return resp, nil
}

// Updates inbound policy, outbound policy, and description for
// the specified security list.
// newName could be "" if you don't want to change the name
// but it's required at leas to have a currentName
// outbound_cidr_policy is the policy for outbound traffic
// from the security list. You can specify one of the following values:
// deny: Packets are dropped. No response is sent.
// reject: Packets are dropped, but a response is sent.
// permit(default): Packets are allowed.
func (c Client) UpdateSecList(
	description string,
	currentName string,
	newName string,
	outbound_cidr_policy string,
	policy string,
) (resp response.SecList, err error) {
	if !c.isAuth() {
		return resp, ErrNotAuth
	}

	if currentName == "" {
		return resp, errors.New("go-oracle-cloud: Empty secure list name")
	}

	if newName == "" {
		newName = currentName
	}

	params := struct {
		Policy               string `json:"policy"`
		Description          string `json:"description,omitempty"`
		Name                 string `json:"name"`
		Outbound_cidr_policy string `json:"outbound_cidr_policy"`
	}{
		Description: description,
		Name: fmt.Sprintf("/Compute-%s/%s/%s",
			c.identify, c.username, newName),
		Outbound_cidr_policy: strings.ToUpper(outbound_cidr_policy),
		Policy:               strings.ToUpper(policy),
	}

	url := fmt.Sprintf("%s/seclist/Compute-%s/%s/%s",
		c.endpoint, c.identify, c.username, currentName)

	if err = request(paramsRequest{
		client: &c.http,
		cookie: c.cookie,
		url:    url,
		body:   &params,
		verb:   "PUT",
		treat:  defaultTreat,
		resp:   &resp,
	}); err != nil {
		return resp, err
	}

	strip(&resp.Account)
	strip(&resp.Name)

	return resp, nil
}
