package packages

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

// PackageManager defines the package manager used by the system
type PackageManager int

const (
	// APTITUDE package manager (Debian, Ubuntu ...)
	APTITUDE PackageManager = iota
	// PACMAN (Arch)
	PACMAN
)

// Packages gives access to the number of available system updates
type Packages struct {
	PrintTemplate string
	Data          data
}

// Data contains the collected information
type data struct {
	Packages string
}

// New returns a instance of type Packages
func New() Packages {
	return Packages{PrintTemplate: "{{ .Packages }}"}
}

// Exec collects the information
func (p *Packages) Exec(pkg PackageManager) error {
	switch pkg {
	case APTITUDE:
		out, err := exec.Command("/usr/lib/update-notifier/apt-check").CombinedOutput()
		if err != nil {
			return errors.Wrap(err, "executing command failed")
		}

		if len(out) == 0 {
			return errors.New("CMD returned empty")
		}

		p.Data.Packages = strings.Split(string(out), ";")[0]
	default:
		return errors.New("Invalid package manager")
	}
	return nil
}

// Print outputs a formatted string using PrintTemplate
func (p *Packages) Print() (string, error) {
	t := template.New("p")

	t, err := t.Parse(p.PrintTemplate)
	if err != nil {
		return "", errors.Wrap(err, "parsing template failed")
	}

	var out bytes.Buffer
	err = t.Execute(&out, p.Data)
	return out.String(), errors.Wrap(err, "executing template failed")
}
