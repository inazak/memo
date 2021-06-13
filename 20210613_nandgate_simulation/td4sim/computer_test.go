package td4sim

import (
  "testing"
  . "github.com/inazak/nandgate"
)

func makeMemory(image [][]int) (memory [][]*Line) {

  for i:=0; i<len(image); i++ {
    line := MakeLines(len(image[i]))
    for j:=0; j<len(image[i]); j++ {
      line[j].State = image[i][j]
    }
    memory = append(memory, line)
  }

  return memory
}

func TestRegister4(t *testing.T) {

  Init()

  d     := MakeLines(4)
  load  := MakeLine()
  clear := MakeLine()
  q     := Register4(d[0],d[1],d[2],d[3], load, clear)

  p := [][]int{
  // d           ld cr nextq
    {1, 1, 1, 0, 1, 0, 0, 0, 0, 0}, //load 1110  out:0000(initial)
    {0, 0, 0, 0, 0, 0, 1, 1, 1, 0}, //nop        out:1110
    {0, 0, 0, 0, 0, 1, 0, 0, 0, 0}, //clear      out:0000
    {0, 1, 1, 0, 1, 0, 0, 0, 0, 0}, //load       out:0000
    {1, 1, 1, 1, 0, 0, 0, 1, 1, 0}, //nop        out:0110
    {1, 1, 1, 1, 1, 1, 0, 0, 0, 0}, //clear&load out:0000
    {0, 0, 0, 0, 0, 0, 1, 1, 1, 1}, //nop        out:1111
  }

  for _, r := range p {

    d[0].State  = r[0]
    d[1].State  = r[1]
    d[2].State  = r[2]
    d[3].State  = r[3]
    load.State  = r[4]
    clear.State = r[5]

    ClockDown()
    Update()
    ClockUp()
    Update()

    if q[0].State != r[6] ||
       q[1].State != r[7] ||
       q[2].State != r[8] ||
       q[3].State != r[9] {
      t.Errorf("Register q want=%v%v%v%v, but=%v%v%v%v",
        r[6],r[7],r[8],r[9], q[0].State,q[1].State,q[2].State,q[3].State)
    }
  }
}

