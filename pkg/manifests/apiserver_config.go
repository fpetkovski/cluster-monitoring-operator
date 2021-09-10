package manifests

import configv1 "github.com/openshift/api/config/v1"

type APIServerConfig configv1.APIServer

var DefaultTLSCiphers = configv1.TLSProfiles[configv1.TLSProfileIntermediateType].Ciphers

func (c *APIServerConfig) GetTLSCiphers() []string {
	if c == nil {
		return DefaultTLSCiphers
	}

	profile := c.Spec.TLSSecurityProfile

	if profile == nil {
		return DefaultTLSCiphers
	}

	if profile.Type != configv1.TLSProfileCustomType {
		if tlsConfig, ok := configv1.TLSProfiles[profile.Type]; ok {
			return tlsConfig.Ciphers
		}
		return DefaultTLSCiphers
	}

	if profile.Custom != nil {
		return profile.Custom.TLSProfileSpec.Ciphers
	}
	return DefaultTLSCiphers

}
