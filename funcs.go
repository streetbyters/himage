package himage

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

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

		i.File.Seek(0, 0)
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

// inDetail ..
func (i *Himage) inDetail() *Himage {
	f, err := os.Open(i.path)
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

	i.File.Seek(0, 0)
	c, _, err := image.DecodeConfig(f)
	if err != nil {
		i.Error = err
		return i
	}
	i.Detail.Width = c.Width
	i.Detail.Height = c.Height

	mime, err := mimetype.DetectFile(i.path)
	if err != nil {
		i.Error = err
		return i
	}
	i.Detail.Mime = mime.String()

	return i
}

// makeTemp ..
func (i *Himage) makeTemp() *Himage {
	if i.Error != nil {
		return i
	}
	//ext := strings.Split(i.Detail.Mime, "/")[0]*/

	name := uuid.New().String()
	if i.dst == "" && i.name != "" {
		name = i.name
	}

	i.tempPath = filepath.Join(os.TempDir(), fmt.Sprintf("%s.jpg", name))
	_, err := os.Create(string(os.PathSeparator) + i.tempPath)
	if err != nil {
		i.Error = err
		return i
	}

	return i
}

// makeQuality ..
func (i *Himage) makeQuality() *Himage {
	i.quality = make(map[string]interface{})
	i.quality["jpg"] = 100
	i.quality["jpeg"] = 100
	i.qJPEG = 100
	i.quality["png"] = png.DefaultCompression
	i.qPNG = png.DefaultCompression

	return i
}

// moveToTemp ..
func (i *Himage) moveToTemp() *Himage {
	if i.tempPath == "" {
		i.makeTemp()
	}

	if i.Error != nil {
		return i
	}

	d, err := os.OpenFile(i.tempPath, os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		i.Error = err
		return i
	}
	defer d.Close()

	if i.path != "" {
		i.movePathToTemp(d)
	} else if i.Multipart != nil {
		i.moveMultipartToTemp(d)
	} else if i.File != nil {
		i.moveFileToTemp(d)
	}

	return i
}

// movePathToTemp ..
func (i *Himage) movePathToTemp(d *os.File) *Himage {
	f, err := os.Open(i.path)
	if err != nil {
		i.Error = err
		return i
	}
	defer f.Close()
	f.Seek(0, 0)
	_, err = io.Copy(d, f)
	if err != nil {
		i.Error = err
		return i
	}

	return i
}

// moveMultipartToTemp ..
func (i *Himage) moveMultipartToTemp(d *os.File) *Himage {
	f, err := i.Multipart.Open()
	if err != nil {
		i.Error = err
		return i
	}
	f.Seek(0, 0)
	if i.Detail.Size <= int64(CHUNK_SIZE) {
		if err := i.bytesCopy(int64(CHUNK_SIZE), f, d); err != nil {
			i.Error = err
			return i
		}
	}

	if err := i.bytesWrite(f, d); err != nil {
		i.Error = err
		return i
	}

	return i
}

func (i *Himage) moveFileToTemp(d *os.File) *Himage {
	i.File.Seek(0, 0)
	if i.Detail.Size <= int64(CHUNK_SIZE) {
		if err := i.bytesCopy(int64(CHUNK_SIZE), i.File, d); err != nil {
			i.Error = err
			return i
		}
	}

	if err := i.bytesWrite(i.File, d); err != nil {
		i.Error = err
		return i
	}

	return i
}

// bytesCopy ..
func (i *Himage) bytesCopy(size int64, r io.Reader, w io.Writer) error {
	buffer := make([]byte, size)
	_, err := r.Read(buffer)
	if err != nil {
		return err
	}

	_, err = w.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

// bytesWrite ..
func (i *Himage) bytesWrite(r io.Reader, w io.Writer) error {
	reading := true
	for reading {
		buffer := make([]byte, CHUNK_SIZE)
		n, err := r.Read(buffer)
		if err != nil {
			if err == io.EOF {
				reading = false
				continue
			} else {
				return err
			}
		}
		_, err = w.Write(buffer[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

// save ..
func (i *Himage) save(im *image.NRGBA) {
	os.Remove(i.tempPath)
	switch i.Detail.Mime {
	case "image/jpeg", "image/jpg":
		if err := imaging.Save(im, i.tempPath, imaging.JPEGQuality(i.qJPEG)); err != nil {
			i.Error = err
		}
		break
	case "image/png":
		if err := imaging.Save(im, i.tempPath, imaging.PNGCompressionLevel(i.qPNG)); err != nil {
			i.Error = err
		}
		break
	}
}
