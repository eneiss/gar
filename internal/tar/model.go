package tar

// Sources for specifications:
// [1]: https://ftp.gnu.org/old-gnu/Manuals/tar-1.12/html_node/tar_123.html
// [2]: https://mort.coffee/home/tar/

import (
	"fmt"
	"github.com/eneiss/gar/internal/utils"
	"time"
)

type typeFlag byte

const (
	RegType   typeFlag = '0' // regular file
	ARegType  typeFlag = 0x0 // backward-compatible variant for regular file, NUL byte
	LinkType  typeFlag = '1' // link (any type)
	SymType   typeFlag = '2' // symlink
	CharType  typeFlag = '3' // character special
	BlockType typeFlag = '4' // block special
	DirType   typeFlag = '5' // directory
	FifoType  typeFlag = '6' // FIFO special
	ContType  typeFlag = '7' // reserved
)

// TODO: parse into usable data structures (especially for integer values?)
type PosixHeader struct {
	// The name, linkname, magic, uname, and gname are null-terminated character
	// strings. All other fileds are zero-filled octal numbers in ASCII. Each
	// numeric field of width w contains w minus 2 digits, a space, and a null,
	// except size, and mtime, which do not contain the trailing null. [1]
	Name     [100]byte
	Mode     [8]byte
	Uid      [8]byte
	Gid      [8]byte
	Size     [12]byte
	Mtime    [12]byte
	Chksum   [8]byte
	Typeflag byte
	Linkname [100]byte
	// NOTE: the following are replaced by Padding [255]byte (source: [2])
	// also makes it simpler
	// magic    [6]byte
	// version  [2]byte
	// uname    [32]byte
	// gname    [32]byte
	// devmajor [8]byte
	// devminor [8]byte
	// prefix   [155]byte
	Padding [255]byte
}

type ParsedPosixHeader struct {
	Name     string
	Mode     [7]byte // TBD, maybe octal string is fine
	Uid      int
	Gid      int
	Size     int64 // NOTE: 8^12 = 2^36 > size of int (32 bits)
	Mtime    int64
	Chksum   int
	Typeflag typeFlag
	Linkname string
	Padding  [255]byte // TBD
}

func BuildRawHeader(raw_input []byte) (*PosixHeader, error) {
	if len(raw_input) < 255 {
		return nil, fmt.Errorf("cannot build tar header from raw input: input too short (expected at least 255 bytes, got %d)", len(raw_input))
	}

	// Source: [1]
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

func (r PosixHeader) ParseRawHeader() *ParsedPosixHeader {
	res := ParsedPosixHeader{
		Name:     string(r.Name[:99]),
		Mode:     [7]byte(r.Mode[:7]),
		Typeflag: typeFlag(r.Typeflag),
		Linkname: string(r.Linkname[:99]),
		Padding:  r.Padding,
	}
	// Processed fields (octal base conversion)
	// Uid
	uid, err := utils.Base8ToBase10(r.Uid[:7])
	if err != nil {
		return nil
	}
	res.Uid = int(uid)
	// Gid
	gid, err := utils.Base8ToBase10(r.Gid[:7])
	if err != nil {
		return nil
	}
	res.Gid = int(gid)
	// Size
	size, err := utils.Base8ToBase10(r.Size[:11])
	if err != nil {
		return nil
	}
	res.Size = int64(size)
	// Mtime
	mtime, err := utils.Base8ToBase10(r.Mtime[:11])
	if err != nil {
		return nil
	}
	res.Mtime = int64(mtime)
	// Chksum
	chksum, err := utils.Base8ToBase10(r.Chksum[:6])
	if err != nil {
		return nil
	}
	res.Chksum = int(chksum)

	return &res
}

// Prints the content of the POSIX header h in a human-friendly way
func (h PosixHeader) Print() {
	fmt.Printf("File name: %s\n", h.Name)
	fmt.Printf("Mode: %s\n", h.Mode[0:8])

	base10Uid, err := utils.Base8ToBase10(h.Uid[:7])
	if err != nil {
		fmt.Printf("Error while displaying header uid: %v", err)
		return
	}
	fmt.Printf("Uid: %d (ASCII: %s)\n", base10Uid, h.Uid)

	base10Gid, err := utils.Base8ToBase10(h.Gid[:7])
	if err != nil {
		fmt.Printf("Error while displaying header gid: %v", err)
		return
	}
	fmt.Printf("Gid: %d (ASCII: %s)\n", base10Gid, h.Gid)

	base10Size, err := utils.Base8ToBase10(h.Size[:11])
	if err != nil {
		fmt.Printf("Error while displaying header size: %v", err)
		return
	}
	fmt.Printf("Size: %d (ASCII: %s)\n", base10Size, h.Size)

	base10Mtime, err := utils.Base8ToBase10(h.Mtime[:11])
	if err != nil {
		fmt.Printf("Error while displaying header mtime: %v", err)
		return
	}
	fmt.Printf("Modification time: %s (ASCII: %s)\n", time.Unix(int64(base10Mtime), 0), h.Mtime)

	base10Checksum, err := utils.Base8ToBase10(h.Chksum[:6]) // last 2 bytes are NUL and SP (space)
	if err != nil {
		fmt.Printf("Error while displaying header chksum: %v", err)
		return
	}
	fmt.Printf("Checksum (base 10 result): %d (ASCII: %s)\n", base10Checksum, h.Chksum)

	fmt.Printf("File type: %c\n", h.Typeflag)
}

// TODO: checksum of header to check for data corruption
// The chksum field is the ASCII representation of the octal value of the
// simple sum of all bytes in the header block. Each 8-bit byte in the header
// is added to an unsigned integer, initialized to zero, the precision of which
// shall be no less than seventeen bits. When calculating the checksum, the
// chksum field is treated as if it were all blanks. [1]
