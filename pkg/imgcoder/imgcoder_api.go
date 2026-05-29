package imgcoder

import (
	"image"
	"image/jpeg"
	"io"
)

// JpegEncoder to support both a software and hardware jpeg encoder
type JpegEncoder interface {
	Encode(w io.Writer, m image.Image, o *jpeg.Options) error
	Decode(r io.Reader) (image.Image, error)
}

type SoftwareJpegEncoder struct {
	DefaultJPGEncodingQuality int
}

// fallback
func GetSoftwareEncoder(jpgEncodingQuality int) JpegEncoder {
	return &SoftwareJpegEncoder{
		DefaultJPGEncodingQuality: jpgEncodingQuality,
	}
}

func (sEncoder *SoftwareJpegEncoder) Encode(w io.Writer, m image.Image, o *jpeg.Options) error {
	if o == nil {
		o = &jpeg.Options{Quality: sEncoder.DefaultJPGEncodingQuality}
	}

	return jpeg.Encode(w, m, o) // Delegate to the jpeg.Encode function
}

func (sEncoder *SoftwareJpegEncoder) Decode(r io.Reader) (image.Image, error) {
	return jpeg.Decode(r)
}
