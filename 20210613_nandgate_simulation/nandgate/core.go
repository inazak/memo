package nandgate

// overview
//
// The simulator is constructed around an object called `Line`.
// `Line` has the state of HI (1) or LO (0), and it also have a function to update own state.
// Circuit takes `Line` as an argument, and returns `Line`.
// There is no Circuit object, for example there is no object of Nand.
// The Nand is function that takes two `Line`s as input and generates one `Line` as output.
// There are MakeLine function, Nand function and Latch function as basic elements.
// All other circuits are defined from this function.
// When all `Line`s are created, they are stored in the `nodelink.Nodes` list.
// 
// Simurator has a `Line` type clock as a package variable, which is used by DFF / DFFC.
// Normally, the clock is updated with the Tick function that combines ClockDown / Up and Update.
// The Update function updates the state of each `Line` of the circuit to the latest.
// The DFF function and the DFFC function hold special data because they hold the clock.
// The clock is the starting point of the state change in this simulation.
// Starting from here, `nandgate.Update` function updates other `Line`s.

const (
  HI = 1
  LO = 0
)

var link  *nodelink
var clock *Line

type Line struct {
  State  int    // HI or LO
  Update func()
}

type nodelink struct {
  Nodes      []*Line
  NextNodes  map[*Line][]*Line
  PrevNodes  map[*Line][]*Line
  StartNodes map[*Line]struct{}
}

func (ln *nodelink) appendNode(l *Line) {
  ln.Nodes = append(ln.Nodes, l)
}

func (ln *nodelink) getNodes() []*Line {
  return ln.Nodes
}

func (ln *nodelink) setStartNode(l *Line) {
  ln.StartNodes[l] = struct{}{}
}

func (ln *nodelink) isStartNode(l *Line) bool {
  _, ok := ln.StartNodes[l]
  return ok
}

func (ln *nodelink) appendLink(from, to *Line) {
  ln.NextNodes[from] = append(ln.NextNodes[from], to)
  ln.PrevNodes[to]   = append(ln.PrevNodes[to], from)
}

func (ln *nodelink) getNextNodes(l *Line) []*Line {
  if _, ok := ln.NextNodes[l] ; ok {
    return ln.NextNodes[l]
  }
  return []*Line{}
}

func (ln *nodelink) getPrevNodes(l *Line) []*Line {
  if _, ok := ln.PrevNodes[l] ; ok {
    return ln.PrevNodes[l]
  }
  return []*Line{}
}


func Init() {
  clock = &Line{
    State:  LO,
    Update: func(){}, //do nothing
  }

  link = &nodelink{}
  link.Nodes      = []*Line{}
  link.StartNodes = make(map[*Line]struct{})
  link.NextNodes  = make(map[*Line][]*Line)
  link.PrevNodes  = make(map[*Line][]*Line)

  link.appendNode(clock)
  link.setStartNode(clock)
}

func ClockUp() {
  clock.State = HI
}

func ClockDown() {
  clock.State = LO
}

func Tick() {
  ClockDown()
  Update()
  ClockUp()
  Update()
}

func MakeLine() *Line {
  out := &Line{
    State:  LO,
    Update: func(){}, //do nothing
  }
  link.appendNode(out)
  return out
}

func MakeLines(size int) []*Line {
  var m []*Line
  for i:=0 ; i<size; i++ {
    m = append(m, MakeLine())
  }
  return m
}

func Connect(from, to *Line) {
  if size := len(link.getPrevNodes(to)) ; size != 0 {
    panic("Connect() cant connect. line already has link.")
  }
  to.Update = func() { //closure
    to.State = from.State
  }
  link.appendLink(from, to)
}


// Nand
//
//  a   b   | out
//  --------|----
//  LO  LO  | HI
//  LO  HI  | HI
//  HI  LO  | HI
//  HI  HI  | LO
//
func Nand(a, b *Line) (out *Line) {
  out = &Line{
    State:  LO,
    Update: func() { //closure
      if a.State == HI && b.State == HI {
        out.State = LO
      } else {
        out.State = HI
      }
    },
  }
  link.appendNode(out)
  link.appendLink(a, out)
  link.appendLink(b, out)
  return out
}


// Latch
//
// a ----+------+
//       | NAND |----+----------------- q
//   +---+------+    |
//   |               +---+------+
//   |                   | NAND |---+
//   |             b ----+------+   |
//   +------------------------------+
//
//  a   b   | q
//  --------|----
//  LO  LO  | HI (dont use)
//  LO  HI  | HI
//  HI  LO  | LO
//  HI  HI  | Keep State
// 
func Latch(a, b *Line) (q *Line) {
  q = &Line{
    State:  LO,
    Update: func() { //closure
      if a.State == LO && b.State == LO { q.State = HI }
      if a.State == LO && b.State == HI { q.State = HI }
      if a.State == HI && b.State == LO { q.State = LO }
      if a.State == HI && b.State == HI { /*keep q.State*/ }
    },
  }
  link.appendNode(q)
  link.appendLink(a, q)
  link.appendLink(b, q)
  return q
}

func Not(a *Line) *Line {
  return Nand(a, a)
}

func And(a, b *Line) *Line {
  return Not(Nand(a, b))
}

func Or(a, b *Line) *Line {
  return Nand(Not(a), Not(b))
}

// DLatch
// if enable is HI(1), then q equals d value.
// if enable is LO(0), then q keep prev value.
// 
//  d   e   | q
//  --------|----
//  LO  LO  | q (latched)
//  HI  LO  | q (latched)
//  LO  HI  | LO
//  HI  HI  | HI
// 
func DLatch(d, e *Line) (q, nq *Line) {
  a := Nand(d, e)
  b := Nand(Not(d), e)
  q  = Latch(a, b)
  nq = Not(q)
  return q, nq
}

func DLatchC(d, e, clear *Line) (q, nq *Line) {
  a := Nand(d, e)
  b := Nand(Not(d), e)
  g := Or(a, clear)
  h := And(b, Not(clear))
  q  = Latch(g, h)
  nq = Not(q)
  return q, nq
}

// DFF has clock
func DFF(in *Line) (q, nq *Line) {
  link.setStartNode(in)
  a, _   := DLatch(in, Not(clock))
  q, nq   = DLatch(a, clock)
  return q, nq
}

func DFFC(in, clear *Line) (q, nq *Line) {
  link.setStartNode(in)
  a, _   := DLatchC(in, Not(clock), clear)
  q, nq   = DLatchC(a, clock, clear)
  return q, nq
}


func Update() {

  var queue []*Line

  // start with StartNodes or node has no PrevNodes
  for _, l := range link.getNodes() {
    if link.isStartNode(l) || len(link.getPrevNodes(l)) == 0 {
      queue = append(queue, l)
    }
  }

  updated := make(map[*Line]struct{})

  for len(queue) > 0 {
    l := queue[0]
    queue = queue[1:]

    // already updated
    if _, ok := updated[l] ; ok {
      continue
    }

    // check if node is updatable, otherwise return to the queue
    if ! link.isStartNode(l) && ! isResolved( updated, link.getPrevNodes(l) ) {
      queue = append(queue, l)
      continue
    }

    l.Update()
    updated[l] = struct{}{}

    // propagated to NextNodes
    for _, t := range link.getNextNodes(l) {
      queue = append(queue, t)
    }
  } //queue loop

}

func isResolved(updated map[*Line]struct{}, lines []*Line) bool {
  for _, l := range lines {
    if _, ok := updated[l] ; ! ok {
      return false
    }
  }
  return true
}


