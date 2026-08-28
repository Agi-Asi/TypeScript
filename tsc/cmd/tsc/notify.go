// Signal-handling helpers for the TypeScript CLI.

package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
)

// notifyContext is like signal.NotifyContext, except on wasip1, where no
// signal can ever be delivered. Go's signal watcher goroutine busy-spins on
// wasip1 (it never gets to idle in single-threaded Wasm), which starves the
// scheduler and prevents e.g. the LSP stdin reader from ever waking. Since
// cancellation can only come from the returned CancelFunc there, an ordinary
// context.WithCancel is equivalent and avoids the busy loop entirely.
func notifyContext(parent context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	if runtime.GOOS == "wasip1" {
		return context.WithCancel(parent)
	}
	return signal.NotifyContext(parent, sigs...)
}
