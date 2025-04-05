package keypad

import (
	"machine"
	"time"
)

const debounce = 30    // ms
const doubleTap = 1000 // ms

const codeAlt = ZX

type Event struct {
	code      Code
	pressed   bool
	timestamp int64
}

type EventType int

const (
	Pressed EventType = 1 << iota
	Double
	Alt
)

func (et EventType) Pressed() bool {
	return et&Pressed > 0
}

func (et EventType) Released() bool {
	return et&Pressed == 0
}

func (et EventType) Single() bool {
	return et&Double == 0
}

func (et EventType) Double() bool {
	return et&Double > 0
}

func (et EventType) Alt() bool {
	return et&Alt > 0
}

func (et EventType) NoAlt() bool {
	return et&Alt == 0
}

type State struct {
	pressed    bool
	debouncing bool
	timestamp  int64
}

type Keypad struct {
	rows []machine.Pin
	cols []machine.Pin

	handlers []func(EventType) bool

	state []State
	alt   bool
	e2    *Event
	e1    *Event
}

// NewKeypad creates a new Keypad instance, and configures input and output pins
// of a keyboard matrix.
func NewKeypad() Keypad {
	rows := []machine.Pin{machine.PA04, machine.PA28, machine.PA16}
	cols := []machine.Pin{machine.PA23, machine.PA22, machine.PA19, machine.PA18, machine.PA01}

	for _, row := range rows {
		row.Configure(machine.PinConfig{Mode: machine.PinInput})
	}
	for _, col := range cols {
		col.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	keypad := Keypad{rows: rows, cols: cols, handlers: nil}
	keypad.clearState()

	return keypad
}

func (k *Keypad) scan() *Event {
	var event *Event

	id := 0
	timestamp := time.Now().UnixMilli()

	for _, row := range k.rows {
		row.Configure(machine.PinConfig{Mode: machine.PinOutput})
		row.Low()

		for _, col := range k.cols {
			state := &k.state[id]
			pressed := !col.Get()

			if pressed == state.pressed {
				if state.debouncing {
					// stop debouncing (false positive)
					state.debouncing = false
					state.timestamp = 0
				}
			} else {
				if state.debouncing {
					if timestamp >= state.timestamp+debounce {
						// stop debouncing
						state.pressed = !state.pressed
						state.debouncing = false
						state.timestamp = 0
						event = &Event{code: Code(id), pressed: state.pressed, timestamp: timestamp}
						break
					}
				} else {
					// start debouncing
					state.debouncing = true
					state.timestamp = timestamp
				}
			}

			id += 1
		}

		row.Configure(machine.PinConfig{Mode: machine.PinInput})

		if event != nil {
			break // exit early
		}
	}

	return event
}

func (k *Keypad) clearState() {
	k.state = make([]State, len(k.rows)*len(k.cols))
	k.alt = false
	k.e2 = nil
	k.e1 = nil
}

// Scan checks whether a key was pressed or released (one or more times), and
// calls the respective event handler. This is different from the internal
// Keypad.scan function that is responsible for checking the physical state of
// all buttons.
//
// Scan also updates the internal alt value, so the ZX key can act like a Fn key
// on a normal keyboard. Each app should have at least one handler, that will
// call peripherals.SoftReset when EventType.Alt and EventType.Released are
// true.
func (k *Keypad) Scan() bool {
	handled := false
	e0 := k.scan()

	if k.handlers != nil && e0 != nil {
		handler := k.handlers[e0.code]
		et := EventType(0)

		switch true {
		case k.e2 != nil && !e0.pressed && !k.e2.pressed && e0.code == k.e2.code && e0.timestamp <= k.e2.timestamp+doubleTap:
			et = Double // Released
		case k.e2 != nil && e0.pressed && k.e2.pressed && e0.code == k.e2.code && e0.timestamp <= k.e2.timestamp+doubleTap:
			et = Double | Pressed
		case !e0.pressed:
			if k.alt && e0.code == codeAlt {
				handler = nil
				k.alt = false
			}
			// et = Released
		case e0.pressed:
			if k.e1 != nil && k.e1.code == codeAlt && k.e1.pressed {
				k.alt = true
			}
			et = Pressed
		}

		if k.alt {
			et |= Alt
		}

		if handler != nil {
			handled = handler(et)
		}

		k.e2 = k.e1
		k.e1 = e0
	}

	return handled
}

// SetHandlers changes the handlers list for a Keypad instance and resets the
// internal state for a Keypad. Ideally, it should be done when the app is idle
// and no key presses are performed.
func (k *Keypad) SetHandlers(handlers []func(EventType)) {
	boolHandlers := make([]func(EventType) bool, len(handlers))

	for i, handler := range handlers {
		if handler != nil {
			boolHandlers[i] = func(et EventType) bool {
				handler(et)
				return true
			}
		}
	}

	k.SetBoolHandlers(boolHandlers)
}

// SetBoolHandlers changes the handlers list for a Keypad instance and resets
// the internal state for a Keypad. Ideally, it should be done when the app is
// idle and no key presses are performed.
//
// Compared to SetHandlers, this function accepts a list of functions returning
// a boolean value, which is useful for conditional updating of app's state.
func (k *Keypad) SetBoolHandlers(handlers []func(EventType) bool) {
	k.handlers = handlers
	k.clearState()
}
