package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Repository database instance
type Repository struct {
	Db *gorm.DB
}

// Init initiates Postgre DB instance
func (p *Repository) Init(conString string) {
	conn, err := gorm.Open("postgres", conString)
	if err != nil {
		fmt.Print(err)
	}
	p.Db = conn
	// defer p.Db.Close()
}

// NewRepository initiates a new repository instance
func NewRepository(config Config) *Repository {
	repository := Repository{}
	conStr := "host=%d port=%d user=%d dbname=%d password=%d sslmode=%d"
	conStr = fmt.Sprintf(conStr, config.Database.Host,
		config.Database.Port, config.Database.User, config.Database.DbName, config.Database.Password, config.Database.SslMode)
	conn, err := gorm.Open("postgres", conStr)
	if err != nil {
		fmt.Print(err)
	}
	repository.Db = conn
	return &repository
}

// Find returns a single record from DB
func (p Repository) Find(theType interface{}, where interface{}) interface{} {
	return p.Db.Find(theType, where)
}

// ByID returns a record for a specific id
func (p Repository) ByID(theType interface{}, id uint) interface{} {
	return p.Db.First(theType, id)
}

// Create a new record
func (p Repository) Create(theType interface{}) {
	p.Db.Create(theType)
}

// Update an existing record
func (p Repository) Update(theType interface{}) {
	p.Db.Update(theType)
}

// UpdateField updates a single field
func (p Repository) UpdateField(fieldName string, theType interface{}) {
	p.Db.Update(fieldName, theType)
}

// Delete a record
func (p Repository) Delete(theType interface{}) {
	p.Db.Delete(theType)
}
