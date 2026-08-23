package archive

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Create an archive",
		Long: `Create a tar or compressed archive of the provided target path.
For example:
	TODO
`,
		Run: archiveRun,
	}

	cmd.Flags().Bool("gzip", false, "Compress the archive with gzip")

	return cmd
}

func archiveRun(cmd *cobra.Command, args []string) {
	fmt.Println("archive called")

	for _, arg := range args {
		fmt.Println(arg)
	}
}
