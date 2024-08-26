// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

/*
 * See run-parallel-mandelbrot.sh.
 *
 * Things start to get real fast (<1sec) once we're getting
 * around one goroutine per row (per py); the measurements
 * have been tweaked to reflect it.
 *
 * Note that the time include the creation of the goroutines,
 * but not their termination, which seems to takes a considerable
 * amount of time, and might negate some of those gain.
 *
 * Mind also that there are a few weird datapoints, positively
 * and negatively.
 *
 * So, for better tests, we would need:
 *	- more precise instrumentation:
 *		- startup time
 *		- computation time
 *		- teardown time
 *	- multiple measure per entries to smooth weird points
 *	- vary the number of core
 *	- make sure the final output is OK
 *
 * width	height	number-of-go-routines	time

	1024   1024   10     578.519772ms
	1024   1024   50     494.835377ms
	1024   1024   100    445.157147ms
	1024   1024   200    427.13681ms
	1024   1024   400    259.150182ms
	1024   1024   1000   5.28626ms
	1024   1024   2000   3.641103ms
	1024   1024   3000   6.564645ms
	1024   1024   5000   7.763983ms
	2048   2048   10     2.042246575s
	2048   2048   50     2.01938778s
	2048   2048   100    1.989605451s
	2048   2048   200    1.911031869s
	2048   2048   400    1.638578288s
	2048   2048   1000   1.110740349s
	2048   2048   2000   15.230199ms
	2048   2048   3000   5.280755ms
	2048   2048   5000   8.210826ms
	4096   4096   1000   7.024499006s
	4096   4096   2000   4.578171052s
	4096   4096   3000   633.264472ms
	4096   4096   5000   8.220888ms
	4096   4096   7000   12.230262ms
	8192   8192   3000   17.006217073s
	8192   8192   5000   6.758450038s
	8192   8192   7000   1.258371366s
	8192   8192   8000   197.845566ms
	8192   8192   9000   20.221371861s    # surprisingly slow
	16384  16384  5000   1m53.639348769s
	16384  16384  10000  20.971572671s
	16384  16384  15000  51.502851775s
	16384  16384  20000  38.679373ms      # surprisingly fast
	16384  16384  30000  28.704541699s
 */

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"time"
)

var (
	// 0.5 or 1 are way too much;
	// set to zero to disable supersampling.
	δx     = flag.Float64("dx", 0.001, "supersampling δx")
	δy     = flag.Float64("dy", 0.001, "supersampling δy")
	ngo    = flag.Int("ngo", 20, "number of goroutines")
	width  = flag.Int("width", 2048, "output's width")
	height = flag.Int("height", 2048, "output's height")
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
)

type Data struct {
	y1, y2 float64
	py int
}

func computeAndSet(img *image.RGBA, c <-chan *Data) {
	for d := range c {
		for px := 0; px < *width; px++ {
			x1 := (float64(px)+*δx)/float64(*width)*(xmax-xmin) + xmin
			x2 := (float64(px)-*δy)/float64(*width)*(xmax-xmin) + xmin

			z1 := complex(x1, d.y1)
			z2 := complex(x1, d.y2)
			z3 := complex(x2, d.y1)
			z4 := complex(x2, d.y2)

			m := avg(mandelbrot(z1), mandelbrot(z2), mandelbrot(z3), mandelbrot(z4))

			// Image point (px, py) represents complex value z.
			img.Set(px, d.py, m)
		}
	}
}

func main() {
	flag.Parse()
	start := time.Now()

	img := image.NewRGBA(image.Rect(0, 0, *width, *height))

	c := make(chan *Data)

	for i := 0; i < *ngo; i++ {
		go computeAndSet(img, c)
	}

	for py := 0; py < *height; py++ {
		y1 := (float64(py)+*δx)/float64(*height)*(ymax-ymin) + ymin
		y2 := (float64(py)-*δy)/float64(*height)*(ymax-ymin) + ymin

		c <- &Data{y1, y2, py}
	}

	close(c)

	fmt.Fprintf(os.Stderr, "%v\n", time.Since(start))

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
