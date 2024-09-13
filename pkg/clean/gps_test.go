package clean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGPSBounds(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("41.87760543823242,-87.62521362304688,41.89404296875,-87.6215591430664")
		assert.InDelta(t, 41.8942, latNorth, 0.0000009)
		assert.InDelta(t, 41.8775, latSouth, 0.0000009)
		assert.InDelta(t, -87.6254, lngWest, 0.0000009)
		assert.InDelta(t, -87.6214, lngEast, 0.0000009)
		assert.NoError(t, err)
	})
	t.Run("FlippedLat", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("41.89404296875,-87.62521362304688,41.87760543823242,-87.6215591430664")
		assert.InDelta(t, 41.8942, latNorth, 0.0000009)
		assert.InDelta(t, 41.8775, latSouth, 0.0000009)
		assert.InDelta(t, -87.6254, lngWest, 0.0000009)
		assert.InDelta(t, -87.6214, lngEast, 0.0000009)
		assert.NoError(t, err)
	})
	t.Run("FlippedLng", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("41.87760543823242,-87.6215591430664,41.89404296875,-87.62521362304688")
		assert.InDelta(t, 41.8942, latNorth, 0.0000009)
		assert.InDelta(t, 41.8775, latSouth, 0.0000009)
		assert.InDelta(t, -87.6254, lngWest, 0.0000009)
		assert.InDelta(t, -87.6214, lngEast, 0.0000009)
		assert.NoError(t, err)
	})
	t.Run("Empty", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("")
		assert.Equal(t, float64(0), latNorth)
		assert.Equal(t, float64(0), lngEast)
		assert.Equal(t, float64(0), latSouth)
		assert.Equal(t, float64(0), lngWest)
		assert.Error(t, err)
	})
	t.Run("One", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("41.87760543823242")
		assert.Equal(t, float64(0), latNorth)
		assert.Equal(t, float64(0), lngEast)
		assert.Equal(t, float64(0), latSouth)
		assert.Equal(t, float64(0), lngWest)
		assert.Error(t, err)
	})
	t.Run("Three", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("41.87760543823242,-87.62521362304688,41.89404296875")
		assert.Equal(t, float64(0), latNorth)
		assert.Equal(t, float64(0), lngEast)
		assert.Equal(t, float64(0), latSouth)
		assert.Equal(t, float64(0), lngWest)
		assert.Error(t, err)
	})
	t.Run("Five", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("41.87760543823242,-87.62521362304688,41.89404296875,-87.6215591430664,41.89404296875")
		assert.Equal(t, float64(0), latNorth)
		assert.Equal(t, float64(0), lngEast)
		assert.Equal(t, float64(0), latSouth)
		assert.Equal(t, float64(0), lngWest)
		assert.Error(t, err)
	})
	t.Run("Invalid", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("95.87760543823242,-197.62521362304688,98.89404296875,-197.6215591430664")
		assert.InDelta(t, 90.0001, latNorth, 0.0000009)
		assert.InDelta(t, 89.9999, latSouth, 0.0000009)
		assert.InDelta(t, -180.0001, lngWest, 0.0000009)
		assert.InDelta(t, -179.9999, lngEast, 0.0000009)
		assert.NoError(t, err)
	})
	t.Run("Invalid2", func(t *testing.T) {
		latNorth, lngEast, latSouth, lngWest, err := GPSBounds("-95.87760543823242,197.62521362304688,-98.89404296875,197.6215591430664")
		assert.InDelta(t, -89.9999, latNorth, 0.0000009)
		assert.InDelta(t, -90.0001, latSouth, 0.0000009)
		assert.InDelta(t, 179.9999, lngWest, 0.0000009)
		assert.InDelta(t, 180.0001, lngEast, 0.0000009)
		assert.NoError(t, err)
	})
}

func TestGPSLatRange(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		latNorth, latSouth, err := GPSLatRange(41.87760543823242, 2)
		assert.InDelta(t, 41.8913, latNorth, 0.0000009)
		assert.InDelta(t, 41.8639, latSouth, 0.0000009)
		assert.NoError(t, err)
	})
	t.Run("Zero", func(t *testing.T) {
		latNorth, latSouth, err := GPSLatRange(0, 2)
		assert.Equal(t, float64(0), latNorth)
		assert.Equal(t, float64(0), latSouth)
		assert.Error(t, err)
	})
}

func TestGPSLngRange(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		lngEast, lngWest, err := GPSLngRange(-87.62521362304688, 2)
		assert.InDelta(t, -87.6389, lngWest, 0.0000009)
		assert.InDelta(t, -87.6116, lngEast, 0.0000009)
		assert.NoError(t, err)
	})
	t.Run("Zero", func(t *testing.T) {
		lngEast, lngWest, err := GPSLngRange(0, 2)
		assert.Equal(t, float64(0), lngEast)
		assert.Equal(t, float64(0), lngWest)
		assert.Error(t, err)
	})
}
