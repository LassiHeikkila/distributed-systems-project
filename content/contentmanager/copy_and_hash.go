package main

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"hash"
	"io"
)

func CopyAndHash(dst io.Writer, src io.Reader, algo crypto.Hash) (int, []byte, error) {
	var h hash.Hash

	switch algo {
	case crypto.MD5:
		h = md5.New()
	case crypto.SHA1:
		h = sha1.New()
	case crypto.SHA256:
		h = sha256.New()
	default:
		return 0, nil, errors.New("unsupported hash")
	}

	n, err := CopyTwice(dst, h, src)
	if err != nil {
		return n, nil, err
	}

	s := h.Sum(nil)
	return n, s, nil
}

const copyBufferSize = 1024

func CopyTwice(dst1, dst2 io.Writer, src io.Reader) (int, error) {
	buf := make([]byte, copyBufferSize)
	totalN := 0

	for {
		rn, rErr := src.Read(buf)
		n1, err1 := dst1.Write(buf[:rn])
		n2, err2 := dst2.Write(buf[:rn])

		if errors.Is(rErr, io.EOF) {
			break
		} else if err1 != nil {
			return totalN, err1
		} else if err2 != nil {
			return totalN, err2
		}

		if n1 < n2 {
			return totalN + n1, errors.New("unbalanced write")
		}
		if n2 < n1 {
			return totalN + n2, errors.New("unbalanced write")
		}

		totalN += n1 // n1 == n2, doesn't matter which one is counted
	}

	return totalN, nil
}
