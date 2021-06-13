package td4sim

import (
  . "github.com/inazak/nandgate"
)

type CPUInfo struct {
  memory         [][]*Line
  registerA      []*Line
  registerB      []*Line
  registerC      []*Line
  programCounter []*Line
  carryFlag      []*Line
  dipsw          []*Line
}

func (c *CPUInfo) Memory() [][]int {
  result := [][]int{}
  for _, lines := range c.memory {
    result = append(result, toIntSlice(lines))
  }
  return result
}

func (c *CPUInfo) RegisterA() []int {
  return toIntSlice(c.registerA)
}

func (c *CPUInfo) RegisterB() []int {
  return toIntSlice(c.registerB)
}

func (c *CPUInfo) RegisterC() []int {
  return toIntSlice(c.registerC)
}

func (c *CPUInfo) ProgramCounter() []int {
  return toIntSlice(c.programCounter)
}

func (c *CPUInfo) CarryFlag() []int {
  return toIntSlice(c.carryFlag)
}

func (c *CPUInfo) DIPSW() []int {
  return toIntSlice(c.dipsw)
}

func (c *CPUInfo) ChangeDIPSW(sw int) {
  if (sw < 0) || (sw >= len(c.dipsw)) {
    return
  }

  if c.dipsw[sw].State == HI {
    c.dipsw[sw].State = LO
  } else {
    c.dipsw[sw].State = HI
  }
}


// ----- utilities -----

func toIntSlice(lines []*Line) []int {

  result := make([]int, len(lines))
  for i, line := range lines {
    result[i] = line.State
  }
  return result
}

func ToString(list []int) string {

  result := ""
  for _, i := range list {
    if i == HI {
      result += "1"
    } else if i == LO {
      result += "0"
    } else {
      result += "-"
    }
  }
  return result
}

// ----- building blocks -----

func Register4(d0,d1,d2,d3, load, clear *Line) []*Line {

  in0   := MakeLine()
  q0, _ := DFFC(in0, clear)
  fb0 := Mux(q0, d0, load)
  Connect(fb0, in0)

  in1   := MakeLine()
  q1, _ := DFFC(in1, clear)
  fb1 := Mux(q1, d1, load)
  Connect(fb1, in1)

  in2   := MakeLine()
  q2, _ := DFFC(in2, clear)
  fb2 := Mux(q2, d2, load)
  Connect(fb2, in2)

  in3   := MakeLine()
  q3, _ := DFFC(in3, clear)
  fb3 := Mux(q3, d3, load)
  Connect(fb3, in3)

  return []*Line{ q0, q1, q2, q3 }
}

func Counter4(d0,d1,d2,d3, load, clear *Line) []*Line {

  in0   := MakeLine()
  q0, _ := DFFC(in0, clear)
  r0  := Not(q0)
  fb0 := Mux(r0, d0, load)
  Connect(fb0, in0)

  in1   := MakeLine()
  q1, _ := DFFC(in1, clear)
  r1  := Xor(q0, q1)
  fb1 := Mux(r1, d1, load)
  Connect(fb1, in1)

  in2   := MakeLine()
  q2, _ := DFFC(in2, clear)
  a1  := And(q1, Xor(q0, q2))
  a2  := And(q2, Not(q1))
  r2  := Or(a1, a2)
  fb2 := Mux(r2, d2, load)
  Connect(fb2, in2)

  in3   := MakeLine()
  q3, _ := DFFC(in3, clear)
  b1  := And(Not(q2), q3)
  b2  := And4(q0, q1, q2, Not(q3))
  b3  := And3(Not(q0), q2, q3)
  b4  := And3(Not(q1), q2, q3)
  r3  := Or4(b1, b2, b3, b4)
  fb3 := Mux(r3, d3, load)
  Connect(fb3, in3)

  return []*Line{ q0, q1, q2, q3 }
}

// ROM is 8bit x 16word memory.
// Looping 8 times after enabling one of the selectors to retrieve memory contents.
// For example, if you want to get the last memory:
//
//  Input          Selector     Memory
//                              +---------------+  <-+
//  (H,H,H,H)      s[0] ---->   | 8bit          |    |
//                      (L)     +---------------+    |
//  [4]in          s[1] ---->   | ....          |    | 16
//    |                 (L)     +---------------+    | word
//    +-> BDec ->  s[n] ---->   | ....          |    |
//                      (L)     +---------------+    |
//                 s[15] --->   | 8bit (Target) |    |
//                      (H)     +---------------+  <-+
//                                 |
//                                 | x 8times loop
//                                 |
//                                 +---> [8]out
// 
func ROM(memory [][]*Line, in []*Line) (out []*Line) {

  out = make([]*Line, 8)
  s := BDec4(in[0], in[1], in[2], in[3])

  for i:=0; i<len(out); i++ {

    out[i] = Or4( Or4( And(s[ 0], memory[ 0][i]), And(s[ 1], memory[ 1][i]),
                       And(s[ 2], memory[ 2][i]), And(s[ 3], memory[ 3][i])),
                  Or4( And(s[ 4], memory[ 4][i]), And(s[ 5], memory[ 5][i]),
                       And(s[ 6], memory[ 6][i]), And(s[ 7], memory[ 7][i])),
                  Or4( And(s[ 8], memory[ 8][i]), And(s[ 9], memory[ 9][i]),
                       And(s[10], memory[10][i]), And(s[11], memory[11][i])),
                  Or4( And(s[12], memory[12][i]), And(s[13], memory[13][i]),
                       And(s[14], memory[14][i]), And(s[15], memory[15][i])))
  }

  return out
}

