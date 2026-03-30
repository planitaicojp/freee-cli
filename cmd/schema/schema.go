package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	cerrors "github.com/planitaicojp/freee-cli/internal/errors"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// schemaNotFoundError is a command lookup error with exit code 3.
type schemaNotFoundError struct {
	name      string
	available []string
}

func (e *schemaNotFoundError) Error() string {
	return fmt.Sprintf("unknown command %q\nhint: available commands: %s", e.name, strings.Join(e.available, ", "))
}

func (e *schemaNotFoundError) ExitCode() int {
	return cerrors.ExitNotFound
}

// FlagSchema describes a single command flag.
type FlagSchema struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

// CommandSchema describes a command and its flags.
type CommandSchema struct {
	Command     string       `json:"command"`
	Description string       `json:"description"`
	Flags       []FlagSchema `json:"flags"`
}

// CommandListItem represents a command in a listing.
type CommandListItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CommandList wraps a slice of CommandListItem for JSON output.
type CommandList struct {
	Commands []CommandListItem `json:"commands"`
}

// NewCmd creates the schema command with the given root command for tree traversal.
func NewCmd(root *cobra.Command) *cobra.Command {
	var showGlobal bool

	cmd := &cobra.Command{
		Use:   "schema [resource] [action]",
		Short: "Show command flag schema",
		Long:  "Display flag schema for any command. Useful for AI agents and script authors.",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, _ := cmd.Flags().GetString("format")
			w := cmd.OutOrStdout()

			target, err := findCommand(root, args)
			if err != nil {
				return err
			}

			// If target has subcommands and we didn't specify an action,
			// list subcommands. Exception: leaf commands with 1 arg go to flags.
			if len(target.Commands()) > 0 && len(args) < 2 {
				items := listSubcommands(target)
				return outputCommandList(w, format, items)
			}

			// Output flag schema
			fullPath := buildCommandPath(root, args)
			flags := extractFlags(target, showGlobal, root)
			schema := CommandSchema{
				Command:     fullPath,
				Description: target.Short,
				Flags:       flags,
			}
			return outputFlagSchema(w, format, schema, showGlobal)
		},
	}

	cmd.Flags().BoolVar(&showGlobal, "show-global", false, "include global flags in output")

	return cmd
}

// findCommand traverses the cobra command tree to find the target command.
func findCommand(root *cobra.Command, args []string) (*cobra.Command, error) {
	cmd := root
	for _, name := range args {
		var found *cobra.Command
		for _, child := range cmd.Commands() {
			if child.Name() == name {
				found = child
				break
			}
		}
		if found == nil {
			var available []string
			for _, child := range cmd.Commands() {
				if !child.Hidden && child.Name() != "help" && child.Name() != "schema" {
					available = append(available, child.Name())
				}
			}
			return nil, &schemaNotFoundError{name: name, available: available}
		}
		cmd = found
	}
	return cmd, nil
}

// listSubcommands returns visible child commands, excluding hidden, help, and schema itself.
func listSubcommands(cmd *cobra.Command) []CommandListItem {
	var items []CommandListItem
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "schema" {
			continue
		}
		items = append(items, CommandListItem{
			Name:        child.Name(),
			Description: child.Short,
		})
	}
	return items
}

// extractFlags collects flag metadata from a command.
// When showGlobal is false, flags owned by root's PersistentFlags are excluded.
func extractFlags(cmd *cobra.Command, showGlobal bool, root *cobra.Command) []FlagSchema {
	rootFlags := make(map[string]bool)
	if !showGlobal {
		root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			rootFlags[f.Name] = true
		})
	}

	var flags []FlagSchema
	seen := make(map[string]bool)

	addFlag := func(f *pflag.Flag) {
		if seen[f.Name] || rootFlags[f.Name] {
			return
		}
		seen[f.Name] = true

		required := false
		if ann := f.Annotations; ann != nil {
			if _, ok := ann[cobra.BashCompOneRequiredFlag]; ok {
				required = true
			}
		}

		defVal := f.DefValue

		flags = append(flags, FlagSchema{
			Name:        f.Name,
			Type:        f.Value.Type(),
			Required:    required,
			Default:     defVal,
			Description: f.Usage,
		})
	}

	cmd.LocalFlags().VisitAll(addFlag)
	cmd.InheritedFlags().VisitAll(addFlag)

	return flags
}

// buildCommandPath constructs the full command path string.
func buildCommandPath(root *cobra.Command, args []string) string {
	parts := []string{root.Name()}
	parts = append(parts, args...)
	return strings.Join(parts, " ")
}

// outputCommandList renders a command listing.
func outputCommandList(w io.Writer, format string, items []CommandListItem) error {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(CommandList{Commands: items})
	case "yaml", "csv":
		return output.New(format).Format(w, items)
	default:
		tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "NAME\tDESCRIPTION")
		for _, item := range items {
			fmt.Fprintf(tw, "%s\t%s\n", item.Name, item.Description)
		}
		return tw.Flush()
	}
}

// FlagSchemaRow is a display-friendly row for table/CSV/YAML output.
type FlagSchemaRow struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    string `json:"required"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

// outputFlagSchema renders a flag schema.
func outputFlagSchema(w io.Writer, format string, schema CommandSchema, showGlobal bool) error {
	if format == "json" {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(schema)
	}

	if format == "yaml" || format == "csv" {
		rows := make([]FlagSchemaRow, len(schema.Flags))
		for i, f := range schema.Flags {
			req := "no"
			if f.Required {
				req = "yes"
			}
			rows[i] = FlagSchemaRow{
				Name:        f.Name,
				Type:        f.Type,
				Required:    req,
				Default:     f.Default,
				Description: f.Description,
			}
		}
		return output.New(format).Format(w, rows)
	}

	// Table format
	fmt.Fprintf(w, "Command: %s\n", schema.Command)
	fmt.Fprintf(w, "Description: %s\n\n", schema.Description)

	if len(schema.Flags) == 0 {
		fmt.Fprintln(w, "No flags.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tTYPE\tREQUIRED\tDEFAULT\tDESCRIPTION")
	for _, f := range schema.Flags {
		req := "no"
		if f.Required {
			req = "yes"
		}
		def := f.Default
		if def == "" || (f.Type == "bool" && def == "false") {
			def = "-"
		}
		fmt.Fprintf(tw, "--%s\t%s\t%s\t%s\t%s\n", f.Name, f.Type, req, def, f.Description)
	}
	tw.Flush()

	if !showGlobal {
		fmt.Fprintf(w, "\nGLOBAL FLAGS (--format, --profile, etc.) are omitted. Use --show-global to include.\n")
	}
	return nil
}
