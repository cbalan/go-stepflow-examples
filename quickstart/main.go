package main

import (
	"context"
	"fmt"
	
	"github.com/cbalan/go-stepflow"
)

func step1(ctx context.Context) error {
	fmt.Println("Executing step 1")
	return nil
}

func aCondition(ctx context.Context) (bool, error) {
	fmt.Println("Return true when ready to proceed.")
	return true, nil
}

func step2(ctx context.Context) error {
	fmt.Println("Executing step 2")
	return nil
}

func main() {
	// Define workflow
	flow, err := stepflow.NewStepFlow(stepflow.Steps().WithName("quickstart.v1").
		Do("step1", step1).
		WaitFor("aCondition", aCondition).
		Do("step2", step2))

	if err != nil {
		panic(err)
	}

	// Execute workflow
	var state []string // Could be loaded from persistent storage
	for !flow.IsCompleted(state) {
		state, err = flow.Apply(context.Background(), state)
		if err != nil {
			panic(err)
		}

		fmt.Println("new state: ", state)
		// Persist state here if needed
	}
}
