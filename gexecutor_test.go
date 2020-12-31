package gexecutor

import (
	"errors"
	"testing"
)

func TestGexecutorErrorStop(t *testing.T) {
	errorMsg := "error stop"
	g := NewGexecutor()
	g.Run(func(g *Gexecutor) error {
		return errors.New(errorMsg)
	})

	if err := g.Wait(); err != nil && err.Error() != errorMsg {
		t.Errorf("error stop")
	}
	g.Quit()
}
