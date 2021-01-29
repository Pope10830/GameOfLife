package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// writePgmImage receives an array of bytes and writes it to a pgm file.
// Note that this function is incomplete. Use the commented-out for loop to receive data from the distributor.
func writePgmImage(p golParams, i ioChans) {
	_ = os.Mkdir("out", os.ModePerm)

	filename := <-i.distributor.filename
	file, ioError := os.Create("out/" + filename + ".pgm")
	check(ioError)
	defer file.Close()

	_, _ = file.WriteString("P5\n")
	//_, _ = file.WriteString("# PGM file writer by pnmmodules (https://github.com/owainkenwayucl/pnmmodules).\n")
	_, _ = file.WriteString(strconv.Itoa(p.imageWidth))
	_, _ = file.WriteString(" ")
	_, _ = file.WriteString(strconv.Itoa(p.imageHeight))
	_, _ = file.WriteString("\n")
	_, _ = file.WriteString(strconv.Itoa(255))
	_, _ = file.WriteString("\n")

	world := make([][]byte, p.imageHeight)
	for i := range world {
		world[i] = make([]byte, p.imageWidth)
	}

	// TODO: write a for-loop to receive the world from the distributor when outputting.
	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			world[y][x] = <- i.distributor.outputVal
		}
	}

	for y := 0; y < p.imageHeight; y++ {
		for x := 0; x < p.imageWidth; x++ {
			_, ioError = file.Write([]byte{world[y][x]})
			check(ioError)
		}
	}

	ioError = file.Sync()
	check(ioError)

	fmt.Println("File", filename, "output done!")
}

// readPgmImage opens a pgm file and sends its data as an array of bytes.
func readPgmImage(p golParams, i ioChans) {
	filename := <-i.distributor.filename
	data, ioError := ioutil.ReadFile("images/" + filename + ".pgm")
	check(ioError)

	fields := strings.Fields(string(data))

	if fields[0] != "P5" {
		panic("Not a pgm file")
	}

	width, _ := strconv.Atoi(fields[1])
	if width != p.imageWidth {
		panic("Incorrect width")
	}

	height, _ := strconv.Atoi(fields[2])
	if height != p.imageHeight {
		panic("Incorrect height")
	}

	maxval, _ := strconv.Atoi(fields[3])
	if maxval != 255 {
		panic("Incorrect maxval/bit depth")
	}

	image := []byte(fields[4])

	for _, b := range image {
		i.distributor.inputVal <- b
	}

	fmt.Println("File", filename, "input done!")
}

func pgmIo(p golParams, i ioChans) {
	for {
		select {
		case command := <-i.distributor.command:
			switch command {
			case ioInput:
				readPgmImage(p, i)
			case ioOutput:
				writePgmImage(p, i)
			case ioCheckIdle:
				i.distributor.idle <- true
			}
		}
	}
}
