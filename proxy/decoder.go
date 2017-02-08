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
// Originally from https://github.com/eawsy/aws-lambda-go-event/blob/master/service/lambda/runtime/event/apigatewayproxyevt/decoder.go
// Changes (kothar - 2017-02):
//   Relocated to go-apex/proxy

package proxy

import "encoding/json"

type requestContextAlias RequestContext

type authorizer map[string]string

// UnmarshalJSON interprets the data as a dynamic map which may carry either a
// Amazon Cognito set of claims or a custom set of attributes. It then choose
// the good one at runtime and fill the authorizer with it.
func (a *authorizer) UnmarshalJSON(data []byte) error {
	var cognito struct {
		Claims *map[string]string
	}
	var custom map[string]string

	err := json.Unmarshal(data, &cognito)
	if err != nil {
		return err
	}

	if cognito.Claims != nil {
		*a = authorizer(*cognito.Claims)
		return nil
	}

	err = json.Unmarshal(data, &custom)
	if err != nil {
		return err
	}

	*a = authorizer(custom)
	return nil
}

type jsonRequestContext struct {
	*requestContextAlias
	Authorizer authorizer
}

// UnmarshalJSON interprets data as a RequestContext with a special authorizer.
// It then leverages type aliasing and struct embedding to fill RequestContext
// with an usual map[string]string.
func (rc *RequestContext) UnmarshalJSON(data []byte) error {
	var jrc jsonRequestContext
	if err := json.Unmarshal(data, &jrc); err != nil {
		return err
	}

	*rc = *(*RequestContext)(jrc.requestContextAlias)
	rc.Authorizer = jrc.Authorizer

	return nil
}

// MarshalJSON reverts the effect of type aliasing and struct embedding used
// during the marshalling step to make the pattern seamless.
func (rc *RequestContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonRequestContext{
		(*requestContextAlias)(rc),
		rc.Authorizer,
	})
}
