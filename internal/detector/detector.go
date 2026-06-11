package detector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Confidence represents the certainty of detection.
type Confidence string

const (
	High Confidence = "High"
	Low  Confidence = "Low"
	None Confidence = "None"
)

// Result holds the outcome of a detection run.
type Result struct {
	Type       string     `json:"type"`
	Confidence Confidence `json:"confidence"`
	Reasons    []string   `json:"reasons"`
}

func fileExists(dir, filename string) bool {
	_, err := os.Stat(filepath.Join(dir, filename))
	return err == nil
}

// Detect scans the given directory and identifies the project type based on exact filenames.
func Detect(projectDir string) *Result {
	hasComposeYml := fileExists(projectDir, "compose.yml")
	hasDockerComposeYml := fileExists(projectDir, "docker-compose.yml")

	hasComposerJson := fileExists(projectDir, "composer.json")
	hasArtisan := fileExists(projectDir, "artisan")

	hasPackageJson := fileExists(projectDir, "package.json")
	hasNextConfigJs := fileExists(projectDir, "next.config.js")
	hasNextConfigTs := fileExists(projectDir, "next.config.ts")

	// Docker Compose
	if hasComposeYml || hasDockerComposeYml {
		var reasons []string
		if hasComposeYml {
			reasons = append(reasons, "Found compose.yml")
		}
		if hasDockerComposeYml {
			reasons = append(reasons, "Found docker-compose.yml")
		}
		return &Result{
			Type:       "docker-compose",
			Confidence: High,
			Reasons:    reasons,
		}
	}

	// Laravel (High)
	if hasComposerJson && hasArtisan {
		return &Result{
			Type:       "laravel",
			Confidence: High,
			Reasons:    []string{"Found composer.json", "Found artisan"},
		}
	}

	// Next.js (High)
	if hasPackageJson && (hasNextConfigJs || hasNextConfigTs) {
		var reasons []string
		reasons = append(reasons, "Found package.json")
		if hasNextConfigJs {
			reasons = append(reasons, "Found next.config.js")
		}
		if hasNextConfigTs {
			reasons = append(reasons, "Found next.config.ts")
		}
		return &Result{
			Type:       "nextjs",
			Confidence: High,
			Reasons:    reasons,
		}
	}

	// Laravel (Low)
	if hasComposerJson {
		return &Result{
			Type:       "laravel",
			Confidence: Low,
			Reasons:    []string{"Found composer.json"},
		}
	}

	// Next.js (Low)
	if hasPackageJson {
		data, err := os.ReadFile(filepath.Join(projectDir, "package.json"))
		if err == nil {
			var pkg map[string]interface{}
			if err := json.Unmarshal(data, &pkg); err == nil {
				deps, _ := pkg["dependencies"].(map[string]interface{})
				if deps != nil {
					if _, ok := deps["next"]; ok {
						return &Result{Type: "nextjs", Confidence: High, Reasons: []string{"Found next in dependencies"}}
					}
					if _, ok := deps["express"]; ok {
						return &Result{Type: "express", Confidence: High, Reasons: []string{"Found express in dependencies"}}
					}
				}

				devDeps, _ := pkg["devDependencies"].(map[string]interface{})
				if devDeps != nil {
					if _, ok := devDeps["next"]; ok {
						return &Result{Type: "nextjs", Confidence: High, Reasons: []string{"Found next in devDependencies"}}
					}
					if _, ok := devDeps["express"]; ok {
						return &Result{Type: "express", Confidence: High, Reasons: []string{"Found express in devDependencies"}}
					}
				}
			}
		}

		return &Result{
			Type:       "express",
			Confidence: Low,
			Reasons:    []string{"Found package.json, but neither next nor express listed"},
		}
	}

	hasGoMod := fileExists(projectDir, "go.mod")
	if hasGoMod {
		b, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
		if err == nil {
			content := string(b)
			if strings.Contains(content, "github.com/gofiber/fiber") {
				return &Result{
					Type:       "fiber",
					Confidence: High,
					Reasons:    []string{"Found github.com/gofiber/fiber in go.mod"},
				}
			}
		}
		
		return &Result{
			Type:       "fiber",
			Confidence: Low,
			Reasons:    []string{"Found go.mod, but gofiber not listed"},
		}
	}

	// Unknown
	return &Result{
		Type:       "unknown",
		Confidence: None,
		Reasons:    []string{"No supported project files found"},
	}
}
