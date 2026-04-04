package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spark",
	Short: "Spark - A CLI for building Fiber v3 projects",
	Long:  `Spark is a CLI tool designed to simplify the development and build process for Fiber v3 projects.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Spark",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Spark CLI v0.0.1")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project with a new module name",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter New Module Name (e.g., github.com/user/repo): ")
		modName, _ := reader.ReadString('\n')
		modName = strings.TrimSpace(modName)

		if modName == "" {
			fmt.Println("Error: Module name cannot be empty")
			return
		}

		oldModule := "github.com/rachmanzz/fiber-starter"

		fmt.Printf("Initializing project with module: %s...\n", modName)

		err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if strings.HasPrefix(d.Name(), ".") || d.Name() == "spark-cli" || d.Name() == "node_modules" || d.Name() == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}
			read, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(read), oldModule) {
				newContents := strings.ReplaceAll(string(read), oldModule, modName)
				info, err := d.Info()
				if err != nil {
					return err
				}
				err = os.WriteFile(path, []byte(newContents), info.Mode())
				if err != nil {
					return err
				}
				fmt.Printf("Updated: %s\n", path)
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error during file replacement: %v\n", err)
			return
		}

		fmt.Println("Running go mod edit...")
		exec.Command("go", "mod", "edit", "-module", modName).Run()
		fmt.Println("Running go mod tidy...")
		exec.Command("go", "mod", "tidy").Run()

		fmt.Printf("Successfully initialized project with module: %s\n", modName)
	},
}

var migrateTo string

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations using tern",
	Run: func(cmd *cobra.Command, args []string) {
		prepareTernEnvironment()

		ternArgs := []string{"migrate"}
		if migrateTo != "" {
			ternArgs = append(ternArgs, "--destination", migrateTo)
		}

		executeTern(ternArgs)
	},
}

var migrateNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prepareTernEnvironment()

		name := args[0]
		ternArgs := []string{"new", name}

		executeTern(ternArgs)
	},
}

func prepareTernEnvironment() {
	ensureTernInstalled()
	ensureTernConfig()

	// Ensure migrations directory exists
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		fmt.Println("Creating migrations directory...")
		if err := os.Mkdir("migrations", 0755); err != nil {
			fmt.Printf("Failed to create migrations directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Set TERN_MIGRATIONS if not already set
	if os.Getenv("TERN_MIGRATIONS") == "" {
		fmt.Println("Setting TERN_MIGRATIONS=migrations")
		os.Setenv("TERN_MIGRATIONS", "migrations")
	}
}

func ensureTernConfig() {
	if _, err := os.Stat("tern.conf"); os.IsNotExist(err) {
		fmt.Println("tern.conf not found. Initializing with tern init...")
		initCmd := exec.Command("tern", "init")
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		if err := initCmd.Run(); err != nil {
			fmt.Printf("Failed to run tern init: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Updating tern.conf with environment variables...")
		configContent := `[database]
host = {{env "DB_HOST"}}
port = {{env "DB_PORT"}}
database = {{env "DB_NAME"}}
user = {{env "DB_USER"}}
password = {{env "DB_PASSWORD"}}
# sslmode = {{env "DB_SSLMODE"}}

[data]
`
		if err := os.WriteFile("tern.conf", []byte(configContent), 0644); err != nil {
			fmt.Printf("Failed to write tern.conf: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("tern.conf updated successfully.")
	}
}

func executeTern(args []string) {
	fmt.Printf("Running: tern %s\n", strings.Join(args, " "))
	runCmd := exec.Command("tern", args...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	// Pass the current environment (including our newly set TERN_MIGRATIONS)
	runCmd.Env = os.Environ()
	if err := runCmd.Run(); err != nil {
		fmt.Printf("Tern command failed: %v\n", err)
	}
}

func ensureTernInstalled() {
	_, err := exec.LookPath("tern")
	if err != nil {
		fmt.Println("tern not found. Installing github.com/jackc/tern/v2@latest...")
		installCmd := exec.Command("go", "install", "github.com/jackc/tern/v2@latest")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Printf("Failed to install tern: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("tern installed successfully.")
	}
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run the application with live reloading using air",
	Run: func(cmd *cobra.Command, args []string) {
		ensureAirInstalled()

		fmt.Println("🚀 Application Start with Air...")
		airCmd := exec.Command("air")
		airCmd.Stdout = os.Stdout
		airCmd.Stderr = os.Stderr
		airCmd.Stdin = os.Stdin
		if err := airCmd.Run(); err != nil {
			fmt.Printf("Air failed: %v\n", err)
		}
	},
}

func ensureAirInstalled() {
	_, err := exec.LookPath("air")
	if err != nil {
		fmt.Println("air not found. Installing github.com/air-verse/air@latest...")
		installCmd := exec.Command("go", "install", "github.com/air-verse/air@latest")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Printf("Failed to install air: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("air installed successfully.")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(devCmd)

	migrateCmd.Flags().StringVarP(&migrateTo, "to", "t", "", "destination migration version")
	migrateCmd.AddCommand(migrateNewCmd)
	rootCmd.AddCommand(migrateCmd)
}
