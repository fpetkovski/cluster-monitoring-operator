package manifests_test

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cluster-monitoring-operator/pkg/manifests"
	"reflect"
	"strings"
	"testing"
)

func TestGetTLSCiphers(t *testing.T) {
	defaultCiphers := manifests.APIServerDefaultTLSCiphers
	defaultTLSVersion := manifests.APIServerDefaultMinTLSVersion

	testCases := []struct {
		name                  string
		config                *manifests.APIServerConfig
		expectedCiphers       []string
		expectedMinTLSVersion configv1.TLSProtocolVersion
	}{
		{
			name:                  "nil config",
			config:                nil,
			expectedCiphers:       defaultCiphers,
			expectedMinTLSVersion: defaultTLSVersion,
		},
		{
			name:                  "nil profile",
			config:                newApiserverConfig(nil),
			expectedCiphers:       defaultCiphers,
			expectedMinTLSVersion: defaultTLSVersion,
		},
		{
			name: "empty profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: "",
			}),
			expectedCiphers:       defaultCiphers,
			expectedMinTLSVersion: defaultTLSVersion,
		},
		{
			name: "invalid profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: "invalid-profile",
			}),
			expectedCiphers:       defaultCiphers,
			expectedMinTLSVersion: defaultTLSVersion,
		},
		{
			name: "old profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileOldType,
			}),
			expectedCiphers:       configv1.TLSProfiles[configv1.TLSProfileOldType].Ciphers,
			expectedMinTLSVersion: configv1.TLSProfiles[configv1.TLSProfileOldType].MinTLSVersion,
		},
		{
			name: "intermediate profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileIntermediateType,
			}),
			expectedCiphers:       defaultCiphers,
			expectedMinTLSVersion: defaultTLSVersion,
		},
		{
			name: "modern profile",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileModernType,
			}),
			expectedCiphers:       configv1.TLSProfiles[configv1.TLSProfileModernType].Ciphers,
			expectedMinTLSVersion: configv1.TLSProfiles[configv1.TLSProfileModernType].MinTLSVersion,
		},
		{
			name: "custom profile without TLS configuration",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
			}),
			expectedCiphers:       defaultCiphers,
			expectedMinTLSVersion: defaultTLSVersion,
		},
		{
			name: "custom profile without configuration",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type:   configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{},
			}),
			expectedCiphers:       defaultCiphers,
			expectedMinTLSVersion: defaultTLSVersion,
		},
		{
			name: "custom profile with configuration",
			config: newApiserverConfig(&configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						Ciphers:       []string{"cipher-1", "cipher-2"},
						MinTLSVersion: configv1.VersionTLS11,
					},
				},
			}),
			expectedCiphers:       []string{"cipher-1", "cipher-2"},
			expectedMinTLSVersion: configv1.VersionTLS11,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actualCiphers := tt.config.GetTLSCiphers()
			if !reflect.DeepEqual(tt.expectedCiphers, actualCiphers) {
				t.Fatalf("invalid ciphers, got %s, want %s", strings.Join(actualCiphers, ", "), strings.Join(tt.expectedCiphers, ", "))
			}

			actualTLSVersion := tt.config.GetMinTLSVersion()
			if tt.expectedMinTLSVersion != actualTLSVersion {
				t.Fatalf("invalid min TLS version, got %s, want %s", actualTLSVersion, tt.expectedMinTLSVersion)
			}
		})
	}
}

func newApiserverConfig(profile *configv1.TLSSecurityProfile) *manifests.APIServerConfig {
	config := manifests.NewAPIServerConfig(&configv1.APIServer{
		Spec: configv1.APIServerSpec{
			TLSSecurityProfile: profile,
		},
	})

	return config
}
