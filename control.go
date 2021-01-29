package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

// getKeyboardCommand sends all keys pressed on the keyboard as runes (characters) on the key chan.
// getKeyboardCommand will NOT work if termbox isn't initialised (in startControlServer)
func getKeyboardCommand(key chan<- rune) {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyCtrlC {
				StopControlServer()
				os.Exit(0)
			} else if key != nil {
				if event.Key != 0 {
					key <- rune(event.Key)
				} else if event.Ch != 0 {
					key <- event.Ch
				}
			}
		}
	}
}

// startControlServer initialises termbox and prints basic information about the game configuration.
func startControlServer(p golParams) {
	e := termbox.Init()
	check(e)

	fmt.Println("Threads:", p.threads)
	fmt.Println("Width:", p.imageWidth)
	fmt.Println("Height:", p.imageHeight)
}

// stopControlServer closes termbox.
// If the program is terminated without closing termbox the terminal window may misbehave.
func StopControlServer() {
	termbox.Close()
}
