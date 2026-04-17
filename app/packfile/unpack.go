package packfile

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"log"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

const (
	objCommit   = 1
	objTree     = 2
	objBlob     = 3
	objOfsDelta = 6
	objRefDelta = 7
)

// Unpack parses a git packfile, writing every non-delta object directly and
// resolving ofs-delta chains against their bases before writing. Ref-deltas
// are not supported.
func Unpack(packFile []byte) {
	if len(packFile) < 12 || string(packFile[:4]) != "PACK" {
		log.Fatalf("invalid packfile header")
	}
	numObjects := int(binary.BigEndian.Uint32(packFile[8:12]))
	r := bytes.NewReader(packFile[12:])

	resolvedDeltasCache := make(map[int]deltaCache)

	for range numObjects {
		ownOffset := 12 + int(r.Size()-int64(r.Len()))
		objType, _ := readObjectHeader(r)

		switch objType {
		case objRefDelta:
			log.Fatalf("Ref delta is not handled")

		case objOfsDelta:
			negOff := readOfsOffset(r)
			instructions := decompressZlib(r)
			ofs := &ofsDeltaObject{
				instructions: instructions,
				baseOffset:   ownOffset - negOff,
			}
			resolvedDeltasCache[ownOffset] = ofs.resolveDelta(packFile, resolvedDeltasCache)

		default:
			content := decompressZlib(r)
			utils.HashObject(content, typeName(objType), utils.HashOptions{WriteHash: true})
			resolvedDeltasCache[ownOffset] = deltaCache{
				objectType:    typeName(objType),
				resolvedDelta: content,
			}
		}
	}
}

// decompressZlib reads a zlib stream from r and returns the decompressed bytes.
func decompressZlib(r io.Reader) []byte {
	zr, err := zlib.NewReader(r)
	if err != nil {
		log.Fatalf("zlib reader: %v", err)
	}
	defer zr.Close()
	data, err := io.ReadAll(zr)
	if err != nil {
		log.Fatalf("zlib read: %v", err)
	}
	return data
}

// readObjectHeader reads a packfile object header from r, returning the type
// and the encoded (inflated-length) size.
func readObjectHeader(r *bytes.Reader) (byte, int) {
	b, err := r.ReadByte()
	if err != nil {
		log.Fatalf("read object header: %v", err)
	}
	objType := (b >> 4) & 0x7
	size := int(b & 0x0f)
	shift := uint(4)
	for b&0x80 != 0 {
		b, err = r.ReadByte()
		if err != nil {
			log.Fatalf("read object size: %v", err)
		}
		size |= int(b&0x7f) << shift
		shift += 7
	}
	return objType, size
}

func typeName(t byte) string {
	switch t {
	case objCommit:
		return "commit"
	case objTree:
		return "tree"
	case objBlob:
		return "blob"
	default:
		return "tag"
	}
}
