//go:build windows && !amd
// +build windows,!amd

package imgcoder

// GetEncoder version that returns a default
// golang software encoder for this build arch
func GetEncoder(jpgEncodingQuality int) JpegEncoder {
	return &SoftwareJpegEncoder{
		DefaultJPGEncodingQuality: jpgEncodingQuality,
	}
}
