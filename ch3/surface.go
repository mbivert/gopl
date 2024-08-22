// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sync"
)

var mu sync.Mutex

const (
	width, height = 800, 600            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func waterdrop(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

// https://mathcurve.com/surfaces.gb/boiteaoeufs/boiteaoeufs.shtml
func eggbox(x, y float64) float64 {
	a := 1./2.
	b := 4.
	return a*(math.Sin(x/b)+math.Sin(y/b))
}

// Formula is correct, but I'm too lazy to make it look good
// https://mathcurve.com/surfaces.gb/paraboloidhyperbolic/paraboloidhyperbolic.shtml
func saddle(x, y float64) float64 {
	a := 1./2.
	b := 1./2.
	return (math.Pow(x/a, 2) - math.Pow(y/b, 2))/ 12
}

func corner(f func(float64, float64) float64, i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func printSVG(w io.Writer, f func(float64, float64) float64) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)

	for i := 0; i < cells; i++ {
Loop:
		for j := 0; j < cells; j++ {
			ax, ay := corner(f, i+1, j)
			bx, by := corner(f, i, j)
			cx, cy := corner(f, i, j+1)
			dx, dy := corner(f, i+1, j+1)

			// 3.1
			for _, x := range []float64{ax, ay, bx, by, cx, cy, dx, dy} {
				if math.IsInf(x, 0) || math.IsNaN(x) {
					continue Loop
				}
			}

			// TODO: use stroke='#%x00%x' to colorize,
			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}

	}

	fmt.Fprintf(w, "</svg>\n")
}

func main() {
	// 3.4
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		switch r.URL.Path {
		case "/waterdrop": printSVG(w, waterdrop)
		case "/eggbox":    printSVG(w, eggbox)
		case "/saddle":    printSVG(w, saddle)
		}

	})

	port := ":8080"
	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
