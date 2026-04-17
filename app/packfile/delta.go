package packfile

import (
	"bytes"
	"io"
	"log"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

type ofsDeltaObject struct {
	instructions []byte
	baseOffset   int
}

type deltaCache struct {
	objectType    string
	resolvedDelta []byte
}

// resolveDelta resolves this ofs-delta by applying its instructions to its
// base object (recursively resolving the base if it is itself an ofs-delta),
// writes the resulting object to .git/objects with its real type, and returns
// the resolved entry. resolvedDeltasCache memoizes offsets already materialized.
func (ofsDelta *ofsDeltaObject) resolveDelta(packFile []byte, resolvedDeltasCache map[int]deltaCache) deltaCache {
	base := getOrReadObject(packFile, ofsDelta.baseOffset, resolvedDeltasCache)
	resolved := deltaCache{
		objectType:    base.objectType,
		resolvedDelta: applyDelta(base.resolvedDelta, ofsDelta.instructions),
	}
	utils.HashObject(resolved.resolvedDelta, resolved.objectType, utils.HashOptions{WriteHash: true})
	return resolved
}

// getOrReadObject returns the resolved content and type of the object at the
// given absolute offset in packFile, consulting the cache first and recursing
// through resolveDelta for ofs-delta bases.
func getOrReadObject(packFile []byte, offset int, resolvedDeltasCache map[int]deltaCache) deltaCache {
	if cached, ok := resolvedDeltasCache[offset]; ok {
		return cached
	}

	r := bytes.NewReader(packFile[offset:])
	objType, _ := readObjectHeader(r)

	if objType == objOfsDelta {
		negOff := readOfsOffset(r)
		instructions := decompressZlib(r)
		ofs := &ofsDeltaObject{
			instructions: instructions,
			baseOffset:   offset - negOff,
		}
		resolved := ofs.resolveDelta(packFile, resolvedDeltasCache)
		resolvedDeltasCache[offset] = resolved
		return resolved
	}

	entry := deltaCache{
		objectType:    typeName(objType),
		resolvedDelta: decompressZlib(r),
	}
	resolvedDeltasCache[offset] = entry
	return entry
}

// readOfsOffset decodes the ofs-delta variable-length base offset (git's
// "offset encoding": each continuation byte bumps the accumulator by 1 before
// shifting in the next 7 bits).
func readOfsOffset(r *bytes.Reader) int {
	b, err := r.ReadByte()
	if err != nil {
		log.Fatalf("read ofs-delta offset: %v", err)
	}
	off := int(b & 0x7f)
	for b&0x80 != 0 {
		b, err = r.ReadByte()
		if err != nil {
			log.Fatalf("read ofs-delta offset: %v", err)
		}
		off++
		off = (off << 7) | int(b&0x7f)
	}
	return off
}

// applyDelta applies git's delta instruction stream to base and returns the
// reconstructed object content. Instructions are either copy-from-base or
// insert-literal opcodes preceded by two size varints (base size, result size).
func applyDelta(base, instructions []byte) []byte {
	r := bytes.NewReader(instructions)
	readVarint := func() int {
		var val int
		var shift uint
		for {
			b, err := r.ReadByte()
			if err != nil {
				log.Fatalf("delta varint: %v", err)
			}
			val |= int(b&0x7f) << shift
			if b&0x80 == 0 {
				return val
			}
			shift += 7
		}
	}
	_ = readVarint() // base size
	_ = readVarint() // result size

	var out []byte
	for {
		op, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("delta op: %v", err)
		}
		if op&0x80 != 0 {
			var off, size int
			for i := range 4 {
				if op&(1<<i) != 0 {
					b, _ := r.ReadByte()
					off |= int(b) << (8 * i)
				}
			}
			for i := range 3 {
				if op&(1<<(4+i)) != 0 {
					b, _ := r.ReadByte()
					size |= int(b) << (8 * i)
				}
			}
			if size == 0 {
				size = 0x10000
			}
			out = append(out, base[off:off+size]...)
		} else if op != 0 {
			buf := make([]byte, op)
			if _, err := io.ReadFull(r, buf); err != nil {
				log.Fatalf("delta literal: %v", err)
			}
			out = append(out, buf...)
		} else {
			log.Fatalf("reserved delta opcode 0")
		}
	}
	return out
}
