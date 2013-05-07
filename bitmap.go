/*
Package bitmap implements a simple type to load, hold and manipulate bitmaps.

*/
package goglutils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type Bitmap struct {
	header        []byte
	dataPos       uint32
	width, height uint32
	imageSize     uint32
	data          []byte
}

// NewBitmap takes a filename and returns a bitmap object that holds the
// data in that file.  If it cannot load the bitmap for whatever reason,
// the Bitmap pointer will be nil.
func NewBitmap(path string) (*Bitmap, error) {
	var b *Bitmap = new(Bitmap)
	if _, err := b.Load(path); err != nil {
		return nil, err
	}
	return b, nil
}

// Load takes a filename and loads the data in the file into
// the Bitmap receiver.  If it cannot load the bitmap, size
// will be 0.
func (b *Bitmap) Load(path string) (int, error) {

	fp, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fp.Close()

	// Get the 54-byte header and load it into the Bitmap object
	b.header = make([]byte, 54)
	if _, err := fp.Read(b.header[0:54]); err != nil {
		return 0, err
	}

	// Parse the header - is this actually a bitmap?
	if b.header[0] != 'B' || b.header[1] != 'M' {
		err := errors.New(fmt.Sprintf("File is not a bitmap: %s", path))
		return 0, err
	}

	// Now read stats from the header
	b.dataPos = binary.LittleEndian.Uint32(b.header[10:14])
	b.imageSize = binary.LittleEndian.Uint32(b.header[2:6])
	b.width = binary.LittleEndian.Uint32(b.header[18:22])
	b.height = binary.LittleEndian.Uint32(b.header[22:26])

	// Sanity checks
	if b.imageSize == 0 {
		b.imageSize = b.width * b.height * 3 // one byte per color
	}
	if b.dataPos == 0 {
		b.dataPos = 54 // BMP header done that way.
	}

	// Allocate appropriate slice to hold the data
	b.data = make([]byte, b.imageSize)
	// Seek to image data offset
	fp.Seek((int64)(b.dataPos), 0)
	size, err := fp.Read(b.data)
	if err != nil {
		return size, err
	}
	return size, nil

}

func (b *Bitmap) Info() {
	fmt.Fprintf(os.Stdout, "Image Header Info\n")
	fmt.Fprintf(os.Stdout, "Width x Height: %d x %d\n", b.width, b.height)
	fmt.Fprintf(os.Stdout, "Image Size: %d\n", b.imageSize)
	fmt.Fprintf(os.Stdout, "Header data:\n")
	for i, val := range b.header {
		fmt.Fprintf(os.Stdout, "%0x ", val)
		if (i+1)%8 == 0 {
			fmt.Fprintf(os.Stdout, "\n")
		}
	}
	// Report the first 16 bytes of data, so we can compare
	// and check that data is correct
	fmt.Println("\n****** DATA ******")
	for i := 0; i < 16; i++ {
		fmt.Fprintf(os.Stdout, "%x ", b.data[i])
		if (i+1)%8 == 0 {
			fmt.Fprintf(os.Stdout, "\n")
		}
	}
}
