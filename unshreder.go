package main

import (
	"fmt"
	"image"
	"image/color"
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

	cols := make(map[int][]color.RGBA)

	for x := 0; x < bounds.Max.X; x++ {

		for y := 0; y < bounds.Max.Y; y++ {

			pixel := img.At(x, y)

			r, g, b, a := pixel.RGBA()

			cols[x] = append(cols[x], color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

			// * a

		}

	}

	var order []int
	usedCols := make(map[int]int)

	//start position
	//column 1
	fmt.Println("cols", len(cols))
	fmt.Println("cols[0]", len(cols[0]))

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

func basicSort(order []int, cols map[int][]color.RGBA, used map[int]int) []int {

	order = append(order, 0)
	used[0] = 0
	tick := 0
	for len(order) != len(cols) {
		if order[tick] == -1 {
			continue
		}

		n, used := findClosestColumn(order[tick], cols, used)
		tick++
		used[n] = n
		order = append(order, n)
		fmt.Println(len(order), len(cols))
	}
	fmt.Println("return order length", len(order))
	//fmt.Println(order)
	return order
}

type columnDiff struct {
	columnNumber int
	r            []float64
	g            []float64
	b            []float64
	lowest       int
}

func findClosestColumn(index int, cols map[int][]color.RGBA, used map[int]int) (int, map[int]int) {

	//var rtn int

	delta := make(map[int]columnDiff)

	//var closer = 0fmt.Println("cols,",len(cols))
	// fmt.Println("delta", len(delta))
	// fmt.Println("detla[0]",len(delta[0]))

	for k := range cols {

		if k == index {
			continue
		}

		if _, ok := used[k]; ok {
			continue
		}

		scores := cols[index]

		cd := columnDiff{}
		cd.columnNumber = k

		for i, p := range cols[k] {

			r, g, b, _ := p.RGBA() //compare each pixel at each y location against each
			if i > len(scores) {
				continue
			}

			sr, sg, sb, _ := scores[i].RGBA()

			//fmt.Println("abs", math.Abs(float64(r-sr))-math.Abs(float64(g-sg))-math.Abs(float64(b-sb)))
			cd.r = append(cd.r, math.Abs(float64(r-sr))-math.Abs(float64(g-sg))-math.Abs(float64(b-sb)))

		}
		if _, ok := used[k]; ok {
			continue
		}else{
			delta[k] = cd
		}
		

	}

	//he lowest array sequence in the diff compared to col diff on.
	fmt.Println("delta ", len(delta))
	fmt.Println("used ", len(used))
	rtn := findLowestDiff(delta)

	return rtn, used
}

func findLowestDiff(cols map[int]columnDiff) int {
	var ra [][]float64

	//PIVOT THE COLS INTO ROWS
	for _, v := range cols {
		ra = append(ra, v.r)
	}

	lowestR := make(map[int]int)
	for i, v := range ra {
		tmp := math.MaxFloat64
		col := -1
		for ii, vv := range v {

			if vv < tmp {
				tmp = vv
				col = ii
			}
		}
		lowestR[i] = col
	}

	rtn := findMostFrequentElement(lowestR)

	return rtn
}

func findMostFrequentElement(input map[int]int) int {

	freq := make(map[int]int)
	for k, _ := range input {

		if val, ok := freq[k]; ok {
			freq[val]++
		} else {
			freq[val] = 1
		}
	}

	maxV := -1

	for _, v := range freq {

		if v > maxV {
			maxV = v

		}

	}

	return maxV
}
