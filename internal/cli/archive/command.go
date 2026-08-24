package archive

import (
	"fmt"

	"github.com/eneiss/gar/internal/utils"
	"github.com/spf13/cobra"
)

var Output string
var BlockingFactor int

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Create an archive",
		Args:  cobra.MinimumNArgs(1),
		Long: `Create a tar or compressed archive of the provided target path.
For example:
	TODO
`,
		RunE: archiveRun,
	}

	cmd.Flags().StringVarP(&Output, "output", "o", "", "Output path where the archive will be written")
	cmd.Flags().BoolP("gzip", "z", false, "Compress the archive with gzip")
	cmd.Flags().IntVarP(&BlockingFactor, "blocking-factor", "b", 20, "Set record size as a multiple of 512 bytes")

	return cmd
}

func archiveRun(cmd *cobra.Command, args []string) error {
	fmt.Println("archive called")

	for _, arg := range args {
		fmt.Println(arg)
	}

	// Check if all provided arguments are existing files
	if err := utils.FilesExist(args); err != nil {
		return err
	}
	return nil
}
