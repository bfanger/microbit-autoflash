# microbit-autoflash

Faster development for the [BBC micro:bit]() on macOS.

Simple utility for macOS that monitors the download folder and flashes files that match `microbit-*.hex`.

The loop:

- Wait for microbit usb drive
- Wait for new hex files in the ~/Downloads folder
- Flash the program
- Quickly unmount
- Remove the file

## Installation

Download the binary from the [release page](https://github.com/bfanger/microbit-autoflash/releases)

or compile from source:

```sh
go get github.com/bfanger/microbit-autoflash
```

## Usage

```sh
microbit-autoflash
```

And code with https://makecode.microbit.org/ or https://python.microbit.org/
