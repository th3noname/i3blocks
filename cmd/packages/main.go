package main

import (
	"fmt"
	"os"

	"github.com/th3noname/i3blocks/blocks/packages"
)

func main() {
	p := packages.New()

	err := p.Exec(packages.APT)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p.PrintTemplate = " {{ .Packages }}"

	s, err := p.Print()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(s)
	fmt.Println(s)
}
