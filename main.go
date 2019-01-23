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

type C64JoyStick struct {
    path   string
    reader *bufio.Reader
    file   *os.File
}

type ButtonEvent struct {
    buttonNumber int
    isDown       bool
}

type StickEvent struct {
    isY      bool
    isCenter bool
    buttonA  bool
    buttonB  bool
}

func (j *JoyStickInputLine) SetByte(character byte) {
    if j.index < BUFLENGTH {
        j.buffer[j.index] = character
        j.index++
    }
}

func (j *JoyStickInputLine) ByteSet(i int) bool {
    return j.buffer[i] == 1
}

func (j *JoyStickInputLine) ButtonNumber() int {
    return int(j.buffer[7])
}

func (j *JoyStickInputLine) ButtonGroup() byte {
    return j.buffer[6]
}

func (j *JoyStickInputLine) StickMatches(group1 byte, group2 byte) bool {
    return j.buffer[4] == group1 && j.buffer[5] == group2
}

func (j *C64JoyStick) Read() {
    stickInputLine := JoyStickInputLine{}

    for i:=0;i<BUFLENGTH;i++ {
        character,err := j.reader.ReadByte()
        if err != nil {
            log.Fatal(err)
        }
        stickInputLine.SetByte(character)
    }

    
    if stickInputLine.ButtonGroup() == 1 {
        fmt.Println("Button pressed")
        event := ButtonEvent{
            buttonNumber: stickInputLine.ButtonNumber(),
            isDown: stickInputLine.ByteSet(4),
        }
        fmt.Println(event)
    } else if stickInputLine.ButtonGroup() == 2 {
        fmt.Println("Joystick moved")
        event := StickEvent{
            isY:      stickInputLine.ByteSet(7),
            isCenter: stickInputLine.StickMatches(0, 0),
            buttonA:  stickInputLine.StickMatches(1, 128),
            buttonB:  stickInputLine.StickMatches(255, 127),
        }
        fmt.Println(event)
    }
}

func (j *C64JoyStick) Close() {
    err := j.file.Close()
    if err != nil {
        log.Fatal(err)
    }
}

// Reads input from joystick 0
func main() {
    controller := OpenController("/dev/input/js0")

    defer controller.Close()
    for {
        controller.Read()
    }
}

func OpenController(path string) C64JoyStick {
    controller := C64JoyStick{path:path}
    file, err := os.Open(controller.path)
    if err != nil {
        log.Fatal(err)
    }

    controller.file = file
    controller.reader = bufio.NewReader(file)

    return controller
}
