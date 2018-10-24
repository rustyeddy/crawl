package moni

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type ConfigLogger struct {
	Output string
	Level  string
	Format string
}

type Configuration struct {
	Addrport string
	Client   bool

	ConfigLogger
	ConfigFile string
	Cli        bool
	Daemon     bool
	Debug      bool
	Depth      int

	Profile bool
	Pubdir  string

	Serve    bool
	StoreDir string
	Tmpldir  string
	Wait     time.Duration
}

type Command struct {
	Command string
	Args    []string
}

var (
	config        *Configuration
	DefaultConfig = Configuration{
		Addrport: ":8888",
		Daemon:   false,
		Depth:    3,
		Profile:  false,
		Pubdir:   "docs",
		Serve:    false,
		StoreDir: "/srv/moni", // or "./.moni"
		Tmpldir:  "app/tmpl",
	}
)

func init() {
	config = &DefaultConfig
}

func GetConfig() *Configuration {
	return config
}

func SetConfig(c *Configuration) {
	config = c // Should merge... probably
}

// SaveFile will write the configuration out to our Storage as JSON.  As of
// now we are using the Storage interface to save our configurations.  Hence,
// we will use a label name and not a filename, although the label name can
// look like a filename..  That is a label can be "config.json" but it can
// NOT have any leading path '/' characters
func (c *Configuration) SaveFile() {
	s := GetStorage()
	if s == nil {
		log.Errorln("Failed saving config ", c.ConfigFile)
		return
	}
	_, err := s.StoreObject("config", c)
	if err != nil {
		log.Errorln("Failed writing configuration", c.ConfigFile, err)
	}
}

// ReadFile fetches our configuration object from our storage container,
// if needed, the object will be converted from JSON to a Go object before
// being returned.
func (c *Configuration) ReadFile() error {
	s := GetStorage() // fatal if nul
	_, err := s.FetchObject("config", c)
	if err != nil {
		return errors.New("failed to read config from store")
	}
	return nil
}
