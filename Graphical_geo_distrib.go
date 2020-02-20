package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
)

const (
	maxEsc = 100
	rMin   = -2.
	rMax   = .5
	iMin   = -1.
	iMax   = 1.
	width  = 750
	red    = 230
	green  = 235
	blue   = 255
)

var nbrOfDistribution = 5
var wg sync.WaitGroup

func mandelbrot(a complex128) float64 {
	i := 0
	for z := a; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = z*z + a
	}
	return float64(maxEsc-i) / maxEsc
}

func compute_mandel_on_range(start int, limit int, height int, scale float64, b *image.NRGBA)  {
	defer wg.Done()
	for x := start; x < limit; x++ {
		for y := 0; y < height; y++ {
			fEsc := mandelbrot(complex(
				float64(x)/scale+rMin,
				float64(y)/scale+iMin))
			b.Set(x, y, color.NRGBA{uint8(red * fEsc),
				uint8(green * fEsc), uint8(blue * fEsc), 255})

		}
	}
}

func main() {
	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	bounds := image.Rect(0, 0, width, height)
	b := image.NewNRGBA(bounds)
	draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)
	runtime.GOMAXPROCS(nbrOfDistribution)
	for distrib_nbr := 0; distrib_nbr < nbrOfDistribution; distrib_nbr+=1 {
		wg.Add(1)
		start := (width / nbrOfDistribution) * distrib_nbr
		limit := start + (width / nbrOfDistribution)
		go compute_mandel_on_range(start, limit, height, scale, b)
	}
	wg.Wait()
	f, err := os.Create("mandelbrot.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = png.Encode(f, b); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}
}