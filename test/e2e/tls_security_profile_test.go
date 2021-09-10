// Copyright 2021 The Cluster Monitoring Operator Authors
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
	"context"
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestTLSSecurityProfileConfiguration(t *testing.T) {
	//type testCase struct {
	//
	//}

	tlsSecurityProfile := configv1.TLSSecurityProfile{
		Type: configv1.TLSProfileOldType,
		Old:  &configv1.OldTLSProfile{},
	}
	setTlsSecurityProfile(t, &tlsSecurityProfile)


}

func setTlsSecurityProfile(t *testing.T, tlsSecurityProfile *configv1.TLSSecurityProfile) {
	ctx := context.Background()
	apiserverConfig, err := f.OpenShiftConfigClient.ConfigV1().APIServers().Get(ctx, "cluster", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	apiserverConfig.Spec.TLSSecurityProfile = tlsSecurityProfile
	if _, err := f.OpenShiftConfigClient.ConfigV1().APIServers().Update(ctx, apiserverConfig, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}
}
