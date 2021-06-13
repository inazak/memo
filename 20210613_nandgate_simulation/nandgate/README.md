# nandgate

I wanted to make a 4bit CPU from NAND, so I write a package that simulates a simple logic circuit.
In this simulator, physical power supply, circuit area, delay time, and other things are not considered at all.

I write [3bitCPU](https://github.com/inazak/cpu3bit) and
[4bitCPU](https://github.com/inazak/td4sim) using this package.


## overview

The simulator is constructed around an object called `Line`.
`Line` has the state of HI (1) or LO (0), and it also have a function to update own state.
Circuit takes `Line` as an argument, and returns `Line`.
There is no Circuit object, for example there is no object of Nand.
The Nand is function that takes two `Line`s as input and generates one `Line` as output.
There are MakeLine function, Nand function and Latch function as basic elements.
All other circuits are defined from this function.
When all `Line`s are created, they are stored in the `nodelink.Nodes` list.

Simurator has a `Line` type clock as a package variable, which is used by DFF / DFFC.
Normally, the clock is updated with the Tick function that combines ClockDown / Up and Update.
The Update function updates the state of each `Line` of the circuit to the latest.
The DFF function and the DFFC function hold special data because they hold the clock.
The clock is the starting point of the state change in this simulation.
Starting from here, `nandgate.Update` function updates other `Line`s.


## test

Try out the sample below. The Init function must be executed.
```
package main

import (
  "fmt"
  . "github.com/inazak/nandgate"
)

func main() {

  Init()

  a := MakeLine()
  b := MakeLine()
  a.State = LO
  b.State = LO
  out := Nand(a, b)

// its mean follows.
//  a ----> +------+
//          | NAND | ----> out 
//  b ----> +------+

  Update()

  fmt.Printf("output is %v\n", out.State)
}
```

The output is as follows.
```
output is 1
```

## test with DFF

Try out the sample below.
Trace and Log functions are included in the package as functions for logging.

```
package main

import (
  . "github.com/inazak/nandgate"
)

func main() {

  Init()

  a    := MakeLine()
  a.State = LO
  q, _ := DFF(Not(a))
  Connect(q, a)

// its mean follows.
//  +------[ NOT ] a <--+
//  |                   |
//  +----> [ DFF ] q ---+

  Trace("output q", "%v", q) // monitor output Line

  for i:=0; i<5; i++ {
    Tick() // clockup/down and Update
    Log("output q") // print tagged Line
  }
}
```

The output is as follows.
```
## output q: 1
## output q: 0
## output q: 1
## output q: 0
## output q: 1
```

