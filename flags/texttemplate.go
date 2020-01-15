package flags

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

type textTemplate struct {
	text string
	t    template.Template
}

func (t *textTemplate) String() string {
	return fmt.Sprintf("%q", t.text)
}

func (t *textTemplate) Set(text string) error {
	p, err := template.New("flag").Parse(text)
	if err != nil {
		return err
	}
	t.t = *p
	t.text = text
	return nil
}

// TextTemplate defines a text template flag with specified name, default value, and usage string.
// The return value is the address of a text template variable that stores the value of the flag.
// If the given template fails to compile an error will be logged and the program will exit with
// exit code 2. (os.Exit(2))
func TextTemplate(name string, value string, usage string) *template.Template {
	t, err := template.New("").Parse(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid template in flag: %v", err)
		os.Exit(2)
	}
	tt := &textTemplate{value, *t}
	flag.Var(tt, name, usage)
	return &tt.t
}

// TextTemplateVar defines a text template flag with specified name, default value, and usage string.
// The argument t points to a text template variable in which to store the value of the flag.
// If the given template fails to compile an error will be logged and the program will exit with
// exit code 2. (os.Exit(2))
func TextTemplateVar(t *template.Template, name string, value string, usage string) {
	p, err := template.New("").Parse(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid template in flag: %v", err)
		os.Exit(2)
	}
	*t = *p
	flag.Var(&textTemplate{value, *t}, name, usage)
}
