// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package capacity

type Level uint8

const (
	LevelZero Level = iota
	LevelMinimal
	LevelReduced
	LevelNormal
	LevelMax
)

const LevelCount = LevelMax + 1

func (v Level) DefaultPercent() int {
	// 0, 25, 75, 100, 125
	return v.ChooseInt([...]int{0, 20, 60, 80, 100})
}

func (v Level) ChooseInt(options [LevelCount]int) int {
	return options[v]
}
