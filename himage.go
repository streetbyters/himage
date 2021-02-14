package himage

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Himage ..
type Himage struct {
	Path      string
	Multipart *multipart.FileHeader
	File      *os.File
	Error     error
	Detail    struct {
		Width  int
		Height int
		Mime   string
		Size   int64
	}
	dst      string
	quality  map[string]interface{}
	qJPEG    int
	qPNG     png.CompressionLevel
	moved    bool
	tempPath string
}

// NewHimageWithPath ..
func NewHimageWithPath(p string) *Himage {
	i := new(Himage)
	i.Path = p
	i.detail().makeQuality()

	return i
}

// NewHimageWithMultipart ..
func NewHimageWithMultipart(f *multipart.FileHeader) *Himage {
	i := new(Himage)
	i.Multipart = f
	i.detail().makeQuality()

	return i
}

// NewHimageWithFile ..
func NewHimageWithFile(f *os.File) *Himage {
	i := new(Himage)
	i.File = f
	i.detail()

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
		break
	case "image/png":
		i.quality["png"] = q.(png.CompressionLevel)
		break
	}
	return i
}

func (i *Himage) makeQuality() *Himage {
	i.quality = make(map[string]interface{})
	i.quality["jpg"] = 100
	i.quality["jpeg"] = 100
	i.quality["png"] = png.DefaultCompression

	return i
}

// detail fetch image details (size, resolutions etc.)
func (i *Himage) detail() *Himage {
	if i.moved {
		i.inDetail()
		return i
	}

	if i.Path != "" {
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

	if i.Path != "" {

	} else if i.Multipart != nil {

	}

	if i.Error == nil {
		i.moved = true
	}
	return i
}

// Resize ..
func (i *Himage) Resize() *Himage {

	return i
}

// Optimize ..
func (i *Himage) Optimize() *Himage {

	return i
}

// inDetail ..
func (i *Himage) inDetail() *Himage {
	f, err := os.Open(i.Path)
	if err != nil {
		i.Error = err
		return i
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		i.Error = err
		return i
	}
	i.Detail.Size = stat.Size()

	c, _, err := image.DecodeConfig(f)
	if err != nil {
		i.Error = err
		return i
	}
	i.Detail.Width = c.Width
	i.Detail.Height = c.Height

	mime, err := mimetype.DetectReader(f)
	if err != nil {
		i.Error = err
		return i
	}
	i.Detail.Mime = mime.String()

	return i
}

func (i *Himage) makeTemp() *Himage {
	if i.Error != nil {
		return i
	}
	ext := strings.Split(i.Detail.Mime, "/")[0]
	i.tempPath = filepath.Join(os.TempDir(), fmt.Sprintf("%s.%s", uuid.New().String(), ext))
	_, err := os.Create(string(os.PathSeparator) + i.tempPath)
	if err != nil {
		i.Error = err
		return i
	}

	return i
}

func (i *Himage) moveToTemp() *Himage {
	if i.Error != nil {
		return i
	}

	return i
}

func (i *Himage) Run() (*Himage, error) {
	if i.tempPath != "" {
		defer os.Remove(i.tempPath)
	}

	return i, i.Error
}
