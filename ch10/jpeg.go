// a elementary subset of magick(1)
package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

var format = flag.String("fmt", "png", "output file format: gif|jpg|png")

// Encoder sometimes takes options, and sometimes not, so we can't
// map to e.g. "gif" -> gif.Encode()
var fmts = map[string]bool{
	"gif": true,
	"jpg": true,
	"png": true,
}

func main() {
	flag.Parse()

	if err := magick(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}

}

func magick(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	switch *format {
	case "gif":
		return gif.Encode(out, img, nil)
	case "jpg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(out, img)
	}
	return fmt.Errorf("Unknown output format '%s'", *format)
}
