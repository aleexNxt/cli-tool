package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aleexNxt/cli-tool/internal/domain"
	"github.com/aleexNxt/cli-tool/internal/infrastructure/config"
)

func newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Konfigurationsverwaltung",
		Long:  "Verwaltet die DevTool-Konfiguration",
	}

	configCmd.AddCommand(newConfigShowCommand())
	configCmd.AddCommand(newConfigSetCommand())
	configCmd.AddCommand(newConfigGetCommand())
	configCmd.AddCommand(newConfigInitCommand())

	return configCmd
}

func newConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Zeigt die aktuelle Konfiguration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configRepo := config.NewConfigRepository()
			cfg, err := configRepo.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Pretty print JSON
			data, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
			return nil
		},
	}
}

func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Setzt einen Konfigurationswert",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			configRepo := config.NewConfigRepository()
			if err := configRepo.Set(key, value); err != nil {
				return fmt.Errorf("failed to set config: %w", err)
			}

			fmt.Printf("✓ Konfiguration gesetzt: %s = %s\n", key, value)
			return nil
		},
	}
}

func newConfigGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Holt einen Konfigurationswert",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			configRepo := config.NewConfigRepository()
			value, err := configRepo.Get(key)
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}

			if value == "" {
				fmt.Printf("%s ist nicht gesetzt\n", key)
			} else {
				fmt.Printf("%s = %s\n", key, value)
			}
			return nil
		},
	}
}

func newConfigInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialisiert eine neue Konfiguration",
		RunE: func(cmd *cobra.Command, args []string) error {
			configRepo := config.NewConfigRepository()

			// Check if config exists
			_, err := configRepo.Load()
			if err == nil {
				fmt.Println("Konfiguration existiert bereits. Überschreiben? (y/n)")
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Abgebrochen.")
					return nil
				}
			}

			// Create default config
			cfg := newDefaultConfig()
			if err := configRepo.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Println("✓ Konfiguration initialisiert")
			fmt.Println("\nStandardwerte:")
			data, _ := json.MarshalIndent(cfg, "", "  ")
			fmt.Println(string(data))

			return nil
		},
	}
}

func newDefaultConfig() *domain.Config {
	return &domain.Config{
		BuildCommand:   "go build -v ./...",
		TestCommand:    "go test -v ./...",
		LintCommand:    "golangci-lint run",
		DeployCommand:  "",
		DockerRegistry: "",
		Environment:    make(map[string]string),
		Timeout:        300,
	}
}
