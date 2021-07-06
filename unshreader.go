package main

import (
	"fmt"
	"image"
)

func unShred(img image.Image, t string, outName string ) {

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
	
	for x:= 0; x < bounds.Max.X; x++ {
		
		for y := 0; y < bounds.Max.Y; y++ {
			
			//fmt.Println("{x,y}",x,y )

			pixel := img.At(x,y)

			r,g,b,a := pixel.RGBA();

			score := r*g*b+a;
			cols[x] += score;
						
		}
		fmt.Println("col: ",x,", value: ",  cols[x],)
	}

	var order []int
	usedCols := make(map[int]int)

	//start position
	//column 1

	basicSort(order, cols, usedCols)

}

func basicSort(order []int, cols map[int]uint32, usedCols map[int]int){

	for k, _ := range cols {

		if used(k, usedCols){
			continue;
		}

		findClosestColumn(k, cols, usedCols);



	}

}

func findClosestColumn(index int, cols map[int]uint32, usedCols map[int]int)(){
	

}

func used(index int, usedCols map[int]int) bool{
	for _, v := range usedCols {
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