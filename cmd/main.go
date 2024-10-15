package main

import (
	"context"
	"path"

	"github.com/spf13/cobra"
	config "github.com/xybor/todennus-config"
	"github.com/xybor/todennus-migration/postgres"
)

var downFlag int
var postgresFlag bool

var rootCommand = &cobra.Command{
	Use:   "todennus-migration",
	Short: "Migrate database of Todennus",
	Run: func(cmd *cobra.Command, args []string) {
		envPaths, err := cmd.Flags().GetStringArray("env")
		if err != nil {
			panic(err)
		}

		migrationPath, err := cmd.Flags().GetString("path")
		if err != nil {
			panic(err)
		}

		config, err := config.Load(envPaths...)
		if err != nil {
			panic(err)
		}

		if postgresFlag {
			postgresPath := generateMigrationPath(migrationPath, "postgres")

			gormDB, err := postgres.Initialize(context.Background(), config)
			if err != nil {
				panic(err)
			}

			db, err := gormDB.DB()
			if err != nil {
				panic(err)
			}

			if downFlag == 0 {
				if err := postgres.Up(context.Background(), db, postgresPath); err != nil {
					panic(err)
				}
			} else {
				if err := postgres.Down(context.Background(), db, postgresPath, downFlag); err != nil {
					panic(err)
				}
			}
		}
	},
}

func main() {
	rootCommand.Flags().StringArray("env", []string{".env"}, "environment file paths")
	rootCommand.Flags().String("path", "./postgres/migration", "migration path")

	rootCommand.Flags().IntVar(&downFlag, "down", 0, "Migrate down with the number of steps")
	rootCommand.Flags().BoolVar(&postgresFlag, "postgres", false, "Migrate postgres database")

	if err := rootCommand.Execute(); err != nil {
		panic(err)
	}
}

func generateMigrationPath(parent, db string) string {
	return path.Join(parent, db, "migration")
}
