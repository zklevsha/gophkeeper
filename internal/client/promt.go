package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

type fn func(string) error

func getYN(label string) string {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func inputSelect(label string, items []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func getInput(label string, validator fn, mask bool) string {

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
