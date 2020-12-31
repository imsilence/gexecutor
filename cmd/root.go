package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/imsilence/gexecutor"
)

func main() {
	go func() {
		<-time.After(15 * time.Second)
		fmt.Println(exec.Command("kill", strconv.Itoa(os.Getpid())).Output())

	}()
	g := gexecutor.NewGexecutor()
	g.Run(func(*gexecutor.Gexecutor) error {
		<-time.After(5 * time.Second)
		return nil
	})

	g.Begin()
	go func() {
		defer g.End()
		for {
			if !g.Next(3 * time.Second) {
				break
			}
			fmt.Println(time.Now())
		}
		// g.AddError(errors.New("xxx"))
	}()

	fmt.Println(g.Wait())
	g.Quit()
}
