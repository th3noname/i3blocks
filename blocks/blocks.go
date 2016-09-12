package blocks

import "fmt"

// Output is the object returned by the blocks
type Output struct {
	FullText            string
	ShortText           string
	Color               string
	Background          string
	Border              string
	MinWidth            string
	Align               string
	Name                string
	Instance            string
	Urgent              bool
	Separator           string
	SeparatorBlockWidth string
	Markup              string
}

func (o *Output) String() string {
	return fmt.Sprintf("%s\n%s\n%s", o.FullText, o.ShortText, o.Color)
}
