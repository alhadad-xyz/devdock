package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/compose"
	"devdock/internal/config"
	"devdock/internal/detector"
	"devdock/internal/env"
	"devdock/internal/utils"
	"devdock/recipes"
)

var (
	initType  string
	initDB    string
	initRedis bool
	initForce bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize DevDock in the current project",
	Run: func(cmd *cobra.Command, args []string) {
		projectDir, _ := ResolveProjectRoot()
		configPath := filepath.Join(projectDir, ".devdock.yml")

		if !initForce {
			if _, err := os.Stat(configPath); err == nil {
				fmt.Println("DevDock is already initialized in this project.")
				fmt.Println("Use --force to overwrite the configuration.")
				os.Exit(1)
			}
		}

		projType := initType
		if projType == "" {
			res := detector.Detect(projectDir)
			if res.Type == "unknown" {
				fmt.Println("DevDock could not detect a supported project type.")
				fmt.Println("To initialize manually, run:")
				fmt.Println("  devdock init --type=laravel")
				fmt.Println("  devdock init --type=nextjs")
				os.Exit(1)
			}
			projType = res.Type
			if res.Confidence == detector.Low {
				fmt.Printf("Detected %s project (Low confidence).\n", res.Type)
				confirm := false
				err := survey.AskOne(&survey.Confirm{
					Message: fmt.Sprintf("Is this a %s project?", res.Type),
					Default: true,
				}, &confirm)
				if err != nil {
					os.Exit(1)
				}
				if !confirm {
					err = survey.AskOne(&survey.Select{
						Message: "Please select the project type:",
						Options: []string{"laravel", "nextjs", "docker-compose"},
					}, &projType)
					if err != nil {
						os.Exit(1)
					}
				}
			} else {
				fmt.Printf("Detected %s project (High confidence).\n", res.Type)
			}
		}

		dirName := filepath.Base(projectDir)
		projName := utils.NormalizeProjectName(dirName)

		cfg := &config.Config{
			Version: "1",
			Project: config.ProjectConfig{
				Name: projName,
				Type: projType,
			},
		}

		if projType == "docker-compose" {
			fmt.Println("Docker Compose project detected. Generating minimal DevDock configuration.")
			writeConfig(configPath, cfg)
			fmt.Println("\nDevDock initialized! DevDock will proxy commands to your existing compose file.")
			fmt.Println("Next step: devdock up")
			return
		}

		// Prompt for services if not provided
		dbService := initDB
		redisService := initRedis

		var recipe *config.Config
		if projType != "docker-compose" {
			r, err := recipes.Load(projType)
			if err == nil {
				recipe = r
			}
		}

		if !cmd.Flags().Changed("db") && !cmd.Flags().Changed("redis") {
			var defaultDB string
			var defaultRedis bool
			if recipe != nil {
				if recipe.Services.Postgres != nil && recipe.Services.Postgres.Enabled {
					defaultDB = "postgres"
				} else if recipe.Services.MySQL != nil && recipe.Services.MySQL.Enabled {
					defaultDB = "mysql"
				}
				if recipe.Services.Redis != nil && recipe.Services.Redis.Enabled {
					defaultRedis = true
				}
			}
			
			var defaultOptions []string
			if defaultDB != "" {
				defaultOptions = append(defaultOptions, defaultDB)
			}
			if defaultRedis {
				defaultOptions = append(defaultOptions, "redis")
			}
			
			var selectedServices []string
			prompt := &survey.MultiSelect{
				Message: "Which services do you want to enable?",
				Options: []string{"postgres", "mysql", "redis"},
				Default: defaultOptions,
			}
			err := survey.AskOne(prompt, &selectedServices)
			if err != nil {
				os.Exit(1)
			}

			for _, s := range selectedServices {
				if s == "postgres" || s == "mysql" {
					dbService = s
				}
				if s == "redis" {
					redisService = true
				}
			}
		}

		if dbService == "postgres" {
			cfg.Services.Postgres = &config.PostgresConfig{
				Enabled: true,
				Version: "15",
				Port:    5432,
			}
		} else if dbService == "mysql" {
			cfg.Services.MySQL = &config.MySQLConfig{
				Enabled: true,
				Version: "8.0",
				Port:    3306,
			}
		}

		if redisService {
			cfg.Services.Redis = &config.RedisConfig{
				Enabled: true,
				Version: "7",
				Port:    6379,
			}
		}

		if recipe != nil {
			cfg.App.Command = recipe.App.Command
			cfg.App.Port = recipe.App.Port
		}

		writeConfig(configPath, cfg)

		composePath := filepath.Join(projectDir, "compose.yml")
		if _, err := os.Stat(composePath); err == nil {
			if !compose.IsDevDockOwned(composePath) && !initForce {
				confirm := false
				err := survey.AskOne(&survey.Confirm{
					Message: "A compose.yml already exists and is not owned by DevDock. Overwrite?",
					Default: false,
				}, &confirm)
				if err != nil || !confirm {
					fmt.Println("Init aborted to protect your compose.yml.")
					os.Exit(1)
				}
			}
		}

		err := compose.Generate(projectDir, cfg)
		if err != nil {
			fmt.Printf("Error generating compose.yml: %v\n", err)
			os.Exit(1)
		}

		err = env.Generate(projectDir, cfg)
		if err != nil {
			fmt.Printf("Error generating environment files: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nDevDock initialized successfully!")
		fmt.Println("Generated .devdock.yml")
		fmt.Println("Generated compose.yml")
		if projType != "laravel" && projType != "nextjs" {
			// Do nothing
		} else {
			fmt.Println("Generated .env and .env.example")
		}
		fmt.Println("\nNext step: devdock up")
	},
}

func writeConfig(path string, cfg *config.Config) {
	b, _ := yaml.Marshal(cfg)
	os.WriteFile(path, b, 0644)
}

func init() {
	initCmd.Flags().StringVar(&initType, "type", "", "Project type (laravel, nextjs, docker-compose)")
	initCmd.Flags().StringVar(&initDB, "db", "", "Database service (postgres, mysql)")
	initCmd.Flags().BoolVar(&initRedis, "redis", false, "Enable Redis")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing configurations")
	rootCmd.AddCommand(initCmd)
}
