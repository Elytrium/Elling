package elling

type SmallTickEvent struct{}

type BigTickEvent struct{}

func DoSmallTick() {
	DispatchEvent(SmallTickEvent{})
}

func DoBigTick() {
	DispatchEvent(BigTickEvent{})
}
