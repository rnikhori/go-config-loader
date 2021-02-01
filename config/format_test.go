package config

import "testing"

func TestDetectFormat(t *testing.T) {
	cases := []struct {
		path    string
		want    Format
		wantErr bool
	}{
		{"config.yaml", FormatYAML, false},
		{"config.yml", FormatYAML, false},
		{"config.YAML", FormatYAML, false},
		{"config.json", FormatJSON, false},
		{"config.txt", "", true},
		{"config", "", true},
	}

	for _, tc := range cases {
		got, err := detectFormat(tc.path)
		if tc.wantErr {
			if err == nil {
				t.Errorf("detectFormat(%q): expected error, got nil", tc.path)
			}
			continue
		}
		if err != nil {
			t.Errorf("detectFormat(%q): unexpected error: %v", tc.path, err)
			continue
		}
		if got != tc.want {
			t.Errorf("detectFormat(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}
