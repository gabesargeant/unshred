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

	m := bounds.Max
	mx := m.X
	step := 9

	cols := make(map[int]uint32, mx)

	for x := 0; x < bounds.Max.X; x++ {

		for y := 0; y < bounds.Max.Y; y++ {

			//fmt.Println("{x,y}",x,y )

			pixel := img.At(x, y)

			//r, _, _, _ := pixel.RGBA()
			//r, g, _, _ := pixel.RGBA()
			r, g, b, _ := pixel.RGBA()
			//r, g, b, a:= pixel.RGBA()
			//fmt.Println("rgba : ", r, " ", g, " ", b)
			if r == 0 {
				r = 1
			}
			if g == 0 {
				g = 1
			}
			if b == 0 {
				b = 1
			}
			score := b // * a
			if y%step == 0 {
				cols[x] += score
			}

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

func basicSort(order []int, cols map[int]uint32, used map[int]int) []int {

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

func findClosestColumn(index int, cols map[int]uint32, used map[int]int) (int, map[int]int) {

	var delta float64 = math.MaxFloat64
	var rtn int
	for k := range cols {

		if k == index {
			continue
		}

		if _, ok := used[k]; ok {
			continue
		}

		d := math.Abs(float64(cols[index] - cols[k]))

		if d <= delta {
			delta = d
			rtn = k
		}

	}

	return rtn, used
}

func prependInt(x []int, y int) []int {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}
