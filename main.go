package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vadimfedulov035/turing-bots/codegen"
	"github.com/vadimfedulov035/turing-bots/robot"
)

// Default values for arguments
const (
	LineLength      = 10
	ProgramFilename = "bot.prog"
	SleepInterval   = time.Second
)

// Argument variables
var (
	lineLength      int
	programFilename string
	sleepInterval   time.Duration
)

// Load passed arguments
func init() {
	const (
		lineLenNote         = "Length of \"endless\" line (min 3, default=10)"
		programFilenameNote = "Program filename (default='bot.prog')"
		sleepIntervalNote   = "Interval between robot operations (default=1s)"
	)
	flag.IntVar(&lineLength, "l", LineLength, lineLenNote)
	flag.StringVar(&programFilename, "g", ProgramFilename, programFilenameNote)
	flag.DurationVar(&sleepInterval, "s", SleepInterval, sleepIntervalNote)
}

func main() {
	flag.Parse()

	// Validate line length early
	if lineLength < 3 {
		const ShortLineMessage = "Error: Line length must be at least 3. Got %d.\n"
		fmt.Fprintf(os.Stderr, ShortLineMessage, lineLength)
		os.Exit(1)
	}

	// Generate universal program for robots
	codegen.GenerateProgram(programFilename, lineLength)

	// Create new robots
	robot1, robot2 := robot.NewRobots(
		lineLength, programFilename, sleepInterval,
	)

	// Channels and WaitGroup
	stopCh := make(chan struct{})     // Channel to signal termination
	completeCh := make(chan struct{}) // Channel to signal completion
	wg := &sync.WaitGroup{}

	// Setup OS signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start robots and wait for execution finish
	wg.Add(2)
	go robot1.Run(wg, stopCh)
	go robot2.Run(wg, stopCh)

	// Goroutine signal grecufull completion in any case
	go func() {
		wg.Wait()
		close(completeCh)
	}()

	// Block until completion or OS signal
	select {
	case <-completeCh:
		fmt.Println("All robots completed their tasks normally.")
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal %v, initiating graceful shutdown...\n", sig)
		// Signal robots to stop
		close(stopCh)
		// Wait for wg.Wait() goroutine
		<-completeCh
		fmt.Println("Shutdown completed gracefully.")
	}
}
