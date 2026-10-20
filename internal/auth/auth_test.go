package auth

import (
	"slices"
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		want    string // provider name, "" for none
		wantErr bool
	}{
		{name: "no provider", cfg: Config{}, want: ""},
		{name: "blank provider", cfg: Config{Provider: "  "}, want: ""},
		{name: "aws preset", cfg: Config{Provider: AWSSSO}, want: AWSSSO},
		{name: "gcloud preset", cfg: Config{Provider: GCloud}, want: GCloud},
		{name: "az preset", cfg: Config{Provider: Az}, want: Az},
		{name: "custom with login", cfg: Config{Provider: Custom, Login: []string{"my-login", "--now"}}, want: Custom},
		{name: "custom without login", cfg: Config{Provider: Custom}, wantErr: true},
		{name: "unknown provider", cfg: Config{Provider: "okta"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Resolve(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
			}
			switch {
			case tt.want == "" && got != nil:
				t.Errorf("Resolve() = %+v, want none", got)
			case tt.want != "" && (got == nil || got.Name != tt.want):
				t.Errorf("Resolve() = %+v, want %s", got, tt.want)
			}
		})
	}
}

func TestResolveLoginCommands(t *testing.T) {
	g, _ := Resolve(Config{Provider: GCloud})
	if want := []string{"gcloud", "auth", "application-default", "login"}; !slices.Equal(g.Login, want) {
		t.Errorf("gcloud login = %v, want %v", g.Login, want)
	}
	over, _ := Resolve(Config{Provider: Az, Login: []string{"az", "login", "--use-device-code"}})
	if want := []string{"az", "login", "--use-device-code"}; !slices.Equal(over.Login, want) {
		t.Errorf("overridden login = %v, want %v", over.Login, want)
	}
	aws, _ := Resolve(Config{Provider: AWSSSO})
	if !aws.UsesSSOSessions() {
		t.Error("aws-sso without a login override should pick an SSO session")
	}
	awsCustom, _ := Resolve(Config{Provider: AWSSSO, Login: []string{"aws", "sso", "login", "--profile", "x"}})
	if awsCustom.UsesSSOSessions() {
		t.Error("aws-sso with a login override should run it as given")
	}
}

func TestPresetsAreNotShared(t *testing.T) {
	a, _ := Resolve(Config{Provider: GCloud, ErrorPatterns: []string{"extra"}})
	b, _ := Resolve(Config{Provider: GCloud})
	if slices.Contains(b.ErrorPatterns, "extra") {
		t.Error("extra patterns leaked into the preset")
	}
	a.ErrorPatterns[0] = "mutated"
	if c, _ := Resolve(Config{Provider: GCloud}); c.ErrorPatterns[0] == "mutated" {
		t.Error("mutating a provider changed the preset")
	}
}

func TestNeedsLogin(t *testing.T) {
	aws, _ := Resolve(Config{Provider: AWSSSO})
	gcp, _ := Resolve(Config{Provider: GCloud, ErrorPatterns: []string{"corp sso expired"}})
	tests := []struct {
		p      *Provider
		output string
		want   bool
	}{
		{aws, "Error: No valid credential sources found", true},
		{aws, "Error: Unable to locate credentials", true},
		{aws, "Error: ExpiredToken: The security token included in the request is expired", true},
		{aws, "Error: could not find default credentials", false},
		{gcp, "google: Could Not Find Default Credentials. See https://cloud.google.com/docs", true},
		{gcp, "Corp SSO expired, log in again", true},
		{gcp, "Error: No valid credential sources found", false},
		{gcp, "Plan: 1 to add", false},
		{nil, "No valid credential sources found", false},
	}
	for _, tt := range tests {
		if got := tt.p.NeedsLogin(tt.output); got != tt.want {
			t.Errorf("NeedsLogin(%q) = %v, want %v", tt.output, got, tt.want)
		}
	}
}
