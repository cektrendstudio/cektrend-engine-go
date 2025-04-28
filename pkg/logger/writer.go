package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utpath"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/uttime"
)

type logWriterObj struct {
	mu          sync.Mutex
	isWriting   bool
	isReady     bool
	filePath    string
	fileName    string
	fileStream  *os.File
	writeQueues []string

	// Path is a log file directory
	Path string
	// FileFormat is a log file name
	FileFormat string
	// Mode is a log writing mode
	Mode Mode
}

type logWriter interface {
	// IsWriting returns the logger writing feature status
	IsWriting() bool
	// StartWriting to enable the logger writing feature
	StartWriting()
	// StopWriting to disable the logger writing feature
	StopWriting()
}

func (ox *logWriterObj) IsWriting() bool {
	return ox.isWriting
}

func (ox *logWriterObj) StopWriting() {
	ox.isWriting = false
}

func (ox *logWriterObj) StartWriting() {
	ox.isWriting = true
}

// private

func (ox *logWriterObj) boot() (err error) {
	if ox.IsWriting() {
		var (
			now = time.Now()

			format      = ox.FileFormat
			formatBinds = map[string]string{
				"d": now.Format("02"),
				"m": now.Format("01"),
				"y": now.Format("2006"),
				"h": now.Format("15"),
				"i": now.Format("04"),
				"s": now.Format("05"),
				"v": "",
			}
		)

		switch ox.Mode {
		case ModeDaily:
			formatBinds["v"] = now.Format("20060102")

		case ModeMonthly:
			formatBinds["v"] = now.Format("200601")

		case ModeYearly:
			formatBinds["v"] = now.Format("2006")

		case ModePermanent:
			formatBinds["v"] = ""
		}

		for k, v := range formatBinds {
			format = strings.ReplaceAll(format, "%"+k, v)
		}

		if ox.fileName != format {
			ox.fileName = format
			ox.filePath = filepath.Join(ox.Path, ox.fileName)

			if !utpath.IsExists(ox.Path) {
				err = os.MkdirAll(ox.Path, os.ModePerm)
				if err != nil {
					return
				}
			}

			err = ox.open()
			if err != nil {
				return
			}
		}
	}

	if !ox.isReady {
		ox.isReady = true

		go func() {
			for {
				time.Sleep(3 * time.Second)
				ox.flush()
			}
		}()
	}

	return
}

func (ox *logWriterObj) open() (err error) {
	if !ox.IsWriting() {
		return nil
	}

	ox.mu.Lock()
	defer ox.mu.Unlock()

	if ox.fileStream != nil {
		_ = ox.fileStream.Close()
	}

	ox.fileStream, err = os.OpenFile(ox.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	return
}

func (ox *logWriterObj) write(msg string) (err error) {
	if !ox.IsWriting() {
		return
	}

	if !ox.isReady {
		err = errors.New("Logger not yet ready")
		return
	}

	if msg == "" {
		return
	}

	ox.mu.Lock()
	defer ox.mu.Unlock()

	ox.writeQueues = append(ox.writeQueues, msg)

	return
}

func (ox *logWriterObj) flush() (err error) {
	if !ox.IsWriting() {
		return nil
	}

	err = ox.boot()
	if err != nil {
		return
	}

	ox.mu.Lock()
	defer ox.mu.Unlock()

	lists := ox.writeQueues
	ox.writeQueues = []string{}

	defer func() {
		if err != nil {
			ox.writeQueues = append(lists, ox.writeQueues...)
		}
	}()

	if len(lists) > 0 {
		for _, v := range lists {
			_, err = ox.fileStream.WriteString(fmt.Sprintf("%s\n", v))
			if err != nil {
				ox.printf("Failed to writing, details: %+v", err)

				errs := ox.open()
				if errs != nil {
					ox.printf("Failed to re-open file %s, details: %+v", ox.filePath, errs)
				}
				return err
			}
		}

		err = ox.fileStream.Sync()
		if err != nil {
			ox.printf("Failed to flushing stream, details: %+v", err)
			return err
		}
	}

	return err
}

func (ox *logWriterObj) printf(msg string, opts ...interface{}) {
	fmt.Printf("[%s] ERR: %s\n", uttime.Format(uttime.DefaultDateTimeFormat, time.Now()), fmt.Sprintf(msg, opts...))
}
