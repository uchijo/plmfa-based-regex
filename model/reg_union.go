package model

import (
	"fmt"

	"github.com/google/uuid"
)

type RegUnion struct {
	Left  RegExp
	Right RegExp
}

func (ru RegUnion) States(startId string) ([]State, string, error) {
	if ru.Left == nil || ru.Right == nil {
		panic("invalid regex.")
	}
	// Left, Right毎にoutputのidが指定されるので、イプシロン動作でいい感じに辻褄合わせる
	states := []State{}
	leftStart, _ := uuid.NewRandom()
	rightStart, _ := uuid.NewRandom()
	start := State{
		Id:    startId,
		IsEnd: false,
		Moves: []Move{
			{
				MType:  Epsilon,
				MoveTo: leftStart.String(),
				Input:  nil,
			},
			{
				MType:  Epsilon,
				MoveTo: rightStart.String(),
				Input:  nil,
			},
		},
	}
	states = append(states, start)

	// TODO: エラー処理ちゃんとする
	leftStates, leftGoal, _ := ru.Left.States(leftStart.String())
	states = append(states, leftStates...)
	rightStates, rightGoal, _ := ru.Right.States(rightStart.String())
	states = append(states, rightStates...)

	goalUUID, _ := uuid.NewRandom()

	leftGoalState := State{
		Id:    leftGoal,
		IsEnd: false,
		Moves: []Move{
			{
				MType:  Epsilon,
				Input:  nil,
				MoveTo: goalUUID.String(),
			},
		},
	}
	rightGoalState := State{
		Id:    rightGoal,
		IsEnd: false,
		Moves: []Move{
			{
				MType:  Epsilon,
				Input:  nil,
				MoveTo: leftGoal,
			},
		},
	}
	states = append(states, leftGoalState, rightGoalState)

	return states, goalUUID.String(), nil
}

func NewUnionFromSlice(contents []RegExp) RegExp {
	length := len(contents)
	if length == 0 {
		return RegSkip{}
	}
	if length == 1 {
		return contents[0]
	}
	if length == 2 {
		return RegUnion{
			Right: contents[0],
			Left:  contents[1],
		}
	}
	fmt.Println(contents[1:])
	return RegUnion{
		Right: contents[0],
		Left: NewUnionFromSlice(contents[1:]),
	}
}
