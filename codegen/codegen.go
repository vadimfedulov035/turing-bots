package codegen

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/vadimfedulov035/turing-bots/robot"
)

const Title = `# Widening ZigZag
# Goal: Find hidden treasure on the field
# Algorithm: move left and right in zigzag checking places; loop forever on successs

`
const CodeWidth = 48

// Write line: code + comment
func writeLine(
	writer io.Writer,
	lineNum int,
	instruction string,
	argument string,
	explanation string,
) {
	// code: instruction + argument
	code := instruction
	if argument != "" {
		code += " " + argument
	}
	code = fmt.Sprintf("%-*s", CodeWidth, code)

	// comment: line number + explanation
	comment := fmt.Sprintf("# %d", lineNum)
	if comment != "" {
		comment += ": " + explanation
	}

	// Write down line
	line := code + comment + "\n"
	writer.Write([]byte(line))
}

// Write title describing the program
func writeTitle(writer io.Writer) {
	writer.Write([]byte(Title))
}

// Write initial loop (start line incrementing)
func writeLoop(writer io.Writer) (lineNum int) {
	writeLine(writer, lineNum, robot.GOTOSTRING, "0", "Endless waiting loop")
	return lineNum + 1
}

// Write header: step comment + move comment
func writeHeader(writer io.Writer, i int, moveIdx int) {
	directions := [2]string{"left", "right"}

	// Build move comment
	var moveComment string
	switch i {
	case 0:
		moveComment = "No movement"
	case 1:
		moveComment = "1 move to " + directions[moveIdx]
	default:
		moveComment = strconv.Itoa(i) + " moves to " + directions[moveIdx]
	}
	// Build step comment
	stepComment := fmt.Sprintf("# -- Step %d: ", i)

	// Write down header
	writer.Write([]byte(stepComment + moveComment + "\n"))
}

// Write single move command (ML or MR) via writeLine()
func writeMove(writer io.Writer, lineNumber int, moveIdx int) int {
	moves := [2]string{robot.MLSTRING, robot.MRSTRING}
	writeLine(writer, lineNumber, moves[moveIdx], "", "")
	return lineNumber + 1
}

// Write base logic check (IFFLAG, GOTO, GOTO) via writeLine()
func writeBase(writer io.Writer, lineNum int) int {
	failTarget := strconv.Itoa(lineNum + 3)

	// IFFLAG
	writeLine(writer, lineNum, robot.IFFLAGSTRING, "", "If found treasure")
	lineNum++

	// GOTO Success
	writeLine(
		writer, lineNum, robot.GOTOSTRING, "0", "Success -> Waiting loop",
	)
	lineNum++

	// GOTO Fail
	writeLine(
		writer, lineNum, robot.GOTOSTRING, failTarget, "Fail -> Next line",
	)
	lineNum++

	// Empty line for readability
	writer.Write([]byte("\n"))

	return lineNum
}

// Write program: title + loop + for(header + moves + base)
func writeProgram(writer io.Writer, lineLength int) {
	writeTitle(writer)

	// Write first code line, set line counter to 1
	lineCount := writeLoop(writer)

	// Iterate up to line length
	for i := range lineLength {
		moveIdx := i % 2                // Even: ML, Odd: MR
		writeHeader(writer, i, moveIdx) // Header is only a comment

		// i step: i movements
		for range i {
			lineCount = writeMove(writer, lineCount, moveIdx)
		}

		// 1 step: 1 base (IFFLAG/GOTO/GOTO block)
		lineCount = writeBase(writer, lineCount)
	}
}

// Create program file and initiates the program writing process.
func GenerateProgram(filename string, lineLength int) {
	// Open file with passed filename
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create '%s': %v", filename, err)
	}
	defer file.Close()

	// Write program
	writeProgram(file, lineLength)

	log.Printf(
		"Generated program with %d steps into '%s'\n", lineLength, filename,
	)
}
