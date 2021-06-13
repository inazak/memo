package nandgate

import (
  "testing"
)

func TestNand3(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  c := MakeLine()
  x := Nand3(a, b, c)

  ps := []struct{
    InputA   int
    InputB   int
    InputC   int
    Expected int
  }{
    { LO, LO, LO, HI },
    { LO, LO, HI, HI },
    { LO, HI, LO, HI },
    { LO, HI, HI, HI },
    { HI, LO, LO, HI },
    { HI, LO, HI, HI },
    { HI, HI, LO, HI },
    { HI, HI, HI, LO },
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    c.State = p.InputC
    Update()
    if x.State != p.Expected {
      t.Errorf("Nand(%v,%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.InputC, p.Expected, x.State)
    }
  }
}

func TestNand4(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  c := MakeLine()
  d := MakeLine()
  x := Nand4(a, b, c, d)

  ps := []struct{
    InputA   int
    InputB   int
    InputC   int
    InputD   int
    Expected int
  }{
    { LO, LO, LO, LO, HI },
    { LO, LO, LO, HI, HI },
    { LO, LO, HI, LO, HI },
    { LO, LO, HI, HI, HI },
    { LO, HI, LO, LO, HI },
    { LO, HI, LO, HI, HI },
    { LO, HI, HI, LO, HI },
    { LO, HI, HI, HI, HI },
    { HI, LO, LO, LO, HI },
    { HI, LO, LO, HI, HI },
    { HI, LO, HI, LO, HI },
    { HI, LO, HI, HI, HI },
    { HI, HI, LO, LO, HI },
    { HI, HI, LO, HI, HI },
    { HI, HI, HI, LO, HI },
    { HI, HI, HI, HI, LO },
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    c.State = p.InputC
    d.State = p.InputD
    Update()
    if x.State != p.Expected {
      t.Errorf("Nand(%v,%v,%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.InputC, p.InputD, p.Expected, x.State)
    }
  }
}

func TestXor(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  x := Xor(a, b)

  ps := []struct{
    InputA   int
    InputB   int
    Expected int
  }{
    { LO, LO, LO },
    { LO, HI, HI },
    { HI, LO, HI },
    { HI, HI, LO },
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    Update()
    if x.State != p.Expected {
      t.Errorf("Xor(%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.Expected, x.State)
    }
  }
}

func TestMux(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  s := MakeLine()
  x := Mux(a, b, s)

  ps := []struct{
    InputA   int
    InputB   int
    InputS   int
    Expected int
  }{
    { LO, LO, LO, LO},
    { HI, LO, LO, HI},
    { LO, HI, LO, LO},
    { HI, HI, LO, HI},
    { LO, LO, HI, LO},
    { LO, HI, HI, HI},
    { HI, LO, HI, LO},
    { HI, HI, HI, HI},
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    s.State = p.InputS
    Update()
    if x.State != p.Expected {
      t.Errorf("Mux(%v,%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.InputS, p.Expected, x.State)
    }
  }
}

func TestMux2(t *testing.T) {

  Init()

  a  := MakeLine()
  b  := MakeLine()
  c  := MakeLine()
  d  := MakeLine()
  s0 := MakeLine()
  s1 := MakeLine()
  x  := Mux2(a, b, c, d, s0, s1)

  ps := []struct{
    InputA   int
    InputB   int
    InputC   int
    InputD   int
    InputS0  int
    InputS1  int
    Expected int
  }{
    // a,  b,  c,  d, s0, s1, out
    { LO, LO, LO, LO, LO, LO, LO},
    { HI, LO, LO, LO, LO, LO, HI},
    { LO, LO, LO, LO, HI, LO, LO},
    { LO, HI, LO, LO, HI, LO, HI},
    { LO, LO, LO, LO, LO, HI, LO},
    { LO, LO, HI, LO, LO, HI, HI},
    { LO, LO, LO, LO, HI, HI, LO},
    { LO, LO, LO, HI, HI, HI, HI},
    { LO, HI, LO, HI, LO, HI, LO},
    { HI, LO, HI, LO, HI, LO, LO},
    { HI, HI, HI, LO, HI, HI, LO},
    { HI, HI, LO, HI, LO, HI, LO},
  }

  for _, p := range ps {
    a.State  = p.InputA
    b.State  = p.InputB
    c.State  = p.InputC
    d.State  = p.InputD
    s0.State = p.InputS0
    s1.State = p.InputS1
    Update()
    if x.State != p.Expected {
      t.Errorf("Mux2(%v,%v,%v,%v,%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.InputC, p.InputD, p.InputS0, p.InputS1, p.Expected, x.State)
    }
  }
}

func TestBDec(t *testing.T) {

  Init()

  s := MakeLine()
  q0,q1 := BDec(s)

  ps := []struct{
    InputS     int
    ExpectedQ0 int
    ExpectedQ1 int
  }{
    { LO, HI, LO },
    { HI, LO, HI },
  }

  for _, p := range ps {
    s.State = p.InputS
    Update()
    if q0.State != p.ExpectedQ0 ||
       q1.State != p.ExpectedQ1 {
      t.Errorf("BDec(%v) expected Q=[%v,%v] got=[%v,%v]",
        p.InputS, p.ExpectedQ0, p.ExpectedQ1, q0.State, q1.State)
    }
  }
}

func TestBDec2(t *testing.T) {

  Init()

  s0 := MakeLine()
  s1 := MakeLine()
  q  := BDec2(s0, s1)

  ps := []struct{
    InputS0    int
    InputS1    int
    ExpectedQ0 int
    ExpectedQ1 int
    ExpectedQ2 int
    ExpectedQ3 int
  }{
    { LO, LO, HI, LO, LO, LO },
    { HI, LO, LO, HI, LO, LO },
    { LO, HI, LO, LO, HI, LO },
    { HI, HI, LO, LO, LO, HI },
  }

  for _, p := range ps {
    s0.State = p.InputS0
    s1.State = p.InputS1
    Update()
    if q[0].State != p.ExpectedQ0 ||
       q[1].State != p.ExpectedQ1 ||
       q[2].State != p.ExpectedQ2 ||
       q[3].State != p.ExpectedQ3 {
      t.Errorf("BDec2(%v,%v) expected Q=[%v,%v,%v,%v] got=[%v,%v,%v,%v]",
        p.InputS0, p.InputS1,
        p.ExpectedQ0, p.ExpectedQ1, p.ExpectedQ2, p.ExpectedQ3,
        q[0].State, q[1].State, q[2].State, q[3].State)
    }
  }
}

func TestBDec3(t *testing.T) {

  Init()

  s0 := MakeLine()
  s1 := MakeLine()
  s2 := MakeLine()
  q  := BDec3(s0,s1,s2)

  ps := []struct{
    InputS0    int
    InputS1    int
    InputS2    int
    ExpectedQ0 int
    ExpectedQ1 int
    ExpectedQ2 int
    ExpectedQ3 int
    ExpectedQ4 int
    ExpectedQ5 int
    ExpectedQ6 int
    ExpectedQ7 int
  }{
    //s0  s1  s2   q0  q1  q2  q3  q4  q5  q6  q7
    { LO, LO, LO,  HI, LO, LO, LO, LO, LO, LO, LO },
    { HI, LO, LO,  LO, HI, LO, LO, LO, LO, LO, LO },
    { LO, HI, LO,  LO, LO, HI, LO, LO, LO, LO, LO },
    { HI, HI, LO,  LO, LO, LO, HI, LO, LO, LO, LO },
    { LO, LO, HI,  LO, LO, LO, LO, HI, LO, LO, LO },
    { HI, LO, HI,  LO, LO, LO, LO, LO, HI, LO, LO },
    { LO, HI, HI,  LO, LO, LO, LO, LO, LO, HI, LO },
    { HI, HI, HI,  LO, LO, LO, LO, LO, LO, LO, HI },
  }

  for _, p := range ps {
    s0.State = p.InputS0
    s1.State = p.InputS1
    s2.State = p.InputS2
    Update()
    if q[0].State != p.ExpectedQ0 ||
       q[1].State != p.ExpectedQ1 ||
       q[2].State != p.ExpectedQ2 ||
       q[3].State != p.ExpectedQ3 ||
       q[4].State != p.ExpectedQ4 ||
       q[5].State != p.ExpectedQ5 ||
       q[6].State != p.ExpectedQ6 ||
       q[7].State != p.ExpectedQ7 {
      t.Errorf("BDec3(%v,%v,%v) expected Q=[%v%v%v%v%v%v%v%v] got=[%v%v%v%v%v%v%v%v]",
        p.InputS0, p.InputS1, p.InputS2,
        p.ExpectedQ0, p.ExpectedQ1, p.ExpectedQ2, p.ExpectedQ3,
        p.ExpectedQ4, p.ExpectedQ5, p.ExpectedQ6, p.ExpectedQ7,
        q[0].State, q[1].State, q[2].State, q[3].State,
        q[4].State, q[5].State, q[6].State, q[7].State)
    }
  }
}

func TestBDec4(t *testing.T) {

  Init()

  s0 := MakeLine()
  s1 := MakeLine()
  s2 := MakeLine()
  s3 := MakeLine()
  q  := BDec4(s0,s1,s2,s3)

  ps := []struct{
    InputS0     int
    InputS1     int
    InputS2     int
    InputS3     int
    ExpectedQ0  int
    ExpectedQ1  int
    ExpectedQ2  int
    ExpectedQ3  int
    ExpectedQ4  int
    ExpectedQ5  int
    ExpectedQ6  int
    ExpectedQ7  int
    ExpectedQ8  int
    ExpectedQ9  int
    ExpectedQ10 int
    ExpectedQ11 int
    ExpectedQ12 int
    ExpectedQ13 int
    ExpectedQ14 int
    ExpectedQ15 int
  }{
    //s0 s1 s2 s3  q0 q1 q2 q3 q4 q5 q6 q7 q8 q9 10 11 12 13 14 15
    { LO,LO,LO,LO, HI,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO },
    { HI,LO,LO,LO, LO,HI,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO },
    { LO,HI,LO,LO, LO,LO,HI,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO },
    { HI,HI,LO,LO, LO,LO,LO,HI,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO },
    { LO,LO,HI,LO, LO,LO,LO,LO,HI,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO },
    { HI,LO,HI,LO, LO,LO,LO,LO,LO,HI,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO },
    { LO,HI,HI,LO, LO,LO,LO,LO,LO,LO,HI,LO,LO,LO,LO,LO,LO,LO,LO,LO },
    { HI,HI,HI,LO, LO,LO,LO,LO,LO,LO,LO,HI,LO,LO,LO,LO,LO,LO,LO,LO },
    { LO,LO,LO,HI, LO,LO,LO,LO,LO,LO,LO,LO,HI,LO,LO,LO,LO,LO,LO,LO },
    { HI,LO,LO,HI, LO,LO,LO,LO,LO,LO,LO,LO,LO,HI,LO,LO,LO,LO,LO,LO },
    { LO,HI,LO,HI, LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,HI,LO,LO,LO,LO,LO },
    { HI,HI,LO,HI, LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,HI,LO,LO,LO,LO },
    { LO,LO,HI,HI, LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,HI,LO,LO,LO },
    { HI,LO,HI,HI, LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,HI,LO,LO },
    { LO,HI,HI,HI, LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,HI,LO },
    { HI,HI,HI,HI, LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,LO,HI },
  }

  for _, p := range ps {
    s0.State = p.InputS0
    s1.State = p.InputS1
    s2.State = p.InputS2
    s3.State = p.InputS3
    Update()
    if q[0].State  != p.ExpectedQ0  ||
       q[1].State  != p.ExpectedQ1  ||
       q[2].State  != p.ExpectedQ2  ||
       q[3].State  != p.ExpectedQ3  ||
       q[4].State  != p.ExpectedQ4  ||
       q[5].State  != p.ExpectedQ5  ||
       q[6].State  != p.ExpectedQ6  ||
       q[7].State  != p.ExpectedQ7  ||
       q[8].State  != p.ExpectedQ8  ||
       q[9].State  != p.ExpectedQ9  ||
       q[10].State != p.ExpectedQ10 ||
       q[11].State != p.ExpectedQ11 ||
       q[12].State != p.ExpectedQ12 ||
       q[13].State != p.ExpectedQ13 ||
       q[14].State != p.ExpectedQ14 ||
       q[15].State != p.ExpectedQ15 {
      t.Errorf("BDec3(%v,%v,%v,%v)", p.InputS0, p.InputS1, p.InputS2, p.InputS3)
    }
  }
}

func TestFAdd(t *testing.T) {
  Init()

  a   := MakeLine()
  b   := MakeLine()
  cin := MakeLine()
  s, c := FAdd(a, b, cin)

  ps := []struct{
    InputA    int
    InputB    int
    InputCin  int
    ExpectedS int
    ExpectedC int
  }{
    { LO, LO, LO,  LO, LO },
    { LO, LO, HI,  HI, LO },
    { LO, HI, LO,  HI, LO },
    { LO, HI, HI,  LO, HI },
    { HI, LO, LO,  HI, LO },
    { HI, LO, HI,  LO, HI },
    { HI, HI, LO,  LO, HI },
    { HI, HI, HI,  HI, HI },
  }

  for _, p := range ps {
    a.State   = p.InputA
    b.State   = p.InputB
    cin.State = p.InputCin
    Update()
    if s.State != p.ExpectedS {
      t.Errorf("Fadd(%v,%v,%v) s.State expected=%v got=%v",
        p.InputA, p.InputB, p.InputCin, p.ExpectedS, s.State)
    }
    if c.State != p.ExpectedC {
      t.Errorf("Fadd(%v,%v,%v) c.State expected=%v got=%v",
        p.InputA, p.InputB, p.InputCin, p.ExpectedC, c.State)
    }
  }
}

