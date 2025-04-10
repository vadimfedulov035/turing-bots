package robot

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// Instruction codes
const (
	ML = iota
	MR
	IFFLAG
	GOTO
)

// Instruction strings
const (
	MLSTRING     = "ML"
	MRSTRING     = "MR"
	IFFLAGSTRING = "IFFLAG"
	GOTOSTRING   = "GOTO"
)

// Field objects
const (
	VoidObject = iota
	RobotObject
	TreasureObject
)

// Field object strings
const (
	VoidObjectString     = "."
	RobotObjectString    = "R"
	TreasureObjectString = "T"
)

// Instruction with code and argument
type Instruction struct {
	Code     int
	Argument int
}

// Instruction implements fmt.Stringer interface
func (instruct Instruction) String() string {
	switch instruct.Code {
	case ML:
		return MLSTRING
	case MR:
		return MRSTRING
	case IFFLAG:
		return IFFLAGSTRING
	case GOTO:
		// For GOTO, include the argument (target line index)
		return fmt.Sprintf("%s %d", GOTOSTRING, instruct.Argument)
	default:
		// Handle unexpected codes gracefully
		return fmt.Sprintf("UNKNOWN(%d)", instruct.Code)
	}
}

// Common robot field with separate lines
type Field struct {
	RobotLines   [][]int
	TreasureLine []int
	LineLength   int
	mu           sync.RWMutex
}

// Field implements fmt.Stringer interface (with implicit locking)
func (f *Field) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	// Tie object with their string representation
	objectSymbols := map[int]string{
		VoidObject:     VoidObjectString,
		RobotObject:    RobotObjectString,
		TreasureObject: TreasureObjectString,
	}
	// Implement custom map()
	mapToString := func(line []int) string {
		symbols := make([]string, len(line))
		for i, obj := range line {
			symbols[i] = objectSymbols[obj]
		}
		return strings.Join(symbols, " ")
	}

	return fmt.Sprintf(
		"Robot1:   [%s]\nRobot2:   [%s]\nTreasure: [%s]\n",
		mapToString(f.RobotLines[0]),
		mapToString(f.RobotLines[1]),
		mapToString(f.TreasureLine),
	)
}

// Robot on the field with instructions to follow
type Robot struct {
	Num           int
	Field         *Field
	FieldPosition int
	Instructions  []Instruction
	SleepInterval time.Duration
}

// Get instructions from file with passed filename
func getInstructions(filename string) []Instruction {
	const parseNumWarning string = "%d line: Failed to parse line number '%s'"
	const readInstructWarning string = "%d line: Got unknown instruction '%s'"
	codes := map[string]int{
		MLSTRING:     ML,
		MRSTRING:     MR,
		IFFLAGSTRING: IFFLAG,
		GOTOSTRING:   GOTO,
	}

	// Get lines to read
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to open and read '%s' file: %v\n", filename, err)
	}
	lines := strings.Split(string(content), "\n")

	// Iterate over lines translating them to instructions
	instructs := make([]Instruction, 0, len(lines))
	for i, line := range lines {
		// Check if comment (pass comment)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Delete comment
		line = strings.Split(line, "#")[0]

		// Trim line spaces, check if empty (pass empty line)
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		// Translate line to instruction
		var instruct Instruction
		code, ok := codes[trimmedLine]

		if ok { // ML, MR, IFFLAG
			instruct.Code = code
			instruct.Argument = 0 // Blank
		} else if strings.HasPrefix(trimmedLine, GOTOSTRING) { // Handle GOTO
			var targetLineNum int
			_, err := fmt.Sscanf(
				strings.TrimSpace(trimmedLine[len(GOTOSTRING):]),
				"%d", &targetLineNum,
			)
			if err != nil || targetLineNum < 0 { // Ensure target is positive
				log.Fatalf(parseNumWarning, i, trimmedLine)
			}
			instruct.Code = GOTO
			instruct.Argument = targetLineNum
		} else { // Unknown instruction
			log.Fatalf(readInstructWarning, i, trimmedLine)
		}

		// Append translated instruction
		instructs = append(instructs, instruct)
	}

	log.Println(instructs)
	return instructs
}

// Generate positions for robots and treasure
func getPositions(lineLength int) (i1 int, i2 int, j int) {
	// Check length constraint
	if lineLength < 3 {
		log.Fatalf("\"Endless\" line should be at least 3 blocks long!")
	}

	// Generate different and *not* adjacent robot positions 'i'
	for {
		i1 = rand.Intn(lineLength)
		i2 = rand.Intn(lineLength)
		// Get absolute difference
		diff := i1 - i2
		if diff < 0 {
			diff = -diff
		}
		// Ensure difference suffices the constraint
		if diff > 1 {
			break
		}
	}

	// Determine min and max positions
	minPos, maxPos := i1, i2
	if i1 > i2 {
		minPos, maxPos = i2, i1
	}

	// Generate treasure position 'j' between minPos and maxPos
	numSpotsBetween := maxPos - minPos - 1 // <- (maxPos - 1) - (minPos + 1) + 1
	j = rand.Intn(numSpotsBetween) +       // [0, numSpotsBetween)
		minPos + 1 // shift to (minPos, maxPos)

	return
}

