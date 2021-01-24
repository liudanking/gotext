package main

import (
	"errors"
	"flag"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/liudanking/goutil/strutil"

	"github.com/fsnotify/fsnotify"
	"github.com/liudanking/gotext/cfg"
	"github.com/liudanking/gotext/ocr"

	log "github.com/liudanking/goutil/logutil"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "~/.gotext_config", "config file")
	flag.Parse()

	config, err := cfg.LoadConfig(configFile)
	if err != nil {
		log.Error("load config from file:[%s] failed:%v", configFile, err)
		return
	}

	log.Notice("load config file OK, config:%+v", config)

	log.Info("start to serve...")
	if err := watchAndServe(config.ServeDir); err != nil {
		log.Error("watchAndServe failed:%v", err)
	}
}

func watchAndServe(serveDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("create watcher failed:%v", err)
		return err
	}
	// defer watcher.Close()

	err = watcher.Add(serveDir)
	if err != nil {
		log.Error("watch serve_dir:[%s] failed:%v", serveDir, err)
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Error("get event from watcher not OK")
				return errors.New("get watcher event not ok")
			}
			log.Debug("get event:%s", event.String())
			if event.Op&fsnotify.Create == fsnotify.Create {
				fn := event.Name //filepath.Join(serveDir, event.Name)
				log.Notice("new file detected:%s", fn)
				if !canOcr(event.Name) {
					log.Info("file:%s can't be ocr, skip", fn)
					continue
				}
				// do ocr work
				GetOCRText(fn)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Error("get error from watcher not OK")
				return errors.New("get watcher error not ok")
			}
			log.Error("get error:%v", err)
		}
	}

	return nil
}

func canOcr(fn string) bool {
	canOcrExtNames := []string{
		".jpg",
		".jpeg",
		".png",
		".webp",
		".bmp",
	}
	extName := strings.ToLower(filepath.Ext(fn))

	return strutil.StringIn(canOcrExtNames, extName)
}

func copyToClipboard(b []byte) error {
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := in.Write(b); err != nil {
		return err
	}

	if err := in.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}

func GetOCRText(fn string) (string, error) {
	content, err := ocr.GetOCRTextWithBaiduAI(fn)
	if err != nil {
		log.Error("GetOCRTextWithBaiduAI failed:%v", err)
		return "", err
	}
	log.Notice("[fn:%s][conent:%s]", fn, content)
	copyToClipboard([]byte(content))
	log.Info("copied to clipboard!")
	return content, nil
}
