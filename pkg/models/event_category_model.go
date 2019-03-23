package models

type EventCategory struct {
	ID          int
	Name        string
	Pattern     string
	InfoURL     string
	ImageURL    string
	Description string
	Priority    int
}

func (ec EventCategory) Validate() map[string]string {
	errors := make(map[string]string)
	//if len(ec.Name) < 1 {
	//	errors["Name"] = "Field is required"
	//} else if err := fieldExists("event_categories", "name", ec.Name, ec.ID); err != nil {
	//	errors["Name"] = err.Error()
	//}
	//if len(ec.Pattern) < 1 {
	//	errors["Pattern"] = "Field is required"
	//} else if err := fieldExists("event_categories", "pattern", ec.Pattern, ec.ID); err != nil {
	//	errors["Pattern"] = err.Error()
	//}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
