package extract

import (
	"bytes"
	"fmt"
	"math"
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

	// TODO: check if file length is a multiple of 512, return otherwise
	archive_len_bytes := len(file)
	if archive_len_bytes%512 != 0 {
		return fmt.Errorf("invalid archive content length: %d (must be a multiple of 512)", archive_len_bytes)
	}

	block_index_512 := 0 // index of the next 512-bytes block to process

	for {
		// Check if we reached the end of the archive with 2 'empty' blocks marking the end of the archive
		if len(file)-block_index_512*512 == 1024 &&
			bytes.Count(file[block_index_512*512:(block_index_512+2)*512], []byte("\x00")) == 512 {

			fmt.Printf("Found end-of-archive blocks, stopping.")
			break
		}
		// TODO: break if invalid file format too or arriving at end of indices

		// Retrieve header info
		header, err := tar.BuildRawHeader(file[block_index_512*512:])
		if err != nil {
			return fmt.Errorf("could not build raw header: %v", err)
		}
		parsedHeader, err := header.Parse()
		if err != nil {
			return fmt.Errorf("could not parse raw header fields: %v", err)
		}

		parsedHeader.Print()

		block_index_512 += 1 // end header, start parsing body

		// TODO: get body
		body_size_blocks := int(math.Ceil(float64(parsedHeader.Size) / 512.0))

		// stop parsing body, go to next file/end-of-archive
		block_index_512 += body_size_blocks
	}
	return nil
}
