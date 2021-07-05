package main

import (
	"flag"
	"fmt"
	"image"

	_ "image/png"
	"os"
)

//Args - input command line args
type Args struct {
	input  *string
	output *string
	shred  *bool
}

func main() {

	fmt.Println("helloworld")
	a := defineFlags()
	flag.Parse()

	//get output

	//get picture

	image, t := getPicture(*a.input)

	if *a.shred {
		shred(image, t, *a.output)
		os.Exit(0)
	}

	unShred(image, t)

}

func getPicture(input string) (image.Image, string) {
	fmt.Println(input)
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Input file doesn't work, or can't be found")
		os.Exit(9)
	}
	defer file.Close()

	image, t, err := image.Decode(file)
	if err != nil {
		fmt.Println("error decoding image file")
		fmt.Println(err)
		os.Exit(9)

	}
	fmt.Println("image type {}", t)
	return image, t
}

func defineFlags() Args {
	a := Args{}
	a.input = flag.String("i", "", "Input file to be unshred")
	a.output = flag.String("o", "./output.png", "Output path and filename, defaults to output.png")
	a.shred = flag.Bool("s", true, "Yes/No : Shred the image, default behavior is to attempt to unshred it")

	return a
}
