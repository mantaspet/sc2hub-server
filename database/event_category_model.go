package database

type EventCategory struct {
	ID       int
	Name     string
	Pattern  string
	InfoURL  string
	ImageURL string
	Priority int
}

func (ec EventCategory) Validate() map[string]string {
	errors := make(map[string]string)
	if len(ec.Name) < 1 {
		errors["Name"] = "Field is required"
	} else if err := fieldExists("event_categories", "name", ec.Name); err != nil {
		errors["Name"] = err.Error()
	}
	if len(ec.Pattern) < 1 {
		errors["Pattern"] = "Field is required"
	} else if err := fieldExists("event_categories", "pattern", ec.Pattern); err != nil {
		errors["Pattern"] = err.Error()
	}
	if ec.Priority == 0 {
		errors["Priority"] = "Field is required"
	} else if ec.Priority < 0 {
		errors["Priority"] = "Priority must be above zero"
	} else if err := fieldExists("event_categories", "priority", ec.Priority); err != nil {
		errors["Priority"] = err.Error()
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
