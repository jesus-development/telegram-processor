package prompt

type Prompter interface {
	StringPrompt(label string) (string, error)
	YesNoPrompt(label string, def bool) (bool, error)
}
