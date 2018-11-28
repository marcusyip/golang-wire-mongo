package migrate

import (
	"github.com/spf13/cobra"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migrate",
		Run: func(cmd *cobra.Command, args []string) {
			job, err := NewMigrationJob()
			if err != nil {
				panic(err)
			}
			job.Run()
		},
	}
	return cmd
}
