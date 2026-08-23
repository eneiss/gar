package extract

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract an archive",
		Long: `Extract the target archive.
For example:
		TODO
`,
		Run: extractRun,
	}

	cmd.Flags().Bool("gzip", false, "Uncompress the archive with gzip")

	return cmd
}

func extractRun(cmd *cobra.Command, args []string) {
	fmt.Println("extract called")

	for _, arg := range args {
		fmt.Println(arg)
	}
}
