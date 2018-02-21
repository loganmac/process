package main

import (
	"fmt"
	"log"

	"github.com/loganmac/process"
)

func main() {
	driver := &Summarizer{}
	processor := process.New(driver)
	driver.PrintHeader("Setting up concert")

	wrapErr(processor.Run("Preparing show", "./test-scripts/noisy-good.sh"))
	wrapErr(processor.Run("Setting up stage", "./test-scripts/good-with-warn.sh"))

	driver.PrintHeader("Let the show begin")

	wrapErr(processor.Run("Opening gates", "./test-scripts/good.sh"))
	wrapErr(processor.Run("Starting show", "./test-scripts/bad.sh"))
	wrapErr(processor.Run("Shouldn't run", "./test-scripts/good.sh"))
}

func wrapErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func clear(s string) string {
	return fmt.Sprintf("\033[2K%s", s)
}

func clearLine() {
	fmt.Print("\033[2K")
}

func printClearLine() {
	fmt.Println("\033[2K")
}

func cursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}

func cursorDown(n int) {
	fmt.Printf("\033[%dB", n)
}
