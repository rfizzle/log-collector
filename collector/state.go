package collector

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

// Create a new state
func newState(path string) (p *State, err error) {
	p = &State{
		Path: path,
		Data: data{
			LastPollTimestamp: time.Now().Add(-1 * time.Hour * 24 * 1).Format(time.RFC3339),
		},
	}
	if fileExists(path) {
		err := p.Load()
		if err != nil {
			return nil, err
		}
	}
	return p, err
}

// Save state
func (p *State) Save() {
	// Marshal to JSON
	file, _ := json.MarshalIndent(p.Data, "", " ")

	// Write to file
	err := ioutil.WriteFile(p.Path, file, 0644)

	if err != nil {
		log.Errorf("error writing state file: %v", err)
	}
}

// Restore state
func (p *State) Load() error {
	// Open our jsonFile
	jsonFile, err := os.Open(p.Path)

	// if os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read the opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Initialize our state struct
	var data data

	// unmarshal our byteArray which contains our
	// jsonFile's content into 'state' which we defined above
	err = json.Unmarshal(byteValue, &data)

	// if json.Unmarshal returns an error then handle it
	if err != nil {
		return err
	}

	// Set data
	p.Data = data

	return nil
}

type State struct {
	Path string
	Data data
}

type data struct {
	LastPollTimestamp string `json:"last_poll_timestamp"`
}

// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
