package main

import (
	"io"
	"net/http"
	"os"
	"syscall"
)

func update(url string) error {
	// TODO: check directory writable (as we need to delete and re-create the executable)
	// TODO: write the new executable in the same directory as the current one

	err := download(url, "/tmp/seismosensor-new")
	if err != nil {
		return err
	}

	procname, err := os.Executable()
	if err != nil {
		return err
	}

	err = os.Remove(procname)
	if err != nil {
		return err
	}

	// Warning, here the executable is vanished, we need to place the new one here immediately

	err = os.Rename("/tmp/seismosensor-new", procname)
	if err != nil {
		// OOOPS! This is a very big issue: the old executable is gone
		return err
	}

	return reboot()
}

func reboot() error {
	procname, err := os.Executable()
	if err != nil {
		return err
	}
	args := []string{procname}
	env := os.Environ()
	return syscall.Exec(procname, args, env) // #nosec G204
}

func download(url string, saveTo string) error {
	out, err := os.Create(saveTo)
	if err != nil {
		return err
	}

	resp, err := http.Get(url) // #nosec G107
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	_ = resp.Body.Close()
	_ = out.Close()
	return nil
}
