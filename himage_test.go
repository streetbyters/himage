package himage

import (
	"errors"
	"image/png"
	"path/filepath"
	"testing"
)

func TestNewHimageWithPath(t *testing.T) {
	hImage := NewHimageWithPath(filepath.Join("test-files", "850x566.png"))
	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("detail mime is not valid"))
	}

	if hImage.Detail.Width != 850 {
		t.Error(errors.New("detail width is not valid"))
	}

	if hImage.Detail.Height != 566 {
		t.Error(errors.New("detail height is not valid"))
	}

	if hImage.Detail.Size == 0 {
		t.Error(errors.New("detail size is not valid"))
	}

	if hImage.qPNG != png.DefaultCompression {
		t.Error(errors.New("png compression is not valid"))
	}

	if len(hImage.quality) == 0 {
		t.Error(errors.New("quality specs is not valid"))
	}
}

func TestNewHimageWithPathNotExistsFile(t *testing.T) {
	hImage := NewHimageWithPath(filepath.Join("test-files", "notfoung.png"))
	if hImage.Error == nil {
		t.Error(errors.New("invalid file open"))
	}
}
