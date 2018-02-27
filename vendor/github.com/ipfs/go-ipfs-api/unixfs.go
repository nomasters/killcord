package shell

import (
	"encoding/json"
	"fmt"
)

type UnixLsObject struct {
	Hash  string
	Size  uint64
	Type  string
	Links []*UnixLsLink
}

type UnixLsLink struct {
	Name string
	Hash string
	Size uint64
	Type string
}

type lsOutput struct {
	Objects map[string]*UnixLsObject
}

// FileList entries at the given path using the UnixFS commands
func (s *Shell) FileList(path string) (*UnixLsObject, error) {
	resp, err := s.newRequest("file/ls", path).Send(s.httpcli)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.Error != nil {
		return nil, resp.Error
	}

	var out lsOutput
	err = json.NewDecoder(resp.Output).Decode(&out)
	if err != nil {
		return nil, err
	}

	for _, object := range out.Objects {
		return object, nil
	}

	return nil, fmt.Errorf("no object in results")
}
