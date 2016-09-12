package main

import (
	"fmt"
	"os"

	"github.com/th3noname/i3blocks/blocks/packages"
)

func main() {
	p := packages.New()

	p.Conf.UrgentValue = 30
	p.Conf.Pkg = packages.APT_HOOK
	p.Conf.PrintTemplate = " {{ .Packages }}"

	err := p.Exec()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	o, err := p.Print()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Print(o.String())
	if o.Urgent {
		os.Exit(33)
	}
}
