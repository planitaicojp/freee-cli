package output

import "io"

// Formatter formats and writes data to a writer.
type Formatter interface {
	Format(w io.Writer, data any) error
}

// Options configures formatter behavior.
type Options struct {
	NoHeader bool
}

// New creates a formatter for the given format name.
func New(format string, opts ...Options) Formatter {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	switch format {
	case "json":
		return &JSONFormatter{}
	case "yaml":
		return &YAMLFormatter{}
	case "csv":
		return &CSVFormatter{Options: o}
	default:
		return &TableFormatter{Options: o}
	}
}
