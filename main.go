package main

import (
	"flag"
	"fmt"
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

	//get picture

	//get output

	if *a.shred {
		fmt.Println("test bool")
		os.Exit(0)
	}

	unShred()
	
}

func getPicture(){

}

func defineFlags() Args {
	a := Args{}
	a.input = flag.String("i", "", "Input file to be unshred")
	a.output = flag.String("o", "./out", "outputDir")
	a.shred = flag.Bool("s", false, "Yes/No : Shred the image, default behavior is to attempt to unshred it")
	
	return a
}
