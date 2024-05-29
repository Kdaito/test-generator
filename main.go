package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Kdaito/test-generator/internal/proceccer/jest"
	"github.com/fatih/color"
)

var src string
var targ string
var typ string
var h string

func init() {
	flag.StringVar(&src, "src", "", "Directory path containing the test cases source to be implemented")
	flag.StringVar(&targ, "targ", "", "Path of the target directory where the test is to be implemented")
	flag.StringVar(&typ, "typ", "", "How to implement test cases")
	flag.StringVar(&h, "", "", "")
}

// 実行例
// go run main.go -src=./test-case -targ=./sample-app/src/services -typ=jest

func main() {
	flag.Parse()

	if (src == "") {
		color.Red("Source directory is undefined. Please select it by -src flag.")
		os.Exit(1)
	}

	if (targ == "") {
		color.Red("Target directory is undefined. Please select it by -targ flag.")
		os.Exit(1)
	}

	if (typ == "") {
		color.Red("The test implementation library is not specified. Please select it by -typ flag.")
		os.Exit(1)
	}

	fmt.Println("Test files are generating...")

	if (typ == "jest") {
		// For Jest
		count, err := jest.Generate(targ, src)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		fmt.Printf("%d test files have been created.", count)
	}

	color.Green("The test file was successfully generated!")
	os.Exit(0)
}
