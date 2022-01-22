package progress

import (
	"strings"
	"text/template"
)

var defaultSpinnerStrings = []string{"|", "/", "-", "\\"}

type Spinner struct {
	strings   []string
	tmpl      *template.Template
	valueFunc func() int
}

func NewSpinner(valueFunc func() int, text string) (*Spinner, error) {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return nil, err
	}
	return &Spinner{
		strings:   defaultSpinnerStrings,
		tmpl:      tmpl,
		valueFunc: valueFunc,
	}, nil
}

func (s *Spinner) String() string {
	sb := &strings.Builder{}
	_ = s.tmpl.Execute(sb, s.strings[s.valueFunc()%len(s.strings)])
	return sb.String()
}
