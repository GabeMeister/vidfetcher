package util

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// Pause pauses program execution for user
func Pause() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Here prints "Here."
func Here() {
	log.Println("Here.")
}

// End prints the word "End." and ends program execution
func End() {
	log.Fatalln("End.")
}

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// WriteLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// InitLoggerAndLogFile creates a log file at filePath location and returns
// the file handle.
func InitLoggerAndLogFile(filePath string) *os.File {
	var file *os.File
	var err error

	if _, err := os.Stat(filePath); err != nil {
		file, err = os.Create(filePath)
	} else {
		file, err = os.OpenFile(filePath, os.O_APPEND, 0644)
	}

	if err != nil {
		log.Fatal("unable to open log file", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	return file
}
