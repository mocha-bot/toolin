package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/mocha-bot/toolin/error_code_generator/templates"
	"gopkg.in/yaml.v2"
)

const (
	packageName = "error_codes"
)

var (
	templateMap = map[string]string{
		"go":     templates.TemplateGo,
		"python": templates.TemplatePython,
		"rust":   templates.TemplateRust,
	}

	languageFileExtMap = map[string]string{
		"go":     "go",
		"python": "py",
		"rust":   "rs",
	}
)

type ErrorDefinition struct {
	Code     string         `yaml:"code"`
	Message  string         `yaml:"message"`
	Category string         `yaml:"category"`
	Context  []ContextField `yaml:"context,omitempty"`
}

type ContextField struct {
	Field       string `yaml:"field"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}

type Contract struct {
	Abbreviation string            `yaml:"abbreviation"`
	Name         string            `yaml:"name"`
	Errors       []ErrorDefinition `yaml:"errors"`
}

func generateCode(contract Contract, lang string) (string, error) {
	tmpl, err := parseLanguageTemplate(lang)
	if err != nil {
		return "", err
	}

	data := struct {
		PackageName         string
		ServiceName         string
		ServiceAbbreviation string
		ServiceErrors       []ErrorDefinition
	}{
		PackageName:         packageName,
		ServiceName:         contract.Name,
		ServiceAbbreviation: contract.Abbreviation,
		ServiceErrors:       contract.Errors,
	}

	t, err := template.New("code").Funcs(template.FuncMap{
		"ToLower": func(s string) string {
			return strings.ToLower(s)
		},
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var output strings.Builder
	err = t.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func parseLanguageTemplate(lang string) (string, error) {
	tmpl, exists := templateMap[lang]
	if !exists {
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	return tmpl, nil
}

func parseLanguageExtFile(lang string) (string, error) {
	fileExtension, exists := languageFileExtMap[lang]
	if !exists {
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	return fileExtension, nil
}

func printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s --contract <path> --language <go|python|rust>\n", os.Args[0])
	fmt.Println("\nOptions:")
	fmt.Println("  --contract   Path to the contract YAML file (default: contract_error.yml)")
	fmt.Println("  --language   Programming language for code generation (go, python, rust)")
}

func main() {
	// Define flags for the contract file and language
	contractFile := flag.String("contract", "contract_error.yml", "Path to the contract YAML file")
	lang := flag.String("language", "", "Programming language for code generation (go, python, rust)")

	// Parse the flags
	flag.Parse()

	// Check for required arguments
	if *lang == "" {
		log.Println("Error: --language is required.")
		printUsage()
		os.Exit(1)
	}

	// Read the YAML file
	yamlFile, err := os.ReadFile(*contractFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var contract Contract
	err = yaml.Unmarshal(yamlFile, &contract)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	code, err := generateCode(contract, *lang)
	if err != nil {
		log.Fatalf("Error generating code: %v", err)
	}

	// Write the generated code to the appropriate file
	fileExtension, err := parseLanguageExtFile(*lang)
	if err != nil {
		log.Fatalf("Error parsing language file extension: %v", err)
	}

	outputFileName := fmt.Sprintf("error_code.%s", fileExtension)

	err = os.WriteFile(outputFileName, []byte(code), 0644)
	if err != nil {
		log.Fatalf("Error writing code file: %v", err)
	}

	fmt.Printf("Code generated successfully in %s\n", outputFileName)
}
