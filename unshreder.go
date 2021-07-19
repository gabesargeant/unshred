package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type column struct {
	columnPosition int
	pixels         []color.RGBA
	used           bool
	rollingDelta   float64
}

func unShred(img image.Image, t string, outName string) {

	fmt.Println("Put te image back together")

	bounds := img.Bounds()
	fmt.Println(bounds)

	cols1 := make(map[int][]color.RGBA)

	cols2 := make([]column, 0)

	for x := 0; x < bounds.Max.X; x++ {

		c := column{}
		c.columnPosition = x
		for y := 0; y < bounds.Max.Y; y++ {

			pixel := img.At(x, y)

			r, g, b, a := pixel.RGBA()

			cols1[x] = append(cols1[x], color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

			c.pixels = append(c.pixels, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

		}
		cols2 = append(cols2, c)

	}

	var order []int
	usedCols1 := make(map[int]int)

	//start position
	//column 1
	fmt.Println("cols", len(cols1))
	fmt.Println("cols[0]", len(cols1[0]))

	order = basicSort(order, cols1, usedCols1, img.Bounds().Max.Y)

	//order = sort2(cols2, img.Bounds())


	// new sort. Pick 1 column. add one more column
	// For each new column added, check it's difference between all current columns and put it in place.

	//

	//

	order = secondSort(order, cols1, usedCols1, img.Bounds().Max.Y)



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

func basicSort(order []int, cols map[int][]color.RGBA, used map[int]int, height int) []int {
	startNo := 20
	order = append(order, startNo)
	used[startNo] = startNo
	tick := startNo
	for len(order) != len(cols) {
		
		n, used := findClosestColumn(order[tick%len(cols)], cols, used, height)
		tick++
		used[n] = n
		order = append(order, n)
		//	fmt.Println(len(order), len(cols))
	}
	fmt.Println("return order length", len(order))
	//fmt.Println(order)
	return order
}


func secondSort(order []int, columns map[int][]color.RGBA, used map[int]int, height int) []int {

	var picture []column

	for k, v := range columns{
		c := column{}
		c.pixels = v
		c.columnPosition = k
		placecolumn(picture, c)

	}

	for _, v := range picture {
		c := v.column
		order = append(order, c)
	}




	return order
}


func placecolumn(picture []column, col column){

	if(len(picture)==0){		
		picture = append(picture, col)
		return;
	}	

	rd := diff{}
	rd.dr = math.MaxFloat64
	rd.dg = math.MaxFloat64
	rd.db = math.MaxFloat64

	for i := 0; i < len(picture); i++{
		position := i
		c := picture[i]

		for j, p := range col.pixels {

			r, g, b, _ := p.RGBA() //get the rgb of the new column			
			sr, sg, sb, _ := c.pixels[j].RGBA()

			//cd.r = append(cd.r, math.Abs(float64(r-sr))+(float64(g-sg))+(float64(b-sb)))
			//fmt.Println(float64(r-sr));
			//cd.totalColumnScore += (float64(i) * math.Abs(float64(r-sr) - float64(g-sg) - float64(b-sb)))
			rd.dr += float64(r - sr)
			rd.dg += float64(g - sg)
			rd.db += float64(b - sb)

		}
		
		
		

		

	}



}

type diff struct {
	pos int
	dr float64
	dg float64
	db float64

}

type columnDiff struct {
	col []color.RGBA
	column           int
	totalColumnScore float64
	tr               float64
	tg               float64
	tb               float64
	r                []float64
	g                []float64
	b                []float64
	lowestCount      int
	tcsa             [400]float64
}

func findClosestColumn(index int, cols map[int][]color.RGBA, used map[int]int, height int) (int, map[int]int) {

	delta := make(map[int]columnDiff)

	for k := range cols {

		if k == index {
			continue
		}

		if _, ok := used[k]; ok {
			continue
		}

		scores := cols[index]

		cd := columnDiff{}
		cd.column = k
		cd.lowestCount = 0
		
		for i, p := range cols[k] {

			r, g, b, _ := p.RGBA() //compare each pixel at each y location against each
			if i > len(scores) {
				continue
			}

			sr, sg, sb, _ := scores[i].RGBA()

			cd.r = append(cd.r, math.Abs(float64(r-sr))+(float64(g-sg))+(float64(b-sb)))
			//fmt.Println(float64(r-sr));
			cd.totalColumnScore += (float64(i) * math.Abs(float64(r-sr) - float64(g-sg) - float64(b-sb)))
			cd.tr += float64(r - sr)
			cd.tg += float64(g - sg)
			cd.tb += float64(b - sb)

		}

		step := height / 400

		for i := 0; i < step; i += step {

			//, p := range cols[k]

			r, g, b, _ := cols[k][i*step].RGBA() //compare each pixel at each y location against each
			if i > len(scores) {
				continue
			}

			sr, sg, sb, _ := scores[i].RGBA()
			
			cd.tr += float64(r - sr)
			cd.tg += float64(g - sg)
			cd.tb += float64(b - sb)

			cd.tcsa[i%step] += (float64(1) * math.Abs(float64(r-sr) + float64(g-sg) + float64(b-sb)))

		}

		if _, ok := used[k]; ok {
			continue
		} else {
			delta[k] = cd
		}

	}

	 lowest := math.MaxFloat64
	lr := math.MaxFloat64
	lg := math.MaxFloat64
	lb := math.MaxFloat64

	var al [400]float64
	for i := 0; i< len(al) ; i++ {
		al[i] = math.MaxFloat64
	}

	rtn := -1
	
	rolling := 0;
	for k, v := range delta {

		if math.Abs(v.totalColumnScore) < math.Abs(lowest) {
			lowest = (math.Abs(v.totalColumnScore))
			//rtn = k
		}

		if math.Abs(v.tr) <= math.Abs(lr) && math.Abs(v.tb) <= lg && math.Abs(v.tb) <= lb {
			lr = v.tr
			lg = v.tg
			lb = v.tb
			rtn = k
		}
		count := 0
		for i := 0; i < len(al) ; i++ {
			if math.Abs(v.tcsa[i]) < math.Abs(al[i]) {
				count++
			}
		}

		if count > rolling {
			rolling = count
			rtn = k
			
		}
		if(count == len(al)){
			//break
		}

	}

	

	//rtn = findLowestDiff(delta)

	return rtn, used
}

func getAKey(m map[int]columnDiff) int {
	for k := range m {
		return k
	}
	return 0
}

func findLowestDiff(cols map[int]columnDiff) int {
	height := 0
	if len(cols) > 1 {
		height = len(cols[getAKey(cols)].r)
	} else {
		fmt.Println("cols = 0", len(cols))

	}

	for i := 0; i < height; i++ {
		lowest := math.MaxFloat64
		columnID := -1
		for id, v := range cols {

			pixelDiffSum := v.r[i]

			if pixelDiffSum < lowest {
				lowest = pixelDiffSum
				columnID = id
			}

		}
		cd := cols[columnID]
		cd.lowestCount++
	}

	score := -1
	closestColumn := -1
	for _, v := range cols {

		if v.lowestCount > score {
			score = v.lowestCount
			closestColumn = v.column

		}

	}

	return closestColumn

	// var ra [][]float64

	// //PIVOT THE COLS INTO ROWS
	// for _, v := range cols {
	// 	ra = append(ra, v.r)
	// }

	// lowestR := make(map[int]int)
	// for i, v := range ra {
	// 	tmp := math.MaxFloat64
	// 	col := -1
	// 	for ii, vv := range v {

	// 		if vv < tmp {
	// 			tmp = vv
	// 			col = ii
	// 		}
	// 	}
	// 	lowestR[i] = col
	// }

	// rtn := findMostFrequentElement(lowestR)

	//return rtn
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
