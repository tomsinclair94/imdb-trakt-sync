package configure

import (
	"errors"
	"fmt"

	"github.com/cecobask/imdb-trakt-sync/cmd"
	"github.com/cecobask/imdb-trakt-sync/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var conf *config.Config
	var confPath string
	command := &cobra.Command{
		Use:   fmt.Sprintf("%s [command]", cmd.CommandNameConfigure),
		Short: "Configure provider credentials and sync options",
		PreRunE: func(c *cobra.Command, args []string) (err error) {
			confPath, err = c.Flags().GetString(cmd.FlagNameConfigFile)
			if err != nil {
				return err
			}
			if conf, err = config.New(confPath, false); err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}
			return conf.Validate()
		},
		RunE: func(c *cobra.Command, args []string) error {
			teaModel, err := config.NewTUI(conf)
			if err != nil {
				return fmt.Errorf("error initializing text-based user interface for the %s command: %w", cmd.CommandNameConfigure, err)
			}
			model, ok := teaModel.(*config.Model)
			if !ok {
				return fmt.Errorf("error type asserting tea.Model to *config.Model")
			}
			if err = model.Err(); err != nil {
				if errors.Is(err, config.ErrUserAborted) {
					return nil
				}
				return fmt.Errorf("error occurred in the config model: %w", err)
			}
			if conf, err = config.NewFromMap(model.Config()); err != nil {
				return fmt.Errorf("error loading config from map: %w", err)
			}
			if err = conf.Validate(); err != nil {
				return fmt.Errorf("error validating config: %w", err)
			}
			return conf.WriteFile(confPath)
		},
	}
	command.Flags().String(cmd.FlagNameConfigFile, cmd.ConfigFileDefault, "path to the config file")
	return command
}
