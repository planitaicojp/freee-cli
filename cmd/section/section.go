package section

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	"github.com/planitaicojp/freee-cli/internal/api"
	"github.com/planitaicojp/freee-cli/internal/model"
	"github.com/planitaicojp/freee-cli/internal/output"
)

var Cmd = &cobra.Command{
	Use:   "section",
	Short: "Manage sections (部門)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	createCmd.Flags().String("name", "", "section name (required)")
	_ = createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("shortcut1", "", "shortcut 1")
	createCmd.Flags().String("shortcut2", "", "shortcut 2")

	updateCmd.Flags().String("name", "", "section name")
	updateCmd.Flags().String("shortcut1", "", "shortcut 1")
	updateCmd.Flags().String("shortcut2", "", "shortcut 2")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List sections",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		freeeAPI := &api.FreeeAPI{Client: client}

		format := cmdutil.GetFormat(cmd)
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		if format != "" && format != "table" {
			var resp any
			if err := freeeAPI.ListSections(client.CompanyID, &resp); err != nil {
				return err
			}
			return output.New(format, opts).Format(os.Stdout, resp)
		}

		var resp model.SectionsResponse
		if err := freeeAPI.ListSections(client.CompanyID, &resp); err != nil {
			return err
		}
		rows := make([]model.SectionRow, len(resp.Sections))
		for i, s := range resp.Sections {
			rows[i] = model.SectionRow{ID: s.ID, Name: s.Name}
		}
		return output.New("table", opts).Format(os.Stdout, rows)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a section",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		body := map[string]any{
			"company_id": client.CompanyID,
			"name":       name,
		}
		if v, _ := cmd.Flags().GetString("shortcut1"); v != "" {
			body["shortcut1"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut2"); v != "" {
			body["shortcut2"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/sections")
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.CreateSection(body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a section",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid section ID: %s", args[0])
		}

		body := map[string]any{
			"company_id": client.CompanyID,
		}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			body["name"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut1"); v != "" {
			body["shortcut1"] = v
		}
		if v, _ := cmd.Flags().GetString("shortcut2"); v != "" {
			body["shortcut2"] = v
		}

		if cmdutil.IsDryRun(cmd) {
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/sections/%d\n", id)
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.UpdateSection(id, body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a section",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid section ID: %s", args[0])
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteSection(client.CompanyID, id)
	},
}
