package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

var copyBuffer = make([]byte, 1024*1024*5)

func doCopy(src, dst string, move, overwrite bool) (err error) {
	info, _ := os.Stat(dst)
	if !overwrite && info != nil {
		return errors.New("le fichier '" + dst + "' existe déjà (" + strconv.Itoa(int(info.Size())) + " octets)")
	}

	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return
	}

	if move {
		if err = os.Rename(src, dst); err == nil {
			return
		}
		/*if le, ok := err.(*os.LinkError); !ok || (le.Err != syscall.Errno(17) && le.Err != syscall.EXDEV) { // ERROR_NOT_SAME_DEVICE
			return
		}*/
	}

	var hsrc, hdst *os.File
	defer func() {
		if hsrc != nil {
			hsrc.Close()
		}
		if hdst != nil {
			hdst.Close()
		}

		if err == nil && move {
			err = os.Remove(src)
		}
	}()

	hsrc, err = os.Open(src)
	if err != nil {
		return
	}

	info, err = hsrc.Stat()
	if err != nil {
		return
	}

	hdst, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return
	}

	if _, err = io.CopyBuffer(hdst, hsrc, copyBuffer); err != nil {
		return
	}

	err = hdst.Sync()
	return
}
