package gron

import (
	"encoding/json"
)

// Parse parses json config into Configuration object
func Parse(jsonConfig []byte) ([]*Task, error) {
	var c []*Task
	err := json.Unmarshal(jsonConfig, &c)
	if err != nil {
		return nil, err
	}

	for key, _ := range c {
		c[key].Processes = make([]*Process, 0)
	}

	return c, nil
}
