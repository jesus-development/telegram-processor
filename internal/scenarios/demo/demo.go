package demo

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"telegram-processor/internal/scenarios"
	"telegram-processor/internal/services/processor"
	"telegram-processor/pkg/cli/prompt"
	"unicode/utf8"
)

const (
	MAX_QUERY_LENGTH = 256
	SEARCH_LIMIT     = 10

	DEFAULT_IMPORT_FILE = "resources/demo/chat-export-30news.json"
)

type (
	demoScenario struct {
		Name        string
		Description string

		Processor processor.MessageProcessor
		Prompter  prompt.Prompter

		Steps []*scenarios.Step
	}
)

// NewDemoScenario if importDb is false - skips import and calc steps
func NewDemoScenario(processor processor.MessageProcessor, prompter prompt.Prompter, importDb bool) *demoScenario {
	return &demoScenario{
		Name:        "Demo",
		Description: "Basic demonstration of all features",
		Processor:   processor,
		Prompter:    prompter,
		Steps: []*scenarios.Step{
			{
				Name:        "Import messages",
				Description: "Import messages to the database from json file",
				Skip:        !importDb,
				Func:        importMessages,
			},
			{
				Name:        "Calc embeddings",
				Description: "Send messages to the OpenAI API and calculate their embeddings",
				Skip:        !importDb,
				Func:        calcEmbeddings,
			},
			{
				Name:        "Search",
				Description: "Search for similar messages in the database",
				Func:        search,
			},
		},
	}
}

func (s *demoScenario) Run(ctx context.Context) error {
	for _, step := range s.Steps {
		if step.Skip {
			continue
		}

		slog.Info("Running step", "name", step.Name, "description", step.Description)
		if err := step.Func(ctx, s, step); err != nil {
			if step.SkipOnFail {
				slog.Warn("Skipping step", "name", step.Name, "error", err)
				continue
			}

			return err
		}

		slog.Info("Finished step", "name", step.Name)
	}

	return nil
}

func importMessages(ctx context.Context, scenario scenarios.Scenario, _ *scenarios.Step) error {
	count, err := scenario.(*demoScenario).Processor.GetCount(ctx)
	if err != nil {
		return fmt.Errorf("s.(*demoScenario).Processor.GetCount -> %w", err)
	}
	if count > 0 {
		yes, err := scenario.(*demoScenario).Prompter.YesNoPrompt("Database is not empty. Do you really want to import more messages?", false)
		if err != nil {
			return fmt.Errorf("s.(*demoScenario).Prompter.YesNoPrompt -> %w", err)
		}
		if !yes {
			return nil
		}
	}

	r, err := os.Open(DEFAULT_IMPORT_FILE)
	if err != nil {
		return fmt.Errorf("os.Open -> %w", err)
	}
	defer r.Close()

	err = scenario.(*demoScenario).Processor.ImportJson(ctx, r)
	if err != nil {
		return fmt.Errorf("s.(*demoScenario).Processor.ImportJson -> %w", err)
	}

	return nil
}

func calcEmbeddings(ctx context.Context, scenario scenarios.Scenario, _ *scenarios.Step) error {
	return scenario.(*demoScenario).Processor.CalculateAndSaveEmbeddings(ctx)
}

func search(ctx context.Context, scenario scenarios.Scenario, _ *scenarios.Step) error {
	for {
		query, err := scenario.(*demoScenario).Prompter.StringPrompt("Search query:")
		if err != nil {
			return fmt.Errorf("s.(*demoScenario).Prompter.StringPrompt -> %w", err)
		}

		if query == "" {
			break
		}

		if utf8.RuneCountInString(query) > MAX_QUERY_LENGTH {
			slog.Warn("Query is too long", "max length", MAX_QUERY_LENGTH)
			continue
		}

		messages, err := scenario.(*demoScenario).Processor.GetClosest(ctx, query, SEARCH_LIMIT)
		if err != nil {
			return fmt.Errorf("s.(*demoScenario).Processor.GetClosest -> %w", err)
		}

		slog.Info("Found messages", "count", len(messages))

		fmt.Println()
		for i, message := range messages {
			fmt.Printf("%d. %f %s\n", i+1, message.Similarity, message.Text)
		}
		fmt.Println()
	}
	return nil
}
