package extract

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/eneiss/gar/internal/tar"
	"github.com/eneiss/gar/internal/utils"
	"github.com/spf13/cobra"
)

var OutputDirPath string

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

	cmd.Flags().StringVarP(&OutputDirPath, "output-dir", "o", "", "Output directory path where the archive will be extracted")
	cmd.Flags().BoolP("gzip", "z", false, "Uncompress the archive with gzip")

	return cmd
}

func extractRun(cmd *cobra.Command, args []string) error {
	file_path := args[0] // already validated by cobra

	// Check if file at provided path exists
	if err := utils.FileExists(file_path); err != nil {
		return err
	}

	// Open archive
	file, err := os.ReadFile(file_path)
	if err != nil {
		return err
	}

	// Change directory to output path, if defined
	if OutputDirPath != "" {
		if err := os.Chdir(OutputDirPath); err != nil {
			return fmt.Errorf("could not use target directory to extract files: %v", err)
		}
		fmt.Printf("Changed dir to specified output dir %s\n", OutputDirPath)
	} else {
		fmt.Printf("Extracting files in current working directory\n")
	}

	archive_len_bytes := len(file)
	// check that archive length is a multiple of 512
	if archive_len_bytes%512 != 0 {
		return fmt.Errorf("invalid archive content length: %d (must be a multiple of 512)", archive_len_bytes)
	}

	block_index_512 := 0 // index of the next 512-bytes block to process

	for { // loop on each 'file' inside the archive
		fmt.Printf("--- Starting iteration on block index %d\n", block_index_512)

		// If there are less than 2 512-byte blocks remaining and we still haven't found an
		// end-of-archive marker before, the archive is invalid
		if len(file)-block_index_512*512 < 1024 {
			return fmt.Errorf("invalid archive file, end-of-archive marker not found before end of file")
		}

		// Check if we reached the end of the archive with 2 'empty' blocks marking the end of the archive
		for _, b := range file[block_index_512*512 : (block_index_512+2)*512] {
			if b != 0 {
				break
			}
			fmt.Printf("Found end-of-archive marker blocks, stopping.\n")
			return nil
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

		// Validate header checksum
		is_checksum_valid, err := parsed_header.ValidateChecksum(file[block_index_512*512 : (block_index_512+1)*512])
		if err != nil {
			return fmt.Errorf("could not validate header checksum: %v", err)
		}

		if !is_checksum_valid {
			return fmt.Errorf("corrupted data, invalid header checksum for file at block %d", block_index_512)
		}

		block_index_512 += 1 // end header, start parsing body

		body_size_blocks := int(math.Ceil(float64(parsed_header.Size) / 512.0))
		body_bytes := file[block_index_512*512 : int64(block_index_512*512)+parsed_header.Size]

		fmt.Printf("512-byte blocks used to store contents of %s in archive: %d\n", parsed_header.Name, body_size_blocks)
		fmt.Printf("Raw body content: %s\n", body_bytes)

		// Handle the extraction of different file types
		switch t := parsed_header.Typeflag; t {
		case tar.RegType:
			// write file to system in current working directory (already chdir to target dir)
			if err = os.WriteFile(parsed_header.Name, body_bytes, os.FileMode(parsed_header.Mode)); err != nil {
				return fmt.Errorf("failed to extract file %s: %v", parsed_header.Name, err)
			}

		case tar.DirType:
			if err = os.MkdirAll(parsed_header.Name, os.FileMode(parsed_header.Mode)); err != nil {
				return fmt.Errorf("failed to create archived directory %s: %v", parsed_header.Name, err)
			}

		default:
			return fmt.Errorf("file type not implemented: %c", t)
		}

		// set correct uid/gid
		if err = os.Chown(parsed_header.Name, parsed_header.Uid, parsed_header.Gid); err != nil {
			return fmt.Errorf("failed to apply uid/gid on file %s: %v", parsed_header.Name, err)
		}
		if err = os.Chtimes(parsed_header.Name, time.Unix(parsed_header.Mtime, 0), time.Unix(parsed_header.Mtime, 0)); err != nil {
			return fmt.Errorf("failed to update access/modification times on file %s: %v", parsed_header.Name, err)
		}

		// stop parsing body, go to next file/end-of-archive
		block_index_512 += body_size_blocks
	}

	return nil
}
