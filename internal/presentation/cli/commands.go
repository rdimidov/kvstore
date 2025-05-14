package cli

import (
	"errors"

	"github.com/rdimidov/kvstore/internal/domain"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kvstore",
		Short: "An in-memory key-value store CLI",
		Long: `kvstore is a command-line application that implements a simple in-memory key-value database.
		
		Supports basic operations:
		- Setting a value by key (set)
		- Retrieving a value by key (get)
		- Deleting a key (delete)
		
		The database is stored entirely in memory and does not persist data to disk.`,
	}
}

func NewGetCmd(getter getter) *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get value by key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key, err := domain.NewKey(args[0])
			if err != nil {
				cmd.PrintErr("invalid key format: " + err.Error() + "\n")
				return
			}
			entry, err := getter.Get(cmd.Context(), key)
			if err != nil {
				if errors.Is(err, domain.ErrKeyNotFound) {
					cmd.PrintErr("key not found\n")
					return
				}
				cmd.PrintErr("error: " + err.Error() + "\n")
				return
			}
			cmd.Print(entry.Value + "\n")
		},
	}
}

func NewSetCmd(setter setter) *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set value by key",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key, err := domain.NewKey(args[0])
			if err != nil {
				cmd.PrintErr("invalid key format: " + err.Error() + "\n")
				return
			}
			value, err := domain.NewValue(args[1])
			if err != nil {
				cmd.PrintErr("invalid value format: " + err.Error() + "\n")
				return
			}
			if err := setter.Set(cmd.Context(), key, value); err != nil {
				cmd.PrintErr("could not save entry: ", err.Error()+"\n")
				return
			}
		},
	}
}

func NewDelCmd(deleter deleter) *cobra.Command {
	return &cobra.Command{
		Use:   "del [key]",
		Short: "Delete value by key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key, err := domain.NewKey(args[0])
			if err != nil {
				cmd.PrintErr("invalid key format: " + err.Error() + "\n")
				return
			}
			if err := deleter.Delete(cmd.Context(), key); err != nil {
				cmd.PrintErr("could not delete entry: ", err.Error()+"\n")
				return
			}
		},
	}
}
