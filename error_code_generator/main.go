package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
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
	CommonErrors []ErrorDefinition `yaml:"common_errors"`
	Service      struct {
		Abbreviation string            `yaml:"abbreviation"`
		Name         string            `yaml:"name"`
		Errors       []ErrorDefinition `yaml:"errors"`
	} `yaml:"service"`
}

func generateCode(contract Contract, lang string) (string, error) {
	var path string
	var tmpl string
	var err error

	templateDirPath := "templates/"

	switch lang {
	case "go":
		path = templateDirPath + "go_template.tmpl"
	case "python":
		path = templateDirPath + "python_template.tmpl"
	case "rust":
		path = templateDirPath + "rust_template.tmpl"
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	tmpl, err = loadTemplate(path)
	if err != nil {
		return "", err
	}

	data := struct {
		PackageName         string
		ServiceName         string
		ServiceAbbreviation string
		CommonErrors        []ErrorDefinition
		ServiceErrors       []ErrorDefinition
	}{
		PackageName:         "mocha_errors",
		ServiceName:         contract.Service.Name,
		ServiceAbbreviation: contract.Service.Abbreviation,
		CommonErrors:        contract.CommonErrors,
		ServiceErrors:       contract.Service.Errors,
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

func loadTemplate(filename string) (string, error) {
	tmpl, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %v", filename, err)
	}
	return string(tmpl), nil
}

func parseLanguageExtFile(lang string) (string, error) {
	languageFileMap := map[string]string{
		"go":     "go",
		"python": "py",
		"rust":   "rs",
	}

	fileExtension, exists := languageFileMap[lang]
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
