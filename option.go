package himage

import (
	"errors"
)

// Anchor is the anchor point for image alignment.
// It has been taken from the repo github.com/disintegration/imaging for easy access to the package.
type Anchor int

// Anchor point positions.
// It has been taken from the repo github.com/disintegration/imaging for easy access to the package.
const (
	Center Anchor = iota
	TopLeft
	Top
	TopRight
	Left
	Right
	BottomLeft
	Bottom
	BottomRight
)

// Resize ..
type Resize struct {
	Anchor         Anchor
	Ratio          int
	Width          int
	Height         int
	WidthOriented  bool
	HeightOriented bool
	Maximize       bool
	Minimize       bool
}

// Valid ..
func (r Resize) Valid() error {
	if (r.Ratio > 0 && r.Width > 0) || (r.Ratio > 0 && r.Height > 0) {
		return errors.New("both ratio and resolution cannot be specified at the same time")
	}

	return nil
}
