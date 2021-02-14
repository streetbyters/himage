package himage

import (
	"errors"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
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

func (i *Himage) SetName(name string) *Himage {
	i.name = name
	return i
}

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

func (i *Himage) makeQuality() *Himage {
	i.quality = make(map[string]interface{})
	i.quality["jpg"] = 100
	i.quality["jpeg"] = 100
	i.qJPEG = 100
	i.quality["png"] = png.DefaultCompression
	i.qPNG = png.DefaultCompression

	return i
}

// detail fetch image details (size, resolutions etc.)
func (i *Himage) detail() *Himage {
	if i.moved {
		i.inDetail()
		return i
	}

	if i.path != "" {
		i.inDetail()
	} else if i.Multipart != nil {
		f, err := i.Multipart.Open()
		if err != nil {
			i.Error = err
			return i
		}
		defer f.Close()
		mime, _ := mimetype.DetectReader(f)
		i.Detail.Mime = mime.String()

		c, _, err := image.DecodeConfig(f)
		if err != nil {
			i.Error = err
			return i
		}
		i.Detail.Width = c.Width
		i.Detail.Height = c.Height

	} else if i.File != nil {
		stat, err := i.File.Stat()
		if err != nil {
			i.Error = err
			return i
		}
		i.Detail.Size = stat.Size()

		mime, _ := mimetype.DetectReader(i.File)
		i.Detail.Mime = mime.String()
		c, _, err := image.DecodeConfig(i.File)
		if err != nil {
			i.Error = err
			return i
		}
		i.Detail.Width = c.Width
		i.Detail.Height = c.Height
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

func (i *Himage) Run() (*Himage, error) {
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
