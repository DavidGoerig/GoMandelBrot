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

type resultsFinal struct {
	x int
	y int
	c color.Color
}

type vectors struct {
	x int
	y int
}

var workerPool = 4
var wg sync.WaitGroup

func mandelbrotCalc(a complex128) float64 {
	i := 0
	for z := a; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = z*z + a
	}
	return float64(maxEsc-i) / maxEsc
}

func farmerWorkInit(posChan <-chan vectors, results chan<- resultsFinal)  {
	scale := width / (rMax - rMin)
	defer wg.Done()
	for pos := range posChan {
		fmt.Println("Sombre pute")
		tempx := pos.x
		tempy := pos.y
		fEsc := mandelbrotCalc(complex(float64(tempx)/scale+rMin, float64(tempy)/scale+iMin))
		color := color.NRGBA{uint8(red * fEsc), uint8(green * fEsc), uint8(blue * fEsc), 255}
		results <- resultsFinal {
			x: tempx,
			y: tempy,
			c: color,
		}
	}
}


func launchWorkers(posChan chan <- vectors) {
	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			posChan <- vectors{x: x, y:y}
			fmt.Println("PosChan", x, y)
		}
	}
	close(posChan)
	wg.Wait()
}

func main() {
	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	bounds := image.Rect(0, 0, width, height)

	results := make(chan resultsFinal)
	posChan := make(chan vectors)
	b := image.NewNRGBA(bounds)

	runtime.GOMAXPROCS(workerPool)
	draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)
	for nb := 0; nb < workerPool; nb += 1 {
		wg.Add(1)
		go farmerWorkInit(posChan, results)
	}
	launchWorkers(posChan)
	close(results)
	for res := range results {
		b.Set(res.x, res.y, res.c)
	}
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