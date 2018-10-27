package moni

import (
	"errors"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	config *Configuration
)

type ConfigLogger struct {
	LevelString  string
	FormatString string
	Logfile      string

	log.Level
	log.Formatter
	Out io.Writer
}

type Configuration struct {
	Addrport string
	Basedir  string
	Client   bool

	ConfigLogger

	ConfigFile string
	Cli        bool
	Daemon     bool
	Debug      bool
	Depth      int

	NoServe bool

	Profile bool
	Pubdir  string

	Storedir string
	Tmpldir  string
	Wait     time.Duration

	*Logerr
}

type Command struct {
	Command string
	Args    []string
}

var (
	DefaultConfig = Configuration{
		Addrport: ":8888",
		Daemon:   false,
		Depth:    3,
		Profile:  false,
		Pubdir:   "docs",
		NoServe:  false,
		Storedir: "/srv/moni", // or "./.moni"
		Tmpldir:  "tmpl",
	}
)

// SetConfig sets and reconfigures the application
func SetConfig(cfg *Configuration) {
	// Set the global config
	config = cfg

	if config.Logerr == nil {
		config.Logerr = NewLogerr("config")
	}
}

// StoreObject will write the configuration out to our Storage as
// configurations.  Hence, we will use a label name and not a
// filename, although the label name can look like a filename..  That
// is a label can be "config.json" but it can NOT have any leading
// path '/' characters
func (c *Configuration) StoreConfig() {
	_, err := storage.StoreObject("config", c)
	if err != nil {
		log.Errorln("Failed writing configuration", c.ConfigFile, err)
	}
}

// ReadFile fetches our configuration object from our storage
// container, if needed, the object will be converted from JSON to a
// Go object before being returned.
func (c *Configuration) FetchConfig() error {
	if _, err := storage.FetchObject("config", c); err != nil {
		return errors.New("failed to read config from store")
	}
	return nil
}
