package manifests_test

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"reflect"
	"strings"
	"testing"
)

func TestGetTLSCiphers(t *testing.T) {
	defaultCiphers := configv1.TLSProfiles[configv1.TLSProfileIntermediateType].Ciphers

	testCases := []struct {
		name            string
		config          *manifests.APIServerConfig
		expectedCiphers []string
	}{
		{
			name: "nil config",
			config: nil,
			expectedCiphers: defaultCiphers,
		},
		{
			name: "nil profile",
			config: newApiserverConfig(nil),
			expectedCiphers: defaultCiphers,
		},
		{
			name: "empty profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: "",
			}),
			expectedCiphers: defaultCiphers,
		},
		{
			name: "invalid profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: "invalid-profile",
			}),
			expectedCiphers: defaultCiphers,
		},
		{
			name: "old profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileOldType,
			}),
			expectedCiphers: configv1.TLSProfiles[configv1.TLSProfileOldType].Ciphers,
		},
		{
			name: "intermediate profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileIntermediateType,
			}),
			expectedCiphers: defaultCiphers,
		},
		{
			name: "modern profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileModernType,
			}),
			expectedCiphers: configv1.TLSProfiles[configv1.TLSProfileModernType].Ciphers,
		},
		{
			name: "custom profile without ciphers",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
			}),
			expectedCiphers: defaultCiphers,
		},
		{
			name: "custom profile with ciphers",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						Ciphers: []string{"cipher-1", "cipher-2"},
					},
				},
			}),
			expectedCiphers: []string{"cipher-1", "cipher-2"},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.config.GetTLSCiphers()
			if !reflect.DeepEqual(tt.expectedCiphers, actual) {
				t.Fatalf("invalid ciphers, got %s, want %s", strings.Join(actual, ", "), strings.Join(tt.expectedCiphers, ", "))
			}
		})
	}
}

func newApiserverConfig(profile *configv1.TLSSecurityProfile) *manifests.APIServerConfig {
	config := manifests.APIServerConfig(configv1.APIServer{
		Spec: configv1.APIServerSpec{
			TLSSecurityProfile: profile,
		},
	})

	return &config
}
