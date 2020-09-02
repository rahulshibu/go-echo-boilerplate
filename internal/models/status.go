package models

// Status - status
type Status struct {
	ID                uint   `gorm:"primary_key" json:"id"`
	StatusName        string `json:"status_name"`
	StatusDescription string `json:"status_description"`
}

//TableName - Set status table name
func (Status) TableName() string {
	return "status"
}
