package compression

import (
	"io"

	bz "github.com/dsnet/compress/bzip2"
)

func Bzip2Compress(data []byte) ([]byte, error) {
	io_writer := io.Writer()
	bz.NewWriter(nil)
}
