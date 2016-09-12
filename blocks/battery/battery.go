package battery

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

// Battery provides information about battery status
type Battery struct {
	PrintTemplate string
	Data          data
}

// Data contains the collected information
type data struct {
	// either Full Charging Discharging
	Status string
	Power  string
	Time   []string
}

// New returns a instance of type Battery
func New() Battery {
	return Battery{PrintTemplate: "{{ .Status }} {{ .Power }} {{ .Time }}"}
}

// Exec collects the information
func (p *Battery) Exec() error {
	out, err := exec.Command("acpi", "-b").CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "executing command failed")
	}

	if len(out) == 0 ||
		bytes.Count(out, []byte("Battery 0: ")) == 0 {
		return errors.New("CMD returned empty or no Battery installed")
	}

	out = bytes.TrimPrefix(out, []byte("Battery 0: "))
	parts := bytes.Split(out, []byte(","))

	p.Data.Status = string(bytes.TrimSpace(parts[0]))
	p.Data.Power = string(bytes.TrimSuffix(bytes.TrimSpace(parts[1]), []byte("%")))
	if len(parts) >= 3 {
		p.Data.Time = strings.Split(string(bytes.TrimSpace(bytes.TrimSuffix(parts[2], []byte("remaining\n")))), ":")
	}

	return nil
}

// Print outputs a formatted string using PrintTemplate
func (p *Battery) Print() (string, error) {
	t := template.New("p")

	t, err := t.Parse(p.PrintTemplate)
	if err != nil {
		return "", errors.Wrap(err, "parsing template failed")
	}

	var out bytes.Buffer
	err = t.Execute(&out, p.Data)
	return out.String(), errors.Wrap(err, "executing template failed")
}
