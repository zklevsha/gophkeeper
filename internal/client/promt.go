package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
)

type fn func(string) error

// getYN returns true if user answered Yes and false otherwise
func getYN(label string) bool {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func inputSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	return result, err
}

func getInput(label string, validator fn, mask bool) (string, error) {

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
	return result, err
}

func getTags(input string) (map[string]string, error) {
	var t map[string]string
	if input != "" {
		err := json.Unmarshal([]byte(input), &t)
		if err != nil {
			return t, fmt.Errorf("cant parse tags: %s", err.Error())
		}
	}
	return t, nil
}
