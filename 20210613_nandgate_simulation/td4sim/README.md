# td4sim

This is 4bit CPU Simulator made using only NAND and Latch. See Japanese book "How to make CPU" (ISBN4-8399-0986-5) for details.
See [nandgate](https://github.com/inazak/nandgate) for details on NAND and Latch implementation.

TD4とは『CPUの創りかた - IC10個のお手軽CPU設計超入門 初歩のデジタル回路動作の基本原理と製作』で
説明されている4bitCPU。[NAND と Latch から始めて](https://github.com/inazak/nandgate) 、
各IC部品の動作を組み立て、最終的にTD4として動くように作った。


## How to use

```
td4sim [OPTIONS]

If there is no option, a demo instruction add 1 to register A and B is loaded.
And the clock advances automatically.

  OPTIONS:
    -load=FILE  ... load memory image text.
    -dipsw=0000 ... load dipsw initial status.
    -manual     ... clock ticking by hand.
```

```
-load=FILE

  プログラムメモリを記述したテキストファイルを読む
  テキストファイルは一行に8bit分、改行区切りで16行分
  8bitは最上位ビットが先頭、最下位ビットが最後になる
  0 または 1 を並べて記述する、スペースは無視される
  #は以降はコメント
  指定しない場合は、A,Bレジスタに交互に1をADDする
  プログラムがロードされる

-dipsw=0000

  ロード時のDIPスイッチの状態を0または1を並べて4bit分記述する
  4bitは最上位ビットが先頭、最下位ビットが最後になる
  指定しない場合は 0000 となる

-manual

  クロックを手動で動かす
  指定しない場合は 1Hz で自動でクロックが動く
```


## Demo

repeat `ADD A,0b0001` and `ADD B,0b0001` endlessly.

![](https://raw.githubusercontent.com/inazak/td4sim/master/misc/sample1.gif)


## Instructions

Memory is 8bit x 16words.
There is also a 4bit DIPSW for input.

```
High     Low      Input 
7 6 5 4  3 2 1 0  dipsw carry | mnemonic
------------------------------|------------
0 0 1 1  d a t a  ----  -     | MOV A,data
0 1 1 1  d a t a  ----  -     | MOV B,data
0 0 0 1  0 0 0 0  ----  -     | MOV A,B
0 1 0 0  0 0 0 0  ----  -     | MOV B,A
0 0 0 0  d a t a  ----  -     | ADD A,data
0 1 0 1  d a t a  ----  -     | ADD B,data
0 0 1 0  0 0 0 0  data  -     | IN  A,data(dipsw)
0 1 1 0  0 0 0 0  data  -     | IN  B,data(dipsw)
1 0 1 1  d a t a  ----  -     | OUT data == MOV C,data
1 0 0 1  0 0 0 0  ----  -     | OUT B    == MOV C,B
1 1 1 1  d a t a  ----  -     | JMP data
1 1 1 0  d a t a  ----  c     | JNC data (jump if carry==0)
```


## Requirements

golang.


