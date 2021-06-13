package cpu3bit

import (
  "testing"
  . "github.com/inazak/nandgate"
)

func TestCounter3(t *testing.T) {

  Init()

  d     := MakeLines(3)
  load  := MakeLine()
  clear := MakeLine()
  q     := Counter3([3]*Line{d[0], d[1], d[2]}, load, clear)

  ps := []struct{
    InputD     []int
    InputLoad  int
    InputClear int
    ExpectedQ  []int
  }{
    { []int{LO,LO,LO}, LO, LO, []int{LO,LO,LO}, }, //count up
    { []int{LO,LO,LO}, LO, LO, []int{HI,LO,LO}, },
    { []int{LO,LO,LO}, LO, LO, []int{LO,HI,LO}, },
    { []int{LO,LO,LO}, LO, LO, []int{HI,HI,LO}, },
    { []int{HI,HI,HI}, LO, LO, []int{LO,LO,HI}, }, //d has no effect
    { []int{HI,HI,HI}, LO, LO, []int{HI,LO,HI}, }, //d has no effect
    { []int{LO,LO,LO}, LO, LO, []int{LO,HI,HI}, },
    { []int{LO,LO,LO}, LO, LO, []int{HI,HI,HI}, },
    { []int{LO,LO,LO}, LO, LO, []int{LO,LO,LO}, },
    { []int{LO,LO,LO}, LO, LO, []int{HI,LO,LO}, },
    { []int{LO,LO,LO}, LO, HI, []int{LO,LO,LO}, }, //clear
    { []int{LO,HI,LO}, LO, HI, []int{LO,LO,LO}, }, //clear and d has no effect
    { []int{LO,LO,LO}, LO, LO, []int{HI,LO,LO}, }, //count up
    { []int{LO,LO,LO}, LO, LO, []int{LO,HI,LO}, },
    { []int{HI,HI,HI}, HI, LO, []int{HI,HI,LO}, }, //count up and load next time
    { []int{LO,LO,LO}, LO, LO, []int{HI,HI,HI}, }, //loaded
    { []int{LO,LO,LO}, LO, LO, []int{LO,LO,LO}, },
  }

  for i, p := range ps {
    d[0].State  = p.InputD[0]
    d[1].State  = p.InputD[1]
    d[2].State  = p.InputD[2]
    load.State  = p.InputLoad
    clear.State = p.InputClear
    Tick()
    Update()

    if q[0].State != p.ExpectedQ[0] ||
       q[1].State != p.ExpectedQ[1] ||
       q[2].State != p.ExpectedQ[2] {
      t.Errorf("Counter3L(%v,%v,%v,%v,%v) [%v] q.State expected=%v%v%v got=%v%v%v",
        p.InputD[0],    p.InputD[1],    p.InputD[2], p.InputLoad, p.InputClear, i,
        p.ExpectedQ[0], p.ExpectedQ[1], p.ExpectedQ[2],
        q[0].State,     q[1].State,     q[2].State)
    }
  }
}

func TestALU(t *testing.T) {

  Init()

  a    := MakeLines(4)
  b    := MakeLines(4)
  mode := MakeLine()

  s, c := ALU([4]*Line{ a[0],a[1],a[2],a[3] },
              [4]*Line{ b[0],b[1],b[2],b[3] },
              mode)

  ps := []struct{
    InputA    []int
    InputB    []int
    InputMode int   //LO:bypass a, HI:add a+b
    ExpectedS []int
    ExpectedC int
  }{
    { []int{LO,LO,LO,LO}, []int{LO,HI,HI,LO}, LO,  []int{LO,LO,LO,LO}, LO, }, //bypass
    { []int{HI,LO,LO,HI}, []int{LO,HI,HI,LO}, LO,  []int{HI,LO,LO,HI}, LO, }, //bypass
    { []int{LO,LO,LO,LO}, []int{LO,LO,LO,LO}, HI,  []int{LO,LO,LO,LO}, LO, }, //add
    { []int{LO,LO,HI,LO}, []int{HI,HI,HI,LO}, HI,  []int{HI,HI,LO,HI}, LO, }, //add
    { []int{HI,HI,HI,HI}, []int{HI,HI,HI,HI}, HI,  []int{LO,HI,HI,HI}, HI, }, //add
  }

  for _, p := range ps {
    for i:=0; i<4; i++ {
      a[i].State = p.InputA[i]
      b[i].State = p.InputB[i]
    }
    mode.State  = p.InputMode
    Update()

    if s[0].State != p.ExpectedS[0] ||
       s[1].State != p.ExpectedS[1] ||
       s[2].State != p.ExpectedS[2] ||
       s[3].State != p.ExpectedS[3] ||
       c.State    != p.ExpectedC {
      t.Errorf("ALU a:%v%v%v%v,b:%v%v%v%v mode:%v expected=s:%v%v%v%v,c:%v got=s:%v%v%v%v,c:%v",
        p.InputA[0],    p.InputA[1],    p.InputA[2],    p.InputA[3],
        p.InputB[0],    p.InputB[1],    p.InputB[2],    p.InputB[3],
        p.InputMode,
        p.ExpectedS[0], p.ExpectedS[1], p.ExpectedS[2], p.ExpectedS[3], p.ExpectedC,
        s[0].State,     s[1].State,     s[2].State,     s[3].State,     c.State)
    }
  }
}

func TestDecoder(t *testing.T) {

  Init()

  i0  := MakeLine()
  i1  := MakeLine()
  i2  := MakeLine()
  cin := MakeLine()

  lra, lrb, lpc, mode, sel := Decoder(i0,i1,i2, cin)

  ps := []struct{
    InputI         []int
    InputCin       int
    ExpectedLoadRA int
    ExpectedLoadRB int
    ExpectedLoadPC int
    ExpectedMode   int
    ExpectedSel    int
    Name           string
  }{
    { []int{LO,HI,LO}, LO,  HI,LO,LO, HI, LO, "add a"},
    { []int{HI,HI,LO}, LO,  HI,LO,LO, LO, LO, "mov a"},
    { []int{LO,LO,HI}, LO,  LO,HI,LO, HI, HI, "add b"},
    { []int{HI,LO,HI}, LO,  LO,HI,LO, LO, HI, "mov b"},
    { []int{LO,HI,HI}, LO,  LO,LO,HI, LO, LO, "jnc c=0"},
    { []int{LO,HI,HI}, HI,  LO,LO,LO, LO, LO, "jnc c=1"},
    { []int{HI,HI,HI}, LO,  LO,LO,HI, LO, LO, "jmp c=0"},
    { []int{HI,HI,HI}, HI,  LO,LO,HI, LO, LO, "jmp c=1"},
  }

  for _, p := range ps {
    i0.State  = p.InputI[0]
    i1.State  = p.InputI[1]
    i2.State  = p.InputI[2]
    cin.State = p.InputCin
    Update()

    if lra.State  != p.ExpectedLoadRA ||
       lrb.State  != p.ExpectedLoadRB ||
       lpc.State  != p.ExpectedLoadPC ||
       mode.State != p.ExpectedMode   ||
       sel.State  != p.ExpectedSel {
      t.Errorf("Decoder [%v] expected=lra:%v,lrb:%v,lpc:%v,mode:%v,sel:%v" +
                                 "got=lra:%v,lrb:%v,lpc:%v,mode:%v,sel:%v", p.Name,
        p.ExpectedLoadRA, p.ExpectedLoadRB, p.ExpectedLoadPC, p.ExpectedMode, p.ExpectedSel,
        lra.State,        lrb.State,        lpc.State,        mode.State,     sel.State)
    }
  }
}


