package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var platform = currentPlatformSpec()
var crnFile = platform.CRNFile

type CRNStore struct {
	mu   sync.RWMutex
	crns []string
}

func NewCRNStore() *CRNStore { return &CRNStore{} }

func (s *CRNStore) Load() error {
	f, err := os.Open(crnFile)
	if err != nil {
		return err
	}
	defer f.Close()

	var list []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if ln := sc.Text(); ln != "" {
			list = append(list, ln)
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	s.crns = list
	s.mu.Unlock()
	return nil
}

func (s *CRNStore) Save(list []string) error {
	f, err := os.Create(crnFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, ln := range list {
		if ln != "" {
			fmt.Fprintln(w, ln)
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}

	s.mu.Lock()
	s.crns = list
	s.mu.Unlock()
	return nil
}

func (s *CRNStore) Snapshot() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.crns))
	copy(out, s.crns)
	return out
}

func (s *CRNStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.crns)
}

func main() {
	runtime.LockOSThread()
	store := NewCRNStore()
	_ = store.Load()

	var (
		loaded   bool
		loadMu   sync.RWMutex
		crnCh    = make(chan string)
		unloadCh = make(chan struct{})
	)

	setLoaded := func(v bool) {
		loadMu.Lock()
		loaded = v
		loadMu.Unlock()
	}
	isLoaded := func() bool {
		loadMu.RLock()
		defer loadMu.RUnlock()
		return loaded
	}

	if platform.Supported {
		// robot goroutine
		go func() {
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			for crn := range crnCh {
				robotgo.TypeStr(crn)
				robotgo.KeyTap("tab")
			}
		}()

		go func() {
			events := hook.Start()
			defer hook.End()
			for ev := range events {
				if platform.MatchesTrigger(ev) && isLoaded() {
					for _, crn := range store.Snapshot() {
						crnCh <- crn
					}
					unloadCh <- struct{}{}
				}
			}
			close(crnCh)
		}()
	}

	a := app.New()
	w := a.NewWindow(appTitle())
	if platform.Supported {
		w.SetTitle(appTitle() + " - " + platform.Name)
	}
	w.Resize(fyne.NewSize(360, 300))

	status := widget.NewLabel("")
	status.Alignment = fyne.TextAlignCenter
	status.Wrapping = fyne.TextWrapWord

	entryBox := container.NewVBox()
	var entries []*widget.Entry

	autoSave := func() {
		var list []string
		for _, e := range entries {
			if t := e.Text; t != "" {
				list = append(list, t)
			}
		}
		if err := store.Save(list); err != nil {
			status.SetText("❌ Save error: " + err.Error())
		}
	}

	var addEntry func()
	addEntry = func() {
		e := widget.NewEntry()
		e.SetPlaceHolder("Enter CRN…")
		e.OnChanged = func(s string) {
			if s != "" && entries[len(entries)-1] == e {
				addEntry()
			}
			if isLoaded() {
				autoSave()
			}
		}
		entries = append(entries, e)
		entryBox.Add(e)
		entryBox.Refresh()
	}

	var loadBtn *widget.Button
	loadBtn = widget.NewButton("Load", func() {
		if !platform.Supported {
			status.SetText(platform.ReadyMessage())
			return
		}

		// gather
		var list []string
		for _, e := range entries {
			if t := e.Text; t != "" {
				list = append(list, t)
			}
		}
		// save & load
		if err := store.Save(list); err != nil {
			status.SetText("❌ Save error: " + err.Error())
			return
		}
		setLoaded(true)
		loadBtn.SetText("Unload")
		status.SetText(fmt.Sprintf("✅ Loaded %d CRNs. Press %s.", len(list), platform.TriggerName))
	})
	if !platform.Supported {
		loadBtn.Disable()
	}

	go func() {
		for range unloadCh {
			setLoaded(false)
			loadBtn.SetText("Load")
			status.SetText(platform.ReadyMessage())
		}
	}()

	if list := store.Snapshot(); len(list) > 0 {
		for _, crn := range list {
			e := widget.NewEntry()
			e.SetText(crn)
			e.OnChanged = func(s string) {
				if s != "" && entries[len(entries)-1] == e {
					addEntry()
				}
				if isLoaded() {
					autoSave()
				}
			}
			entries = append(entries, e)
			entryBox.Add(e)
		}
	}
	addEntry()

	status.SetText(platform.ReadyMessage())

	w.SetContent(container.NewVBox(
		status,
		loadBtn,
		entryBox,
	))

	w.ShowAndRun()
}
