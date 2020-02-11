package formatter

import "fmt"

type SimpleFormatter struct {
	Key   string
	Value interface{}
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
