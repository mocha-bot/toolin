package templates

var (
	TemplateGo = `
    package {{.PackageName}}

    import "fmt"

    // ServiceAbbreviation represents the abbreviation for the {{.ServiceName}} service
    const ServiceAbbreviation = "{{.ServiceAbbreviation}}"

    // {{.ServiceName}} service-specific error codes
    const (
    {{ range .ServiceErrors }}
        Err{{.Code}} = ServiceAbbreviation + "_{{.Code}}"
    {{ end }}
    )

    // ErrorMessages holds the mapping of error codes to human-readable messages
    var ErrorMessages = map[string]string{
    {{ range .ServiceErrors }}
        Err{{.Code}}: "{{.Message}}",
    {{ end }}
    }

    // ErrorCategories holds the mapping of error codes to their categories
    var ErrorCategories = map[string]string{
    {{ range .ServiceErrors }}
        Err{{.Code}}: "{{.Category}}",
    {{ end }}
    }

    // GenerateDynamicError creates a formatted error message based on the error code and context
    func GenerateDynamicError(code string, context map[string]interface{}) string {
        messageTemplate, exists := ErrorMessages[code]
        if !exists {
            return "Unknown error code: " + code
        }

        // Replace placeholders with actual context values
        for _, value := range context {
            messageTemplate = fmt.Sprintf(messageTemplate, value)
        }

        return messageTemplate
    }
  `
)
