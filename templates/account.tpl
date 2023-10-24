accounts: {
{{- range .}}
    {{ .Name}} {
        jetstream: enabled
        users: [
        {{- range .Users}}
            { nkey: {{ .Nkey }} }
        {{- end }}
        ]
    }
{{- end }}
}
