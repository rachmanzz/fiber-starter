package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
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
	Short: "Run database migrations using goose",
	Run: func(cmd *cobra.Command, args []string) {
		prepareGooseEnvironment()

		gooseArgs := []string{"up"}
		if migrateTo != "" {
			gooseArgs = append(gooseArgs, migrateTo)
		}

		executeGoose(gooseArgs)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the most recently applied migration",
	Run: func(cmd *cobra.Command, args []string) {
		prepareGooseEnvironment()

		executeGoose([]string{"down"})
	},
}

var migrateNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prepareGooseEnvironment()

		executeGooseNoDB([]string{"create", args[0], "sql"})
	},
}

func prepareGooseEnvironment() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}
	ensureGooseInstalled()

	// Ensure migrations directory exists
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		fmt.Println("Creating migrations directory...")
		if err := os.Mkdir("migrations", 0755); err != nil {
			fmt.Printf("Failed to create migrations directory: %v\n", err)
			os.Exit(1)
		}
	}
}

func buildDSN() string {
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     net.JoinHostPort(os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: "sslmode=" + sslMode,
	}
	return u.String()
}

func executeGoose(args []string) {
	gooseArgs := append([]string{"-dir", "migrations", "postgres", buildDSN()}, args...)
	display := "goose -dir migrations postgres <dsn> " + strings.Join(args, " ")
	runGoose(gooseArgs, display)
}

func executeGooseNoDB(args []string) {
	gooseArgs := append([]string{"-dir", "migrations"}, args...)
	display := "goose -dir migrations " + strings.Join(args, " ")
	runGoose(gooseArgs, display)
}

func runGoose(args []string, display string) {
	fmt.Printf("Running: %s\n", display)
	runCmd := exec.Command("goose", args...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Env = os.Environ()
	if err := runCmd.Run(); err != nil {
		fmt.Printf("Goose command failed: %v\n", err)
		os.Exit(1)
	}
}

func ensureGooseInstalled() {
	_, err := exec.LookPath("goose")
	if err != nil {
		fmt.Println("goose not found. Installing github.com/pressly/goose/v3@latest...")
		installCmd := exec.Command("go", "install", "github.com/pressly/goose/v3/cmd/goose@latest")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Printf("Failed to install goose: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("goose installed successfully.")
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
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateNewCmd)
	rootCmd.AddCommand(migrateCmd)
}
