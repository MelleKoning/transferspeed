//go:build !windows && amd64
// +build !windows,amd64

package imgcoder_test

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/MelleKoning/transferspeed/pkg/imgcoder"
)

func TestLibJpegEncoder(t *testing.T) {
	img, err := LoadJpegImageLibJpeg("../../data/lionface.jpg")
	if err != nil {
		t.Fatalf("unable to load image: %v", err)
	}

	start := time.Now()

	for n := 0; n < 5; n++ {
		buffer := new(bytes.Buffer) // Create a new buffer for each iteration

		encoder := imgcoder.GetEncoder(85)
		if err := encoder.Encode(buffer, img, nil); err != nil {
			t.Fatal()
		}

		assert.NotNil(t, buffer) // buffer contains the jpeg now
	}

	fmt.Printf("LibJpegEncoder took %s", time.Since(start))
}

func TestDecodeImage(t *testing.T) {
	file, err := os.Open("../../data/lionface.jpg")
	if err != nil {
		t.Fatal()
	}
	defer file.Close()

	coder := imgcoder.LibJpegEncoder{}

	var img image.Image

	if img, err = coder.Decode(file); err != nil {
		t.Fatal()
	}

	assert.NotNil(t, img)
}

func BenchmarkLibJpegEncoder(b *testing.B) {
	img, err := LoadJpegImageLibJpeg("../data/lionface.jpg")
	if err != nil {
		b.Fatalf("unable to load image: %v", err)
	}

	b.ResetTimer()

	encoder := imgcoder.GetEncoder(85)

	for i := 0; i < b.N; i++ {
		buffer := new(bytes.Buffer) // Create a new buffer for each iteration
		if err := encoder.Encode(buffer, img, nil); err != nil {
			b.Fatal()
		}
	}

	ms := float64(b.Elapsed()/time.Duration(b.N)) / 1e6
	b.ReportMetric(ms, "ms/op")
}

func BenchmarkLibJpegDecoder(b *testing.B) {
	file, err := os.Open("../../data/lionface.jpg")
	if err != nil {
		b.Fatal()
	}
	defer file.Close()

	b.ResetTimer()

	coder := imgcoder.LibJpegEncoder{}

	for i := 0; i < b.N; i++ {
		var img image.Image

		_, _ = file.Seek(0, 0) // gomnd:nolint

		if img, err = coder.Decode(file); err != nil {
			b.Fatalf("%v", err)
		}

		assert.NotNil(b, img)
	}

	ms := float64(b.Elapsed()/time.Duration(b.N)) / 1e6
	b.ReportMetric(ms, "ms/op")
}

func LoadJpegImageLibJpeg(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := imgcoder.LibJpegEncoder{}

	img, err := decoder.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
