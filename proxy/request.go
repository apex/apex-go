//
// Copyright 2017 Alsanium, SAS. or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Adapted from https://github.com/eawsy/aws-lambda-go-net/blob/master/service/lambda/runtime/net/apigatewayproxy/server.go
// Changes (kothar - 2017-02):
//   - Relocated to go-apex/proxy
//   - All code not related to constructing an http.Request removed
//   - Remaining code placed in buildRequest function

package proxy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/apex/go-apex"
)

// Constructs an http.Request object from a proxyEvent
func buildRequest(proxyEvent *Event, ctx *apex.Context) (*http.Request, error) {
	u, err := url.Parse(proxyEvent.Path)
	if err != nil {
		return nil, fmt.Errorf("Parse request path: %s", err)
	}
	q := u.Query()
	for k, v := range proxyEvent.QueryStringParameters {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	dec := proxyEvent.Body
	if proxyEvent.IsBase64Encoded {
		data, err := base64.StdEncoding.DecodeString(dec)
		if err != nil {
			return nil, fmt.Errorf("Decode base64 request body: %s", err)
		}
		dec = string(data)
	}

	req, err := http.NewRequest(proxyEvent.HTTPMethod, u.String(), strings.NewReader(dec))
	if err != nil {
		return nil, fmt.Errorf("Create request: %s", err)
	}

	proxyEvent.Body = "... truncated"

	for k, v := range proxyEvent.Headers {
		req.Header.Set(k, v)
	}
	hbody, err := json.Marshal(proxyEvent)
	if err != nil {
		return nil, fmt.Errorf("Marshal proxy event: %s", err)
	}
	req.Header.Set("X-ApiGatewayProxy-Event", string(hbody))
	if ctx != nil {
		req.Header.Set("X-ApiGatewayProxy-Context", string(ctx.ClientContext))
	}

	req.Host = proxyEvent.Headers["Host"]

	return req, nil
}
