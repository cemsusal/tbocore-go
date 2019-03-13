package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Repository database instance
type Repository struct {
	db *gorm.DB
}

// Init initiates Postgre DB instance
func (p Repository) Init(conString string) {
	conn, err := gorm.Open("postgres", conString)
	if err != nil {
		fmt.Print(err)
	}
	p.db = conn
	defer p.db.Close()
}

// Find returns a single record from DB
func (p Repository) Find(theType interface{}, where interface{}) interface{} {
	return p.db.Find(theType, where)
}

// ByID returns a record for a specific id
func (p Repository) ByID(theType interface{}, id uint) interface{} {
	return p.db.First(theType, id)
}

// Create a new record
func (p Repository) Create(theType interface{}) {
	p.db.Create(theType)
}

// Update an existing record
func (p Repository) Update(theType interface{}) {
	p.db.Update(theType)
}

// UpdateField updates a single field
func (p Repository) UpdateField(fieldName string, theType interface{}) {
	p.db.Update(fieldName, theType)
}

// Delete a record
func (p Repository) Delete(theType interface{}) {
	p.db.Delete(theType)
}
