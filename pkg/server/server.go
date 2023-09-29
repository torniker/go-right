package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/slayer/autorestart"
	"github.com/torniker/go-right/env"
	"gopkg.in/fsnotify.v1"
)

var spin = spinner.New(spinner.CharSets[6], 100*time.Millisecond)

// Start echo web server
func Start(e *echo.Echo, shouldWatch bool) {
	if shouldWatch {
		binFilename := filepath.Base(env.Get(env.MainGoPathKey))
		// suffix := filepath.Base(dir)
		autorestart.WatchFilename = fmt.Sprintf("%s/build/%s", env.Get(env.BasePathKey), binFilename)
		restart := autorestart.GetNotifier()
		go func() {
			<-restart
			fmt.Println("\033[36mI will restart shortly\033[00m")
			spin.Stop()
		}()
		if fileExists(autorestart.WatchFilename) {
			autorestart.StartWatcher()
			go watch()
		} else {
			log.Warn("watcher did not start")
		}
		if err := e.Start(env.Get(env.PortKey)); err != nil {
			log.Infof("shutting down the server, error: %s", err)
		}
	} else {
		if err := e.Start(env.Get(env.PortKey)); err != nil {
			log.Infof("shutting down the server, error: %s", err)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Watch file changes and restart server
func watch() {
	ctx := context.Background()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if filepath.Ext(event.Name) == ".go" {
						spin.Color("blue")
						spin.Suffix = "\033[33m Building... \033[00m"
						spin.Start()
						rebuild()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(ctx, err)
			}
		}
	}()

	err = filepath.Walk(env.Get(env.BasePathKey),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(path, "/node_modules/") {
				return nil
			}
			if !info.IsDir() && filepath.Ext(info.Name()) == ".go" {
				err = watcher.Add(filepath.Dir(path))
				if err != nil {
					log.Error(ctx, err)
				}
			}
			return nil
		},
	)
	if err != nil {
		log.Error(ctx, err)
	}
	<-done
}

func rebuild() {
	cmd := exec.Command("go", "build", "-o", fmt.Sprintf("%s/build", env.Get(env.BasePathKey)), fmt.Sprintf("%s/%s", env.Get(env.BasePathKey), env.Get(env.MainGoPathKey)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
