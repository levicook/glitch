package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
)

const (
	watchFlags = fsnotify.FSN_CREATE |
		fsnotify.FSN_DELETE |
		fsnotify.FSN_MODIFY
)

var (
	watcher  *fsnotify.Watcher
	watched  = make(map[string]struct{})
	exitCode = make(chan int)
	rootPath string

	emptyStruct = struct{}{}
	hasSuffix   = strings.HasSuffix
	contains    = strings.Contains
	sprintf     = fmt.Sprintf

	buildQueued = false
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func clearScrollBuffer() {
	print("\033c")
}

func build() {

	{
		log.Println("glitch: building")
		cmd := exec.Command("go", "build")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = rootPath
		if err := cmd.Run(); err != nil {
			return
		}
	}

	{
		log.Println("glitch: build OK - vetting")
		cmd := exec.Command("go", "vet", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = rootPath
		if err := cmd.Run(); err != nil {
			return
		}
	}

	{
		log.Println("glitch: vet OK - testing")
		cmd := exec.Command("go", "test", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = rootPath
		if err := cmd.Run(); err != nil {
			return
		}
	}

	{
		log.Println("glitch: test OK - installing")
		cmd := exec.Command("go", "install")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = rootPath
		if err := cmd.Run(); err != nil {
			return
		}
	}

	log.Println("glitch: install OK - waiting for next build event")
}

func maybeQueueBuild(path string) {
	buildQueued = hasSuffix(path, ".go")
}

func handleCreate(path string) {
	watch(path)
	maybeQueueBuild(path)
}

func handleDelete(path string) {
	if _, watching := watched[path]; watching {
		_ = watcher.RemoveWatch(path)
		delete(watched, path)
	}
	maybeQueueBuild(path)
}

func handleModify(path string) {
	maybeQueueBuild(path)
}

func handleEvent(ev *fsnotify.FileEvent) {
	if len(ev.Name) > 0 {
		switch {
		case ev.IsCreate():
			handleCreate(ev.Name)
		case ev.IsDelete():
			handleDelete(ev.Name)
		case ev.IsModify():
			handleModify(ev.Name)
		}
	}
}

var (
	gitSuffix   = sprintf("%v.git", string(filepath.Separator))
	gitContains = sprintf("%v%v", gitSuffix, string(filepath.Separator))
)

func watch(dir string) {

	if _, watching := watched[dir]; watching {
		return
	}

	walker := func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if hasSuffix(path, gitSuffix) || contains(path, gitContains) {
			return nil
		}

		if fileInfo.IsDir() {
			if err = watcher.WatchFlags(path, watchFlags); err == nil {
				watched[path] = emptyStruct
			}
		}

		return err
	}

	_ = filepath.Walk(dir, walker)
}

func periodicallyLogWatchedPaths() {
	logWatchedPaths := func() {
		log.Printf("glitch: watching: %v paths", len(watched))
		for path, _ := range watched {
			log.Println("glitch: watching:", path)
		}
	}

	logWatchedPaths()
	for _ = range time.Tick(5 * time.Second) {
		logWatchedPaths()
	}
}

func periodicallyLogBuildStatus() {
	logBuildStatus := func() {
		log.Println("glitch: buildQueued:", buildQueued)
	}

	logBuildStatus()
	for _ = range time.Tick(300 * time.Millisecond) {
		logBuildStatus()
	}
}

func runEventLoop() {
	for {
		select {
		case ev := <-watcher.Event:
			handleEvent(ev)
		case err := <-watcher.Error:
			panicIf(err)
		}
	}
}

func runBuildLoop() {
	consumeBuildQueue := func() {
		if buildQueued {
			buildQueued = false
			clearScrollBuffer()
			build()
		}
	}

	for _ = range time.Tick(1 * time.Second) {
		go consumeBuildQueue()
	}
}

func main() {
	wd, err := os.Getwd()
	panicIf(err)
	rootPath = wd

	w, err := fsnotify.NewWatcher()
	panicIf(err)
	watcher = w
	defer watcher.Close()

	//go periodicallyLogWatchedPaths()
	//go periodicallyLogBuildStatus()
	go runEventLoop()
	go runBuildLoop()

	watch(rootPath)
	os.Exit(<-exitCode)
}
