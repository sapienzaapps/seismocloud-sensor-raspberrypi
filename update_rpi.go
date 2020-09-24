package main

import (
	"context"
	"crypto/md5" // #nosec G501 used only for integrity check and esp8266 compatibility
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/config"
)

func checkUpdate() error {
	// TODO: check directory writable (as we need to delete and re-create the executable)
	// TODO: write the new executable in the same directory as the current one

	updateExists, err := download("/tmp/seismosensor-new")
	if err != nil {
		return err
	} else if !updateExists {
		return nil
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
	return syscall.Exec(procname, args, env) // #nosec G204 no risk of file inclusion
}

func download(saveTo string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build HTTP Request
	url := fmt.Sprintf("%s/firmware/raspi/%s/bin", config.FirmwareServer, cfg.GetDeviceID())
	log.Debugf("Checking for update with URL %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("x-version", config.Version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotModified {
		// No updates
		return false, nil
	} else if resp.StatusCode != 200 {
		// Server error
		return false, fmt.Errorf("invalid HTTP response status: %d", resp.StatusCode)
	}

	// Save new firmware
	out, err := os.Create(filepath.Clean(saveTo))
	if err != nil {
		return false, err
	}

	_, err = io.Copy(out, resp.Body)
	_ = out.Close()
	if err != nil {
		return false, err
	}

	// Checking MD5 sum
	fp, err := os.Open(filepath.Clean(saveTo))
	if err != nil {
		return false, err
	}

	chash := md5.New() // #nosec G401 used only for integrity check and esp8266 compatibility

	if _, err := io.Copy(chash, fp); err != nil {
		_ = fp.Close()
		return false, err
	}
	_ = fp.Close()

	localHash := hex.EncodeToString(chash.Sum(nil))
	if !strings.EqualFold(localHash, resp.Header.Get("x-md5")) {
		_ = os.Remove(saveTo)
		return false, errors.New("corrupted file")
	}

	// Save and MD5 check OK
	return true, nil
}
