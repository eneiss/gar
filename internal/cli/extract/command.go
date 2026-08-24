package extract

import (
	// "fmt"
	"os"

	"github.com/eneiss/gar/internal/tar"
	"github.com/eneiss/gar/internal/utils"
	"github.com/spf13/cobra"
)

var Output string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract an archive",
		Args:  cobra.ExactArgs(1),
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
	file_path := args[0] // already validated by cobra

	// Check if file at provided path exists
	if err := utils.FileExists(file_path); err != nil {
		return err
	}

	file, err := os.ReadFile(file_path)
	if err != nil {
		return err
	}

	header, err := tar.BuildHeader(file)
	if err != nil {
		return err
	}

	header.Print()

	return nil
}
