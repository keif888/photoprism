package video

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormats(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := NewFormats("")
		assert.Empty(t, list)
	})
	t.Run("Single", func(t *testing.T) {
		list := NewFormats("magicyuv")
		assert.Len(t, list, 1)
		assert.True(t, list.Contains("magicyuv"))
	})
	t.Run("CommaSeparated", func(t *testing.T) {
		list := NewFormats("magicyuv, m8rg ,M8RA")
		assert.Len(t, list, 3)
		assert.True(t, list.Contains("magicyuv"))
		assert.True(t, list.Contains("m8rg"))
		assert.True(t, list.Contains("m8ra"))
	})
	t.Run("MixedContainerAndCodec", func(t *testing.T) {
		list := NewFormats("avi, magicyuv")
		assert.Len(t, list, 2)
		assert.True(t, list.Contains("avi"))
		assert.True(t, list.Contains("magicyuv"))
	})
}

func TestFormats_Contains(t *testing.T) {
	list := NewFormats("magicyuv, m8rg")

	t.Run("Match", func(t *testing.T) {
		assert.True(t, list.Contains("magicyuv"))
	})
	t.Run("CaseInsensitive", func(t *testing.T) {
		assert.True(t, list.Contains("MagicYUV"))
		assert.True(t, list.Contains("M8RG"))
	})
	t.Run("Trimmed", func(t *testing.T) {
		assert.True(t, list.Contains("  magicyuv  "))
		assert.True(t, list.Contains(`"magicyuv"`))
		assert.True(t, list.Contains(".magicyuv"))
	})
	t.Run("NoMatch", func(t *testing.T) {
		assert.False(t, list.Contains("avc1"))
	})
	t.Run("Empty", func(t *testing.T) {
		assert.False(t, list.Contains(""))
	})
	t.Run("EmptyList", func(t *testing.T) {
		empty := NewFormats("")
		assert.False(t, empty.Contains("magicyuv"))
	})
	t.Run("NoArgs", func(t *testing.T) {
		assert.False(t, list.Contains())
	})
	t.Run("MultipleArgsFirstMatches", func(t *testing.T) {
		assert.True(t, list.Contains("magicyuv", "avi"))
	})
	t.Run("MultipleArgsSecondMatches", func(t *testing.T) {
		assert.True(t, list.Contains("avc1", "magicyuv"))
	})
	t.Run("MultipleArgsWithEmpty", func(t *testing.T) {
		assert.True(t, list.Contains("", "magicyuv"))
		assert.False(t, list.Contains("", "avc1"))
	})
	t.Run("MultipleArgsNoMatch", func(t *testing.T) {
		assert.False(t, list.Contains("avc1", "hevc", "mp4"))
	})
}

func TestFormats_Allow(t *testing.T) {
	list := NewFormats("magicyuv")

	assert.False(t, list.Allow("magicyuv"))
	assert.True(t, list.Allow("avc1"))
	assert.True(t, list.Allow(""))
}

func TestFormats_Set(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := make(Formats)
		list.Set("")
		assert.Empty(t, list)
	})
	t.Run("AddsToExisting", func(t *testing.T) {
		list := NewFormats("magicyuv")
		list.Set("m8rg, AVI")
		assert.Len(t, list, 3)
		assert.True(t, list.Contains("magicyuv"))
		assert.True(t, list.Contains("m8rg"))
		assert.True(t, list.Contains("avi"))
	})
	t.Run("Duplicate", func(t *testing.T) {
		list := NewFormats("magicyuv")
		list.Set("MAGICYUV, magicyuv")
		assert.Len(t, list, 1)
	})
}

func TestFormats_Add(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		list := make(Formats)
		list.Add("MagicYUV")
		assert.Len(t, list, 1)
		assert.True(t, list.Contains("magicyuv"))
	})
	t.Run("Empty", func(t *testing.T) {
		list := make(Formats)
		list.Add("")
		list.Add("   ")
		list.Add(".")
		assert.Empty(t, list)
	})
	t.Run("Idempotent", func(t *testing.T) {
		list := make(Formats)
		list.Add("avi")
		list.Add("avi")
		assert.Len(t, list, 1)
	})
}

func TestFormats_String(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := NewFormats("")
		assert.Equal(t, "", list.String())
	})
	t.Run("Sorted", func(t *testing.T) {
		list := NewFormats("magicyuv, m8rg, m8ra")
		assert.Equal(t, "m8ra, m8rg, magicyuv", list.String())
	})
}
