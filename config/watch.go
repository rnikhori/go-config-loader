package config

import (
	"fmt"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

// debounceInterval coalesces bursts of filesystem events (editors
// frequently emit several events for a single logical save, e.g. a
// temp-file write followed by a rename) into a single reload.
const debounceInterval = 100 * time.Millisecond

// Watch loads path into out, then watches it for changes using fsnotify.
// On each change it reloads path into out — reapplying env and, if args
// is non-nil, flag overlays, exactly as LoadWithFlags would — and calls
// onChange with the result (nil on a successful reload, non-nil on a
// parse/validation error or a watcher failure).
//
// Watch returns a stop function that stops watching and releases the
// underlying OS resources. After stop is called, onChange will not be
// invoked again, though a reload already in progress may still finish.
func Watch(path string, out interface{}, args []string, onChange func(error)) (stop func() error, err error) {
	if err := LoadWithFlags(path, out, args); err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("config: create watcher: %w", err)
	}

	dir := filepath.Dir(path)
	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("config: watch %s: %w", dir, err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		watcher.Close()
		return nil, fmt.Errorf("config: resolve path %s: %w", path, err)
	}

	var stopped atomic.Bool
	done := make(chan struct{})

	go runWatchLoop(watcher, absPath, path, args, out, onChange, &stopped, done)

	stop = func() error {
		stopped.Store(true)
		close(done)
		return watcher.Close()
	}

	return stop, nil
}

func runWatchLoop(
	watcher *fsnotify.Watcher,
	absPath, path string,
	args []string,
	out interface{},
	onChange func(error),
	stopped *atomic.Bool,
	done chan struct{},
) {
	var timer *time.Timer

	reload := func() {
		if stopped.Load() {
			return
		}
		onChange(LoadWithFlags(path, out, args))
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			eventAbsPath, err := filepath.Abs(event.Name)
			if err != nil || eventAbsPath != absPath {
				continue
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounceInterval, reload)

		case watchErr, ok := <-watcher.Errors:
			if !ok {
				return
			}
			if !stopped.Load() {
				onChange(fmt.Errorf("config: watcher error: %w", watchErr))
			}

		case <-done:
			if timer != nil {
				timer.Stop()
			}
			return
		}
	}
}
