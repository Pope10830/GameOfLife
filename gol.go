package main

import (
	"fmt"
	"strconv"
	"strings"
	"os"
	"time"
)

// Returns the number of live neighbours of a given cell
func getNumberOfLiveNeighbours(xCoord int, yCoord int, world [][]byte, p golParams) int {
	var numberOfLiveNeighbours int = 0
	var checkY, checkX int
	for y := -1; y < 2; y++ {
		for x := -1; x < 2; x++ {
			if (!(y == 0 && x == 0)) {
				// Loop Y To Other Edge
				if (yCoord + y < 0) {
					checkY = yCoord + y + p.imageHeight
				} else if (yCoord + y >= p.imageHeight) {
					checkY = yCoord + y - p.imageHeight
				} else {
					checkY = yCoord + y
				}

				// Loop X To Other Edge
				if (xCoord + x < 0) {
					checkX = xCoord + x + p.imageWidth
				} else if (xCoord + x >= p.imageWidth) {
					checkX = xCoord + x - p.imageWidth
				} else {
					checkX = xCoord + x
				}

				if (world[checkY][checkX] == 0xFF) {
					numberOfLiveNeighbours++
				}
			}
		}
	}

	return numberOfLiveNeighbours
}

// Returns true if cell should be switched between alive and dead
func toSwitch(x int, y int, world [][]byte, p golParams) bool {
	if (world[y][x] == 0xFF) {
		if (getNumberOfLiveNeighbours(x, y, world, p) < 2) {
			return true
		} else if (getNumberOfLiveNeighbours(x, y, world, p) > 3) {
			return true
		} else if (getNumberOfLiveNeighbours(x, y, world, p) == 2 || getNumberOfLiveNeighbours(x, y, world, p) == 3) {
			return false
		} else {
			return false
		}
	} else {
		if (getNumberOfLiveNeighbours(x, y, world, p) == 3) {
			return true
		} else {
			return false
		}
	}
}

// Calculates the new value of a given cell
func calculateCell(world [][]byte, cellVar cell, p golParams) uint8 {
	if (toSwitch(cellVar.x, cellVar.y, world, p)) {
		return world[cellVar.y][cellVar.x] ^ 0xFF
	} else {
		return world[cellVar.y][cellVar.x]
	}
}

// Main function of go routine workers
func workerFunction(world [][]byte, c cellChannels, p golParams, i int) {
	var resultVar resultType

	resultVar.id = i
	c.idle <- i // Send the id of the thread to the idle channel initially

	for {
		resultVar.cell = <- c.coordinates // Read coordinates of cell to be calculated
		resultVar.cellValue = calculateCell(world, resultVar.cell, p)
		c.result <- resultVar	// Send the result back to the main thread
		c.idle <- i	// Send the id of this thread to the idle channel when done
	}
}

func outputToPgm(p golParams, d distributorChans, world [][]byte, turn int) {
	var filename string = ("File " + strconv.Itoa(p.imageWidth) + "x" + strconv.Itoa(p.imageHeight) + " - " + strconv.Itoa(turn))

	d.io.command <- ioOutput
	d.io.filename <- filename

	for y:=0; y < p.imageHeight; y++ {
		for x:=0;x<p.imageWidth;x++ {
			d.io.outputVal <- world[y][x]
		}
	}

	// Make sure that the Io has finished any output before exiting.
	d.io.command <- ioCheckIdle
	<-d.io.idle
}

func getLiveCells(p golParams, world [][]byte) []cell {
	var aliveCells []cell
	// Go through the world and append the cells that are still alive.
	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			if world[y][x] != 0 {
				aliveCells = append(aliveCells, cell{x: x, y: y})
			}
		}
	}

	return aliveCells
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p golParams, d distributorChans, alive chan []cell, keyChan <-chan rune) {
	// Create the 2D slice to store the world.
	world := make([][]byte, p.imageHeight)
	oldWorld := make([][]byte, p.imageHeight)
	for i := range world {
		world[i] = make([]byte, p.imageWidth)
		oldWorld[i] = make([]byte, p.imageWidth)
	}

	// Request the io goroutine to read in the image with the given filename.
	d.io.command <- ioInput
	d.io.filename <- strings.Join([]string{strconv.Itoa(p.imageWidth), strconv.Itoa(p.imageHeight)}, "x")

	// The io goroutine sends the requested image byte by byte, in rows.
	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			val := <-d.io.inputVal
			if val != 0 {
				fmt.Println("Alive cell at", x, y)
				world[y][x] = val
			}
		}
	}

	// Create the channels
	var cellChans cellChannels
	cellChans.coordinates = make(chan cell, p.threads/2)
	cellChans.result = make(chan resultType, p.threads/2)
	cellChans.idle = make(chan int, p.threads/2)

	// Create p.thread threads
	for i := 0; i < p.threads; i++ {
		go workerFunction(oldWorld, cellChans, p, i)
	}

	ticker := time.NewTicker(2000 * time.Millisecond)

	// Calculate the new state of Game of Life after the given number of turns.
	for turns := 0; turns < p.turns; turns++ {
		// Create a copy of the world at the start of the turn
		for y := 0; y < p.imageHeight; y++ {
			for x:= 0; x < p.imageWidth; x++ {
				oldWorld[y][x] = world[y][x]
			}
		}

		// Keeps track of the number of results sent back to main thread
		var resultsRecieved int = 0

		// Loop through each cell
		for y := 0; y < p.imageHeight; y++ {
			for x := 0; x < p.imageWidth; x++ {
				var cellVar cell
				cellVar.x = x
				cellVar.y = y

				// Loop through reading channels until the coordinates are sent
				// If a result is read we still want to wait for the coordinates
				Loop:
					for {
						select {
							case keyboardInput := <- keyChan: //case there is something on the rune channel
								if keyboardInput == rune('s') {
									outputToPgm(p, d, oldWorld, turns)
								} else if keyboardInput == rune('q') {
									outputToPgm(p, d, oldWorld, turns)
									fmt.Println("Program Terminating")
								 	os.Exit(0)
								} else if keyboardInput == rune('p') {
									fmt.Println("Program Paused. Turn: " + strconv.Itoa(turns))
									for {
										key := <- keyChan
										if key == rune('p') {
											fmt.Println("Continuing")
											break
										}
									}
								}
							case <- cellChans.idle:
								cellChans.coordinates <- cellVar
								break Loop // Break For Loop
							case result := <- cellChans.result:
								world[result.cell.y][result.cell.x] = result.cellValue
								resultsRecieved++
							case <- ticker.C:
								fmt.Println("Number Of Alive Cells: " + strconv.Itoa(len(getLiveCells(p, oldWorld))))
						}
					}
			}
		}

		// Check all threads have returned results
		for {
			if resultsRecieved == p.imageHeight * p.imageWidth {
				break
			} else {
				result := <- cellChans.result
				world[result.cell.y][result.cell.x] = result.cellValue
				resultsRecieved++
			}
		}

	}

	// Create an empty slice to store coordinates of cells that are still alive after p.turns are done.
	var finalAlive []cell
	finalAlive = getLiveCells(p, world)

	outputToPgm(p, d, world, p.turns)

	// Return the coordinates of cells that are still alive.
	alive <- finalAlive
}
