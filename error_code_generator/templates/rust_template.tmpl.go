package templates

var (
	TemplateRust = `
    pub const SERVICE_ABBREVIATION: &str = "{{.ServiceAbbreviation}}";

    // {{.ServiceName}} service-specific error codes
    pub mod {{.ServiceName | ToLower}} {
    {{ range .ServiceErrors }}
        pub const {{.Code}}: &str = concat!(SERVICE_ABBREVIATION, "_{{.Code}}");
    {{ end }}
    }

    pub static ERROR_MESSAGES: phf::Map<&'static str, &'static str> = phf_map! {
    {{ range .ServiceErrors }}
        {{.Code}} => "{{.Message}}",
    {{ end }}
    };

    pub static ERROR_CATEGORIES: phf::Map<&'static str, &'static str> = phf_map! {
    {{ range .ServiceErrors }}
        {{.Code}} => "{{.Category}}",
    {{ end }}
    };

    `
)