// Create two robots on the same field with a dot between them
func NewRobots(
	lineLength int, filename string, sleepInterval time.Duration,
) (*Robot, *Robot) {
	// Get robots and treasure positions
	i1, i2, j := getPositions(lineLength)

	// Create common robot field with separate lines
	field := &Field{
		RobotLines:   make([][]int, 2),
		TreasureLine: make([]int, lineLength),
		LineLength:   lineLength,
	}
	field.RobotLines[0] = make([]int, lineLength)
	field.RobotLines[1] = make([]int, lineLength)

	// Set robots
	instructions := getInstructions(filename)
	robot1 := &Robot{
		Num:           0,
		Field:         field,
		FieldPosition: i1,
		Instructions:  instructions,
		SleepInterval: sleepInterval,
	}
	robot1.Field.RobotLines[0][robot1.FieldPosition] = RobotObject
	robot2 := &Robot{
		Num:           1,
		Field:         field,
		FieldPosition: i2,
		Instructions:  instructions,
		SleepInterval: sleepInterval,
	}
	robot2.Field.RobotLines[1][robot2.FieldPosition] = RobotObject

	// Set treasure
	field.TreasureLine[j] = TreasureObject

	log.Printf("Robots at %d and %d, Treasure at %d\n", i1, i2, j)
	log.Printf("Initial Field State:\n%s\n", field)
	return robot1, robot2
}

// Move one step left on the field, print field
func (r *Robot) MoveLeft(i int) int {
	time.Sleep(r.SleepInterval)

	// Simulate endless line with wrapping
	r.Field.mu.Lock()
	r.Field.RobotLines[r.Num][r.FieldPosition] = VoidObject
	if r.FieldPosition > 0 {
		r.FieldPosition = r.FieldPosition - 1 // Move left
	} else {
		r.FieldPosition = r.Field.LineLength - 1 // Wrap around
	}
	r.Field.RobotLines[r.Num][r.FieldPosition] = RobotObject
	r.Field.mu.Unlock()

	log.Printf("[Robot %d] %d line: MoveLeft\n", r.Num+1, i)

	return i + 1
}

// Move one step right on the field, print field
func (r *Robot) MoveRight(i int) int {
	time.Sleep(r.SleepInterval)

	// Simulate endless line with wrapping
	r.Field.mu.Lock()
	r.Field.RobotLines[r.Num][r.FieldPosition] = VoidObject
	r.FieldPosition = (r.FieldPosition + 1) % r.Field.LineLength
	r.Field.RobotLines[r.Num][r.FieldPosition] = RobotObject
	r.Field.mu.Unlock()

	log.Printf("[Robot %d] %d line: MoveRight\n", r.Num+1, i)

	return i + 1
}

// Check if treasure below current position, if no jump over instruction
func (r *Robot) IfFlag(i int) int {
	time.Sleep(r.SleepInterval)

	// Check current position fast
	r.Field.mu.RLock()
	onTreasure := r.Field.TreasureLine[r.FieldPosition] == TreasureObject
	r.Field.mu.RUnlock()

	log.Printf("[Robot %d] %d line: IfFlag: ", r.Num+1, i)
	if onTreasure {
		log.Printf("%s -> TRUE (Treasure found!) -> Executing next line", IFFLAGSTRING)
		i += 1
	} else { // If no treasure, skip the *next* instruction
		log.Printf("%s -> FALSE (No treasure) -> Skipping next line", IFFLAGSTRING)
		i += 2
	}

	return i
}

// Return index to go to and log it
func (r *Robot) GoTo(i, lineToGoTo int) int {
	log.Printf("[Robot %d] %d line: Go to %d\n", r.Num+1, i, lineToGoTo)
	return lineToGoTo
}

// Perform i-th instruction
func (r *Robot) Perform(i int) int {
	// Set current instruction
	instruct := r.Instructions[i]

	// Call method based on the code
	switch instruct.Code {
	case ML:
		i = r.MoveLeft(i)
	case MR:
		i = r.MoveRight(i)
	case IFFLAG: // Branching based on success by resetting i
		i = r.IfFlag(i)
	case GOTO: // Recursive Perform with GoTo argument (instruction number)
		i = r.GoTo(i, instruct.Argument)
	}

	// Print field (.String() method locks mutex for reading)
	fmt.Printf("[Robot %d] Field:\n%s\n", r.Num+1, r.Field)
	return i
}

func (r *Robot) Run(wg *sync.WaitGroup, stopCh chan struct{}) {
	defer wg.Done()

	// Loop until all instructions are performed or signal received
	for i := 1; i >= 0 && i < len(r.Instructions); {
		select {
		case <-stopCh:
			log.Printf("Robot %d: Received stop signal.\n", r.Num+1)
			return
		default:
			i = r.Perform(i)
			// 0-th instruction indicates treasure was found
			if i == 0 {
				log.Printf("Robot %d: Found the treasure!\n", r.Num+1)
				return
			}
		}
	}
	log.Printf("Robot %d: Loop escape.\n", r.Num+1)
}
