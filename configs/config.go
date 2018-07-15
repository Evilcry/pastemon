package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	"github.com/asggo/store"
)

type Keyword struct {
	Keyword string
	Prefix  string
}

type Regex struct {
	Regex    string
	Compiled *regexp.Regexp
	Prefix   string
	Match    string
}

type Config struct {
	Keys     map[string]time.Time
	Ds       *store.Store
	Keywords []*Keyword // A list of keywords to search for in the data.
	Regexes  []*Regex   // A list of regular expressions to test against data.
	Buckets  []string   `json:"buckets"`       // List of buckets we need to create.
	DbFile   string     `json:"database_file"` // File to use for the Store database.
	MaxSize  int        `json:"max_size"`      // Do not save files larger than this many bytes.
	MaxTime  int        `json:"max_time"`      // Max time, in seconds, to store previously downloaded keys.
	Sleep    int        // Time, in seconds, to wait between each run.
	Save     bool
}

// NewConfig ctor
// new configuration instance
func NewConfig() Config {
	var c Config

	data, err := ioutil.ReadFile("configs/config.json")
	if err != nil {
		log.Fatal("[-] Could not read config file.")
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		log.Fatal("[-] Could not parse config file.")
	}

	c.Keys = make(map[string]time.Time)

	// Compile our regular expressions
	for i, _ := range c.Regexes {
		r := c.Regexes[i]
		r.Compiled = regexp.MustCompile(r.Regex)
	}

	return c
}
