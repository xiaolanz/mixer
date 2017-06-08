// Copyright 2017 Istio Authors
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

package serviceControl

import (
	"net/http"
	"os"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	servicecontrol "google.golang.org/api/servicecontrol/v1"
	"istio.io/mixer/pkg/adapter"
	"fmt"
	"encoding/json"
)

func createAPIClient(logger adapter.Logger, clientSecretFile string) (*servicecontrol.Service, error) {
	logger.Infof("Creating service control client...\n")
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Transport: http.DefaultTransport})

	bytes, err := ioutil.ReadFile(clientSecretFile)
	if err != nil {
		return nil, err
	}

	o, err := google.ConfigFromJSON(bytes, servicecontrol.CloudPlatformScope, servicecontrol.ServicecontrolScope)
	logger.Infof("Created oauth config %v\n", o)

	// TODO need authorize for the first time.
//	ts, err := authorize(ctx, *o)
//	if err != nil {
//		return nil, err
//	}

//	t, err := ts.Token()

	t, err := o.Exchange(ctx, "4/-wzh6lqur-rLcj_hI1B1hMWTYGETVANetyYB0R43U6U")
	if err != nil {
		return nil, err
	}

	httpClient := o.Client(ctx, t)
	s, err := servicecontrol.New(httpClient)
	logger.Infof("Created service control client")
	return s, err
}

func authorize(ctx context.Context, config oauth2.Config) (oauth2.TokenSource, error) {
	authURL := config.AuthCodeURL("")

	showURL(authURL)

	code, err := obtainCode()

	if err != nil {
		return nil, err
	}

	token, err := config.Exchange(ctx, code)

	if err != nil {
		return nil, err
	}
  showToken(token)
	return config.TokenSource(ctx, token), nil
}

func showURL(url string) {
	fmt.Printf("Please go to >> %s << in order to authorize the prodreview tool.\n\nCopy the code displayed there, and then paste the code here: ", url)
	os.Stdout.Sync()
}

func obtainCode() (string, error) {
	var code string
	_, err := fmt.Scanln(&code)
	return code, err
}

func showToken(token *oauth2.Token) error {
	jt, err := json.Marshal(token)
	if err != nil {
		return err
	}
	fmt.Println("Save this token in a secure place and give it to prodreview_tool via --credentials the next time.\n\n")
	fmt.Println(string(jt))
	return nil
}
