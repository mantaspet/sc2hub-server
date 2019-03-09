package models

type EventCategory struct {
	ID       int
	Name     string
	Pattern  string
	InfoURL  string
	ImageURL string
	Order    int
}

func (ec EventCategory) Validate() map[string]string {
	errors := make(map[string]string)
	if len(ec.Name) < 1 {
		errors["Name"] = "Field is required"
	}
	if len(ec.Pattern) < 1 {
		errors["Pattern"] = "Field is required"
	}
	if ec.Order == 0 {
		errors["Order"] = "Field is required"
	}
	if ec.Order < 0 {
		errors["Order"] = "Order must be above zero"
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
