package main

import (
	"fmt"
	"os"

	"github.com/th3noname/i3blocks/blocks/battery"
)

func main() {
	b := battery.New()

	err := b.Exec()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.PrintTemplate = `{{ if eq .Status "Full" -}}
		 {{ .Power }}%
		{{- else if eq .Status "Charging" -}}
		 {{ .Power }}% {{ index .Time 0 }}:{{ index .Time 1 }}
		{{- else if eq .Status "Discharging" -}}
		{{- if le .Power "75" -}}
		
		{{- else if le .Power "50" -}}
		
		{{- else if le .Power "25" -}}
		
		{{- else if le .Power "5" -}}
		
		{{- else -}}
		
		{{- end }} {{ .Power }}% {{ index .Time 0 }}:{{ index .Time 1 }}
		{{- end }}`

	s, err := b.Print()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(s)
	fmt.Println(s)
}
