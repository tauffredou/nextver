package formatter

import (
	"fmt"
	"os"
	"text/template"
)

type SimpleFormatter struct {
	Key   string
	Value interface{}
}

func (f *SimpleFormatter) Template(text string) error {
	tpl, err := template.New("Release").Parse(text)
	if err != nil {
		return err
	}
	return tpl.Execute(os.Stdout, f.Value)
}

func (f *SimpleFormatter) Json() {
	fmt.Printf("{\"%s\":\"%s\"}\n", f.Key, f.Value)
}

func (f *SimpleFormatter) Yaml() {
	fmt.Printf("%s: \"%s\"\n", f.Key, f.Value)
}

func (f *SimpleFormatter) Console() {
	fmt.Println(f.Value)
}
