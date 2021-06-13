package main

import (
  "flag"
  "fmt"
  "os"
  "bufio"
  "time"
  "github.com/nsf/termbox-go"
  "github.com/inazak/td4sim"
)

const (
  fgColor   = termbox.ColorWhite
  bgColor   = termbox.ColorBlack
  fgEmColor = termbox.ColorBlack
  bgEmColor = termbox.ColorWhite
)


var display = []string{
        //01234567890123456789012345678901234567890123456789
/*  0 */ "                TD4 Simulator                     ",
/*  1 */ "                                                  ",
/*  2 */ " Register A [####]       Address    Memory        ",
/*  3 */ "                               0    [########]    ",
/*  4 */ " Register B [####]             1    [########]    ",
/*  5 */ "                               2    [########]    ",
/*  6 */ " Register C [####]             3    [########]    ",
/*  7 */ "                               4    [########]    ",
/*  8 */ " Carry Flag [#]                5    [########]    ",
/*  9 */ "                               6    [########]    ",
/* 10 */ " Program Counter               7    [########]    ",
/* 11 */ " [####]                        8    [########]    ",
/* 12 */ "                               9    [########]    ",
/* 13 */ " DIP Switch [####]            10    [########]    ",
/* 14 */ "                              11    [########]    ",
/* 15 */ "                              12    [########]    ",
/* 17 */ "                              13    [########]    ",
/* 18 */ "                              14    [########]    ",
/* 19 */ "                              15    [########]    ",
/* 20 */ " [auto mode (default)]                            ",
/* 21 */ "   q=quit                                         ",
/* 22 */ " [manual mode]                                    ",
/* 23 */ "   q=quit, c=clock, 0,1,2,3=change dipsw          ",

}

var memoryimage = [][]int{
  //D0           D7
  {1,0,0,0,0,0,0,0}, // ADD A,0b0001
  {1,0,0,0,1,0,1,0}, // ADD B,0b0001
  {1,0,0,0,0,0,0,0},
  {1,0,0,0,1,0,1,0},
  {1,0,0,0,0,0,0,0}, //4
  {1,0,0,0,1,0,1,0},
  {1,0,0,0,0,0,0,0},
  {1,0,0,0,1,0,1,0},
  {1,0,0,0,0,0,0,0}, //8
  {1,0,0,0,1,0,1,0},
  {1,0,0,0,0,0,0,0},
  {1,0,0,0,1,0,1,0},
  {1,0,0,0,0,0,0,0}, //12
  {1,0,0,0,1,0,1,0},
  {1,0,0,0,0,0,0,0},
  {1,0,0,0,1,0,1,0},
}

var dipswimage = []int{ 0,0,0,0 }

var info *td4sim.CPUInfo

var file = flag.String("load", "", "textfile representing memory image")
var dipsw = flag.String("dipsw",  "", "initial dipsw setting (order3-0)")
var manual = flag.Bool("manual", false, "clocking by hand")

func main() {
  flag.Parse()

  // load text file
  if *file != "" {
    var err error
    memoryimage, err = loadMemoryImageText(*file)
    if err != nil {
      fmt.Printf("%v", err)
      return
    }
  }

  if *dipsw != "" {
    dipswimage = convert(*dipsw)
    if len(dipswimage) != 4 {
      fmt.Printf("parameter dipsw error")
      return
    }
    //reverse elements
    for i,j := 0,len(dipswimage)-1; i < j; i,j = i+1,j-1 {
      dipswimage[i], dipswimage[j] = dipswimage[j], dipswimage[i]
    }
  }

  td4sim.Initialize()
  info = td4sim.MakeTD4(memoryimage, dipswimage)

  // termbox init and event-loop
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

  //auto clocking
  if ! *manual {
    go func(){
      for {
        select {
        case <- time.After(time.Millisecond * 1000):
          td4sim.TickTock()
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
        case ev.Ch == 'c':
          if *manual {
            td4sim.TickTock() // clockdown/up and update
            render()
          }
        case ev.Ch == '0':
          if *manual {
            info.ChangeDIPSW(0)
            render()
          }
        case ev.Ch == '1':
          if *manual {
            info.ChangeDIPSW(1)
            render()
          }
        case ev.Ch == '2':
          if *manual {
            info.ChangeDIPSW(2)
            render()
          }
        case ev.Ch == '3':
          if *manual {
            info.ChangeDIPSW(3)
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
  setBinaryText(13, 2, ToRunes(td4sim.ToString(info.RegisterA())))
  setBinaryText(13, 4, ToRunes(td4sim.ToString(info.RegisterB())))
  setBinaryText(13, 6, ToRunes(td4sim.ToString(info.RegisterC())))
  setBinaryText(13, 8, ToRunes(td4sim.ToString(info.CarryFlag())))
  setBinaryText(2, 11, ToRunes(td4sim.ToString(info.ProgramCounter())))
  setBinaryText(13,13, ToRunes(td4sim.ToString(info.DIPSW())))
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
    runes := ToRunes(td4sim.ToString(w))
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
      //reversed order
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


