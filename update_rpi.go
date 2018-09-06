package main

import (
	"github.com/kardianos/osext"
	"io"
	"net/http"
	"os"
	"syscall"
)

func updateStage1(url string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	// TODO: check directory writable
	// TODO: check executability in directory (ie. if FS is mounted with noexec or SELinux is blocking something)
	// TODO: check if the executable itself is writable

	err = download(url, wd+string(os.PathSeparator)+"seismosensor-new")
	if err != nil {
		return err
	}

	procname, err := osext.Executable()
	if err != nil {
		return err
	}

	args := []string{"seismosensor-new", "-stage2update", procname}
	env := os.Environ()
	return syscall.Exec(wd+string(os.PathSeparator)+"seismosensor-new", args, env)
}

func updateStage2(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	in, err := os.Open(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	in.Close()
	out.Close()

	args := []string{filename}
	env := os.Environ()
	return syscall.Exec(filename, args, env)
}

func download(url string, saveTo string) error {
	out, err := os.Create(saveTo)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()
	out.Close()
	return nil
}
