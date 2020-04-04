package thereum

import "time"

// Rough notes on the types of delays I eventually want

// Delayer wraps the delay method for pausing block commitment
type Delayer interface {
	Delay()
}

type ConstantTimeDelay struct {
	delay time.Duration
}

type AdjustedTimeDelay struct {
}

type PauseDelay struct {
}

func (d *PauseDelay) RunSwitch(input <-chan struct{}) {

}

func (d *PauseDelay) Delay() {

}
