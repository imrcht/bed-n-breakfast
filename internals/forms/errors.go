package forms

type errors map[string][]string

// * Add appends the message to errors in that field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// * Get returns first message of errors of that field
func (e errors) Get(field string) string {
	es := e[field]

	if len(es) == 0 {
		return ""
	}

	return es[0]
}
