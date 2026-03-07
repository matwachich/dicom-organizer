package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

var copyBuffer = make([]byte, 1024*1024*5)

func doCopy(src, dst string, move, overwrite bool) (err error) {
	stat, _ := os.Stat(dst)
	if !overwrite && stat != nil {
		return errors.New("le fichier destination existe déjà (" + strconv.Itoa(int(stat.Size())) + " octets)")
	}

	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return
	}

	if move {
		if overwrite {
			os.Remove(dst)
		}

		if err = os.Rename(src, dst); err == nil {
			return
		}
		if le, ok := err.(*os.LinkError); !ok || le.Err != syscall.EXDEV {
			return
		}
	}

	hsrc, err := os.Open(src)
	if err != nil {
		return
	}
	defer hsrc.Close()

	info, err := hsrc.Stat()
	if err != nil {
		return
	}

	hdst, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return
	}

	err = func() error {
		defer hdst.Close()
		if _, err := io.CopyBuffer(hdst, hsrc, copyBuffer); err != nil {
			return err
		}
		return hdst.Sync()
	}()
	if err != nil {
		return
	}

	if move {
		err = os.Remove(src)
	}
	return
}
