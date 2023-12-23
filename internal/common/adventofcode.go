package common

import (
	"errors"
	"fmt"
)

// BaseExercise is the base struct for all exercises.
type BaseExercise struct{}

// One is the first part of the exercise.
func (c BaseExercise) One(instr string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// Two is the second part of the exercise.
func (c BaseExercise) Two(instr string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

// Vis is the visualization of the exercise.
func (c BaseExercise) Vis(instr string, outdir string) error {
	return errors.New("not implemented")
}

// Close is called when the exercise is done.
func Close() {
	if recover() != nil {
		fmt.Printf("PANIC: %v\n", recover())
	}
}
