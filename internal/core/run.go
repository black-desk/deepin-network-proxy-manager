package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/black-desk/deepin-network-proxy-manager/internal/core/monitor"
	"github.com/black-desk/deepin-network-proxy-manager/internal/core/repeater"
	"github.com/black-desk/deepin-network-proxy-manager/internal/core/rulemanager"
	"github.com/black-desk/deepin-network-proxy-manager/internal/core/watcher"
	. "github.com/black-desk/lib/go/errwrap"
)

func (c *Core) Run() (err error) {
	defer Wrap(&err, "Error occurs while running the core.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	c.start()

	go func() {
		select {
		case <-c.ctx.Done():
		case <-sigChan:
		}
		c.cancel()
	}()

	return c.pool.Wait()
}

func (c *Core) start() {
	c.pool.Go(c.runWatcher)
	c.pool.Go(c.runMonitor)
	c.pool.Go(c.runRuleManager)

	if c.cfg.Repeater == nil {
		return
	}

	c.pool.Go(c.runRepeater)

	return
}

func (c *Core) runMonitor() (err error) {
	defer c.cancel()

	var m *monitor.Monitor
	m, err = injectedMonitor(c)
	if err != nil {
		return
	}

	err = m.Run()
	return
}

func (c *Core) runRuleManager() (err error) {
	defer c.cancel()

	var r *rulemanager.RuleManager
	r, err = injectedRuleManager(c)
	if err != nil {
		return
	}

	err = r.Run()
	return
}

func (c *Core) runRepeater() (err error) {
	defer c.cancel()

	var r *repeater.Repeater
	r, err = injectedRepeater(c)

	if err != nil {
		return
	}

	err = r.Run()
	return
}

func (c *Core) runWatcher() (err error) {
	defer c.cancel()

	var w *watcher.Watcher
	w, err = injectedWatcher(c)

	if err != nil {
		return
	}

	err = w.Run()
	return
}
