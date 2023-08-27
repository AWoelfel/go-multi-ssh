package connection

import (
	"fmt"
	"github.com/muesli/termenv"
)

type ClientContext struct {
	Namespace string
	Pod       string
	Container string
	Col       termenv.ANSI256Color
}

func (c *ClientContext) Label() string {
	return fmt.Sprintf("%s/%s:%s", c.Namespace, c.Pod, c.Container)
}

func (c *ClientContext) Color() termenv.ANSI256Color {
	return c.Col
}
