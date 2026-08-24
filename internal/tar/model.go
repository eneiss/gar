package tar

import (
	"fmt"
	"github.com/eneiss/gar/internal/utils"
)

type PosixHeader struct {
	Name     [100]byte
	Mode     [8]byte
	Uid      [8]byte
	Gid      [8]byte
	Size     [12]byte
	Mtime    [12]byte
	Chksum   [8]byte
	Typeflag byte
	Linkname [100]byte
	// NOTE: the following are replaced by padding [255]byte
	// magic    [6]byte
	// version  [2]byte
	// uname    [32]byte
	// gname    [32]byte
	// devmajor [8]byte
	// devminor [8]byte
	// prefix   [155]byte
	Padding [255]byte
}

func BuildHeader(raw_input []byte) (*PosixHeader, error) {
	if len(raw_input) < 255 {
		return nil, fmt.Errorf("cannot build tar header from raw input: input too short (expected at least 255 bytes, got %d)", len(raw_input))
	}

	header := PosixHeader{
		Name:     [100]byte(raw_input[0:100]),
		Mode:     [8]byte(raw_input[100:108]),
		Uid:      [8]byte(raw_input[108:116]),  // base 8, last byte is 0x00
		Gid:      [8]byte(raw_input[116:124]),  // base 8, last byte is 0x00
		Size:     [12]byte(raw_input[124:136]), // base 8, last byte is 0x00
		Mtime:    [12]byte(raw_input[136:148]),
		Chksum:   [8]byte(raw_input[148:156]),
		Typeflag: byte(raw_input[156]),
		Linkname: [100]byte(raw_input[157:257]),
		Padding:  [255]byte(raw_input[257:512]),
	}

	return &header, nil
}

// Prints the content of the POSIX header h in a human-friendly way
func (h PosixHeader) Print() {
	fmt.Printf("File name: %s\n", h.Name)
	fmt.Printf("Mode: %s\n", h.Mode[0:8])
	fmt.Printf("Uid: %d\n", utils.Base8ToBase10(h.Uid[:7]))
	fmt.Printf("Gid: %d\n", utils.Base8ToBase10(h.Gid[:7]))
	fmt.Printf("Size: %d bytes\n", utils.Base8ToBase10(h.Size[:11]))
}
