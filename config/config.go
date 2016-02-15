package config

import (
	"encoding/json"

	"github.com/melkorm/gron/runner"
)

// Parse parses json config into Configuration object
func Parse(jsonConfig []byte) ([]*runner.Task, error) {
	var c []*runner.Task
	err := json.Unmarshal(jsonConfig, &c)
	if err != nil {
		return nil, err
	}

	for key, _ := range c {
		c[key].Processes = make([]*runner.Process, 0)
	}

	return c, nil
}
