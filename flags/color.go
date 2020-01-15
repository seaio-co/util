package flags

import (
	"flag"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

type hexColor struct {
	color.Color
}

func (c *hexColor) String() string {
	r, g, b, a := c.Color.RGBA()
	return fmt.Sprintf("rgba(%d, %d, %d, %d)", r, g, b, a)
}

func (c *hexColor) Set(s string) error {
	if strings.HasPrefix(s, "#") {
		s = s[1:]
	}
	if len(s) == 3 {
		s = fmt.Sprintf("%c0%c0%c0", s[0], s[1], s[2])
	}
	if len(s) != 6 {
		return fmt.Errorf("color should be 3 or 6 hex digits")
	}
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return fmt.Errorf("not hexadecimal: %v", err)
	}
	c.Color = &color.RGBA{
		R: uint8(n >> 16),
		G: uint8(n >> 8),
		B: uint8(n),
		A: 0xff,
	}
	return nil
}

// HexColor defines a hex color flag with specified name, default value, and usage string.
// The return value is the address of an RGBA color variable that stores the value of the flag.
func HexColor(name string, value color.Color, usage string) color.Color {
	c := &hexColor{value}
	flag.Var(c, name, usage)
	return c
}

// HexColorVar defines a hex color flag with specified name, default value, and usage string.
// The argument c points to an RGBA color variable in which to store the value of the flag.
func HexColorVar(c *color.Color, name string, value color.Color, usage string) {
	p := &hexColor{value}
	*c = p
	flag.Var(p, name, usage)
}
