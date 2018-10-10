package signal

import (
	"sync"
)

// Emitter holds information about the subscribers
type Emitter struct {
	objectSubscribers []chan string
	signalSubscribers map[string][]chan string
	emitterMutex      sync.Mutex
}

// Subscribe returns a new channel to which all object signals will be sent
func (e *Emitter) Subscribe() chan string {
	e.emitterMutex.Lock()
	defer e.emitterMutex.Unlock()

	if e.objectSubscribers == nil {
		e.objectSubscribers = make([]chan string, 0)
	}

	ch := make(chan string, 1)
	e.objectSubscribers = append(e.objectSubscribers, ch)

	return ch
}

// SubscribeSignal returns a channel to which the specified signal will
// be sent
func (e *Emitter) SubscribeSignal(sig string) chan string {
	e.emitterMutex.Lock()
	defer e.emitterMutex.Unlock()

	if e.signalSubscribers == nil {
		e.signalSubscribers = make(map[string][]chan string, 0)
	}

	if e.signalSubscribers[sig] == nil {
		e.signalSubscribers[sig] = make([]chan string, 0)
	}

	ch := make(chan string, 1)
	e.signalSubscribers[sig] = append(e.signalSubscribers[sig], ch)

	return ch
}

// Unsubscribe() removes a channel from subscriptions and closes it
func (e *Emitter) Unsubscribe(ch chan string) {
	e.emitterMutex.Lock()
	defer e.emitterMutex.Unlock()

	if e.objectSubscribers != nil {
		for i, v := range e.objectSubscribers {
			if v == ch {
				e.objectSubscribers = append(e.objectSubscribers[:i], e.objectSubscribers[i+1:]...)
				close(ch)
				return
			}
		}
	}

	if e.signalSubscribers != nil {
		for sig := range e.signalSubscribers {
			for i, v := range e.signalSubscribers[sig] {
				if v == ch {
					e.signalSubscribers[sig] = append(e.signalSubscribers[sig][:i], e.signalSubscribers[sig][i+1:]...)
					close(ch)
					return
				}
			}
		}
	}
}

// Emit() emits a signal by the object in a goroutine
func (e *Emitter) Emit(signal string) {
	go e.emit(signal)
}

// emit sends a signal to all subscribers in a thread-safe, locking way
func (e *Emitter) emit(signal string) {
	e.emitterMutex.Lock()
	defer e.emitterMutex.Unlock()

	for _, ch := range e.objectSubscribers {
		ch <- signal
	}

	if (e.signalSubscribers != nil) && (e.signalSubscribers[signal] != nil) {
		for _, ch := range e.signalSubscribers[signal] {
			ch <- signal
		}
	}
}
