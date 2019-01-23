package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

const BUFLENGTH = 8

type JoyStickInputLine struct {
    buffer  [BUFLENGTH]byte
    index   int
}

func (j *JoyStickInputLine) SetByte(character byte) {
    if j.index < BUFLENGTH {
        j.buffer[j.index] = character
        j.index++
    }
}

var reader *bufio.Reader

func read() {
    stickInputLine := JoyStickInputLine{}

    for i:=0;i<BUFLENGTH;i++ {
        character,err := reader.ReadByte()
        if err != nil {
            log.Fatal(err)
        }
        stickInputLine.SetByte(character)
    }

    fmt.Println("Got a new line")
    fmt.Println(stickInputLine.buffer)
}

// Reads input from joystick 0
func main() {
    file, err := os.Open("/dev/input/js0")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    reader = bufio.NewReader(file)
    for {
        read()
    }
}
