package baselib

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

type Compressor struct {
	Quality int
	Encoder *png.Encoder
}

const (
	pngExt  = "png"
	jpegExt = "jpeg"
)

func New(q int, compressionLvl int) *Compressor {
	return &Compressor{
		Quality: q,
		Encoder: &png.Encoder{
			CompressionLevel: png.CompressionLevel(compressionLvl),
		},
	}
}

func (c *Compressor) Compress(path string, alias string) (string, error) {
	const op = "lib.compressor.compress"

	in, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", op, err)
	}
	defer in.Close()

	img, format, err := image.Decode(in)
	fmt.Println(format)

	if err != nil {
		return "", fmt.Errorf("failed to decode file %s: %w", op, err)
	}
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	outPath := filepath.Join(dir, alias+ext)

	out, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", op, err)
	}

	switch format {
	case pngExt:
		err = c.Encoder.Encode(out, img)
		if err != nil {
			return "", fmt.Errorf("failed to encode file %s: %w", op, err)
		}
	case jpegExt:
		if err = jpeg.Encode(out, img, &jpeg.Options{Quality: c.Quality}); err != nil {
			return "", fmt.Errorf("failed to encode file %s: %w", op, err)
		}
	}

	return outPath, nil
}
