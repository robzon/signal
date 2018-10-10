Signal
======

Simple implementation of signal subscription/emission for Go.

Rough usage example:

```go
import "github.com/robzon/signal"

type Dog struct {
  signal.Emitter
  name string
}

func (d *Dog) SetName(name string) {
  d.name = name
  d.Emit("name_changed") // emit a signal to all subscribers
}

func main() {
  d := &Dog{}

  ch := d.Subscribe() // subsribe to all signals on the object
  // ... OR ...
  ch := d.SubscribeSignal("name_changed") // Subscribe to a specific signal

  go func(ch chan string) {
    for sig := range ch {
      fmt.Println(sig, "received")
    }
  }(ch)

  // ...

  d.Unsubscribe(ch) // remove subscription and close the channel
}
```
