package cpu3bit

import (
  . "github.com/inazak/nandgate"
)

// overview
// 
// This is 3bit CPU Emurator made using only NAND and Latch.
// The following objects appear in this program.
//
// RegisterA,B    : 4bit register.
// ProgramCounter : 3bit counter.
// Selector       : unit to select one of the two inputs.
// Memory         : 7bit * 8word memory.
// ALU            : unit performing bypass or addition.
// Decoder        : unit that break up an instruction into separate signals.
// 
// The process is as follows.
//
// The clock goes on and the counter is updated.
// Lower 4bits of memory indicated by program counter are sent to ALU.
//
//   [ProgramCounter(3bit)] --> [Memory(7bit*8word)] --(4bit)--> [ALU]
//
// Rest 3bits of memory indicated by program counter are sent to Decoder.
//
//   [ProgramCounter(3bit)] --> [Memory(7bit*8word)] --(3bit)--> [Decoder]
//
// The decoder sends 1bit to the selector which determines the register to read.
//
//   [Decoder] --(1bit)--> [selector]
// 
// The selector uses input 1bit to send the value of either register to the ALU.
//
//   [RegisterA(4bit)] --> [selector] --(4bit)--> [ALU]
//   [RegisterB(4bit)] -->
//
// The decoder sends 1bit to the ALU which determines the processing.
//
//   [Decoder] --(1bit)--> [ALU]
//
// ALU sends operation result to register, program counter, and carry flag.
// The value of carry flag sends to the decoder.
//
//   [ALU] --(4bit)--> [RegisterA]
//                 --> [RegisterB]
//                 --> [ProgramCounter]
//         --(1bit)--> [CarryFlag] --(1bit)--> [Decoder]
//
// The decoder sends 1bit to registers and program counter.
// It is determine whether to upddate itself with the value sent from ALU.
//
//   [Decoder] --(1bit)--> load [RegisterA]
//             --(1bit)--> load [RegisterB]
//             --(1bit)--> load [ProgramCounter]
//

// --------- circuit parts ----------

// 3bit counter, loadable
func Counter3(d [3]*Line, load, clear *Line) (q [3]*Line) {
  r0    := MakeLine()
  q0, _ := DFFC(r0, clear)
  fb0   := Not(q0)
  in0   := Mux(fb0, d[0], load)
  Connect(in0, r0)

  r1    := MakeLine()
  q1, _ := DFFC(r1, clear)
  fb1   := Xor(q0, q1)
  in1   := Mux(fb1, d[1], load)
  Connect(in1, r1)

  r2    := MakeLine()
  q2, _ := DFFC(r2, clear)
  a1    := And(q1, Xor(q0, q2))
  a2    := And(q2, Not(q1))
  fb2   := Or(a1, a2)
  in2   := Mux(fb2, d[2], load)
  Connect(in2, r2)

  return [3]*Line{ q0,q1,q2 }
}

// 4bit Register
func Register4(d [4]*Line, load, clear *Line) (q [4]*Line) {

  r0    := MakeLine()
  q0, _ := DFFC(r0, clear)
  i0        := Mux(q0, d[0], load)
  Connect(i0, r0)

  r1    := MakeLine()
  q1, _ := DFFC(r1, clear)
  i1        := Mux(q1, d[1], load)
  Connect(i1, r1)

  r2    := MakeLine()
  q2, _ := DFFC(r2, clear)
  i2        := Mux(q2, d[2], load)
  Connect(i2, r2)

  r3    := MakeLine()
  q3, _ := DFFC(r3, clear)
  i3        := Mux(q3, d[3], load)
  Connect(i3, r3)

  return [4]*Line{ q0,q1,q2,q3 }
}

// --------- 3bit cpu structures ----------

// ALU:
// When the mode is low, it is bypassed(s=a),
// and when it is high, it is addition(s=a+b).
func ALU(a, b [4]*Line, mode *Line) (s [4]*Line, c *Line) {
  c0 := MakeLine()
  s0, c1 := FAdd(a[0], And(b[0], mode), c0)
  s1, c2 := FAdd(a[1], And(b[1], mode), c1)
  s2, c3 := FAdd(a[2], And(b[2], mode), c2)
  s3, c  := FAdd(a[3], And(b[3], mode), c3)
  return [4]*Line{ s0,s1,s2,s3 }, c
}

func CFlag(d *Line) *Line {
  r    := MakeLine()
  q, _ := DFF(r)
  Connect(d, r)
  return q
}

// Selector:
// When sel is lo, output q is a, and when sel is hi, output q is b.
func Selecter(a, b [4]*Line, sel *Line) (q [4]*Line) {
  q0 := Mux(a[0], b[0], sel)
  q1 := Mux(a[1], b[1], sel)
  q2 := Mux(a[2], b[2], sel)
  q3 := Mux(a[3], b[3], sel)
  return [4]*Line{ q0,q1,q2,q3 }
}

// Decoder:
// takes 3bit input and carryin.
//
//  instructions                       0:mov    0:A
//  bit           | output             1:add    1:B
//  2  1  0  cin  | loadA loadB loadPC ALCmode  Sel      
//  --------------|------------------------------------ 
//  0  1  0  x    | 1     0     0      1        0       ADD A, X
//  0  1  1  x    | 1     0     0      0        0       MOV A, X
//  1  0  0  x    | 0     1     0      1        1       ADD B, X
//  1  0  1  x    | 0     1     0      0        1       MOV B, X
//  1  1  0  0    | 0     0     1      0        x       JNC X
//  1  1  1  x    | 0     0     1      0        x       JMP X
//  0  0  0  x    | x     x     x      x        x       NOP
//  0  0  1  x    | x     x     x      x        x       NOP
//
func Decoder(i0, i1, i2, cin *Line) (lda, ldb, ldpc, mode, sel *Line) {

  lda  = And(Not(i2),    i1)
  ldb  = And(    i2, Not(i1))
  ldpc = Or(And3(i2, i1, Not(cin)), And3(i2, i1, i0))
  mode = Or(And3(Not(i2), i1, Not(i0)), And3(i2, Not(i1), Not(i0)))
  sel  = And(i2, Not(i1))

  return
}

