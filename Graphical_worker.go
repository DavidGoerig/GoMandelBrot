/*

 @author:  David Goerig
 @id:      djg53
 @module:  Concurrency and Parallelism - CO890
 @asses:   assess 4- Go Worker and geometric
           distribution comparaison
*/

/*
** thread pool / worker implementation of mandelbrot
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
/*
** vector used to launch worker
*/
type vectors struct {
	x int
	y int
}

//	global wait group variable
var wg sync.WaitGroup

/*
*	param:	a complex123
*	desc:	manderlbort calculation
*	return:	result
*/
func mandelbrotCalc(a complex128) float64 {
	i := 0
	for z := a; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = z*z + a
	}
	return float64(maxEsc-i) / maxEsc
}

/*
*	param:	posChan <-chan: channel in order to have the position to compute
*	desc:	get by the posChan the position to compute, and compute it
*	return:	/
*/
func farmerWorkInit(posChan <-chan vectors, b *image.NRGBA)  {
	scale := width / (rMax - rMin)
	defer wg.Done()
	for pos := range posChan {
		tempx := pos.x
		tempy := pos.y
		fEsc := mandelbrotCalc(complex(float64(tempx)/scale+rMin, float64(tempy)/scale+iMin))
		color := color.NRGBA{uint8(red * fEsc), uint8(green * fEsc), uint8(blue * fEsc), 255}
		b.Set(tempx, tempy, color)
	}
}

/*
*	param:	posChan <- vectors: the channel in order to send pos to the worker
*	desc:	launch workers by giving them work
*	return:	/
*/
func launchWorkers(posChan chan <- vectors) {
	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			posChan <- vectors{x: x, y:y}
		}
	}
	close(posChan)
}

/*
*	param:	workerPool int: sice of the pool
*	desc:	main farmer / thread pool function, create farmer first, then send work
*	return:	computed image
*/
func farmers(workerPool int) *image.NRGBA{
	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	bounds := image.Rect(0, 0, width, height)
	posChan := make(chan vectors)
	b := image.NewNRGBA(bounds)

	runtime.GOMAXPROCS(workerPool)
	// Image creation
	draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)
	// Workers initialisation
	for nb := 0; nb < workerPool; nb += 1 {
		wg.Add(1)
		go farmerWorkInit(posChan, b)
	}
	// send data to workers through the pos channel
	launchWorkers(posChan)
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
	workerPool, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err, "Nbr of distribution")
		os.Exit(0)
	}
	var b = farmers(workerPool)
	print_in_png(b)
}