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

		for y := 0; y < 200; y++ {

			pixel := img.At(x, y)

			r, g, b, a := pixel.RGBA()

			cols[x] = append(cols[x], color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

		}

	}

	var order []int
	usedCols := make(map[int]int)

	//start position
	//column 1
	fmt.Println("cols",len(cols))
	fmt.Println("cols[0]",len(cols[0]))
	

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

func findClosestColumn(index int, cols map[int][]color.RGBA, used map[int]int) (int, map[int]int) {

	//var rtn int

	delta := make(map[int][]map[string]float64)

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
		deltaArray := make([]map[string]float64, len(cols[0]))

		for i, p := range cols[k] {
			scoreDelta := make(map[string]float64)

			r, g, b, _ := p.RGBA()  //compare each pixel at each y location against each 
			sr,sg,sb, _ := scores[i].RGBA()
			// scoreDelta = append(scoreDelta, math.Abs(float64(r-sr)))
			// scoreDelta = append(scoreDelta, math.Abs(float64(g-sg)))
			// scoreDelta = append(scoreDelta, math.Abs(float64(b-sb)))

			scoreDelta["r"] = math.Abs(float64(r-sr))
			scoreDelta["g"] = math.Abs(float64(g-sg))
			scoreDelta["b"] = math.Abs(float64(b-sb))

			
			deltaArray = append(deltaArray, scoreDelta)
			
			delta[i] = deltaArray

		}


	}
	fmt.Println(len(used))
	// fmt.Println("cols,",len(cols))
	// fmt.Println("delta", len(delta))
	// fmt.Println("detla[0]",len(delta[0]))

	//closest := -1
	//max64 := math.MaxFloat64

	// for i, v := range delta {
		
	// 	//fmt.Println(i,v)




	// }
	

	return index, used
}

