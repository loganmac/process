package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

/****************************
* SUMMARIZER PROCESS DRIVER *
*****************************/

var (
	// coloring and style
	errorFmt         = color.New(color.FgHiRed)
	successFmt       = color.New(color.FgHiGreen)
	failureFmt       = color.New(color.FgHiRed).Add(color.Bold)
	stackFmt         = color.New(color.FgHiRed).Add(color.Faint)
	errorHeaderFmt   = color.New(color.FgHiRed).Add(color.Bold).Add(color.ReverseVideo)
	successHeaderFmt = color.New(color.FgHiGreen).Add(color.Bold)
	headerFmt        = color.New(color.FgHiMagenta).Add(color.Bold)
	spinnerFmt       = color.New(color.FgHiYellow).Add(color.Bold)
	spinnerPromptFmt = color.New(color.FgHiYellow).Add(color.Bold).Add(color.Underline)
	// indentation
	headerIndent = strings.Repeat(" ", 2)
	taskIndent   = strings.Repeat(" ", 4)
	outputIndent = strings.Repeat(" ", 6)

	termWidth, termHeight, _ = terminal.GetSize(0)
	lineWidth                = termWidth - 7
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
	if (len(task) + 1) > lineWidth {
		task = task[:lineWidth-3] + "..."
	}
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
		msg = msg[:lineWidth-3] + "..."
	}
	fmt.Println(clear(outputIndent + msg))
}

// HandleErr is called when there is a line of error
func (d *Summarizer) HandleErr(msg string) {
	d.reprintSpinner()
	d.stack = append(d.stack, msg)
	cursorUp(1)
	if (len(msg) + 1) > lineWidth {
		msg = msg[:lineWidth-3] + "..."
	}
	errorFmt.Println(clear(outputIndent + msg))
}

// HandleSuccess is called when a process exits successfully
func (d *Summarizer) HandleSuccess() {
	cursorUp(2)
	str := fmt.Sprintf("%s%s %s", taskIndent, taskSuccess, d.taskName)
	successFmt.Println(clear(str))
	clearLine()
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

// PrintSuccess is used to print task success messages for this driver
func (d *Summarizer) PrintSuccess(text string) {
	printClearLine()
	successHeaderFmt.Println(clear(headerIndent + text))
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
