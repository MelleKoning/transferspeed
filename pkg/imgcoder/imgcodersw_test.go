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

func TestSoftwareEncoder(t *testing.T) {
	img, err := LoadJpegImage("../../data/lionface.jpg")
	if err != nil {
		t.Fatalf("unable to load image: %v", err)
	}

	start := time.Now()

	encoder := imgcoder.GetSoftwareEncoder(85)

	for n := 0; n < 5; n++ {
		buffer := new(bytes.Buffer) // Create a new buffer for each iteration
		if err := encoder.Encode(buffer, img, nil); err != nil {
			t.Fatal()
		}

		assert.NotNil(t, buffer) // buffer contains the jpeg now
	}

	fmt.Printf("SoftwareEncoder took %s", time.Since(start))
}

func BenchmarkSoftwareJpegEncoder(b *testing.B) {
	// we are calling GetSoftwareEncoder instead of GetEncoder
	// otherwise we will get the hardwareencoder!
	encoder := imgcoder.GetSoftwareEncoder(85)

	img, err := LoadJpegImage("../../data/lionface.jpg")
	if err != nil {
		b.Fatalf("unable to load image: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buffer := new(bytes.Buffer) // Create a new buffer for each iteration
		if err := encoder.Encode(buffer, img, nil); err != nil {
			b.Fatal("Encode failed:", err)
		}
	}

	ms := float64(b.Elapsed()/time.Duration(b.N)) / 1e6
	b.ReportMetric(ms, "ms/op")
}

func BenchmarkSoftwareJpegDecoder(b *testing.B) {
	file, err := os.Open("../../data/lionface.jpg")
	if err != nil {
		b.Fatal()
	}
	defer file.Close()

	b.ResetTimer()

	coder := imgcoder.GetSoftwareEncoder(85)

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

func LoadJpegImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := imgcoder.SoftwareJpegEncoder{}

	img, err := decoder.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
