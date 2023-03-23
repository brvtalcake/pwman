package compression

import (
	"io"

	bz "github.com/dsnet/compress/bzip2"
)

func GetEncryptedContent(r *bz.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func CompressEncryptedContent(w *bz.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}
