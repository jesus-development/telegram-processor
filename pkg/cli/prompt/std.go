package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// StdPrompter doesn't handle backspace
type StdPrompter struct {
	r io.Reader
	w io.Writer
}

func NewStdPrompter(in io.Reader, out io.Writer) *StdPrompter {
	return &StdPrompter{
		r: in,
		w: out,
	}
}

func (p *StdPrompter) StringPrompt(label string) (string, error) {
	r := bufio.NewReader(os.Stdin)

	fmt.Fprint(p.w, label+" ")
	s, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("r.ReadString -> %w", err)
	}

	return strings.TrimSpace(s), nil
}

func (p *StdPrompter) YesNoPrompt(label string, def bool) (bool, error) {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(p.r)

	for {
		fmt.Fprintf(p.w, "%s (%s) ", label, choices)
		s, err := r.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("r.ReadString -> %w", err)
		}

		s = strings.TrimSpace(s)
		if s == "" {
			return def, nil
		}

		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true, nil
		}
		if s == "n" || s == "no" {
			return false, nil
		}
	}
}
