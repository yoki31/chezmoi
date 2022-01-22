package progress

import (
	"io"
	"strings"
	"text/template"
)

type ReaderStatsTemplater struct {
	readerStatsFunc func() ReaderStats
	tmpl            *template.Template
}

func NewReaderStatsTemplater(readerStatsFunc ReaderStatsFunc, text string) (*ReaderStatsTemplater, error) {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return nil, err
	}
	return &ReaderStatsTemplater{
		readerStatsFunc: readerStatsFunc,
		tmpl:            tmpl,
	}, nil
}

func (f *ReaderStatsTemplater) Execute(w io.Writer) error {
	return f.tmpl.Execute(w, f.readerStatsFunc())
}

func (f ReaderStatsTemplater) String() string {
	sb := &strings.Builder{}
	if err := f.Execute(sb); err != nil {
		return ""
	}
	return sb.String()
}
