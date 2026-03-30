package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

type TableFormatter struct {
	Options Options
}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		_, err := fmt.Fprintf(w, "%v\n", data)
		return err
	}

	if val.Len() == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()

	headers := make([]string, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Name
		}
		headers[i] = strings.ToUpper(name)
	}
	if !f.Options.NoHeader {
		if _, err := fmt.Fprintln(tw, strings.Join(headers, "\t")); err != nil {
			return err
		}
	}

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}
		fields := make([]string, row.NumField())
		for j := 0; j < row.NumField(); j++ {
			jsonTag := elemType.Field(j).Tag.Get("json")
			fields[j] = formatValue(row.Field(j), jsonTag)
		}
		if _, err := fmt.Fprintln(tw, strings.Join(fields, "\t")); err != nil {
			return err
		}
	}

	return tw.Flush()
}

// formatValue formats a reflect value for table display.
func formatValue(v reflect.Value, jsonTag string) string {
	tagName := strings.Split(jsonTag, ",")[0]
	isIDOrCode := tagName == "id" || strings.HasSuffix(tagName, "_id") || strings.HasSuffix(tagName, "_code")

	if (v.Kind() == reflect.Int || v.Kind() == reflect.Int64) && !isIDOrCode {
		return formatAmount(v.Int())
	}
	if v.Kind() == reflect.String && tagName == "status" {
		return StatusLabel(v.String())
	}
	return fmt.Sprintf("%v", v.Interface())
}

// formatAmount formats an integer with comma separators.
func formatAmount(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	if neg {
		return "-" + string(result)
	}
	return string(result)
}
