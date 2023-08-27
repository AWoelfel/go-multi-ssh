package connection

import (
	"fmt"
	"github.com/muesli/termenv"
	"strings"
)

type ClientContext struct {
	ID        string
	Container string
	Col       termenv.ANSI256Color
}

func (c *ClientContext) Label() string {
	return fmt.Sprintf("%s (%s)", c.ID[:12], strings.TrimPrefix(c.Container, "/"))
}

func (c *ClientContext) Color() termenv.ANSI256Color {
	return c.Col
}
