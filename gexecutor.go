package gexecutor

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Runner 运行者
type Runner func(g *Gexecutor) error

// Gexecutor Goroutine执行器
type Gexecutor struct {
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	errCh chan error
}

// NewGexecutor 构造器
func NewGexecutor() *Gexecutor {
	ctx, cancel := context.WithCancel(context.Background())

	return &Gexecutor{
		ctx:    ctx,
		cancel: cancel,
		errCh:  make(chan error, 1),
	}
}

// Done 获取停止信号
func (g *Gexecutor) Done() <-chan struct{} {
	return g.ctx.Done()
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
	if err == nil {
		return
	}
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
	case <-g.Done():
		return false
	}
}

// Wait 等待终端结束信号
func (g *Gexecutor) Wait() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-interrupt:
		return nil
	case err := <-g.errCh:
		return err
	}
}

// Quit 退出执行器
func (g *Gexecutor) Quit() {
	g.cancel()
	g.wg.Wait()
}
