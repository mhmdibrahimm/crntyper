package main

import (
	"path/filepath"
	"testing"

	hook "github.com/robotn/gohook"
)

func TestPlatformSpecForSupportedSystems(t *testing.T) {
	home := filepath.Join("Users", "test")

	tests := []struct {
		goos         string
		name         string
		backspaceRaw uint16
	}{
		{goos: "darwin", name: "macOS", backspaceRaw: rawcodeMacBackspace},
		{goos: "windows", name: "Windows", backspaceRaw: rawcodeWindowsBackspace},
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			spec := platformSpecFor(tt.goos, home)
			if !spec.Supported {
				t.Fatalf("expected %s to be supported", tt.goos)
			}
			if spec.Name != tt.name {
				t.Fatalf("expected name %q, got %q", tt.name, spec.Name)
			}
			if spec.BackspaceRaw != tt.backspaceRaw {
				t.Fatalf("expected raw code %d, got %d", tt.backspaceRaw, spec.BackspaceRaw)
			}
			if spec.CRNFile != filepath.Join(home, "Documents", "crns.txt") {
				t.Fatalf("unexpected CRN file: %s", spec.CRNFile)
			}
		})
	}
}

func TestPlatformSpecForUnsupportedSystem(t *testing.T) {
	tests := []string{"linux", "freebsd"}

	for _, goos := range tests {
		t.Run(goos, func(t *testing.T) {
			spec := platformSpecFor(goos, "")
			if spec.Supported {
				t.Fatalf("expected %s to be unsupported", goos)
			}
			if spec.CRNFile != "crns.txt" {
				t.Fatalf("expected local fallback CRN file, got %s", spec.CRNFile)
			}
		})
	}
}

func TestMatchesTrigger(t *testing.T) {
	spec := platformSpecFor("darwin", "/Users/test")

	tests := []struct {
		name  string
		event hook.Event
		want  bool
	}{
		{
			name:  "key down with backspace keycode",
			event: hook.Event{Kind: hook.KeyDown, Keycode: keycodeBackspace},
			want:  true,
		},
		{
			name:  "key down with platform raw code",
			event: hook.Event{Kind: hook.KeyDown, Rawcode: rawcodeMacBackspace},
			want:  true,
		},
		{
			name:  "key down with ascii backspace",
			event: hook.Event{Kind: hook.KeyDown, Keychar: '\b'},
			want:  true,
		},
		{
			name:  "key up ignored",
			event: hook.Event{Kind: hook.KeyUp, Keycode: keycodeBackspace},
			want:  false,
		},
		{
			name:  "other key ignored",
			event: hook.Event{Kind: hook.KeyDown, Keycode: 42, Rawcode: 42, Keychar: 'x'},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := spec.MatchesTrigger(tt.event); got != tt.want {
				t.Fatalf("MatchesTrigger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsupportedPlatformDoesNotMatchTrigger(t *testing.T) {
	spec := platformSpecFor("linux", "/home/test")

	if spec.MatchesTrigger(hook.Event{Kind: hook.KeyDown, Keycode: keycodeBackspace}) {
		t.Fatal("expected unsupported platform to ignore trigger")
	}
}
