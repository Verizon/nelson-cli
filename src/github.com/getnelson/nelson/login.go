//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//:
//:   Licensed under the Apache License, Version 2.0 (the "License");
//:   you may not use this file except in compliance with the License.
//:   You may obtain a copy of the License at
//:
//:       http://www.apache.org/licenses/LICENSE-2.0
//:
//:   Unless required by applicable law or agreed to in writing, software
//:   distributed under the License is distributed on an "AS IS" BASIS,
//:   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//:   See the License for the specific language governing permissions and
//:   limitations under the License.
//:
//: ----------------------------------------------------------------------------
package main

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

type CreateSessionRequest struct {
	AccessToken string `json:"access_token"`
}

// { "session_token": "xxx", "expires_at": 12345 }
type Session struct {
	SessionToken string `json:"session_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

///////////////////////////// CLI ENTRYPOINT ////////////////////////////////

func Login(client *gorequest.SuperAgent, githubToken string, nelsonHost string, disableTLS bool) []error {
	baseURL := createEndpointURL(nelsonHost, !disableTLS)
	sess, errs := createSession(client, githubToken, baseURL)
	if errs != nil {
		return errs
	}
	writeConfigFile(sess, baseURL, defaultConfigPath()) // TIM: side-effect, discarding errors seems wrong
	return nil
}

///////////////////////////// INTERNALS ////////////////////////////////

func createEndpointURL(host string, useTLS bool) string {
	u := "://" + host
	if useTLS {
		return "https" + u
	} else {
		return "http" + u
	}
}

/* TODO: any error handling here... would be nice */
func createSession(client *gorequest.SuperAgent, githubToken string, baseURL string) (Session, []error) {
	ver := CreateSessionRequest{AccessToken: githubToken}
	url := baseURL + "/auth/github"
	_, bytes, errs := client.
		Post(url).
		Set("User-Agent", UserAgentString(globalBuildVersion)).
		Send(ver).
		SetCurlCommand(globalEnableCurl).
		SetDebug(globalEnableDebug).
		Timeout(GetTimeout(globalTimeoutSeconds)).
		EndBytes()

	if len(errs) > 0 {
		return Session{}, errs
	}

	var result Session
	if err := json.Unmarshal(bytes, &result); err != nil {
		return Session{}, []error{err}
	}

	return result, nil
}
