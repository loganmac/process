package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/loganmac/process"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	driver := &Summarizer{}
	processor := process.New(driver)
	driver.PrintHeader("Setting up concert")

	processor.Run("Preparing show", "./test-scripts/noisy-good.sh")
	processor.Run("Setting up stage", "./test-scripts/good-with-warn.sh")

	driver.PrintHeader("Let the show begin")

	processor.Run("Opening gates", "./test-scripts/good.sh")
	processor.Run("Starting show", "./test-scripts/bad.sh")
	processor.Run("Shouldn't run", "./test-scripts/good.sh")
}

/****************************
* SUMMARIZER PROCESS DRIVER *
*****************************/

var (
	// coloring and style
	errorFmt         = color.New(color.FgHiRed)
	successFmt       = color.New(color.FgHiGreen).Add(color.Bold)
	failureFmt       = color.New(color.FgHiRed).Add(color.Bold)
	stackFmt         = color.New(color.FgHiRed).Add(color.Faint)
	errorHeaderFmt   = color.New(color.FgHiRed).Add(color.Bold).Add(color.ReverseVideo)
	headerFmt        = color.New(color.FgHiMagenta).Add(color.Bold)
	spinnerFmt       = color.New(color.FgHiYellow).Add(color.Bold)
	spinnerPromptFmt = color.New(color.FgHiYellow).Add(color.Bold).Add(color.Underline)
	// indentation
	headerIndent = strings.Repeat(" ", 2)
	taskIndent   = strings.Repeat(" ", 4)
	outputIndent = strings.Repeat(" ", 6)
	// terminal sizing
	termWidth, termHeight, _ = terminal.GetSize(0)
	lineWidth                = termWidth - 7
	// spinner
	taskSpinner       = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	taskSpinnerLength = len(taskSpinner)
	taskSuccess       = "✓"
	taskFailure       = "✖"
)

// Summarizer is a driver that implements `process.driver`
type Summarizer struct {
	taskName        string
	stack           []string
	spinnerPosition int
	lastSpin        time.Time
}

// Initialize is called when a display driver is first attached
func (d *Summarizer) Initialize(task string) {
	d.taskName = task
	d.stack = []string{}
	d.spinnerPosition = 0
	d.lastSpin = time.Now()
	d.printSpinner()
}

// HandleNothing is called when there was no output yet
func (d *Summarizer) HandleNothing() {
	d.reprintSpinner()
	time.Sleep(50 * time.Millisecond)
}

// HandleOut is called when there is a line of output
func (d *Summarizer) HandleOut(msg string) {
	d.reprintSpinner()
	d.stack = append(d.stack, msg)
	cursorUp(1)
	if (len(msg) + 1) > lineWidth {
		msg = msg[:lineWidth]
	}
	fmt.Println(clear(outputIndent + msg))
}

// HandleErr is called when there is a line of error
func (d *Summarizer) HandleErr(msg string) {
	d.reprintSpinner()
	d.stack = append(d.stack, msg)
	cursorUp(1)
	if (len(msg) + 1) > lineWidth {
		msg = msg[:lineWidth]
	}
	errorFmt.Println(clear(outputIndent + msg))
}

// HandleSuccess is called when a process exits successfully
func (d *Summarizer) HandleSuccess() {
	cursorUp(2)
	str := fmt.Sprintf("%s%s %s", taskIndent, taskSuccess, d.taskName)
	successFmt.Println(clear(str))
}

// HandleFailure is called when a process exits with a bad exit code
func (d *Summarizer) HandleFailure() {
	cursorUp(2)
	str := fmt.Sprintf("%s%s %s", taskIndent, taskFailure, d.taskName)
	failureFmt.Println(clear(str))
	printClearLine()
	fmt.Print(clear(headerIndent))
	errStr := fmt.Sprintf("Error executing task '%s':\n", d.taskName)
	errorHeaderFmt.Println(clear(errStr))
	for _, msg := range d.stack {
		stackFmt.Println(clear(taskIndent + msg))
	}
	printClearLine()
	os.Exit(1)
}

// PrintHeader is used to print task headers for this driver
func (d *Summarizer) PrintHeader(text string) {
	printClearLine()
	headerFmt.Println(clear(headerIndent + text))
	printClearLine()
}

// Prints and updates the spinner
func (d *Summarizer) reprintSpinner() {
	cursorUp(2)
	d.printSpinner()
}

// Prints and updates the spinner
func (d *Summarizer) printSpinner() {
	fmt.Print(clear(taskIndent))
	spinnerFmt.Print(getSpinner(d.spinnerPosition) + " ")
	spinnerPromptFmt.Println(d.taskName + "\n")
	if time.Since(d.lastSpin) > (50 * time.Millisecond) {
		d.spinnerPosition++
		d.lastSpin = time.Now()
	}
}

func getSpinner(pos int) string {
	return taskSpinner[pos%taskSpinnerLength]
}

func clear(s string) string {
	return fmt.Sprintf("\033[2K%s", s)
}

func printClearLine() {
	fmt.Println("\033[2K")
}

func cursorUp(n int) {
	fmt.Printf("\033[%dA", n)
}
