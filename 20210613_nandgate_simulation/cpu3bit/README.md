# cpu3bit

3bit CPU Simulator made using only NAND and Latch. For fun.
See [nandgate](https://github.com/inazak/nandgate) for details on NAND and Latch implementation.


## How to use

```
cpu3bit [OPTIONS]

If there is no option, a demo 3x3 instruction is loaded.
And the clock advances automatically.

  OPTIONS:
    -load=FILE  ... load memory image text.
    -manual     ... clock ticking by hand.
```


## Demo

calculate 3 * 3 = 9

```
101 0000   # 0: mov B 0
011 1101   # 1: mov A 13
100 0011   # 2: add B 3
010 0001   # 3: add A 1
110 0010   # 4: jnc 2
111 0101   # 5: jmp 5 (end)
000 0000   # Register B is 1001 (9)
000 0000   #
```

![](https://raw.githubusercontent.com/inazak/cpu3bit/master/misc/sample2.gif)


## Installation

```
$ go get github.com/inazak/cpu3bit/cmd/cpu3bit
```


## Diagram

![](https://raw.githubusercontent.com/inazak/cpu3bit/master/misc/sample1.png)


## Instructions

Memory is 7bits and 8words.
Lower 4bits is immidiate data.
Higher 3bits is operator.

```
High     Low
6  5  4  3-0  carry  | mnemonic
---------------------|----------------
0  1  0  Imd    x    | ADD A, Imd
0  1  1  Imd    x    | MOV A, Imd
1  0  0  Imd    x    | ADD B, Imd
1  0  1  Imd    x    | MOV B, Imd
1  1  0  Imd    0    | JNC Imd (jump if carry=0)
1  1  1  Imd    x    | JMP Imd
0  0  0  Imd    x    | NOP
0  0  1  Imd    x    | NOP
```


## Requirements

golang.


