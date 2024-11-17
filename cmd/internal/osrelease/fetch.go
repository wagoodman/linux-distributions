package osrelease

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mitchellh/mapstructure"
	"github.com/scylladb/go-set/strset"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Fetch(dest string) ([]string, error) {
	fmt.Println("Starting generation for OS release catalog...")

	tmpDir, err := os.MkdirTemp("", "os-release-repo")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	repoURL := "https://github.com/which-distro/os-release.git"
	fmt.Println("Cloning repository:", repoURL)
	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to clone repo: %w", err)
	}

	// walk the directory and process os-release files
	var osReleaseInfos []Info
	ids := strset.New()
	fmt.Println("Parsing os-release files:")
	skipDirs := strset.New(".git")
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// skip paths that are not processable
		basename := filepath.Base(info.Name())
		if skipDirs.Has(basename) || strings.HasPrefix(basename, ".") || info.Size() == 0 {
			return nil
		}

		pathRelToTmp, err := filepath.Rel(tmpDir, path)
		if err != nil {
			pathRelToTmp = path
		}

		isDiscontinued := strings.Contains(path, "/discontinued/")
		osr, err := parseOSReleaseFile(path, isDiscontinued)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}
		if osr != nil && osr.ID != "" {
			ids.Add(osr.ID)
			osReleaseInfos = append(osReleaseInfos, *osr)
			fmt.Printf("  cataloged %s --> %s\n", pathRelToTmp, osr.String())
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking through repo: %w", err)
	}

	fmt.Printf("Cataloged %d os-release entries\n", len(osReleaseInfos))

	// write out to dest file

	fh, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	enc := json.NewEncoder(fh)
	enc.SetIndent("", "  ")

	if err := enc.Encode(osReleaseInfos); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	idsSlice := ids.List()
	sort.Strings(idsSlice)

	return idsSlice, fh.Close()
}

func parseOSReleaseFile(path string, isDiscontinued bool) (*Info, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]interface{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := trimQuotes(parts[1])

		if key == "ID_LIKE" {
			data[key] = strings.Fields(value)
		} else {
			data[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	info := &Info{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   info,
		TagName:  "field",
	})
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(data); err != nil {
		return nil, err
	}

	if info.ID == "" {
		return nil, nil
	}

	info.Discontinued = isDiscontinued

	if info.VersionID != "" {
		fields := strings.Split(info.VersionID, ".")
		info.MajorVersion = fields[0]
		if len(fields) > 1 {
			info.MinorVersion = fields[1]
		}
	}

	return info, nil
}

func trimQuotes(s string) string {
	return strings.Trim(s, `"`)
}