// ROM:
// take memory[A][B] has the following layout,
// and decide which record to return according to input 3bit.
//
//     B[0 .... .... 6]
//      +-------------+
// A[0] | memory data |
//      +-------------+
// A[1] | memory data |
//      +-------------+
//      |    ....     |
//      +-------------+
// A[7] | memory data |
//      +-------------+
//
func ROM(memory [][]*Line, in [3]*Line) (out [7]*Line) {

  s := BDec3(in[0], in[1], in[2])

  for i:=0; i<len(out); i++ {

    out[i] = Or( Or4( And(s[ 0], memory[ 0][i]), And(s[ 1], memory[ 1][i]),
                      And(s[ 2], memory[ 2][i]), And(s[ 3], memory[ 3][i])),
                 Or4( And(s[ 4], memory[ 4][i]), And(s[ 5], memory[ 5][i]),
                      And(s[ 6], memory[ 6][i]), And(s[ 7], memory[ 7][i])))
  }

  return out
}


// --------- utilities ----------

func lineToInt(lines ...*Line) []int {
  s := []int{}
  for _, line := range lines {
    s = append(s, line.State)
  }
  return s
}

func ToString(list []int) string {

  s := ""
  for _, i := range list {
    if i == HI {
      s += "1"
    } else if i == LO {
      s += "0"
    } else {
      s += "-"
    }
  }
  return s
}

// --------- computer ----------

type CPUInfo struct {
  memory         [][]*Line
  registerA      []*Line
  registerB      []*Line
  programCounter []*Line
  carryFlag      []*Line
}

func (c *CPUInfo) Memory() [][]int {
  r := [][]int{}
  for _, m := range c.memory {
    r = append(r, lineToInt(m...))
  }
  return r
}

func (c *CPUInfo) RegisterA() []int {
  return lineToInt(c.registerA...)
}

func (c *CPUInfo) RegisterB() []int {
  return lineToInt(c.registerB...)
}

func (c *CPUInfo) ProgramCounter() []int {
  return lineToInt(c.programCounter...)
}

func (c *CPUInfo) CarryFlag() []int {
  return lineToInt(c.carryFlag...)
}


// MakeComputer:
// The specification of 3bit CPU is as follows
//
// - Memory is 7bits and 8words.
// - Lower 4bits is immidiate data.
// - Higher 3bits is operator.
//
// Instructions
//  High     Low
//  6  5  4  3-0  carry  | mnemonic
//  ---------------------|----------------
//  0  1  0  Imd    x    | ADD A, Imd
//  0  1  1  Imd    x    | MOV A, Imd
//  1  0  0  Imd    x    | ADD B, Imd
//  1  0  1  Imd    x    | MOV B, Imd
//  1  1  0  Imd    0    | JNC Imd (jump if carry=0)
//  1  1  1  Imd    x    | JMP Imd
//  0  0  0  Imd    x    | NOP
//  0  0  1  Imd    x    | NOP
//
func MakeComputer(image [][]int) *CPUInfo {

  info := &CPUInfo{}

  //load memory image
  for i:=0; i<len(image); i++ {
    line := MakeLines(len(image[i]))
    for j:=0; j<len(image[i]); j++ {
      line[j].State = image[i][j]
    }
    info.memory = append(info.memory, line)
  }

  // reset line
  clear := MakeLine()

  // register A
  RAin   := MakeLines(4)
  RAload := MakeLine()
  RAout  := Register4([4]*Line{RAin[0],RAin[1],RAin[2],RAin[3]}, RAload, clear)

  // register B
  RBin   := MakeLines(4)
  RBload := MakeLine()
  RBout  := Register4([4]*Line{RBin[0],RBin[1],RBin[2],RBin[3]}, RBload, clear)

  // program counter
  PCin   := MakeLines(3)
  PCload := MakeLine()
  PCout  := Counter3([3]*Line{PCin[0],PCin[1],PCin[2]}, PCload, clear)

  // carry flag
  Cin  := MakeLine()
  Cout := CFlag(Cin)

  //make ROM
  ROMout := ROM(info.memory, PCout)

  // decoder
  isLoadRA, isLoadRB, isLoadPC, isMode, isSEL := Decoder(ROMout[4],ROMout[5],ROMout[6],Cout)
  Connect(isLoadRA, RAload)
  Connect(isLoadRB, RBload)
  Connect(isLoadPC, PCload)

  // register selecter
  SELout := Selecter(RAout, RBout, isSEL)

  // ALU
  ALUout, Carry := ALU([4]*Line{ ROMout[0],ROMout[1],ROMout[2],ROMout[3] }, SELout, isMode)
  Connect(Carry, Cin)

  // connect from ALU output to register and program counter
  for i:=0; i<4; i++ {
    Connect(ALUout[i], RAin[i])
    Connect(ALUout[i], RBin[i])
  }
  for i:=0; i<3; i++ {
    Connect(ALUout[i], PCin[i])
  }

  // return value
  info.registerA      = RAout[0:4]
  info.registerB      = RBout[0:4]
  info.programCounter = PCout[0:3]
  info.carryFlag      = []*Line{ Cout }

  return info
}

func Initialize() {
  Init() //nandgate.Init()
}

func TickTock() {
  Tick() //nandgate.Tick()
}



