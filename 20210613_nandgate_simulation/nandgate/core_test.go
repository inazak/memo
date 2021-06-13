package nandgate

import (
  "testing"
)

func TestNand(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  x := Nand(a, b)

  ps := []struct{
    InputA   int
    InputB   int
    Expected int
  }{
    { LO, LO, HI },
    { LO, HI, HI },
    { HI, LO, HI },
    { HI, HI, LO },
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    Update()
    if x.State != p.Expected {
      t.Errorf("Nand(%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.Expected, x.State)
    }
  }
}

func TestLatch(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  x := Latch(a, b)

  ps := []struct{
    InputA   int
    InputB   int
    Expected int
  }{
    { LO, HI, HI },
    { HI, LO, LO },
    { HI, HI, LO },
    { LO, HI, HI },
    { HI, HI, HI },
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    Update()
    if x.State != p.Expected {
      t.Errorf("Latch(%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.Expected, x.State)
    }
  }
}

func TestNot(t *testing.T) {

  Init()

  a := MakeLine()
  x := Not(a)

  ps := []struct{
    Input    int
    Expected int
  }{
    { HI, LO },
    { LO, HI },
    { HI, LO },
    { HI, LO },
  }

  for _, p := range ps {
    a.State = p.Input
    Update()
    if x.State != p.Expected {
      t.Errorf("Not(%v) expected=%v got=%v",
        p.Input, p.Expected, x.State)
    }
  }
}

func TestAnd(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  x := And(a, b)

  ps := []struct{
    InputA   int
    InputB   int
    Expected int
  }{
    { LO, LO, LO },
    { LO, HI, LO },
    { HI, LO, LO },
    { HI, HI, HI },
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    Update()
    if x.State != p.Expected {
      t.Errorf("And(%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.Expected, x.State)
    }
  }
}

func TestOr(t *testing.T) {

  Init()

  a := MakeLine()
  b := MakeLine()
  x := Or(a, b)

  ps := []struct{
    InputA   int
    InputB   int
    Expected int
  }{
    { LO, LO, LO },
    { LO, HI, HI },
    { HI, LO, HI },
    { HI, HI, HI },
  }

  for _, p := range ps {
    a.State = p.InputA
    b.State = p.InputB
    Update()
    if x.State != p.Expected {
      t.Errorf("Or(%v,%v) expected=%v got=%v",
        p.InputA, p.InputB, p.Expected, x.State)
    }
  }
}

func TestDFF(t *testing.T) {

  Init()

  a := MakeLine()
  q, nq := DFF(a)

  ps := []struct{
    Input      int
    ExpectedQ  int
    ExpectedNQ int
  }{
    { LO, LO, HI },
    { HI, HI, LO },
    { HI, HI, LO },
    { LO, LO, HI },
    { LO, LO, HI },
    { HI, HI, LO },
  }

  //initialize
  a.State = LO
  Tick()
  Update()

  for i, p := range ps {
    prevQ  := q.State
    prevNQ := nq.State

    a.State = p.Input
    ClockDown()
    Update() //keep q/nq state and store input state

    if q.State != prevQ {
      t.Errorf("DFF q.State after ClockDown[%v] expected=%v got=%v",
        i, prevQ, q.State)
    }
    if nq.State != prevNQ {
      t.Errorf("DFF nq.State after ClockDown[%v] expected=%v got=%v",
        i, prevNQ, nq.State)
    }

    a.State = HI //test changing state at this point has no effect

    ClockUp()
    Update() //update state

    if q.State != p.ExpectedQ {
      t.Errorf("DFF q.State after ClockUp[%v] expected=%v got=%v",
        i, p.ExpectedQ, q.State)
    }
    if nq.State != p.ExpectedNQ {
      t.Errorf("DFF nq.State after ClockUp[%v] expected=%v got=%v",
        i, p.ExpectedNQ, nq.State)
    }
  }
}

func TestDFFLoop(t *testing.T) {

  Init()

  a    := MakeLine()
  q, _ := DFF(Not(a))
  Connect(q, a)
  a.State = LO

// its mean follows.
//  +------[ NOT ] a <--+
//  |                   |
//  +----> [ DFF ] q ---+

  for i, p := range []int{ HI, LO, HI, LO, HI, } {

    Tick() // clockup-down and Update
    if q.State != p {
      t.Errorf("DFFLoop q.State after Tick[%v] expected=%v got=%v",
        i, p, q.State)
    }
  }
}

func TestDFFC(t *testing.T) {

  Init()

  a        := MakeLine()
  clear    := MakeLine()
  q, nq := DFFC(a, clear)

  ps := []struct{
    Input      int
    InputClear int
    ExpectedQ  int
    ExpectedNQ int
  }{
    { LO, LO, LO, HI },
    { HI, LO, HI, LO },
    { HI, LO, HI, LO },
    { LO, LO, LO, HI },
    { LO, LO, LO, HI },
    { HI, LO, HI, LO },
    { HI, HI, LO, HI },
    { LO, HI, LO, HI },
    { HI, HI, LO, HI },
    { HI, HI, LO, HI },
  }

  //initialize
  a.State     = LO
  clear.State = LO
  Tick()
  Update()

  for i, p := range ps {
    prevQ  := q.State
    prevNQ := nq.State

    a.State     = p.Input
    clear.State = p.InputClear

    if p.InputClear == HI {
      prevQ  = LO //expect clearing
      prevNQ = HI
    }

    ClockDown()
    Update() //keep q/nq state and store input state

    if q.State != prevQ {
      t.Errorf("DFFC q.State after ClockDown[%v] expected=%v got=%v",
        i, prevQ, q.State)
    }
    if nq.State != prevNQ {
      t.Errorf("DFFC nq.State after ClockDown[%v] expected=%v got=%v",
        i, prevNQ, nq.State)
    }

    a.State = HI //test changing state at this point has no effect

    ClockUp()
    Update() //update state

    if q.State != p.ExpectedQ {
      t.Errorf("DFFC q.State after ClockUp[%v] expected=%v got=%v",
        i, p.ExpectedQ, q.State)
    }
    if nq.State != p.ExpectedNQ {
      t.Errorf("DFFC nq.State after ClockUp[%v] expected=%v got=%v",
        i, p.ExpectedNQ, nq.State)
    }
  }
}


