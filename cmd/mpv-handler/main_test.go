package main

import (
	"reflect"
	"testing"
)

func TestParseExtraArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace only",
			input: "   \t  ",
			want:  nil,
		},
		{
			name:  "space separated flags",
			input: "--fs --profile=fast",
			want:  []string{"--fs", "--profile=fast"},
		},
		{
			name:  "quoted argument with spaces",
			input: `--title="hello world" --fs`,
			want:  []string{`--title=hello world`, "--fs"},
		},
		{
			name:  "quoted standalone value",
			input: `--script-opts "profile=fast mode=cinema"`,
			want:  []string{"--script-opts", "profile=fast mode=cinema"},
		},
		{
			name:  "escaped quote inside quoted text",
			input: `--term-playing-msg="title: ""demo"""`,
			want:  []string{`--term-playing-msg=title: "demo"`},
		},
		{
			name:  "backslashes before quote are preserved correctly",
			input: `--msg="C:\\temp\\file"`,
			want:  []string{`--msg=C:\\temp\\file`},
		},
		{
			name:  "tabs also split arguments",
			input: "--fs\t--no-border",
			want:  []string{"--fs", "--no-border"},
		},
		{
			name:    "unterminated quote returns error",
			input:   `"unterminated`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseExtraArgs(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseExtraArgs(%q) error = nil, want error", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseExtraArgs(%q) unexpected error: %v", tt.input, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseExtraArgs(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}
