package eol

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read(file string) (map[string][]ReleaseInfo, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	var releasesByID map[string][]ReleaseInfo
	dec := json.NewDecoder(fh)
	if err := dec.Decode(&releasesByID); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return releasesByID, fh.Close()
}
