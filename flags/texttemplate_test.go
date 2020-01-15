package flags

import (
	"bytes"
	"testing"
)

func TestTextTemplate(t *testing.T) {
	data := struct{ Name string }{"Francesc"}
	tests := []struct {
		text    string
		out     string
		invalid bool
	}{
		{text: "{{.Name}}", out: "Francesc"},
		{text: "", out: ""},
		{text: "{{", invalid: true},
		{text: "Hello, {{.Name}}", out: "Hello, Francesc"},
	}

	for _, test := range tests {
		var tt textTemplate
		if err := tt.Set(test.text); err != nil {
			if !test.invalid {
				t.Errorf("parsing %s failed unexpectedly: %v", test.text, err)
			}
			continue
		}
		if test.invalid {
			t.Errorf("parsing %s should have failed", test.text)
			continue
		}

		if tt.text != test.text {
			t.Errorf("the text template should be %s; got %s", test.text, tt.text)
		}

		var buf bytes.Buffer
		if err := tt.t.Execute(&buf, data); err != nil {
			t.Errorf("unexpected error executing %s: %v", test.text, err)
			continue
		}
		if got := buf.String(); got != test.out {
			t.Errorf("expected text was %s; got %s", test.out, got)
			continue
		}
	}
}
