package tar

// Sources for specifications:
// [1]: https://ftp.gnu.org/old-gnu/Manuals/tar-1.12/html_node/tar_123.html
// [2]: https://mort.coffee/home/tar/

import (
	"bytes"
	"fmt"
	"github.com/eneiss/gar/internal/utils"
	"time"
)

type TypeFlag byte

const (
	RegType   TypeFlag = '0' // regular file
	ARegType  TypeFlag = 0x0 // backward-compatible variant for regular file, NUL byte
	LinkType  TypeFlag = '1' // link (any type)
	SymType   TypeFlag = '2' // symlink
	CharType  TypeFlag = '3' // character special
	BlockType TypeFlag = '4' // block special
	DirType   TypeFlag = '5' // directory
	FifoType  TypeFlag = '6' // FIFO special
	ContType  TypeFlag = '7' // reserved
)

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
	Mode     int
	Uid      int
	Gid      int
	Size     int64 // NOTE: 8^12 = 2^36 > size of int (32 bits)
	Mtime    int64
	Chksum   int
	Typeflag TypeFlag
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

func (r *PosixHeader) Parse() (*ParsedPosixHeader, error) {
	res := ParsedPosixHeader{
		// NOTE: os.WriteFile behaves strangely when null bytes are present in the
		// target file name, so we trim them
		// https://github.com/golang/go/issues/24195
		Name:     string(bytes.Trim(r.Name[:99], "\x00")),
		Typeflag: TypeFlag(r.Typeflag),
		Linkname: string(bytes.Trim(r.Linkname[:99], "\x00")),
		Padding:  r.Padding,
	}
	// Processed fields (octal base conversion)
	// Mode
	mode, err := utils.Base8ToBase10(r.Mode[4:7]) // TODO: find what to do with first bytes
	if err != nil {
		return nil, fmt.Errorf("cloud not parse octal value for header section mode: %v", err)
	}
	res.Mode = int(mode)
	// Uid
	uid, err := utils.Base8ToBase10(r.Uid[:7])
	if err != nil {
		return nil, fmt.Errorf("cloud not parse octal value for header section uid: %v", err)
	}
	res.Uid = int(uid)
	// Gid
	gid, err := utils.Base8ToBase10(r.Gid[:7])
	if err != nil {
		return nil, fmt.Errorf("cloud not parse octal value for header section gid: %v", err)
	}
	res.Gid = int(gid)
	// Size
	size, err := utils.Base8ToBase10(r.Size[:11])
	if err != nil {
		return nil, fmt.Errorf("cloud not parse octal value for header section size: %v", err)
	}
	res.Size = int64(size)
	// Mtime
	mtime, err := utils.Base8ToBase10(r.Mtime[:11])
	if err != nil {
		return nil, fmt.Errorf("cloud not parse octal value for header section mtime: %v", err)
	}
	res.Mtime = int64(mtime)
	// Chksum
	chksum, err := utils.Base8ToBase10(r.Chksum[:6])
	if err != nil {
		return nil, fmt.Errorf("cloud not parse octal value for header section chksum: %v", err)
	}
	res.Chksum = int(chksum)

	return &res, nil
}

// Prints the content of the parsed POSIX header h in a human-friendly way
func (h ParsedPosixHeader) Print() {
	fmt.Printf("File name: %s\n", h.Name)
	fmt.Printf("Mode: %o\n", h.Mode)
	fmt.Printf("Uid: %d\n", h.Uid)
	fmt.Printf("Gid: %d\n", h.Gid)
	fmt.Printf("Size: %d\n", h.Size)
	fmt.Printf("Modification time: %s\n", time.Unix(h.Mtime, 0))
	fmt.Printf("Checksum (base 10 result): %d\n", h.Chksum)
	fmt.Printf("File type: %c\n", h.Typeflag)
}

// Computes the checksum of the raw POSIX header bytes header_bytes and
// compares it to the checksum value in the parsed header h.
// Returns true if the computed checksum matches h's checksum field value,
// false otherwise. Return an error if the header bytes array isn't the right
// length.
// When computing the checksum, the checksum field value is assumed to be
// blank. ([1])
func (h ParsedPosixHeader) ValidateChecksum(header_bytes []byte) (bool, error) {
	if len(header_bytes) != 512 {
		return false, fmt.Errorf("invalid header length: got %d, wanted 512", len(header_bytes))
	}

	sum := 0
	for i, b := range header_bytes {
		if i < 148 || i >= 156 { // filter out checksum bytes in computation
			sum += int(b)
		} else {
			// > When calculating the checksum, the chksum field is treated as if it were all blanks. [1]
			// Apparently blank litterally meant ' ' (SP character, 0x20)
			// https://github.com/gitGNU/gnu_tar/blob/da8d0659a6fe8faf76b3a3275cf1f403e78edb1f/src/list.c#L370
			sum += ' '
		}
	}

	fmt.Printf("Computed checksum (base 10): %d, header checksum field value: %d\n", sum, h.Chksum)

	return h.Chksum == sum, nil
}
