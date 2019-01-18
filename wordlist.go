package main

import (
	"bufio"
	"os"
)

type WordlistInput struct {
	data     [][]byte
	position int
}

func NewWordlistInput(filePath string) (*WordlistInput, error) {
	var wl WordlistInput
	wl.position = -1
	valid, err := wl.validFile(filePath)
	if err != nil {
		return &wl, err
	}
	if valid {
		err = wl.readFile(filePath)
	}
	return &wl, err
}

//Next will increment the cursor position, and return a boolean telling if there's words left in the list
func (w *WordlistInput) Next() bool {
	w.position++
	if w.position >= len(w.data) {
		return false
	}
	return true
}

//Value returns the value from wordlist at current cursor position
func (w *WordlistInput) Value() []byte {
	return w.data[w.position]
}

//Total returns the size of wordlist
func (w *WordlistInput) Total() int {
	return len(w.data)
}

//validFile checks that the wordlist file exists and can be read
func (w *WordlistInput) validFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	f.Close()
	return true, nil
}

//readFile reads the file line by line to a byte slice
func (w *WordlistInput) readFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var data [][]byte
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		data = append(data, []byte(reader.Text()))
	}
	w.data = data
	return reader.Err()
}
