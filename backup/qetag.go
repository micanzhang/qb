package backup

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
)

func QEtag(filepath string) (etag string, err error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return
	}

	f, err := os.Open(filepath)
	if err != nil {
		return
	}

	size := stat.Size()
	blockCnt := blockCount(size)
	sha1Buf := make([]byte, 0, 21)

	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = calSha1(sha1Buf, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(f, BLOCK_SIZE)
			sha1BlockBuf, err = calSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = calSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return

}

const (
	BLOCK_BITS = 22 // Indicate that the blocksize is 4M
	BLOCK_SIZE = 1 << BLOCK_BITS
)

func blockCount(fsize int64) int {
	return int((fsize + (BLOCK_SIZE - 1)) >> BLOCK_BITS)
}

func calSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}
