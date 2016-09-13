package battery

import (
	"bytes"
	"fmt"
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
	BatteryID     int
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
func (b *Battery) Exec() error {
	out, err := exec.Command("acpi", "-b").CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "command execution failed")
	}

	if len(out) == 0 {
		return errors.New("command returned empty")
	}

	return b.parseACPI(out)
}

// Print outputs a formatted string using PrintTemplate
func (b *Battery) Print() (blocks.Output, error) {
	t := template.New("b")

	t, err := t.Parse(b.Conf.PrintTemplate)
	if err != nil {
		return blocks.Output{}, errors.Wrap(err, "parsing template failed")
	}

	var out bytes.Buffer
	err = t.Execute(&out, b.Data)

	return blocks.Output{
			ShortText: out.String(),
			FullText:  out.String(),
			Urgent:    b.Data.Power <= b.Conf.UrgentValue,
			Color:     b.Conf.Color,
		},
		errors.Wrap(err, "executing template failed")
}

func (b *Battery) parseACPI(d []byte) error {
	batteryPrefix := []byte(fmt.Sprintf("Battery %d: ", b.Conf.BatteryID))
	
	batteries := bytes.Split(d, []byte("\n"))
	
	for _, battery := range batteries {
		i := bytes.Index(battery, batteryPrefix)
		if i == -1 {
			continue 
		}
		
		parts := bytes.Split(battery[i + len(batteryPrefix):], []byte(","))
		
		for c, _ := range parts {
			parts[c] = bytes.TrimSpace(parts[c])
		}
		
		b.Data.Status = string(parts[0])
		
		var err error
		b.Data.Power, err = strconv.Atoi(
			string(parts[1][:len(parts[1]) - 1]),
		)
		if err != nil {
			return errors.Wrap(err, "converting power value failed")
		}
		
		if len(parts) >= 3 {
			b.Data.Time = strings.Split(string(parts[2][:8]), ":")
		}
		
		return nil
	}
	
	return errors.New("specified battery id does not exist")
}