// IC74HC74 is D-type flip-flop
func IC74HC74 (d, nclear *Line) (q, nq *Line) {
  in   := MakeLine()
  q, nq = DFFC(in, Not(nclear))
  Connect(d, in)

  return q, nq
}

// IC74HC153 is Selector
func IC74HC153(c10,c11,c12,c13, c20,c21,c22,c23, a, b *Line) (y1, y2 *Line) {

  y1 = Mux2(c10,c11,c12,c13, a, b)
  y2 = Mux2(c20,c21,c22,c23, a, b)

  return y1, y2
}

// IC74HC283 is 4bit full adder
func IC74HC283(a0,a1,a2,a3, b0,b1,b2,b3, c0 *Line) (s [4]*Line, c4 *Line) {

  s0, c1 := FAdd(a0, b0, c0)
  s1, c2 := FAdd(a1, b1, c1)
  s2, c3 := FAdd(a2, b2, c2)
  s3, c4 := FAdd(a3, b3, c3)

  s = [4]*Line{s0,s1,s2,s3}
  return s, c4
}

func IC74HC161R(d0,d1,d2,d3, nload, nclear *Line) (q []*Line) {

  q = Register4(d0,d1,d2,d3, Not(nload), Not(nclear))
  return q
}

func IC74HC161C(d0,d1,d2,d3, nload, nclear *Line) (q []*Line) {

  q = Counter4(d0,d1,d2,d3, Not(nload), Not(nclear))
  return q
}

func Decoder(op0,op1,op2,op3, ncflg *Line) (sela, selb *Line, nload [4]*Line) {
  hi := MakeLine()
  hi.State = 1

  sela = Or(op0, op3)
  selb = op1

  nload[0] = Or(op2, op3)
  nload[1] = Or(op3, Nand3(hi, hi, op2))
  nload[2] = Nand3(hi, Nand3(hi, hi, op2), op3)
  nload[3] = Nand3(op3, op2, Or(op0, ncflg))

  return sela, selb, nload
}

// ----- setup -----

func MakeTD4(image [][]int, dipsw []int) *CPUInfo {

  info := &CPUInfo{}

  //load memory image
  for i:=0; i<len(image); i++ {
    line := MakeLines(len(image[i]))
    for j:=0; j<len(image[i]); j++ {
      line[j].State = image[i][j]
    }
    info.memory = append(info.memory, line)
  }

  low := MakeLine()

  //make DIPSW
  info.dipsw = MakeLines(4)
  for i:=0; i<len(dipsw); i++ {
    info.dipsw[i].State = dipsw[i]
  }

  //register pins
  RAin := MakeLines(4)
  RBin := MakeLines(4)
  RCin := MakeLines(4)
  PCin := MakeLines(4)
  RAnotLD := MakeLine()
  RBnotLD := MakeLine()
  RCnotLD := MakeLine()
  PCnotLD := MakeLine()
  clear  := MakeLine()
  nclear := Not(clear)

  //make register
  RAout := IC74HC161R(RAin[0],RAin[1],RAin[2],RAin[3], RAnotLD, nclear)
  RBout := IC74HC161R(RBin[0],RBin[1],RBin[2],RBin[3], RBnotLD, nclear)
  RCout := IC74HC161R(RCin[0],RCin[1],RCin[2],RCin[3], RCnotLD, nclear)
  PCout := IC74HC161C(PCin[0],PCin[1],PCin[2],PCin[3], PCnotLD, nclear)

  //make ROM
  ROMout := ROM(info.memory, PCout)

  //selecter pins
  SelRA := MakeLine()
  SelRB := MakeLine()

  //make selecter
  A0in,A1in := IC74HC153(RAout[0], RBout[0], info.dipsw[0], low,
                         RAout[1], RBout[1], info.dipsw[1], low, SelRA, SelRB )
  A2in,A3in := IC74HC153(RAout[2], RBout[2], info.dipsw[2], low,
                         RAout[3], RBout[3], info.dipsw[3], low, SelRA, SelRB )

  //make adder
  Sigma, Carry := IC74HC283(A0in,A1in,A2in,A3in,
                            ROMout[0],ROMout[1],ROMout[2],ROMout[3], low)

  Connect(Sigma[0], RAin[0])
  Connect(Sigma[0], RBin[0])
  Connect(Sigma[0], RCin[0])
  Connect(Sigma[0], PCin[0])

  Connect(Sigma[1], RAin[1])
  Connect(Sigma[1], RBin[1])
  Connect(Sigma[1], RCin[1])
  Connect(Sigma[1], PCin[1])

  Connect(Sigma[2], RAin[2])
  Connect(Sigma[2], RBin[2])
  Connect(Sigma[2], RCin[2])
  Connect(Sigma[2], PCin[2])

  Connect(Sigma[3], RAin[3])
  Connect(Sigma[3], RBin[3])
  Connect(Sigma[3], RCin[3])
  Connect(Sigma[3], PCin[3])

  //make cflag
  CFlag, NotCFlag := IC74HC74(Carry, nclear)

  //make decoder
  DecSelRA, DecSelRB, DecNotLoad := Decoder(ROMout[4],ROMout[5],ROMout[6],ROMout[7], NotCFlag)

  Connect(DecSelRA, SelRA)
  Connect(DecSelRB, SelRB)
  Connect(DecNotLoad[0], RAnotLD)
  Connect(DecNotLoad[1], RBnotLD)
  Connect(DecNotLoad[2], RCnotLD)
  Connect(DecNotLoad[3], PCnotLD)

  //set info
  info.registerA = RAout
  info.registerB = RBout
  info.registerC = RCout
  info.programCounter = PCout
  info.carryFlag = []*Line{ CFlag }

  return info
}


func Initialize() {
  Init() //nandgate.Init()
}

func TickTock() {
  Tick() //nandgate.Tick()
}


