package process

import (
	"bufio"
	"log"
	"os/exec"
	"sync"
)

/******************
* PROCESS MANAGER *
*******************/

// Processor takes a driver and runs external commands
type Processor struct {
	stdout  chan string
	stderr  chan string
	success chan bool
	failure chan bool
	done    chan bool
	driver  Driver
}

// Driver is the interface that you configure the process with
// has handlers for each type of output event
type Driver interface {
	Initialize(string)
	HandleNothing()
	HandleOut(string)
	HandleErr(string)
	HandleSuccess()
	HandleFailure()
}

// New returns a new processor
func New(driver Driver) *Processor {
	stdout := make(chan string)
	stderr := make(chan string)
	success := make(chan bool)
	failure := make(chan bool)
	done := make(chan bool)
	return &Processor{
		stdout: stdout, stderr: stderr, success: success,
		failure: failure, done: done, driver: driver}
}

// Run runs an external command and does something with the errors and output
func (p *Processor) Run(task, cmdName string, cmdArgs ...string) {
	// attach the driver to respond to output
	go p.attachDriver(task)
	// run the command
	cmd := exec.Command(cmdName, cmdArgs...)

	// set stdout and stderr to a pipe
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("error creating pipe for stdout of command: %v", err)
	}
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("error creating pipe for stderr of command: %v", err)
	}

	// create scanners for each pipe
	outScanner := bufio.NewScanner(cmdOut)
	errScanner := bufio.NewScanner(cmdErr)

	// create a waitgroup for each pipe listener
	var wg sync.WaitGroup
	wg.Add(2)

	// spawn a listener for stdout
	go func() {
		for outScanner.Scan() {
			p.stdout <- outScanner.Text()
		}
		wg.Done()
	}()

	// spawn a listener for stderr
	go func() {
		for errScanner.Scan() {
			p.stderr <- errScanner.Text()
		}
		wg.Done()
	}()

	// start the command
	if err = cmd.Start(); err != nil {
		log.Fatalf("error starting command: %v", err)
	}

	// wait for the listeners to drain their scanners
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			p.failure <- true
		} else {
			log.Fatalf("error waiting on command: %v", err)
		}
	} else {
		// signal that the process completed successfully
		p.success <- true
	}
	// wait for the driver to handle all of the output
	<-p.done
}

func (p *Processor) attachDriver(task string) {
	p.driver.Initialize(task)
jump:
	for {
		select {
		case out := <-p.stdout:
			p.driver.HandleOut(out)
		case err := <-p.stderr:
			p.driver.HandleErr(err)
		case _ = <-p.success:
			p.driver.HandleSuccess()
			break jump
		case _ = <-p.failure:
			p.driver.HandleFailure()
			break jump
		default:
			p.driver.HandleNothing()
		}
	}
	p.done <- true
}

