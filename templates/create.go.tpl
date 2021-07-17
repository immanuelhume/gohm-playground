func ({{ toReceiver .Name }} *{{ .Name }}) Create({{ toParams .Fields }}) {{ .Package }}.{{ .Name }} {
	