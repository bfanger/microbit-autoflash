package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	defer fmt.Println("Bye")
	if err := microbitConnect(); err != nil {
		log.Fatalf("could not connect to micro:bit: %v", err)
	}
	fmt.Println("micro:bit connected")
	for {
		file, err := downloadAdded()
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
		if err := microbitConnect(); err != nil {
			log.Fatalf("could not connect to micro:bit: %v", err)
		}
		name := path.Base(file)
		name = name[9 : len(name)-4]
		fmt.Printf("Flashing: %s\n", name)
		source, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("could not open file: %v", err)
		}

		if err := ioutil.WriteFile("/Volumes/MICROBIT/program.hex", source, 0775); err != nil {
			log.Fatalf("could not write file: %v", err)
		}
		cmd := exec.Command("/usr/sbin/diskutil", "unmount", "/Volumes/MICROBIT")
		if err := cmd.Run(); err != nil {
			fmt.Printf("could unmount write file: %v\n", err)
		}
		if err := os.Remove(file); err != nil {
			log.Fatalf("could not remove file: %v", err)
		}
	}
}

// microbitConnect waits for the microbit to connect
func microbitConnect() error {
	for {
		connected, err := isMicrobitConnected()
		if err != nil {
			return err
		}
		if connected {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func isMicrobitConnected() (bool, error) {
	if _, err := os.Stat("/Volumes/MICROBIT"); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func downloadAdded() (string, error) {
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()

	type downloadAddedEvent struct {
		name string
		err  error
	}

	result := make(chan downloadAddedEvent)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Create && strings.HasSuffix(event.Name, ".hex") && strings.HasPrefix(path.Base(event.Name), "microbit-") {
					result <- downloadAddedEvent{name: event.Name, err: nil}
					return
				}

			case err := <-watcher.Errors:
				result <- downloadAddedEvent{err: err}
				return
			}
		}
	}()

	if watcher.Add(path.Join(os.Getenv("HOME"), "Downloads")); err != nil {
		return "", fmt.Errorf("failed to add watcher: %v", err)
	}
	unpack := <-result
	close(result)

	return unpack.name, unpack.err
}
