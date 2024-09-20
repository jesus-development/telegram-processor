package scenarios

import "context"

type (
	// Scenario is kind of saga
	Scenario interface {
		Run(ctx context.Context) error
	}

	Step struct {
		Name        string
		Description string
		Skip        bool
		SkipOnFail  bool
		Func        func(ctx context.Context, scenario Scenario, step *Step) error
	}
)
