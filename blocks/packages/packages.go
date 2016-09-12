package packages

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/th3noname/i3blocks/blocks"
)

// PackageManager defines the package manager used by the system
type PackageManager int

const (
	NULL PackageManager = iota
	// APT package manager (Debian, Ubuntu ...)
	// This should work on all Debian based distributions out of
	// the box.
	APT
	// APT_HOOK is faster than the APT version above, but depends
	// on a change in the apt configuration
	//
	// The line
	//     DPkg::Post-Invoke-Success {"/usr/lib/update-notifier/apt-check 2> /var/opt/updates";};
	// needs to be added either into /etc/apt/apt.conf or one of
	// the existing files or a new file in /etc/apt/apt.conf.d/.
	APT_HOOK
	// PACMAN (Arch)
	PACMAN
)

// Packages gives access to the number of available system updates
type Packages struct {
	Conf configuration
	Data data
}

type data struct {
	Packages string
}

type configuration struct {
	PrintTemplate string
	UrgentValue   string
	Color         string
	Pkg           PackageManager
}

// New returns a instance of type Packages
func New() Packages {
	return Packages{
		Conf: configuration{
			PrintTemplate: "{{ .Status }} {{ .Power }} {{ .Time }}",
		},
	}
}

// Exec collects the information
// Check the PackageManager constant to learn about supported package managers.
func (p *Packages) Exec() error {
	switch p.Conf.Pkg {
	case APT:
		out, err := exec.Command("/usr/lib/update-notifier/apt-check").CombinedOutput()
		if err != nil {
			return errors.Wrap(err, "executing command failed")
		}

		p.Data.Packages, err = parseAPT(out)
		return errors.Wrap(err, "parsing output failed")

	case APT_HOOK:
		out, err := ioutil.ReadFile("/var/opt/updates")
		if err != nil {
			return errors.Wrap(err, "reading updates file failed")
		}

		p.Data.Packages, err = parseAPT(out)
		return errors.Wrap(err, "parsing output failed")

	default:
		return errors.New("Invalid package manager")
	}
}

// Print outputs a formatted string using PrintTemplate
func (p *Packages) Print() (blocks.Output, error) {
	t := template.New("p")

	t, err := t.Parse(p.Conf.PrintTemplate)
	if err != nil {
		return blocks.Output{}, errors.Wrap(err, "parsing template failed")
	}

	var out bytes.Buffer
	err = t.Execute(&out, p.Data)

	return blocks.Output{
			ShortText: out.String(),
			FullText:  out.String(),
			Urgent:    p.Data.Packages >= p.Conf.UrgentValue,
			Color:     p.Conf.Color,
		},
		errors.Wrap(err, "executing template failed")
}

func parseAPT(b []byte) (string, error) {
	if len(b) == 0 {
		return "", errors.New("empty data body")
	}

	return strings.Split(string(b), ";")[0], nil
}
