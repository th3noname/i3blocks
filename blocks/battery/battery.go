package battery

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/th3noname/i3blocks/blocks"
)

// Battery provides information about battery status
type Battery struct {
	Conf configuration
	Data data
}

type data struct {
	// either Full, Charging, Discharging
	Status string

	// in percent
	Power int

	// slice containing the remaining time (h:m:s)
	Time []string
}

type configuration struct {
	PrintTemplate string
	UrgentValue   int
	Color         string
}

// New returns a instance of type Battery
func New() Battery {
	return Battery{
		Conf: configuration{
			PrintTemplate: "{{ .Status }} {{ .Power }} {{ .Time }}",
		},
	}
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

	parts := bytes.Split(bytes.TrimPrefix(out, []byte("Battery 0: ")), []byte(","))

	p.Data.Status = string(bytes.TrimSpace(parts[0]))
	p.Data.Power, err = strconv.Atoi(string(bytes.TrimSuffix(bytes.TrimSpace(parts[1]), []byte("%"))))
	if err != nil {
		return errors.Wrap(err, "converting power value failed")
	}

	if len(parts) >= 3 {
		p.Data.Time = strings.Split(string(bytes.TrimSpace(bytes.TrimSuffix(parts[2], []byte("remaining\n")))), ":")
	}

	return nil
}

// Print outputs a formatted string using PrintTemplate
func (p *Battery) Print() (blocks.Output, error) {
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
			Urgent:    p.Data.Power <= p.Conf.UrgentValue,
			Color:     p.Conf.Color,
		},
		errors.Wrap(err, "executing template failed")
}
