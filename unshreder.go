package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)

func unShred(img image.Image, t string, outName string) {

	fmt.Println("Put te image back together")

	bounds := img.Bounds()
	fmt.Println(bounds)

	//m := bounds.Max
	//mx := m.X
	//step := m.Y / 30

	cols := make(map[int][]uint32)

	for x := 0; x < bounds.Max.X; x++ {

		for y := 0; y < bounds.Max.Y; y++ {

			pixel := img.At(x, y)

			r, g, b, _ := pixel.RGBA()

			if r == 0 {
				r = 1
			}
			if g == 0 {
				g = 1
			}
			if b == 0 {
				b = 1
			}
			//use greyscale for the score
			//score := ((float64(r)*0.21)*(float64(b)*0.07) + (float64(g) * 0.72)) / 3 // * a
			score := r + b
			cols[x] = append(cols[x], uint32(score))

		}

	}

	var order []int
	usedCols := make(map[int]int)

	//start position
	//column 1

	order = basicSort(order, cols, usedCols)

	//writeimage(order, image)
	//outputCols := make([]int, mx)

	out := image.NewRGBA(bounds)

	for i, x := range order {
		//fmt.Println("final x ", x)
		for y := 0; y < bounds.Max.Y; y++ {

			out.Set(i, y, img.At(x, y))

		}
	}

	outfile, _ := os.Create(outName)
	defer outfile.Close()

	png.Encode(outfile, out)

}

func basicSort(order []int, cols map[int][]uint32, used map[int]int) []int {

	order = append(order, 0)
	used[0] = 0
	tick := 0
	for len(order) != len(cols) {
		n, used := findClosestColumn(order[tick], cols, used)
		tick++
		used[n] = n
		order = append(order, n)
		//fmt.Println(len(order), len(cols))
	}
	fmt.Println("return order length", len(order))
	//fmt.Println(order)
	return order
}

func findClosestColumn(index int, cols map[int][]uint32, used map[int]int) (int, map[int]int) {

	var rtn int
	
	var delta []float64
	var closer = 0

	for k := range cols {

		if k == index {
			continue
		}

		if _, ok := used[k]; ok {
			continue
		}

		scores := cols[index]
		scores2 := cols[k]

		var scoreDelta []float64

		for i := range scores {
			//fmt.Println(i);
			scoreDelta = append(scoreDelta, math.Abs(float64(scores[i]-scores2[i])))
		}

		if len(delta) == 0 {
		 	for i:=0; i < len(scoreDelta); i++{
				 delta = append(delta, math.MaxFloat64)
			 }
		}

		clsr := 0
		for i := range scoreDelta {

			if scoreDelta[i] <= delta[i] {
				//fmt.Println(scoreDelta[])
				clsr++
			}

		}

		if clsr > closer {
			delta = scoreDelta
			rtn = k
			closer = clsr
		}

	}
	fmt.Println(closer);

	return rtn, used
}

func prependInt(x []int, y int) []int {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}
