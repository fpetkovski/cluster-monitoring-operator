// Copyright 2020 The Cluster Monitoring Operator Authors
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

package e2e

import (
	"bytes"
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"strings"

	"github.com/Jeffail/gabs"
)

func getThanosRules(body []byte, expGroupName, expRuleName string) error {
	j, err := gabs.ParseJSON([]byte(body))
	if err != nil {
		return err
	}

	groups, err := j.Path("data.groups").Children()
	if err != nil {
		return err
	}

	for i := 0; i < len(groups); i++ {
		groupName := groups[i].S("name").Data().(string)
		if groupName != expGroupName {
			continue
		}

		rules, err := groups[i].Path("rules").Children()
		if err != nil {
			return err
		}

		for j := 0; j < len(rules); j++ {
			ruleName := rules[j].S("name").Data().(string)
			if ruleName == expRuleName {
				return nil
			}
		}
	}
	return fmt.Errorf("'%s' alert not found in '%s' group", expRuleName, expGroupName)
}

// startPortForward initiates a port forwarding connection to a pod on the localhost interface.
//
// startPortForward blocks until the port forwarding proxy server is ready to receive connections.
func startPortForward(config *rest.Config, scheme string, name string, ns string, port string) error {
	roundTripper, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", ns, name)
	hostIP := strings.TrimLeft(config.Host, "htps:/")
	serverURL := url.URL{Scheme: scheme, Path: path, Host: hostIP}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	forwarder, err := portforward.New(dialer, []string{port}, stopChan, readyChan, out, errOut)
	if err != nil {
		return err
	}

	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			panic(err)
		}
	}()

	<-readyChan
	return nil
}