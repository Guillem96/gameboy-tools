# Game Boy Tools

Game Boy tools written in Go.

This repository contains a set of tools to interact with the Game Boy system.

## Software Design

This code is thought to work with any micro controller in the market (arduino, raspberry pi, etc.), for this
reason, we have created some interfaces that allow us to abstract from the underlying controller.

Hence, we have an abstract Game Boy cartridge pin ([`GameBoyPin`](gbproxy/base.go)) and an abstract
Game Boy Proxy ([`GameBoyProxy`](gbproxy/base.go)). The later is backed up on the abstract pins so it
can implement the following methods:

- `SelectAddress`: Writes to the A0-A15 pins the desired address we want to select in the Game Boy cart.
- `Read` and `Write`: These methods read the byte from the selected address and writes a specified value in the 
selected address. Both methods use the D0-D7 pins.
- `SetReadMode` & `SetWriteMode`: A0-A15 are always set as output mode, because the micro controller is the 
one in charge of managing the selected addresses (this might change depending on the use case). On the other hand,
the D0-D7 pins are bidirectional (we have to write or read from them). Therefore this two methods will help you
prepare the pins before a read or a write operation respectively.

Going back to the `GameBoyPin`, it defines these methods: 

- `Read`: Returns `true` if the pin is in High state.
- `High`: Sets the pin to high state.
- `Low`: Sets the pin to low state.
- `Input`: Sets the pin to input mode.
- `Output`: Sets the pin to output mode.

> Remember that the interfaces described here are abstract and do not provide any functionality, aside
from the code "schema". Therefore if you want to make this work with your own micro controller you'll have to implement
this two interfaces yourself.

By default we provide a Raspberry Pi proxy implementation.

## Hardware

Obviously to use this repository you need a Game Boy or a Game Boy color.

Also, to make use of this tools I recommend this two amazing pieces of hardware:

- [Breakout & ROM Game Boy Cart PCB](https://stacksmashing.gumroad.com/l/gbcart) by [stacksmashing](https://www.youtube.com/c/stacksmashing)
- [Cartridge Breakout Board](https://www.tindie.com/products/driptronics/cartridge-breakout-board-for-gameboy/)

> Now we are only supporting Raspberry Pi, but in a future we are planning to 
release Arduino and Raspberry Pi Zero support.