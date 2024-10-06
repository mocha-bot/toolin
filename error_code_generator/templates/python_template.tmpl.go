package templates

var (
	TemplatePython = `
    # Generated Python Code

    class ErrorCodes:
        PACKAGE_NAME = "{{ .PackageName }}"
        SERVICE_NAME = "{{ .ServiceName }}"
        SERVICE_ABBREVIATION = "{{ .ServiceAbbreviation }}"

        SERVICE_ERRORS = {
            {{ range .ServiceErrors }}
            "{{ .Code }}": "{{ .Message }}",
            {{ end }}
        }
    `
)
