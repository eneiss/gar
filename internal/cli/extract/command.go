package extract

import (
	"fmt"

	"github.com/eneiss/gar/internal/utils"
	"github.com/spf13/cobra"
)

var Output string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract an archive",
		Args:  cobra.MinimumNArgs(1),
		Long: `Extract the target archive.
For example:
		TODO
`,
		RunE: extractRun,
	}

	cmd.Flags().StringVarP(&Output, "output", "o", "", "Output path where the archive will be extracted")
	cmd.Flags().BoolP("gzip", "z", false, "Uncompress the archive with gzip")

	return cmd
}

func extractRun(cmd *cobra.Command, args []string) error {
	fmt.Println("extract called")

	for _, arg := range args {
		fmt.Println(arg)
	}

	// Check if all provided arguments are existing files
	if err := utils.FilesExist(args); err != nil {
		return err
	}
	return nil
}
