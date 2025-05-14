package cli

import (
	"context"

	"github.com/spf13/cobra"
)

type CLI struct {
	root *cobra.Command
}

func New(app appController) *CLI {
	cobra.EnableCaseInsensitive = true
	cli := &CLI{root: NewRootCmd()}

	// register sub-commands
	cli.root.AddCommand(
		NewGetCmd(app),
		NewSetCmd(app),
		NewDelCmd(app),
	)
	return cli
}

func (c *CLI) Execute(ctx context.Context) error {
	return c.root.ExecuteContext(ctx)
}
