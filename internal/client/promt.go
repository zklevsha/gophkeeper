package client

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

type fn func(string) error

func promptGetInput(label string, validator fn, mask bool) string {

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     label,
		Templates: templates,
		Validate:  promptui.ValidateFunc(validator),
	}

	if mask {
		prompt.Mask = '*'
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}
