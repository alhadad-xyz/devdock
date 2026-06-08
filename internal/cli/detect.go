package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"devdock/internal/detector"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect the project type in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		projectDir, _ := ResolveProjectRoot()
		res := detector.Detect(projectDir)

		if jsonFlag {
			b, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(b))
			return
		}

		if res.Type == "unknown" {
			fmt.Println("DevDock could not detect a supported project type.")
			fmt.Println()
			fmt.Println("Supported types in this version: laravel, nextjs, docker-compose")
			fmt.Println()
			fmt.Println("To initialize manually, run:")
			fmt.Println("  devdock init --type=laravel")
			fmt.Println("  devdock init --type=nextjs")
			os.Exit(1)
		}

		fmt.Printf("Detected project type: %s\n", res.Type)
		fmt.Printf("Confidence: %s\n", res.Confidence)
		fmt.Println("Reasons:")
		for _, r := range res.Reasons {
			fmt.Printf("  - %s\n", r)
		}

		if res.Confidence == detector.Low {
			fmt.Println()
			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Is this a %s project?", res.Type),
				Default: true,
			}
			err := survey.AskOne(prompt, &confirm)
			if err != nil {
				os.Exit(1)
			}

			if !confirm {
				var selectedType string
				selectPrompt := &survey.Select{
					Message: "Please select the project type:",
					Options: []string{"laravel", "nextjs", "docker-compose"},
				}
				err = survey.AskOne(selectPrompt, &selectedType)
				if err != nil {
					os.Exit(1)
				}
				res.Type = selectedType
				fmt.Printf("\nSelected project type: %s\n", res.Type)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
