// Basic Driver for the Adafruit I2C Backpack 7-Segment Display
// Adapted from the adafruit C++ driver for Arduino

package i2c7Seg

import (
	"github.com/d2r2/go-i2c"
)

const (
	HT16K33_CMD_BRIGHTNESS  byte = 0xE0 //< I2C register for BRIGHTNESS setting
	HT16K33_BLINK_CMD       byte = 0x80 ///< I2C register for BLINK setting
	HT16K33_BLINK_DISPLAYON byte = 0x01 ///< I2C value for steady on
	HT16K33_BLINK_OFF       byte = 0    ///< I2C value for steady off
	HT16K33_BLINK_2HZ       byte = 1    ///< I2C value for 2 Hz blink
	HT16K33_BLINK_1HZ       byte = 2    ///< I2C value for 1 Hz blink
	HT16K33_BLINK_HALFHZ    byte = 3    ///< I2C value for 0.5 Hz blink
)

var (
	sevenSegFontTable = [...]byte{
		0b00000000, // (space)
		0b10000110, // !
		0b00100010, // "
		0b01111110, // #
		0b01101101, // $
		0b11010010, // %
		0b01000110, // &
		0b00100000, // '
		0b00101001, // (
		0b00001011, // )
		0b00100001, // *
		0b01110000, // +
		0b00010000, // ,
		0b01000000, // -
		0b10000000, // .
		0b01010010, // /
		0b00111111, // 0
		0b00000110, // 1
		0b01011011, // 2
		0b01001111, // 3
		0b01100110, // 4
		0b01101101, // 5
		0b01111101, // 6
		0b00000111, // 7
		0b01111111, // 8
		0b01101111, // 9
		0b00001001, // :
		0b00001101, // ;
		0b01100001, // <
		0b01001000, // =
		0b01000011, // >
		0b11010011, // ?
		0b01011111, // @
		0b01110111, // A
		0b01111100, // B
		0b00111001, // C
		0b01011110, // D
		0b01111001, // E
		0b01110001, // F
		0b00111101, // G
		0b01110110, // H
		0b00110000, // I
		0b00011110, // J
		0b01110101, // K
		0b00111000, // L
		0b00010101, // M
		0b00110111, // N
		0b00111111, // O
		0b01110011, // P
		0b01101011, // Q
		0b00110011, // R
		0b01101101, // S
		0b01111000, // T
		0b00111110, // U
		0b00111110, // V
		0b00101010, // W
		0b01110110, // X
		0b01101110, // Y
		0b01011011, // Z
		0b00111001, // [
		0b01100100, //
		0b00001111, // ]
		0b00100011, // ^
		0b00001000, // _
		0b00000010, // `
		0b01011111, // a
		0b01111100, // b
		0b01011000, // c
		0b01011110, // d
		0b01111011, // e
		0b01110001, // f
		0b01101111, // g
		0b01110100, // h
		0b00010000, // i
		0b00001100, // j
		0b01110101, // k
		0b00110000, // l
		0b00010100, // m
		0b01010100, // n
		0b01011100, // o
		0b01110011, // p
		0b01100111, // q
		0b01010000, // r
		0b01101101, // s
		0b01111000, // t
		0b00011100, // u
		0b00011100, // v
		0b00010100, // w
		0b01110110, // x
		0b01101110, // y
		0b01011011, // z
		0b01000110, // {
		0b00110000, // |
		0b01110000, // }
		0b00000001, // ~
		0b00000000, // del
	}
)

type SevenSegI2C struct {
	i2c           *i2c.I2C
	postion       byte
	displaybuffer [5]uint8
}

func NewSevenSegI2C(address byte, bus int) (*SevenSegI2C, error) {
	i2c, err := i2c.NewI2C(address, bus)
	if err != nil {
		return nil, err
	}
	new7seg := &SevenSegI2C{
		i2c:     i2c,
		postion: 0,
	}
	new7seg.Begin()
	return new7seg, nil
}
func (self *SevenSegI2C) Begin() error {

	// turn on oscillator
	var buffer [1]byte
	buffer[0] = 0x21
	_, err := self.i2c.WriteBytes(buffer[:])
	if err != nil {
		return err
	}

	// internal RAM powers up with garbage/random values.
	// ensure internal RAM is cleared before turning on display
	// this ensures that no garbage pixels show up on the display
	// when it is turned on.
	self.Clear()
	self.WriteDisplay()

	self.BlinkRate(HT16K33_BLINK_OFF)

	self.SetBrightness(15) // max brightness

	return nil
}
func (self *SevenSegI2C) Clear() {
	self.displaybuffer = [5]uint8{0, 0, 0, 0, 0}
}
func (self *SevenSegI2C) WriteDisplay() error {
	var buffer [17]byte

	buffer[0] = 0x00 // start at address $00

	for i := 0; i < 5; i++ {
		buffer[1+2*i] = byte(self.displaybuffer[i] & 0xFF)
		// buffer[2+2*i] = byte(self.displaybuffer[i] >> 8)
	}
	_, err := self.i2c.WriteBytes(buffer[:])
	return err
}
func (self *SevenSegI2C) SetBrightness(brightness byte) error {
	if brightness > 15 {
		brightness = 15 // limit to max brightness
	}
	var buffer [1]byte
	buffer[0] = HT16K33_CMD_BRIGHTNESS | brightness
	_, err := self.i2c.WriteBytes(buffer[:])
	return err
}
func (self *SevenSegI2C) BlinkRate(rate byte) error {
	if rate > 3 {
		rate = 0 // turn off if not sure
	}
	var buffer [1]byte
	buffer[0] = HT16K33_BLINK_CMD | HT16K33_BLINK_DISPLAYON | (rate << 1)
	_, err := self.i2c.WriteBytes(buffer[:])
	return err
}
func (self *SevenSegI2C) WriteDigitRaw(d byte, bitmask byte) {
	if d > 4 {
		return
	}
	self.displaybuffer[d] = bitmask
}
func (self *SevenSegI2C) DrawColon(visible bool) {
	if visible {
		self.displaybuffer[2] = 0x2
	} else {
		self.displaybuffer[2] = 0
	}
}
func (self *SevenSegI2C) WriteAsciiChar(d byte, ascii byte, dp bool) {
	if d > 4 {
		return
	}
	var letter byte = sevenSegFontTable[ascii-32]
	var dpMask byte = 0x00
	if dp {
		dpMask = 0x80
	}
	self.WriteDigitRaw(d, letter|dpMask)
}
