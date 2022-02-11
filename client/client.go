// Copyright (c) 2021 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Client struct {
	client  *http.Client
	timeout int
	baseUrl *url.URL
}

type ClientOption func(*Client) error

func NewClient(options ...ClientOption) (*Client, error) {
	client := Client{
		client:  &http.Client{},
		timeout: 15,
		baseUrl: new(url.URL),
	}

	for _, option := range options {
		if err := option(&client); err != nil {
			return nil, err
		}
	}

	return &client, nil
}

func EnableServerValidation(enableServerValidation bool) func(*Client) error {
	return func(c *Client) error {
		c.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !enableServerValidation},
		}

		return nil
	}
}

func SetAuthentication(user, password string) func(*Client) error {
	return func(c *Client) error {
		c.baseUrl.User = url.UserPassword(user, password)
		return nil
	}
}

func SetBaseUrl(baseUrl string) func(*Client) error {
	return func(c *Client) error {
		u, err := url.Parse(baseUrl)
		if err != nil {
			return err
		}

		u.User = c.baseUrl.User
		c.baseUrl = u

		return nil
	}
}

func (c *Client) Timeout() int {
	return c.timeout
}

func (c *Client) DoRequest(method, reqUrlPath string, body []byte) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeout)*time.Second)
	defer cancel()

	reqUrl := *c.baseUrl // copy baseUrl
	reqUrl.Path = path.Join(reqUrl.Path, reqUrlPath)

	req, err := http.NewRequestWithContext(ctx, method, reqUrl.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	return c.client.Do(req)
}
