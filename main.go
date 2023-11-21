package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"os"
	"time"
)

type imgData struct {
	size          int
	scale         float64
	cx, cy        float64
	maxIterations int
}

var palette = make([]color.Color, 0)

func main() {
	fractal(os.Stdout)
}

func fractal(out io.Writer) {

	nthreadsPtr := flag.Uint("t", 1, "Número de hilos a usar")
	sizePtr := flag.Uint("s", 256, "Tamaño del GIF")
	nframesPtr := flag.Uint("n", 256, "Número de frames del GIF")
	periodPtr := flag.Uint("p", 10, "Periodo de oscilación de la fase")

	flag.Parse()

	nthreads := int(*nthreadsPtr)
	size := int(*sizePtr)
	nframes := int(*nframesPtr)
	period := int(*periodPtr)

	delay := period * 100 / nframes
	dphase := float64(2.0*math.Pi) / float64(nframes)

	palette = append(palette, color.Black)
	for i := 1; i <= 0xFF; i++ {
		palette = append(palette, color.RGBA{0, uint8(i), 0, 0xFF})
	}

	mod := 0.75
	phase := math.Pi/2 + 0.3

	maxIterations := 256
	scale := 1.0 / (float64(size) / 2)

	anim := gif.GIF{LoopCount: nframes}

	start := time.Now()
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, size, size)
		img := image.NewPaletted(rect, palette)

		cx := mod * math.Cos(phase)
		cy := mod * math.Sin(phase)

		var data imgData
		data.cx = cx
		data.cy = cy
		data.size = size
		data.scale = scale
		data.maxIterations = maxIterations

		var ch = make(chan bool, nthreads)

		for j := 0; j < nthreads; j++ {
			ch <- true
		}

		x := 0
		for range ch {
			go renderColumn(x, img, data, ch)
			x++
			if x >= size {
				break
			}
		}

		phase += dphase
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)

	fmt.Fprintln(os.Stderr, time.Since(start).Seconds())
}

func renderColumn(x int, img *image.Paletted, data imgData, ch chan<- bool) {
	for y := 0; y < data.size; y++ {
		px := float64(x-data.size/2) * data.scale
		py := float64(y-data.size/2) * data.scale

		iterations := computeIterations(px, py, data.cx, data.cy, data.maxIterations)

		index := uint8(iterations)
		if data.maxIterations != len(palette) {
			index = uint8((float64(iterations) / float64(data.maxIterations)) * float64((len(palette))))
			index = max(0, index-1)
		}
		img.SetColorIndex(x, y, index)
	}

	ch <- true
}

func computeIterations(x, y, cx, cy float64, maxIteration int) int {
	zx := x
	zy := y
	iterations := 0

	for modSquared(zx, zy) <= 4.0 && iterations < maxIteration-1 {
		zx, zy = computeNext(zx, zy, cx, cy)
		iterations++
	}

	return iterations
}

func computeNext(currX, currY, cx, cy float64) (x, y float64) {
	// z_n = z_n-1^2 + c
	x = currX*currX - currY*currY + cx
	y = 2.0*currX*currY + cy

	return x, y
}

func modSquared(x, y float64) float64 {
	return x*x + y*y
}
