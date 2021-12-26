// store - the file where records are stored.

package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	// define the encoding used to persist record sizes and index entries in.
	enc = binary.BigEndian
)

const (
	// define the number of bytes used to store the record's length.
	lenWidth = 8
)

// a simple wrapper around a file with two APIs to append and read bytes
// to/from the file.
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

// newStore creates a store for the given file. It gets the current size, in
// case we are recreating the store from a file with existing data (like after)
// a restart.
func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Append persists the given bytes to the store. Writes the length of the
// record, so we know how many bytes to return. Writes to the buffered writer
// to reduce system calls/improve performance. Returns the number of bytes
// written, and their position in the store; used by the segment to create an
// index entry for the record.
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

// Read returns the record stored at the given position. First it flushes the
// writer buffer to ensure all records are written to disk, then it gets the
// size (in bytes) for the whole record, finally reading and then returning
// the record.
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

// ReadAt reads len(p) bytes in p, beginning at offset. It implements
// io.ReaderAt on the 'store' type.
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, nil
	}
	return s.File.ReadAt(p, off)
}

// Close persists and buffered data before closing the file.
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
