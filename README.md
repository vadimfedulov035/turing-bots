# Turing Bots - Robot Meeting Problem

A classic computer science problem solved in Golang.

## The Problem

Imagine an infinitely long strip of white cells, like a tape.
*   Two robots (Robot 1 and Robot 2) are placed on different cells.
*   Exactly one cell *between* the robots is black (the "treasure"). All other cells are white.
*   The robots don't know their initial positions or the distance to the black cell.
*   You need to write a *single program* that, when loaded onto *both* robots, will guarantee they eventually meet at the black cell.

## Allowed Commands

The robot programs consist of a sequence of instructions. Each robot executes the program line by line.

1.  `ML`: Move one cell to the Left. Takes 1 second. Proceeds to the next line.
2.  `MR`: Move one cell to the Right. Takes 1 second. Proceeds to the next line.
3.  `IF FLAG`: Check if the robot is currently on the black cell ("found treasure").
    *   If YES (on the black cell): Proceed to the *next* line. Takes 1 second.
    *   If NO (on a white cell): Skip the next line and proceed to the line *after* that. Takes 1 second.
4.  `GOTO N`: Jump *immediately* to line number `N` (0-indexed based on the actual instructions loaded, ignoring comments/blank lines). Takes 0 seconds.

## The Solution: Expanding Search

Since the robots must use the same program and don't know where the treasure is, they need a systematic way to search. The implemented solution uses an expanding search pattern:

1.  **Check Current Location:** First, check if the robot starts on the "treasure" (ensure "treasure" is found even if generation went wrong).
2.  **Expand Search:** The robot alternates directions (Right, Left, then Right...) and increases the number of steps taken in each direction (1 step, then 2 steps, then 3 steps, and so on).
3.  **Why it Works:** Because the treasure is between the robots, this expanding search pattern guarantees that one of the robots will eventually step onto the treasure cell. Since both robots run the same program, they will both find it.
4.  **Stopping:** When `IF FLAG` detects the treasure, the next instruction executed is `GOTO 0`. Line 0 contains another `GOTO 0`, creating an infinite waiting loop interpreted as a success, as 0 line is not accessible if condition was not met.

The `codegen.go` component automatically generates a program file (`bot.prog` by default) implementing this expanding search algorithm.

## Project Structure

*   `main.go`: The main application entry point. Handles command-line arguments, initializes the simulation environment, creates the robots, starts their execution in separate goroutines, and manages graceful shutdown.
*   `robot/robots.go`: Defines the `Robot`, `Field`, and `Instruction` types. Contains the logic for robot movement (`MoveLeft`, `MoveRight`), condition checking (`IfFlag`), jumps (`GoTo`), instruction parsing from the program file, and the core `Run` loop for robot execution. It also simulates the shared "infinite" line (using wrapping on a finite array).
*   `codegen/codegen.go`: Contains the logic to automatically generate the robot program file (`bot.prog`) based on the expanding search algorithm.

## Prerequisites

*   Go compiler (version 1.22 or later recommended).

## How to Build

Navigate to the project's root directory in your terminal and run:

```bash
go build .
```

This will create an executable file (e.g., `turing-bots` on Linux/macOS or `turing-bots.exe` on Windows).

## How to Run

Execute the compiled program from your terminal. You can customize the simulation using command-line flags:

```bash
./turing-bots [flags]
```

### Command-Line Arguments

*   `-l <int>`: Specifies the length of the simulated line. The "infinite" line wraps around this length. Must be at least 3.
    *   Default: `10`
*   `-g <string>`: The filename for the generated robot program code.
    *   Default: `bot.prog`
*   `-s <duration>`: The sleep interval between each robot action (e.g., `1s`, `500ms`, `100ns`). This controls the simulation speed.
    *   Default: `1s` (1000ms)

### Examples

*   Run with default settings:
    ```bash
    ./turing-bots
    ```
*   Run with a line length of 20 and a faster simulation speed:
    ```bash
    ./turing-bots -l 20 -s 250ms
    ```
*   Generate the program to a different file:
    ```bash
    ./turing-bots -g my_robot_program.txt
    ```

## Simulation Output

The program will:
1.  Log the randomly chosen initial positions for Robot 1, Robot 2, and the Treasure.
2.  Log the initial state of the field.
3.  Print the state of the field after each robot takes an action (ML, MR, IFFLAG...). Each action is logged with the robot number and the instruction line executed.
4.  Indicate when a robot finds the treasure (`IF FLAG -> TRUE`).
5.  The simulation ends when both robots have found the treasure and entered their `GOTO 0` loop, or if you interrupt it (e.g., with Ctrl+C).

Example Field Output:

```
[Robot 1] Field:
Robot1:   [. . . R . . . . . .]
Robot2:   [. . . . . . . R . .]
Treasure: [. . . . . T . . . .]

[Robot 2] Field:
Robot1:   [. . . R . . . . . .]
Robot2:   [. . . . . . R . . .]
Treasure: [. . . . . T . . . .]

... and so on ...
```

## License

---

[![License](https://img.shields.io/badge/license-GPLv3-blue.svg)](#)  
**License**: GNU General Public License v3.0  
