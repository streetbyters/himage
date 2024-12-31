package himage

import (
	"errors"
	"github.com/disintegration/imaging"
	"image"
	"image/png"
	"mime/multipart"
	"os"
)

// Himage ..
type Himage struct {
	Multipart *multipart.FileHeader
	File      *os.File
	Error     error
	Detail    struct {
		Width  int
		Height int
		Mime   string
		Size   int64
	}
	path         string
	dst          string
	quality      map[string]interface{}
	qJPEG        int
	qPNG         png.CompressionLevel
	moved        bool
	resized      bool
	optimized    bool
	tempPath     string
	name         string
	removeOrigin bool
}

// NewHimageWithPath ..
func NewHimageWithPath(p string) *Himage {
	i := new(Himage)
	i.path = p
	i.removeOrigin = false
	i.detail().makeQuality()

	return i
}

// NewHimageWithMultipart ..
func NewHimageWithMultipart(f *multipart.FileHeader) *Himage {
	i := new(Himage)
	i.Multipart = f
	i.removeOrigin = false
	i.detail().makeQuality()
	return i
}

// NewHimageWithFile ..
func NewHimageWithFile(f *os.File) *Himage {
	i := new(Himage)
	i.File = f
	i.removeOrigin = false
	i.detail().makeQuality()
	return i
}

// SetName ..
func (i *Himage) SetName(name string) *Himage {
	i.name = name
	return i
}

// RemoveOrigin ..
func (i *Himage) RemoveOrigin(val bool) *Himage {
	i.removeOrigin = val
	return i
}

// SetDestination ..
func (i *Himage) SetDestination(p string) *Himage {
	i.dst = p
	return i
}

// SetQuality ..
func (i *Himage) SetQuality(q interface{}) *Himage {
	switch i.Detail.Mime {
	case "image/jpg", "image/jpeg":
		i.quality["jpg"] = q.(int)
		i.quality["jpeg"] = q.(int)
		i.qJPEG = q.(int)
		break
	case "image/png":
		i.quality["png"] = q.(png.CompressionLevel)
		i.qPNG = q.(png.CompressionLevel)
		break
	}
	return i
}

// Move ..
func (i *Himage) Move() *Himage {
	if i.Error != nil {
		return i
	}

	if i.dst == "" {
		i.Error = errors.New("destination path is nil")
		return i
	}

	i.moveToTemp()

	if i.Error == nil {
		i.moved = true
	}
	return i
}

// Resize ..
func (i *Himage) Resize(option Resize) *Himage {
	if !i.moved {
		i.Move()
	}

	if i.Error != nil {
		return i
	}

	if err := option.Valid(); err != nil {
		i.Error = err
		return i
	}

	src, err := imaging.Open(i.tempPath)
	if err != nil {
		i.Error = err
		return i
	}

	width := option.Width
	height := option.Height

	if option.Ratio > 0 {
		if option.WidthOriented {
			if option.Maximize {
				width = option.Width + (option.Width / option.Ratio)
			} else {
				width = option.Width - (option.Width / option.Ratio)
			}
		}

		if option.HeightOriented {
			if option.Maximize {
				width = option.Height + (option.Height / option.Ratio)
			} else {
				width = option.Height - (option.Height / option.Ratio)
			}
		}
	}

	var im *image.NRGBA

	if option.Anchor > 0 {
		im = imaging.Fill(src, width, height, imaging.Anchor(option.Anchor), imaging.Lanczos)
	} else {
		im = imaging.Resize(src, width, height, imaging.Lanczos)
	}

	i.save(im)

	if i.Error != nil {
		i.resized = true
	}

	return i
}

// Run ..
func (i *Himage) Finish() (*Himage, error) {
	if i.tempPath != "" {
		defer os.Remove(i.tempPath)
	}

	if i.path != "" {
		i.Error = os.Remove(i.path)
	} else if i.File != nil {
		i.File.Close()
		i.Error = os.Remove(i.File.Name())
	}

	return i, i.Error
}
