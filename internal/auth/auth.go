// Package auth describes how LazyTF helps a user log in to their cloud: a
// login command to run, and the Terraform output that means credentials are
// missing or expired. AWS SSO, gcloud and az are presets; a custom command
// covers anything else. With no provider configured LazyTF does nothing about
// credentials.
package auth

import (
	"fmt"
	"strings"
)

// Provider names accepted in the config.
const (
	AWSSSO = "aws-sso"
	GCloud = "gcloud"
	Az     = "az"
	Custom = "custom"
)

// Config is the `auth` section of config.yml or .lazytf.yml.
type Config struct {
	// Provider is one of aws-sso, gcloud, az, custom, or empty for none.
	Provider string `yaml:"provider,omitempty"`
	// Login is the command and arguments to run for a custom provider. It
	// overrides the preset's command for the others.
	Login []string `yaml:"login,omitempty"`
	// ErrorPatterns are extra substrings of Terraform output (matched without
	// regard to case) that mean the user must log in.
	ErrorPatterns []string `yaml:"error_patterns,omitempty"`
}

// Provider is a resolved auth provider.
type Provider struct {
	Name string
	// Label names the cloud in prompts: "AWS", "Google Cloud", "Azure".
	Label string
	// Login is the argv of the login command. Empty for aws-sso, which picks
	// an SSO session from ~/.aws/config first.
	Login         []string
	ErrorPatterns []string
}

var presets = map[string]Provider{
	AWSSSO: {
		Name:  AWSSSO,
		Label: "AWS",
		ErrorPatterns: []string{
			"No valid credential sources found",
			"Unable to locate credentials",
			"SSO session associated with this profile has expired",
			"ExpiredToken",
		},
	},
	GCloud: {
		Name:  GCloud,
		Label: "Google Cloud",
		Login: []string{"gcloud", "auth", "application-default", "login"},
		ErrorPatterns: []string{
			"could not find default credentials",
			"Reauthentication required",
			"invalid_grant",
			"oauth2: cannot fetch token",
		},
	},
	Az: {
		Name:  Az,
		Label: "Azure",
		Login: []string{"az", "login"},
		ErrorPatterns: []string{
			"Please run 'az login'",
			"AADSTS70043",
			"AADSTS700082",
			"Interactive authentication is needed",
		},
	},
}

// Resolve turns a config into a provider. It returns nil, nil when no
// provider is configured, and an error for an unknown provider or a custom
// provider without a login command.
func Resolve(c Config) (*Provider, error) {
	name := strings.TrimSpace(c.Provider)
	if name == "" {
		return nil, nil
	}
	var p Provider
	if name == Custom {
		if len(c.Login) == 0 {
			return nil, fmt.Errorf("auth provider %q needs a login command", Custom)
		}
		p = Provider{Name: Custom, Label: "Your"}
	} else {
		preset, ok := presets[name]
		if !ok {
			return nil, fmt.Errorf("unknown auth provider %q (want %s, %s, %s or %s)", name, AWSSSO, GCloud, Az, Custom)
		}
		p = preset
		p.ErrorPatterns = append([]string{}, preset.ErrorPatterns...)
	}
	if len(c.Login) > 0 {
		p.Login = append([]string{}, c.Login...)
	}
	p.ErrorPatterns = append(p.ErrorPatterns, c.ErrorPatterns...)
	return &p, nil
}

// NeedsLogin reports whether Terraform output looks like missing or expired
// credentials for this provider.
func (p *Provider) NeedsLogin(output string) bool {
	if p == nil {
		return false
	}
	lower := strings.ToLower(output)
	for _, pattern := range p.ErrorPatterns {
		if pattern != "" && strings.Contains(lower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// UsesSSOSessions reports whether login first picks an AWS SSO session.
func (p *Provider) UsesSSOSessions() bool {
	return p != nil && p.Name == AWSSSO && len(p.Login) == 0
}
