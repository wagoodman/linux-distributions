package eol

import (
	"encoding/json"
	"fmt"
)

type ReleaseInfo struct {
	APIID             string           `json:"apiID"`
	OSReleaseID       string           `json:"osReleaseID"`
	Cycle             string           `json:"cycle,omitempty"`
	Codename          string           `json:"codename,omitempty"`
	LTS               *ConditionalDate `json:"lts,omitempty"`
	ReleaseDate       string           `json:"releaseDate,omitempty"`
	EOL               *ConditionalDate `json:"eol,omitempty"`
	Latest            string           `json:"latest,omitempty"`
	LatestReleaseDate string           `json:"latestReleaseDate,omitempty"`
	Support           *ConditionalDate `json:"support,omitempty"`
	ExtendedSupport   *ConditionalDate `json:"extendedSupport,omitempty"`
	Discontinued      *ConditionalDate `json:"discontinued,omitempty"`
	Link              string           `json:"link,omitempty"`
}

// types derived from https://endoflife.date/docs/api

type ConditionalDate struct {
	Evaluated bool   `json:"evaluated"`
	Until     string `json:"until,omitempty"`
}

func (r *ReleaseInfo) UnmarshalJSON(data []byte) error {
	type Alias ReleaseInfo
	aux := &struct {
		LTS             any `json:"lts"`
		EOL             any `json:"eol"`
		Support         any `json:"support"`
		ExtendedSupport any `json:"extendedSupport"`
		Discontinued    any `json:"discontinued"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal Release: %w", err)
	}

	handleField := func(value any) *ConditionalDate {
		if value == nil {
			return nil
		}
		switch v := value.(type) {
		case bool:
			if v {
				return &ConditionalDate{Evaluated: true}
			}
			return &ConditionalDate{Evaluated: false}
		case string:
			return &ConditionalDate{Evaluated: true, Until: v}
		default:
			panic(fmt.Errorf("unexpected type: %T", value))
		}
	}

	r.LTS = handleField(aux.LTS)
	r.EOL = handleField(aux.EOL)
	r.Support = handleField(aux.Support)
	r.ExtendedSupport = handleField(aux.ExtendedSupport)
	r.Discontinued = handleField(aux.Discontinued)

	return nil
}
