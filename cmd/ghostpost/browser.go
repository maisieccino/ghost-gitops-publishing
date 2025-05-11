// cmd/ghostpost/browser.go

package main

import (
	"os/exec"
	"runtime"
)

func launchBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default: // linux, *bsd
		return exec.Command("xdg-open", url).Start()
	}
}
