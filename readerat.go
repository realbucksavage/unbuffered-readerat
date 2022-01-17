package unbufra

import (
	"errors"
	"io"
	"io/ioutil"
)

var (
	ErrInvalidOffset = errors.New("invalid offset")
)

type unbufferedReaderAt struct {
	reader io.Reader
	read   int64
}

func (u *unbufferedReaderAt) ReadAt(p []byte, offset int64) (int, error) {
	if offset < u.read {
		return 0, ErrInvalidOffset
	}

	diff := offset - u.read
	written, err := io.CopyN(ioutil.Discard, u.reader, diff)
	u.read += written

	if err != nil {
		return 0, err
	}

	n, err := u.reader.Read(p)
	u.read += int64(n)

	return n, err
}

func NewUnbufferedReader(reader io.Reader) io.ReaderAt {
	return &unbufferedReaderAt{reader: reader}
}
