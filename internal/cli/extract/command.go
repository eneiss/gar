package extract

import (
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

	archive_len_bytes := len(file)
	// check that archive length is a multiple of 512
	if archive_len_bytes%512 != 0 {
		return fmt.Errorf("invalid archive content length: %d (must be a multiple of 512)", archive_len_bytes)
	}

	block_index_512 := 0 // index of the next 512-bytes block to process

	for { // loop on each 'file' inside the archive
		fmt.Printf("--- Starting iteration on block index %d\n", block_index_512)

		// Check if we reached the end of the archive with 2 'empty' blocks marking the end of the archive
		if len(file)-block_index_512*512 == 1024 {
			// Check if all bytes are null
			for i, b := range file[block_index_512*512 : (block_index_512+2)*512] {
				if b != 0 {
					return fmt.Errorf("non-zero end of archive byte at index %s, invalid file", block_index_512*512+i)
				}
			}
			fmt.Printf("Found end-of-archive blocks, stopping.")
			break
		}

		// Retrieve header info
		header, err := tar.BuildRawHeader(file[block_index_512*512:])
		if err != nil {
			return fmt.Errorf("could not build raw header: %v", err)
		}
		parsed_header, err := header.Parse()
		if err != nil {
			return fmt.Errorf("could not parse raw header fields: %v", err)
		}

		parsed_header.Print()

		block_index_512 += 1 // end header, start parsing body

		// TODO: get body
		body_size_blocks := int(math.Ceil(float64(parsed_header.Size) / 512.0))

		fmt.Printf("512-byte blocks used to store contents of %s in archive: %d\n", parsed_header.Name, body_size_blocks)

		fmt.Printf("Raw body content: %s\n", file[block_index_512*512:int64(block_index_512*512)+parsed_header.Size])

		// stop parsing body, go to next file/end-of-archive
		block_index_512 += body_size_blocks
	}

	return nil
}
