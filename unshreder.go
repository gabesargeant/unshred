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

	cols := make(map[int]uint32, mx)
	//fill cols with numbers
	// for i := range cols {
	// 	cols[i] = i
	// }

	for x := 0; x < bounds.Max.X; x++ {

		for y := 0; y < bounds.Max.Y; y++ {

			//fmt.Println("{x,y}",x,y )

			pixel := img.At(x, y)

			r, g, b, a := pixel.RGBA()

			score := r*g*b + a
			cols[x] += score

		}
		fmt.Println("col: ", x, ", value: ", cols[x])
	}

	var order []int
	usedCols := make(map[int]int)

	//start position
	//column 1

	order = basicSort(order, cols, usedCols)


	//writeimage(order, image)
	//outputCols := make([]int, mx)

	out := image.NewRGBA(bounds)

	for i, x := range order{
		
		for y := 0; y < bounds.Max.Y; y++ {

			out.Set(i, y, img.At(x,y))

		}
	}

	outfile, _ := os.Create(outName)
	defer outfile.Close()

	png.Encode(outfile, out)



}

func basicSort(order []int, cols map[int]uint32, used map[int]int) []int{

	order = append(order, 0)
	used[0]=0;

	for len(order) != len(cols) {
		n, used := findClosestColumn(order[len(order)-1], cols, used)
		used[n] = n
		order = append(order, n)
		fmt.Println(len(order), len(cols))
	}
	return order
}

func findClosestColumn(index int, cols map[int]uint32, used map[int]int) (int, map[int]int) {

	var smaller int = 0
	var larger int = 0

	for k := range cols {
		fmt.Print(k);
		if usedColumn(index, used) {
			continue
		}

		if cols[index] < cols[k] && cols[k] <= cols[smaller] {
			smaller = k
		}

		if cols[index] > cols[k] && cols[k] >= cols[larger] {
			larger = k
		}

	}
	largerDelta := math.Abs(float64(cols[index] - cols[larger]))

	smallerDelta := math.Abs(float64(cols[index] - cols[smaller]))

	if smallerDelta < largerDelta {
		return smaller,used
	}
	return larger,used
}

func usedColumn(index int, used map[int]int) bool {
	for _, v := range used {
		if v == index {
			return true
		}
	}
	return false
}

func prependInt(x []int, y int) []int {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}
