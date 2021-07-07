package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func shred(img image.Image, t string, outName string) {

	fmt.Println("shred an image")

	bounds := img.Bounds()
	fmt.Println(bounds)

	m := bounds.Max
	mx := m.X

	cols := make([]int, mx)

	for i := range cols {
		cols[i] = i
	}

	//fmt.Println("cols ", cols)

	rand.Seed(time.Now().UnixNano())
	for i := len(cols) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		cols[i], cols[j] = cols[j], cols[i]
	}

	out := image.NewRGBA(bounds)

	for i, x := range cols {
		fmt.Println(x)
		for y := 0; y < bounds.Max.Y; y++ {

			out.Set(i, y, img.At(x, y))

		}

	}

	outfile, _ := os.Create(outName)
	defer outfile.Close()

	png.Encode(outfile, out)

}
