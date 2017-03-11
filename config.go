package gron

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

// Conf is main interface for configuration
type Conf interface {
	Parse(conf io.Reader) ([]*Job, error)
}

// JSONConf config is type to allow JSON configuration
type JSONConf struct {
}

// Parse parses json config into slice of Tasks to run
func (jsonConf *JSONConf) Parse(conf io.Reader) ([]*Job, error) {
	buf, err := ioutil.ReadAll(conf)
	if err != nil {
		log.Fatal(err)
	}

	var c []*Job
	err = json.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
