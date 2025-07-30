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

func quickStartStepFlow() (stepflow.StepFlow, error) {
	return stepflow.NewStepFlow(stepflow.Steps().WithName("quickstart.v1").
		Do("step1", step1).
		WaitFor("aCondition", aCondition).
		Do("step2", step2))
}

func main() {

	// Workflow definition.
	flow, err := quickStartStepFlow()
	if err != nil {
		panic(err)
	}

	// Workflow execution.
	var state []string
	i := 0
	for !flow.IsCompleted(state) {
		// Could load state from persistent storage.
		fmt.Printf("[%d] old state: %v \n", i, state)

		// Apply workflow on the old state
		state, err = flow.Apply(context.Background(), state)
		if err != nil {
			panic(err)
		}

		// Could save state to persistent storage.
		fmt.Printf("[%d] new state: %v \n\n", i, state)

		i++
	}
}
