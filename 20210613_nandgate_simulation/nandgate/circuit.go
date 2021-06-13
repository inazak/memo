package nandgate

func Nand3(a, b, c *Line) *Line {
  x := Nand(a, b)
  return Nand(c, Nand(x, x))
}

func Nand4(a, b, c, d *Line) *Line {
  x := Nand(a, b)
  y := Nand(c, d)
  return Nand(Nand(x, x), Nand(y, y))
}

func And3(a, b, c *Line) *Line {
  return And(c, And(a, b))
}

func Or3(a, b, c *Line) *Line {
  return Or(c, Or(a, b))
}

func And4(a, b, c, d *Line) *Line {
  return And(d, And(c, And(a, b)))
}

func Or4(a, b, c, d *Line) *Line {
  return Or(d, Or(c, Or(a, b)))
}

func Xor(a, b *Line) *Line {
  c := Nand(a, b)
  d := Nand(a, c)
  e := Nand(b, c)
  return Nand(d, e)
}


// Mux is Multiplexor
// if sel is LO, then output a, otherwise output b.
//
//  sel | out
//  LO  | a
//  HI  | b
//
//  a   b   sel | out
//  ------------|----
//  LO  LO  LO  | LO
//  LO  HI  LO  | LO
//  HI  LO  LO  | HI
//  HI  HI  LO  | HI
//  LO  LO  HI  | LO
//  LO  HI  HI  | HI
//  HI  LO  HI  | LO
//  HI  HI  HI  | HI
//
func Mux(a, b, sel *Line) *Line {
  return Or( And(a, Not(sel)), And(b, sel) )
}

// Mux2 is Multiplexor has 2bit selector
//
//  s0   s1   | out
//  LO   LO   | a
//  HI   LO   | b
//  LO   HI   | c
//  HI   HI   | d
//
func Mux2(a, b, c, d, s0, s1 *Line) *Line {
  return Or(And(Mux(a, b, s0), Not(s1)), And(Mux(c, d, s0), s1))
}

// BDec is Binary Decoder
//
//   sel | q0  q1
//   LO  | HI  LO
//   HI  | LO  HI
//
func BDec(sel *Line) (q0, q1 *Line) {
  q0 = Not(sel)
  q1 = MakeLine()
  Connect(sel, q1)
  return q0, q1
}

// BDec2 is Binary Decoder has 2bit selector
//
//  s0  s1 | q0  q1  q2  q3
//  LO  LO | HI  LO  LO  LO
//  HI  LO | LO  HI  LO  LO
//  LO  HI | LO  LO  HI  LO
//  HI  HI | LO  LO  LO  HI
//
func BDec2(s0, s1 *Line) []*Line {
  q0 := And(Not(s0), Not(s1))
  q1 := And(    s0 , Not(s1))
  q2 := And(Not(s0),     s1 )
  q3 := And(    s0 ,     s1 )
  return []*Line{q0, q1, q2, q3}
}

// BDec3 is Binary Decoder has 3bit selector
//
//  s0  s1  s2 | q0  q1  q2  q3  q4  q5  q6  q7
//  LO  LO  LO | HI  LO  LO  LO  LO  LO  LO  LO
//  HI  LO  LO | LO  HI  LO  LO  LO  LO  LO  LO
//  LO  HI  LO | LO  LO  HI  LO  LO  LO  LO  LO
//  HI  HI  LO | LO  LO  LO  HI  LO  LO  LO  LO
//  LO  LO  HI | LO  LO  LO  LO  HI  LO  LO  LO
//  HI  LO  HI | LO  LO  LO  LO  LO  HI  LO  LO
//  LO  HI  HI | LO  LO  LO  LO  LO  LO  HI  LO
//  HI  HI  HI | LO  LO  LO  LO  LO  LO  LO  HI
//
func BDec3(s0,s1,s2 *Line) []*Line {
  q := make([]*Line, 8)
  q[0] = And3(Not(s0), Not(s1), Not(s2))
  q[1] = And3(    s0 , Not(s1), Not(s2))
  q[2] = And3(Not(s0),     s1 , Not(s2))
  q[3] = And3(    s0 ,     s1 , Not(s2))
  q[4] = And3(Not(s0), Not(s1),     s2 )
  q[5] = And3(    s0 , Not(s1),     s2 )
  q[6] = And3(Not(s0),     s1 ,     s2 )
  q[7] = And3(    s0 ,     s1 ,     s2 )
  return q
}

func BDec4(s0,s1,s2,s3 *Line) []*Line {
  q := make([]*Line, 16)
  q[0]  = And4(Not(s0), Not(s1), Not(s2), Not(s3))
  q[1]  = And4(    s0 , Not(s1), Not(s2), Not(s3))
  q[2]  = And4(Not(s0),     s1 , Not(s2), Not(s3))
  q[3]  = And4(    s0 ,     s1 , Not(s2), Not(s3))
  q[4]  = And4(Not(s0), Not(s1),     s2 , Not(s3))
  q[5]  = And4(    s0 , Not(s1),     s2 , Not(s3))
  q[6]  = And4(Not(s0),     s1 ,     s2 , Not(s3))
  q[7]  = And4(    s0 ,     s1 ,     s2 , Not(s3))
  q[8]  = And4(Not(s0), Not(s1), Not(s2),     s3 )
  q[9]  = And4(    s0 , Not(s1), Not(s2),     s3 )
  q[10] = And4(Not(s0),     s1 , Not(s2),     s3 )
  q[11] = And4(    s0 ,     s1 , Not(s2),     s3 )
  q[12] = And4(Not(s0), Not(s1),     s2 ,     s3 )
  q[13] = And4(    s0 , Not(s1),     s2 ,     s3 )
  q[14] = And4(Not(s0),     s1 ,     s2 ,     s3 )
  q[15] = And4(    s0 ,     s1 ,     s2 ,     s3 )
  return q
}

// FAdd is 1bit FullAdder
func FAdd(a, b, cin *Line) (s, c *Line) {
  w := Nand(a, b)
  x := Nand(b, cin)
  y := Nand(a, cin)
  c  = Or(Not(w), Or(Not(x), Not(y)))
  s  = Xor(Xor(a, b), cin)
  return
}