func TestCounter4(t *testing.T) {

  Init()

  d     := MakeLines(4)
  load  := MakeLine()
  clear := MakeLine()
  q     := Counter4(d[0],d[1],d[2],d[3], load, clear)

  p := [][]int{
  // d           ld cr nextq
    {0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //initial missed swing
    {0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
    {0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
    {0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
    {0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
    {0, 0, 0, 0, 0, 0, 1, 0, 1, 0},
    {0, 0, 0, 0, 0, 0, 0, 1, 1, 0},
    {0, 0, 0, 0, 0, 0, 1, 1, 1, 0},
    {0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
    {0, 0, 0, 0, 0, 0, 1, 0, 0, 1},
    {0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
    {0, 0, 0, 0, 0, 0, 1, 1, 0, 1},
    {0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
    {0, 0, 0, 0, 0, 0, 1, 0, 1, 1},
    {0, 0, 0, 0, 0, 0, 0, 1, 1, 1},
    {0, 0, 0, 0, 0, 0, 1, 1, 1, 1},
    {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
    {1, 1, 0, 0, 1, 0, 1, 0, 0, 0}, //load next
    {0, 0, 0, 0, 0, 0, 1, 1, 0, 0}, //loaded
    {0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, //nop countup
  }

  for _, r := range p {

    d[0].State  = r[0]
    d[1].State  = r[1]
    d[2].State  = r[2]
    d[3].State  = r[3]
    load.State  = r[4]
    clear.State = r[5]

    ClockDown()
    Update()
    ClockUp()
    Update()

    if q[0].State != r[6] ||
       q[1].State != r[7] ||
       q[2].State != r[8] ||
       q[3].State != r[9] {
      t.Errorf("Counter q want=%v%v%v%v, but=%v%v%v%v",
        r[6],r[7],r[8],r[9], q[0].State,q[1].State,q[2].State,q[3].State)
    }
  }
}

func TestROM(t *testing.T) {

  Init()

  info := &CPUInfo{}
  info.memory = makeMemory([][]int{
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,1,0,1,1,1,1,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
  })

  sel := MakeLines(4)
  out := ROM(info.memory, sel)

  sel[0].State = 0
  sel[1].State = 1
  sel[2].State = 0
  sel[3].State = 0

  Update()

  if out[0].State != 0 { t.Errorf("out[0] want=%v, but=%v", 0, out[0].State) }
  if out[1].State != 1 { t.Errorf("out[1] want=%v, but=%v", 1, out[1].State) }
  if out[2].State != 0 { t.Errorf("out[2] want=%v, but=%v", 0, out[2].State) }
  if out[3].State != 1 { t.Errorf("out[3] want=%v, but=%v", 1, out[3].State) }
  if out[4].State != 1 { t.Errorf("out[4] want=%v, but=%v", 1, out[4].State) }
  if out[5].State != 1 { t.Errorf("out[5] want=%v, but=%v", 1, out[5].State) }
  if out[6].State != 1 { t.Errorf("out[6] want=%v, but=%v", 1, out[6].State) }
  if out[7].State != 0 { t.Errorf("out[7] want=%v, but=%v", 0, out[7].State) }

}

func TestIC74HC161R(t *testing.T) {

  Init()
  in := MakeLines(4)
  nl := MakeLine()
  nc := MakeLine()

  out := IC74HC161R(in[0],in[1],in[2],in[3], nl, nc)

  Connect(out[0], in[0])
  Connect(out[1], in[1])
  Connect(out[2], in[2])
  Connect(out[3], in[3])

  nl.State = 0 //load on
  nc.State = 1 //clear off

  in[0].State = 1
  in[1].State = 1
  in[2].State = 0
  in[3].State = 0

  ClockDown()
  Update()

  if out[0].State != 0 &&
     out[1].State != 0 &&
     out[2].State != 0 &&
     out[3].State != 0 {
    t.Errorf("want=0000 but=%v%v%v%v", out[0].State,out[1].State,out[2].State,out[3].State)
  }

  ClockUp() //regist
  Update()

  if out[0].State != 0 &&
     out[1].State != 0 &&
     out[2].State != 0 &&
     out[3].State != 0 {
    t.Errorf("want=0000 but=%v%v%v%v", out[0].State,out[1].State,out[2].State,out[3].State)
  }

  ClockDown()
  Update()

  if out[0].State != 0 &&
     out[1].State != 0 &&
     out[2].State != 0 &&
     out[3].State != 0 {
    t.Errorf("want=0000 but=%v%v%v%v", out[0].State,out[1].State,out[2].State,out[3].State)
  }

  ClockUp() //output
  Update()

  if out[0].State != 1 &&
     out[1].State != 1 &&
     out[2].State != 0 &&
     out[3].State != 0 {
    t.Errorf("want=1100 but=%v%v%v%v", out[0].State,out[1].State,out[2].State,out[3].State)
  }
}

func TestDecoder(t *testing.T) {
  Init()

  op   := MakeLines(4)
  cflg := MakeLine()

  sa, sb, nload := Decoder(op[0],op[1],op[2],op[3], Not(cflg))

  p := [][]int{
  // op              sel   nload
  // 3  2  1  0  cf  b  a  0  1  2  3
    {0, 0, 0, 0, 0,  0, 0, 0, 1, 1, 1}, //Add A,Im
    {0, 0, 0, 1, 0,  0, 1, 0, 1, 1, 1}, //MOV A,B
    {0, 0, 1, 0, 0,  1, 0, 0, 1, 1, 1}, //IN  A
    {0, 0, 1, 1, 0,  1, 1, 0, 1, 1, 1}, //MOV A,Im
    {0, 1, 0, 0, 0,  0, 0, 1, 0, 1, 1}, //MOV B,A
    {0, 1, 0, 1, 0,  0, 1, 1, 0, 1, 1}, //Add B,Im
    {0, 1, 1, 0, 0,  1, 0, 1, 0, 1, 1}, //IN  B
    {0, 1, 1, 1, 0,  1, 1, 1, 0, 1, 1}, //MOV B,Im
    {1, 0, 0, 1, 0,  0, 1, 1, 1, 0, 1}, //OUT B
    {1, 0, 1, 1, 0,  1, 1, 1, 1, 0, 1}, //OUT Im
    {1, 1, 1, 0, 0,  1, 1, 1, 1, 1, 0}, //JNC(C=0) Jump
    {1, 1, 1, 0, 1,  1, 1, 1, 1, 1, 1}, //JNC(C=1)
    {1, 1, 1, 1, 0,  1, 1, 1, 1, 1, 0}, //JMP
  }

  for _, r := range p {
    op[0].State = r[3]
    op[1].State = r[2]
    op[2].State = r[1]
    op[3].State = r[0]
    cflg.State  = r[4]
    Update()

    if sa.State != r[6] ||
       sb.State != r[5] ||
       nload[0].State != r[7] ||
       nload[1].State != r[8] ||
       nload[2].State != r[9] ||
       nload[3].State != r[10] {
         t.Errorf("Decoder op=%v%v%v%v cf=%v want=%v%v:%v%v%v%v but=%v%v:%v%v%v%v",
           op[0].State,op[1].State,op[2].State,op[3].State, cflg.State,
           r[6], r[5], r[7],r[8],r[9],r[10],
           sa.State, sb.State, nload[0].State,nload[1].State,nload[2].State,nload[3].State)
    }
  }
}

func TestMakeTD4(t *testing.T) {

  Init()
  info := MakeTD4([][]int{
  //D0             D7
    {1,0,0,0,0,0,0,0}, // add a, 1
    {1,0,0,0,0,0,0,0},
    {1,0,0,0,0,0,0,0},
    {1,0,0,0,0,0,0,0},
    {1,0,0,0,0,0,0,0}, //4
    {1,0,0,0,0,0,0,0},
    {1,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0}, //8
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0}, //12
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
    {0,0,0,0,0,0,0,0},
  },
  []int{ 0,0,0,0 }) //dipsw

  info.dipsw[0].State = 0
  info.dipsw[1].State = 0
  info.dipsw[2].State = 0
  info.dipsw[3].State = 0

  for i:=0; i<6; i++ {
    Tick()

    ra := 0
    for j:=0; j<4; j++ {
      ra += info.registerA[j].State << uint(j)
    }

    if i != ra {
      t.Errorf("TD4 RegA want=%v but=%v", i, ra)
    }
  }
}



