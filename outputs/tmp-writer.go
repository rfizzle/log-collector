// Modified version of lumberjack (https://github.com/natefinch/lumberjack)
package outputs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// ensure we always implement io.WriteCloser
var _ io.WriteCloser = (*TmpWriter)(nil)

type TmpWriter struct {
	size         int64
	file         *os.File
	fileOpen     bool
	mu           sync.Mutex
	previousFile *os.File
	WriteCount   int

	millCh    chan bool
	startMill sync.Once
}

var (
	// os_Stat exists so it can be mocked out by tests.
	osStat = os.Stat
)

// Write implements io.Writer.
func (l *TmpWriter) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.write(p)
}

func (l *TmpWriter) WriteString(s string) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.write([]byte(s))
}

func (l *TmpWriter) write(p []byte) (n int, err error) {
	if l.file == nil {
		if err = l.openExistingOrNew(); err != nil {
			return 0, err
		}
	}

	if len(p) == 0 {
		return 0, nil
	}

	// Append newline to byte
	pStringWithNewline := string(p) + "\n"
	l.WriteCount += 1

	// Write new line
	n, err = l.file.Write([]byte(pStringWithNewline))
	l.size += int64(n)

	return n, err
}

// Close implements io.Closer, and closes the current logfile.
func (l *TmpWriter) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

// close closes the file if it is open.
func (l *TmpWriter) close() error {
	if l.file == nil {
		l.fileOpen = false
		return nil
	}
	err := l.file.Close()
	l.previousFile = l.file
	l.file = nil
	l.fileOpen = false
	return err
}

// Rotate causes TmpWriter to close the existing log file and immediately create a
// new one.
func (l *TmpWriter) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

// Returns current file size
func (l *TmpWriter) Size() int64 {
	return l.size
}

// rotate closes the current file and opens a new file
func (l *TmpWriter) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.openNew(); err != nil {
		return err
	}
	return nil
}

// openNew opens a new log file for writing. This methods assumes the last file has already been closed.
func (l *TmpWriter) openNew() error {
	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	f, err := ioutil.TempFile(os.TempDir(), randomStringWithLength(32))
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.file = f
	l.size = 0
	l.fileOpen = true
	l.WriteCount = 0
	return nil
}

// openExistingOrNew opens the logfile if it exists.  If there is no such file, a new file is created.
func (l *TmpWriter) openExistingOrNew() error {
	filename := l.filename()
	info, err := osStat(filename)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
		return l.openNew()
	}
	l.file = file
	l.size = info.Size()
	l.fileOpen = true
	l.WriteCount = 0
	return nil
}

// filename generates the name of the logfile from the current time.
func (l *TmpWriter) filename() string {
	file, _ := ioutil.TempFile(os.TempDir(), randomStringWithLength(32))
	name := file.Name()
	_ = file.Close()
	return name
}

// Return current log file pointer
func (l *TmpWriter) CurrentFile() *os.File {
	return l.file
}

// Return previous log file pointer
func (l *TmpWriter) PreviousFile() *os.File {
	return l.previousFile
}

// Delete current log file and update struct
func (l *TmpWriter) DeleteCurrentFile() (err error) {
	if l.file != nil && fileExists(l.file.Name()) {
		if err = os.Remove(l.file.Name()); err != nil {
			return err
		}
	}
	l.file = nil
	l.fileOpen = false
	return nil
}

// Delete previous log file and update struct
func (l *TmpWriter) DeletePreviousFile() (err error) {
	if l.previousFile != nil && fileExists(l.previousFile.Name()) {
		if err = os.Remove(l.previousFile.Name()); err != nil {
			return err
		}
	}
	l.previousFile = nil
	return nil
}

// Close open file and cleanup any current or previous log files
func (l *TmpWriter) Exit() (err error) {
	if l.fileOpen {
		if err = l.close(); err != nil {
			return err
		}
	}
	if l.file != nil {
		if err = l.DeleteCurrentFile(); err != nil {
			return err
		}
	}
	if l.previousFile != nil {
		if err = l.DeletePreviousFile(); err != nil {
			return err
		}
	}

	return nil
}
