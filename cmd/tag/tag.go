package tag

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
	Use:   "tag",
	Short: "Manage tags (メモタグ)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)

	createCmd.Flags().String("name", "", "tag name (required)")
	_ = createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("shortcut1", "", "shortcut 1")
	createCmd.Flags().String("shortcut2", "", "shortcut 2")

	updateCmd.Flags().String("name", "", "tag name")
	updateCmd.Flags().String("shortcut1", "", "shortcut 1")
	updateCmd.Flags().String("shortcut2", "", "shortcut 2")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags",
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
			if err := freeeAPI.ListTags(client.CompanyID, &resp); err != nil {
				return err
			}
			return output.New(format, opts).Format(os.Stdout, resp)
		}

		var resp model.TagsResponse
		if err := freeeAPI.ListTags(client.CompanyID, &resp); err != nil {
			return err
		}
		rows := make([]model.TagRow, len(resp.Tags))
		for i, t := range resp.Tags {
			rows[i] = model.TagRow{ID: t.ID, Name: t.Name}
		}
		return output.New("table", opts).Format(os.Stdout, rows)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a tag",
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
			fmt.Fprintln(os.Stderr, "[dry-run] POST /api/1/tags")
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.CreateTag(body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a tag",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid tag ID: %s", args[0])
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
			fmt.Fprintf(os.Stderr, "[dry-run] PUT /api/1/tags/%d\n", id)
			return output.New("json").Format(os.Stdout, body)
		}

		freeeAPI := &api.FreeeAPI{Client: client}
		var resp any
		if err := freeeAPI.UpdateTag(id, body, &resp); err != nil {
			return err
		}
		opts := output.Options{NoHeader: cmdutil.IsNoHeader(cmd)}
		return output.New(cmdutil.GetFormat(cmd), opts).Format(os.Stdout, resp)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a tag",
	Args:  cmdutil.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid tag ID: %s", args[0])
		}
		freeeAPI := &api.FreeeAPI{Client: client}
		return freeeAPI.DeleteTag(client.CompanyID, id)
	},
}
