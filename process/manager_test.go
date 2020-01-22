package process

import (
	"context"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestManager(t *testing.T) {

	manager := NewManager(3*time.Second, func(logType OutputType, line string, process *Process) {
		color.Green("%s[%s] => %s\n", process.GetName(), logType, line)
	})
	manager.AddProgram("test", "/bin/sleep 1", 5, "")
	manager.AddProgram("prometheus", "/bin/echo Hello", 1, "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	manager.Watch(ctx)
}
