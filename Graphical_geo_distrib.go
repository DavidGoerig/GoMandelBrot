/*

 @author:  David Goerig
 @id:      djg53
 @module:  Concurrency and Parallelism - CO890
 @asses:   assess 4- Go Worker and geometric
           distribution comparaison
*/

/*
** geometric distribution representation of mandelbrot
*/
package main
/*
** needed imports
*/
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
	"strconv"
)

/*
** const used for the mandelbro creation
*/
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

//	global wait group variable
var wg sync.WaitGroup


/*
*	param:	a complex123
*	desc:	manderlbort calculation
*	return:	result
*/
func mandelbrot(a complex128) float64 {
	i := 0
	for z := a; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = z*z + a
	}
	return float64(maxEsc-i) / maxEsc
}

/*
*	param:	start: where starting the computation
			limit: where to stop the computation
			height: picture height
			scale: the scale
			b: the image
*	desc:	calc on range mandelbort
*	return:	/
*/

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

/*
*	param:	nbrOfDistribution int: nbr of intervals / in how many is divided the work
*	desc:	divide in an equal way the work in threads
*	return:	the image
*/
func geo_distrib_algo(nbrOfDistribution int) *image.NRGBA{
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
	return b
}

/*
*	param:	b image to print
*	desc:	print in a png obtained image
*	return:	/
*/
func print_in_png(b *image.NRGBA) {
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

/*
*	desc:	entry point
*/
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: build, and pass as argument in the number of equal distribution.")
		os.Exit(1)
	}
	nbrOfDistribution, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err, "Nbr of distribution")
		os.Exit(0)
	}
	var b = geo_distrib_algo(nbrOfDistribution)
	print_in_png(b)
}