package mermaidgen

import (
	"context"
	"github.com/cbalan/go-stepflow"
)

func noop(_ context.Context) error {
	return nil
}

func alwaysTrue(_ context.Context) (bool, error) {
	return true, nil
}

func alwaysFalse(_ context.Context) (bool, error) {
	return false, nil
}

func alwaysRetry(_ context.Context, _ error) (bool, error) {
	return true, nil
}

//go:generate go run gen.go
func DemoSteps() *stepflow.StepsSpec {
	basicSteps := stepflow.Steps().Retry("demo", alwaysRetry, stepflow.Steps().
		Do("importantActivity", noop).
		LoopUntil("visitAllItems", alwaysFalse, stepflow.Steps().
			Do("step1", noop).
			WaitFor("step1ToBeDone", alwaysTrue)))

	return stepflow.Steps().
		WithName("deploy.v1").
		Do("getPreActualState", noop).
		WaitFor("preActualState", alwaysTrue).
		Case("shouldDeploy", alwaysTrue, stepflow.Steps().
			Do("createDeployRequest", noop).
			WaitFor("acceptedDeployRequest", alwaysTrue).
			WaitFor("monitorDeploymentProgress", alwaysTrue).
			Do("getPostActualState", noop).
			WaitFor("postActualState", alwaysTrue).
			Do("validatePostDeployActualState", noop)).
		Steps("basicSteps", basicSteps)

}
