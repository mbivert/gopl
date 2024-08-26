// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

var (
	// 0.5 or 1 are way too much;
	// set to zero to disable supersampling.
	δx = flag.Float64("dx", 0.001, "supersampling δx")
	δy = flag.Float64("dy", 0.001, "supersampling δy")
)

func main() {
	flag.Parse()

	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height; py++ {
		y1 := (float64(py)+*δx)/height*(ymax-ymin) + ymin
		y2 := (float64(py)-*δy)/height*(ymax-ymin) + ymin

		for px := 0; px < width; px++ {
			x1 := (float64(px)+*δx)/width*(xmax-xmin) + xmin
			x2 := (float64(px)-*δy)/width*(xmax-xmin) + xmin

			z1 := complex(x1, y1)
			z2 := complex(x1, y2)
			z3 := complex(x2, y1)
			z4 := complex(x2, y2)

			m := avg(mandelbrot(z1), mandelbrot(z2), mandelbrot(z3), mandelbrot(z4))

			// Image point (px, py) represents complex value z.
			img.Set(px, py, m)
		}
	}

	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func avg(x, y, z, t color.Color) color.Color {
	// uint32 should be enough, but as the extra space is somewhat
	// "reserved" already (could be used), let's get bigger.
	var r, g, b, a uint64

	for _, u := range []color.Color{x, y, z, t} {
		ru, gu, bu, au := u.RGBA()
		r += uint64(ru)
		g += uint64(gu)
		b += uint64(bu)
		a += uint64(au)
	}

	return color.RGBA{uint8(r / 4), uint8(g / 4), uint8(b / 4), uint8(a / 4)}
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	var v complex128

	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{255 - contrast*n, 255 - contrast*(n/2), 0, 0xFF}
		}
	}

	return color.RGBA{10, 165, 0, 0xFF}
}
