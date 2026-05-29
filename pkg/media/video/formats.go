package video

import (
	"sort"
	"strings"

	"github.com/photoprism/photoprism/pkg/clean"
)

// Formats represents a set of video container and codec names for allow/exclude
// lookups. Keys are stored lowercased and matched case-insensitively.
type Formats map[string]bool

// NewFormats creates a Formats set from a comma-separated string of container
// and/or codec names. Returns an empty set when the input is empty.
func NewFormats(formats ...string) Formats {
	n := strings.Count(strings.Join(formats, clean.FormatSep), clean.FormatSep)
	list := make(Formats, n+8)

	for i := range formats {
		list.Set(formats[i])
	}

	return list
}

// Contains reports whether any of the given format names is in the list.
// Empty values are skipped so callers can pass both a codec and a container
// even when only one is known.
func (b Formats) Contains(formats ...string) bool {
	if len(b) == 0 || len(formats) == 0 {
		return false
	}

	for _, format := range formats {
		if format = clean.Format(format); format == "" {
			continue
		} else if _, ok := b[format]; ok {
			return true
		}
	}

	return false
}

// Allow reports whether the given format is NOT in the list.
func (b Formats) Allow(format string) bool {
	return !b.Contains(format)
}

// Set adds a comma-separated list of format names to the list.
func (b Formats) Set(formats string) {
	if formats == "" {
		return
	}

	for c := range strings.SplitSeq(formats, clean.FormatSep) {
		b.Add(c)
	}
}

// Add adds a format name to the list after normalizing it.
func (b Formats) Add(format string) {
	format = clean.Format(format)
	if format == "" {
		return
	}

	b[format] = true
}

// String returns the list as a comma-separated string in alphabetical order.
func (b Formats) String() string {
	if len(b) == 0 {
		return ""
	}

	list := make([]string, 0, len(b))
	for s := range b {
		list = append(list, s)
	}

	sort.Strings(list)
	return strings.Join(list, ", ")
}
