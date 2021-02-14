package himage

import (
	"bytes"
	"errors"
	"github.com/disintegration/imaging"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_detail(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "2200x1467.png")
	hImage.detail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 2200 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 1467 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_detailOSFile(t *testing.T) {
	hImage := new(Himage)
	hImage.File, _ = os.Open(filepath.Join("test-files", "2200x1467.png"))
	defer hImage.File.Close()
	hImage.detail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 2200 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 1467 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_inDetail(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "2200x1467.png")
	hImage.inDetail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 2200 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 1467 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_inDetailNotExistsFile(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "notexist.png")
	hImage.inDetail()

	if hImage.Error == nil {
		t.Error(errors.New("erroneous reading operation"))
	}
}

func Test_makeTemp(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "2200x1467.png")
	hImage.inDetail()

	hImage.makeTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	_, err := os.Stat(hImage.tempPath)
	if os.IsNotExist(err) {
		t.Error(errors.New("temp file does not created"))
	}
}

func Test_moveToTempLargeFileAndWithoutTemp(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "2200x1467.png")
	hImage.inDetail()

	hImage.moveToTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	_, err := os.Stat(hImage.tempPath)
	if os.IsNotExist(err) {
		t.Error(errors.New("temp file does not created"))
	}
	p := hImage.tempPath

	hImage = new(Himage)
	hImage.path = p
	hImage.inDetail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 2200 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 1467 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_moveToTempLargeFile(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "2200x1467.png")
	hImage.inDetail()

	hImage.makeTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	hImage.moveToTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	_, err := os.Stat(hImage.tempPath)
	if os.IsNotExist(err) {
		t.Error(errors.New("temp file does not created"))
	}
	p := hImage.tempPath

	hImage = new(Himage)
	hImage.path = p
	hImage.inDetail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 2200 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 1467 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_moveToTempSmallFileAndWithoutTemp(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "10x10.png")
	hImage.inDetail()

	hImage.moveToTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	_, err := os.Stat(hImage.tempPath)
	if os.IsNotExist(err) {
		t.Error(errors.New("temp file does not created"))
	}
	p := hImage.tempPath

	hImage = new(Himage)
	hImage.path = p
	hImage.inDetail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 10 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 10 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_moveToTempSmallFile(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "10x10.png")
	hImage.inDetail()

	hImage.makeTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	hImage.moveToTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	_, err := os.Stat(hImage.tempPath)
	if os.IsNotExist(err) {
		t.Error(errors.New("temp file does not created"))
	}
	p := hImage.tempPath

	hImage = new(Himage)
	hImage.path = p
	hImage.inDetail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 10 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 10 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_moveToTempLargeOSFileAndWithoutTemp(t *testing.T) {
	hImage := new(Himage)
	osFile, _ := os.Open(filepath.Join("test-files", "2200x1467.png"))
	defer osFile.Close()
	hImage.File = osFile
	hImage.detail()

	hImage.moveToTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	_, err := os.Stat(hImage.tempPath)
	if os.IsNotExist(err) {
		t.Error(errors.New("temp file does not created"))
	}
	p := hImage.tempPath

	hImage = new(Himage)
	hImage.path = p
	hImage.inDetail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 2200 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 1467 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_moveToTempSmallOSFile(t *testing.T) {
	hImage := new(Himage)
	osFile, _ := os.Open(filepath.Join("test-files", "10x10.png"))
	defer osFile.Close()
	hImage.File = osFile
	hImage.detail()

	hImage.makeTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	hImage.moveToTemp()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	_, err := os.Stat(hImage.tempPath)
	if os.IsNotExist(err) {
		t.Error(errors.New("temp file does not created"))
	}
	p := hImage.tempPath

	hImage = new(Himage)
	hImage.path = p
	hImage.inDetail()

	if hImage.Error != nil {
		t.Error(hImage.Error)
	}

	if hImage.Detail.Width != 10 {
		t.Error(errors.New("width is not valid"))
	}

	if hImage.Detail.Height != 10 {
		t.Error(errors.New("height is not valid"))
	}

	if hImage.Detail.Mime != "image/png" {
		t.Error(errors.New("mime type is not valid"))
	}
}

func Test_bytesCopy(t *testing.T) {
	hImage := new(Himage)

	size := 32 * 1024
	readerBytes := make([]byte, size)
	for i := 0; i < size; i++ {
		readerBytes[i] = 1
	}

	writesBytes := make([]byte, 0)
	buff := bytes.NewBuffer(writesBytes)

	err := hImage.bytesCopy(int64(size), bytes.NewReader(readerBytes), buff)
	if err != nil {
		t.Error(err)
	}

	if buff.Len() != size {
		t.Error(errors.New("erroneous writing operation"))
	}
}

func Test_bytesCopyReaderSizeNotValid(t *testing.T) {
	hImage := new(Himage)

	size := 0
	readerBytes := make([]byte, size)
	for i := 0; i < size; i++ {
		readerBytes[i] = 1
	}

	writesBytes := make([]byte, 0)
	buff := bytes.NewBuffer(writesBytes)

	err := hImage.bytesCopy(int64(32*1024), bytes.NewReader(readerBytes), buff)
	if err == nil {
		t.Error(errors.New("erroneous reading operation"))
	}

	if buff.Len() != 0 {
		t.Error(errors.New("erroneous writing operation"))
	}
}

func Test_bytesCopyWriterSizeNotValid(t *testing.T) {
	hImage := new(Himage)

	size := 32 * 1024
	readerBytes := make([]byte, size)
	for i := 0; i < size; i++ {
		readerBytes[i] = 1
	}

	f, _ := ioutil.TempFile(string(os.PathSeparator)+"tmp", "himage")
	f.Close()

	err := hImage.bytesCopy(int64(size), bytes.NewReader(readerBytes), f)
	if err == nil {
		t.Error(errors.New("erroneous writing operation"))
	}
}

func Test_bytesWrite(t *testing.T) {
	hImage := new(Himage)

	size := 32 * 1024 * 1024
	readerBytes := make([]byte, size)
	for i := 0; i < size; i++ {
		readerBytes[i] = 1
	}

	f, _ := ioutil.TempFile(string(os.PathSeparator)+"tmp", "himage")
	defer f.Close()

	err := hImage.bytesWrite(bytes.NewReader(readerBytes), f)
	if err != nil {
		t.Error(err)
	}
	f.Sync()
	st, _ := f.Stat()

	if st.Size() != int64(size) {
		t.Error(errors.New("erroneous writing operation"))
	}
}

func Test_bytesWriteWriteError(t *testing.T) {
	hImage := new(Himage)

	size := 32 * 1024 * 1024
	readerBytes := make([]byte, size)
	for i := 0; i < size; i++ {
		readerBytes[i] = 1
	}

	f, _ := ioutil.TempFile(string(os.PathSeparator)+"tmp", "himage")
	f.Close()

	err := hImage.bytesWrite(bytes.NewReader(readerBytes), f)
	if err == nil {
		t.Error(errors.New("erroneous writing operation"))
	}
}

func Test_bytesWriteReadError(t *testing.T) {
	hImage := new(Himage)

	f, _ := ioutil.TempFile(string(os.PathSeparator)+"tmp", "himage")
	f.Close()

	f2, _ := ioutil.TempFile(string(os.PathSeparator)+"tmp", "himage")
	defer f.Close()

	err := hImage.bytesWrite(f, f2)
	if err == nil {
		t.Error(errors.New("erroneous reading operation"))
	}
}

func Test_save(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "2200x1467.png")
	hImage.detail().makeQuality().SetQuality(png.BestCompression)

	hImage.moveToTemp()

	src, _ := imaging.Open(hImage.path)
	fill := imaging.Resize(src, 100, 100, imaging.Lanczos)
	hImage.save(fill)
	if hImage.Error != nil {
		t.Error(hImage.Error)
	}
}

func Test_saveWithJPGFile(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "640x426.jpeg")
	hImage.detail().makeQuality().SetQuality(50)

	hImage.moveToTemp()

	src, _ := imaging.Open(hImage.path)
	fill := imaging.Resize(src, 100, 100, imaging.Lanczos)
	hImage.save(fill)
	if hImage.Error != nil {
		t.Error(hImage.Error)
	}
}

func Test_saveWithJPGFileInvalidTempPath(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "640x426.jpeg")
	hImage.detail().makeQuality().SetQuality(50)

	src, _ := imaging.Open(hImage.path)
	fill := imaging.Resize(src, 100, 100, imaging.Lanczos)
	hImage.save(fill)
	if hImage.Error == nil {
		t.Error(errors.New("erroneous writing operation"))
	}
}

func Test_saveInvalidTempPath(t *testing.T) {
	hImage := new(Himage)
	hImage.path = filepath.Join("test-files", "10x10.png")
	hImage.detail().makeQuality()

	src, _ := imaging.Open(hImage.path)
	fill := imaging.Resize(src, 100, 100, imaging.Lanczos)
	hImage.save(fill)
	if hImage.Error == nil {
		t.Error(errors.New("erroneous writing operation"))
	}
}
