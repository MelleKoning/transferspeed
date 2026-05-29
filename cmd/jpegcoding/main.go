package main

/*
#cgo LDFLAGS: -lturbojpeg
#include <turbojpeg.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

func main() {
	inputPath := "./data/lionface.jpg"
	outputPath := "output.jpg"

	// 1. Read image from disk
	imgData, err := os.ReadFile(inputPath)

	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Create a C pointer to the image data
	cImgData := C.CBytes(imgData)
	defer C.free(cImgData)

	// 2. Initialize Decompressor
	// Note: Modern TurboJPEG uses tjInitDecompress() without arguments
	decompHandle := C.tjInitDecompress()
	if decompHandle == nil {
		log.Fatal("Unable to create decompress handle")
	}
	defer C.tjDestroy(decompHandle) // Fixed: TurboJPEG uses unified tjDestroy

	// 3. Extract Header
	var width, height, jpegSubsamp C.int
	if C.tjDecompressHeader2(decompHandle, (*C.uchar)(cImgData), C.ulong(len(imgData)), &width, &height, &jpegSubsamp) < 0 {
		log.Fatal("Failed to decompress header")
	}
	fmt.Printf("Image size: %dx%d\n", width, height)

	// 4. Decode image to RGB bitmap
	bufferSize := int(width * height * 3)
	buffer := C.malloc(C.size_t(bufferSize))
	defer C.free(buffer)

	// Fixed: Passed required size and matched types
	if C.tjDecompress2(decompHandle, (*C.uchar)(cImgData), C.ulong(len(imgData)), (*C.uchar)(buffer), width, 0, height, C.TJPF_RGB, 0) < 0 {
		log.Fatal("Decompression failed")
	}

	// 5. Draw Square Box in the Middle Quarter
	pixels := (*[1 << 30]byte)(buffer)[:bufferSize:bufferSize]

	startX := int(width) / 4
	endX := (3 * int(width)) / 4
	startY := int(height) / 4
	endY := (3 * int(height)) / 4

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			offset := (y*int(width) + x) * 3
			pixels[offset] = 255 // Red
			pixels[offset+1] = 0 // Green
			pixels[offset+2] = 0 // Blue
		}
	}

	// 6. Initialize Compressor
	compHandle := C.tjInitCompress()
	if compHandle == nil {
		log.Fatal("Unable to create compress handle")
	}
	defer C.tjDestroy(compHandle) // Fixed: Unified tjDestroy

	// 7. Encode Image
	var compressedBuf *C.uchar = nil
	var compressedSize C.ulong = 0

	// TurboJPEG manages allocating this buffer dynamically if we pass a pointer to a null pointer.
	// Fixed: Removed TJFLAG_ICCR (not standard) and passed the required subsampling factor.
	res := C.tjCompress2(
		compHandle,
		(*C.uchar)(buffer),
		width,
		0, // pitch
		height,
		C.TJPF_RGB,
		&compressedBuf,
		&compressedSize,
		jpegSubsamp, // Preserve original subsampling
		90,          // Quality
		0,           // Flags
	)

	if res < 0 {
		log.Fatal("Compression failed")
	}
	// TurboJPEG allocated compressedBuf on the C heap, we must free it using tjFree
	defer C.tjFree((*C.uchar)(unsafe.Pointer(compressedBuf)))

	// Convert C buffer to Go slice
	outputSlice := C.GoBytes(unsafe.Pointer(compressedBuf), C.int(compressedSize))

	// 8. Write back to disk
	err = os.WriteFile(outputPath, outputSlice, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Println("Image processed and saved successfully.")
}
