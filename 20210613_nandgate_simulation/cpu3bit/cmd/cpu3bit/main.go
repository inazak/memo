package main

import (
  "flag"
  "fmt"
  "os"
  "bufio"
  "time"
  "github.com/nsf/termbox-go"
  "github.com/inazak/cpu3bit"
)

const (
  fgColor = termbox.ColorWhite
  bgColor = termbox.ColorBlack
  fgEmColor = termbox.ColorBlack
  bgEmColor = termbox.ColorWhite
)

var display = []string{
        //01234567890123456789012345678901234567890123456789
/*  0 */ "                3bit CPU Demo                     ",
/*  1 */ "                                                  ",
/*  2 */ " Register A [####]       Address    Memory        ",
/*  3 */ "                               0    [#######]     ",
/*  4 */ " Register B [####]             1    [#######]     ",
/*  5 */ "                               2    [#######]     ",
/*  6 */ " Carry Flag [#]                3    [#######]     ",
/*  7 */ "                               4    [#######]     ",
/*  8 */ " Program Counter               5    [#######]     ",
/*  9 */ " [###]                         6    [#######]     ",
/* 10 */ "                               7    [#######]     ",
/* 11 */ " q: quit                          bit6 ... 0      ",
/* 12 */ " t: tick(manual mode only)                        ",
/* 13 */ "                                                  ",
/* 14 */ "  ########## 3bit cpu Instructions ##########     ",
/* 15 */ "                                                  ",
/* 16 */ "  6  5  4  3-0  carry  | mnemonic                 ",
/* 17 */ "  ---------------------|----------------          ",
/* 18 */ "  0  1  0  Imd    x    | ADD A, Imd               ",
/* 19 */ "  0  1  1  Imd    x    | MOV A, Imd               ",
/* 20 */ "  1  0  0  Imd    x    | ADD B, Imd               ",
/* 21 */ "  1  0  1  Imd    x    | MOV B, Imd               ",
/* 22 */ "  1  1  0  Imd    0    | JNC Imd (jump if carry=0)",
/* 23 */ "  1  1  1  Imd    x    | JMP Imd                  ",
/* 24 */ "  0  0  0  Imd    x    | NOP                      ",
/* 25 */ "  0  0  1  Imd    x    | NOP                      ",
}

// sample image is 3*3=9
var memoryimage = [][]int{
  //bit0 .. bit6 
  {0,0,0,0,1,0,1}, //0: mov b,0
  {1,0,1,1,1,1,0}, //1: mov a,13
  {1,1,0,0,0,0,1}, //2: add b,3
  {1,0,0,0,0,1,0}, //3: add a,1
  {0,1,0,0,0,1,1}, //4: jnc 2
  {1,0,1,0,1,1,1}, //5: jmp 5
  {0,0,0,0,0,0,0}, // Register B is 1001 (9)
  {0,0,0,0,0,0,0},
}

var load   = flag.String("load", "", "textfile representing memory image")
var manual = flag.Bool("manual", false, "ticking by hand")

var info *cpu3bit.CPUInfo

func main() {
  flag.Parse()

  // load text file
  if *load != "" {
    var err error
    memoryimage, err = loadMemoryImageText(*load)
    if err != nil {
      fmt.Printf("%v", err)
      return
    }
  }

  // initialize
  cpu3bit.Initialize()
  info = cpu3bit.MakeComputer(memoryimage)

  err := termbox.Init()
  if err != nil {
    panic(err)
  }
  defer termbox.Close()

  eventQueue := make(chan termbox.Event)
  go func(){
    for {
      eventQueue <- termbox.PollEvent()
    }
  }()

  render()

  //auto ticking
  if ! *manual {
    go func(){
      for {
        select {
        case <- time.After(time.Millisecond * 1000):
          cpu3bit.TickTock()
          render()
        }
      }
    }()
  }

  for {
    select {
    case ev := <-eventQueue:
      if ev.Type == termbox.EventKey {
        switch {
        case ev.Ch == 't':
          if *manual {
            cpu3bit.TickTock() // clockdown/up and update
            render()
          }
        case ev.Ch == 'q' || ev.Key == termbox.KeyEsc:
          return
        }
      }
    }
  }
}

func render() {
  termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
  //title
  setText(0, 0, fgEmColor, bgEmColor, display[0])
  //other
  for i:=1; i<len(display); i++ {
    setText(0, i, fgColor, bgColor, display[i])
  }
  //text update
  setBinaryText(13, 2, ToRunes(cpu3bit.ToString(info.RegisterA())))
  setBinaryText(13, 4, ToRunes(cpu3bit.ToString(info.RegisterB())))
  setBinaryText(13, 6, ToRunes(cpu3bit.ToString(info.CarryFlag())))
  setBinaryText( 2, 9, ToRunes(cpu3bit.ToString(info.ProgramCounter())))
  setMemoryText()
  setAddrArrow()

  //reflesh
  termbox.Flush()
}


func setText(x, y int, fg, bg termbox.Attribute, msg string) {
  for _, c := range msg {
	  termbox.SetCell(x, y, c, fg, bg)
    x++
  }
}

func ToRunes(s string) []rune {
  runes := []rune{}
  for _, r := range s {
    runes = append(runes, r)
  }
  //reverse order
  for i,j := 0,len(runes)-1; i<j; i,j = i+1,j-1 {
    runes[i], runes[j] = runes[j], runes[i]
  }
  return runes
}

func setBinaryText(x, y int, runes []rune) {
  for  i, r := range runes {
    if r == '1' {
      termbox.SetCell(x+i, y, r, fgEmColor, bgEmColor)
    } else {
      termbox.SetCell(x+i, y, r, fgColor, bgColor)
    }
  }
}

func setMemoryText() {
  for i, w := range info.Memory() {
    runes := ToRunes(cpu3bit.ToString(w))
    setBinaryText(37, 3+i, runes)
  }
}

func setAddrArrow() {
  y := 0
  p := info.ProgramCounter()
  for i:=0; i<len(p); i++ {
    y += p[i] << uint(i)
  }
  termbox.SetCell(28, 3+y, '>', fgColor, bgColor)
}


// load memory image [8][7]int from textfile
func loadMemoryImageText(filename string) (image [][]int, err error) {

  f, err := os.Open(filename)
  if err != nil {
    return image, err
  }
  defer f.Close()

  s := bufio.NewScanner(f)
  for s.Scan() {
    if a := convert(s.Text()); len(a) != 0 {
    //reverse the order
    for i,j := 0,len(a)-1; i < j; i,j = i+1,j-1 {
      a[i], a[j] = a[j], a[i]
    }
    image = append(image, a)
    }
  }

  if s.Err() != nil {
    return image, s.Err()
  }

  return image, nil
}

func convert(s string) (result []int) {
  for _, r := range s {
    switch r {
    case '#': return result
    case '0': result = append(result, 0)
    case '1': result = append(result, 1)
    }
  }
  return result
}

