//go:build !windows && amd64
// +build !windows,amd64

package imgcoder

import (
	"image"
	"image/jpeg"
	"io"

	libjpeg "github.com/pixiv/go-libjpeg/jpeg"
)

type LibJpegEncoder struct {
	jpgEncodingQuality int
}

// GetEncoder version that returns a libjpeg hardware encoder
// for this build arch
func GetEncoder(jpgEncodingQuality int) JpegEncoder {
	return &LibJpegEncoder{
		jpgEncodingQuality: jpgEncodingQuality,
	}
}

func (libJpegEncoder *LibJpegEncoder) Encode(w io.Writer, m image.Image, o *jpeg.Options) error {
	// forward the requested Quality
	var options *libjpeg.EncoderOptions
	if o != nil {
		options = &libjpeg.EncoderOptions{Quality: o.Quality}
	} else {
		// the libjpeg needs default Quality options passed in
		options = &libjpeg.EncoderOptions{Quality: libJpegEncoder.jpgEncodingQuality}
	}

	libjpeg.
	return libjpeg.Encode(w, m, options)
}

func (LibJpegEncoder) Decode(r io.Reader) (image.Image, error) {
	// forward the requested Quality
	options := &libjpeg.DecoderOptions{
		DCTMethod: libjpeg.DCTFloat,
	}

	return libjpeg.Decode(r, options)
}
