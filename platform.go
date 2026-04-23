package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	hook "github.com/robotn/gohook"
)

const (
	keycodeBackspace        uint16 = 0x000E
	rawcodeMacBackspace     uint16 = 0x0033
	rawcodeWindowsBackspace uint16 = 0x0008
)

type platformSpec struct {
	GOOS         string
	Name         string
	Supported    bool
	CRNFile      string
	TriggerName  string
	BackspaceRaw uint16
}

func currentPlatformSpec() platformSpec {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("⚠️  could not get home dir, using local cwd:", err)
	}

	return platformSpecFor(runtime.GOOS, home)
}

func platformSpecFor(goos, home string) platformSpec {
	crnFile := "crns.txt"
	if home != "" {
		crnFile = filepath.Join(home, "Documents", "crns.txt")
	}

	switch goos {
	case "darwin":
		return platformSpec{
			GOOS:         goos,
			Name:         "macOS",
			Supported:    true,
			CRNFile:      crnFile,
			TriggerName:  "Backspace",
			BackspaceRaw: rawcodeMacBackspace,
		}
	case "windows":
		return platformSpec{
			GOOS:         goos,
			Name:         "Windows",
			Supported:    true,
			CRNFile:      crnFile,
			TriggerName:  "Backspace",
			BackspaceRaw: rawcodeWindowsBackspace,
		}
	default:
		return platformSpec{
			GOOS:        goos,
			Name:        goos,
			Supported:   false,
			CRNFile:     crnFile,
			TriggerName: "Backspace",
		}
	}
}

func (s platformSpec) ReadyMessage() string {
	if s.Supported {
		return fmt.Sprintf("%s detected. Press Load to import %s", s.Name, s.CRNFile)
	}

	return fmt.Sprintf("%s is not supported. CRN Typer supports macOS and Windows.", s.Name)
}

func (s platformSpec) MatchesTrigger(ev hook.Event) bool {
	if !s.Supported {
		return false
	}

	if ev.Kind != hook.KeyDown {
		return false
	}

	// gohook reports Backspace differently across OS backends, so accept the
	// virtual keycode plus the platform raw code and common control chars.
	return ev.Keycode == keycodeBackspace ||
		(s.BackspaceRaw != 0 && ev.Rawcode == s.BackspaceRaw) ||
		ev.Keychar == '\b' ||
		ev.Keychar == 127
}
