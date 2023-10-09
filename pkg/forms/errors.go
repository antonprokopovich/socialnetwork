package forms

type formErrors map[string][]string

// Add adds error messages for a given field to the map.
func (e formErrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves the first error message for a given field from the map
func (e formErrors) Get(field string) string {
	eMsg := e[field]
	if len(eMsg) == 0 {
		return ""
	}
	return eMsg[0]
}
