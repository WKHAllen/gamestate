package gamestate

import (
	"reflect"
	"testing"
)

// Assert a condition
func assert(value bool, t *testing.T, err string) {
	if !value {
		t.Errorf(err)
		panic(err)
	}
}

// Assert an equality
func assertEq[T any](left T, right T, t *testing.T) {
	if !reflect.DeepEqual(left, right) {
		t.Errorf("Assertion error: %#v != %#v", left, right)
		panic("Assertion error")
	}
}

// Assert no error occurred
func assertNoErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf(err.Error())
		panic(err.Error())
	}
}

func TestTransitions(t *testing.T) {
	ch := make(chan int, 5)

	action1 := func() any { ch <- 1; return 1 }
	action2 := func() any { ch <- 2; return 2 }
	action3 := func() any { ch <- 3; return 3 }
	action4 := func() any { ch <- 4; return 4 }
	action5 := func() any { ch <- 5; return 5 }

	state1 := NewGameState(action1)
	state2 := NewGameState(action2)
	state3 := NewGameState(action3)
	state4 := NewGameState(action4)
	state5 := NewGameState(action5)

	state1.Transition(state3)
	state3.Transition(state4)
	state4.Transition(state2)
	state2.Transition(state5)

	lastRet, err := RunGame(state1, state5)
	assertNoErr(err, t)
	assertEq(lastRet, 5, t)

	assertEq(<-ch, 1, t)
	assertEq(<-ch, 3, t)
	assertEq(<-ch, 4, t)
	assertEq(<-ch, 2, t)
	assertEq(<-ch, 5, t)
}

func TestTransitionsOn(t *testing.T) {
	ch := make(chan int, 5)

	action1 := func() any { ch <- 1; return 1 }
	action2 := func() any { ch <- 2; return 2 }
	action3 := func() any { ch <- 3; return 3 }
	action4 := func() any { ch <- 4; return 4 }
	action5 := func() any { ch <- 5; return 5 }

	state1 := NewGameState(action1)
	state2 := NewGameState(action2)
	state3 := NewGameState(action3)
	state4 := NewGameState(action4)
	state5 := NewGameState(action5)

	state1.TransitionOn(state3, 1)
	state3.TransitionOn(state4, 3)
	state4.TransitionOn(state2, 4)
	state2.TransitionOn(state5, 2)

	lastRet, err := RunGame(state1, state5)
	assertNoErr(err, t)
	assertEq(lastRet, 5, t)

	assertEq(<-ch, 1, t)
	assertEq(<-ch, 3, t)
	assertEq(<-ch, 4, t)
	assertEq(<-ch, 2, t)
	assertEq(<-ch, 5, t)
}

func TestGetSetGameValues(t *testing.T) {
	SetGameValue("total", 0)

	action1 := func() any {
		SetGameValue("number1", 1)
		total, ok := GetGameValue("total")
		assert(ok, t, "could not get total value")
		SetGameValue("total", total.(int) + 1)
		return 1
	}
	action2 := func() any {
		SetGameValue("number2", 2)
		total, ok := GetGameValue("total")
		assert(ok, t, "could not get total value")
		SetGameValue("total", total.(int) + 2)
		return 2
	}
	action3 := func() any {
		SetGameValue("number3", 3)
		total, ok := GetGameValue("total")
		assert(ok, t, "could not get total value")
		SetGameValue("total", total.(int) + 3)
		return 3
	}
	action4 := func() any {
		SetGameValue("number4", 4)
		total, ok := GetGameValue("total")
		assert(ok, t, "could not get total value")
		SetGameValue("total", total.(int) + 4)
		return 4
	}
	action5 := func() any {
		SetGameValue("number5", 5)
		total, ok := GetGameValue("total")
		assert(ok, t, "could not get total value")
		SetGameValue("total", total.(int) + 5)
		return 5
	}

	state1 := NewGameState(action1)
	state2 := NewGameState(action2)
	state3 := NewGameState(action3)
	state4 := NewGameState(action4)
	state5 := NewGameState(action5)

	state1.Transition(state3)
	state3.Transition(state4)
	state4.Transition(state2)
	state2.Transition(state5)

	lastRet, err := RunGame(state1, state5)
	assertNoErr(err, t)
	assertEq(lastRet, 5, t)

	total, ok := GetGameValue("total")
	assert(ok, t, "could not get total value")
	assertEq(total.(int), 15, t)

	number1, ok := GetGameValue("number1")
	assert(ok, t, "could not get number1 value")
	assertEq(number1.(int), 1, t)
	number2, ok := GetGameValue("number2")
	assert(ok, t, "could not get number2 value")
	assertEq(number2.(int), 2, t)
	number3, ok := GetGameValue("number3")
	assert(ok, t, "could not get number3 value")
	assertEq(number3.(int), 3, t)
	number4, ok := GetGameValue("number4")
	assert(ok, t, "could not get number4 value")
	assertEq(number4.(int), 4, t)
	number5, ok := GetGameValue("number5")
	assert(ok, t, "could not get number5 value")
	assertEq(number5.(int), 5, t)

	DeleteGameValue("total")
	_, ok = GetGameValue("total")
	assert(!ok, t, "could not delete total value")
}

func TestMultipleEndings(t *testing.T) {
	start := func() any { return 0 }
	middle := func() any {
		if ending, ok := GetGameValue("ending"); ok {
			return ending
		} else {
			return 0
		}
	}
	end1 := func() any { return -1 }
	end2 := func() any { return -2 }
	end3 := func() any { return -3 }

	startState := NewGameState(start)
	middleState := NewGameState(middle)
	endState1 := NewGameState(end1)
	endState2 := NewGameState(end2)
	endState3 := NewGameState(end3)

	startState.Transition(middleState)
	middleState.TransitionOn(endState1, 1)
	middleState.TransitionOn(endState2, 2)
	middleState.TransitionOn(endState3, 3)

	SetGameValue("ending", 1)
	lastRet := RunGameToEnd(startState)
	assertEq(lastRet, -1, t)

	SetGameValue("ending", 2)
	lastRet = RunGameToEnd(startState)
	assertEq(lastRet, -2, t)

	SetGameValue("ending", 3)
	lastRet = RunGameToEnd(startState)
	assertEq(lastRet, -3, t)
}
