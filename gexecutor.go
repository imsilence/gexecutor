package gexecutor

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Runner 运行者
type Runner func(g *Gexecutor) error

// SignalStop 退出信号
type SignalStop struct{}

// Gexecutor Goroutine执行器
type Gexecutor struct {
	wg sync.WaitGroup

	stopCh chan SignalStop
	errCh  chan error
}

// NewGexecutor 构造器
func NewGexecutor() *Gexecutor {
	return &Gexecutor{
		stopCh: make(chan SignalStop, 1),
		errCh:  make(chan error, 1),
	}
}

// StopCh 获取停止信号
func (g *Gexecutor) StopCh() chan SignalStop {
	return g.stopCh
}

// Begin 开始
func (g *Gexecutor) Begin() {
	g.wg.Add(1)
}

// End 结束
func (g *Gexecutor) End() {
	g.wg.Done()
}

// AddError 运行结束
func (g *Gexecutor) AddError(err error) {
	g.errCh <- err
}

// Run 运行
func (g *Gexecutor) Run(runner Runner) {
	g.Begin()
	go func() {
		defer g.End()
		g.AddError(runner(g))
	}()
}

// Next 执行
func (g *Gexecutor) Next(d time.Duration) bool {
	select {
	case <-time.After(d):
		return true
	case <-g.stopCh:
		return false
	}
}

// Wait 等待终端结束信号
func (g *Gexecutor) Wait() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-interrupt:
			return nil
		case err := <-g.errCh:
			if err != nil {
				return err
			}
		}
	}
}

// Quit 退出执行器
func (g *Gexecutor) Quit() {
	close(g.stopCh)
	g.wg.Wait()
}
