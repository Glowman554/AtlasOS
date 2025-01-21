package main

import (
	"os"

	"github.com/foxboron/sbctl"
	"github.com/foxboron/sbctl/config"
	"github.com/foxboron/sbctl/logging"
	"github.com/foxboron/sbctl/lsm"
	"github.com/spf13/cobra"
)

var removeFileCmd = &cobra.Command{
	Use: "remove-file",
	Aliases: []string{
		"rm-file",
		"rm",
	},
	Short: "Remove file from database",
	RunE: func(cmd *cobra.Command, args []string) error {
		state := cmd.Context().Value(stateDataKey{}).(*config.State)

		if state.Config.Landlock {
			if err := lsm.Restrict(); err != nil {
				return err
			}
		}

		if len(args) < 1 {
			logging.Println("Need to specify file")
			os.Exit(1)
		}
		files, err := sbctl.ReadFileDatabase(state.Fs, state.Config.FilesDb)
		if err != nil {
			return err
		}
		if _, ok := files[args[0]]; !ok {
			logging.Print("File %s doesn't exist in database!\n", args[0])
			os.Exit(1)
		}
		delete(files, args[0])
		if err := sbctl.WriteFileDatabase(state.Fs, state.Config.FilesDb, files); err != nil {
			return err
		}
		logging.Print("Removed %s from the database.\n", args[0])
		return nil
	},
}

func init() {
	CliCommands = append(CliCommands, cliCommand{
		Cmd: removeFileCmd,
	})
}
