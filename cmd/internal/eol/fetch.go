package eol

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var ErrNotFound = fmt.Errorf("not found")

var translateIDs = map[string]string{
	"amzn":          "amazon-linux",
	"ol":            "oracle-linux",
	"pop":           "pop-os",
	"rocky":         "rocky-linux",
	"opensuse-leap": "opensuse",
	"sled":          "sles",
}

func Fetch(dest string, ids []string) error {

	fmt.Println("Fetching EOL data for operating systems")

	var releasesByID = make(map[string][]ReleaseInfo)
	for _, osReleaseID := range ids {
		releases, err := fetchAndParse(osReleaseID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				continue
			}
			return fmt.Errorf("failed to fetch data for %s: %w", osReleaseID, err)
		}

		fmt.Printf("   ...found %d releases for %s\n", len(releases), osReleaseID)
		releasesByID[osReleaseID] = releases
	}

	fh, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}

	enc := json.NewEncoder(fh)
	enc.SetIndent("", "  ")

	if err := enc.Encode(releasesByID); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return fh.Close()
}

func fetchAndParse(id string) ([]ReleaseInfo, error) {
	osReleaseID := id
	if translated, ok := translateIDs[id]; ok {
		id = translated
	}

	url := "https://endoflife.date/api/" + id + ".json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var releases []ReleaseInfo
	if err := json.Unmarshal(data, &releases); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	for i := range releases {
		releases[i].OSReleaseID = osReleaseID
		releases[i].APIID = id
	}

	return releases, nil
}
