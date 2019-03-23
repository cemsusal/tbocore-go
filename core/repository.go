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
func NewRepository(config *Config) *Repository {
	repository := Repository{}
	conStr := "host=" + config.Database.Host + " port=" + config.Database.Port + " user=" + config.Database.User + " dbname=" + config.Database.DbName + " password=" + config.Database.Password + " sslmode=" + config.Database.SslMode
	// conStr := "host=62.75.148.168 port=5432 user=tmmaster dbname=tmuat password=Dky1qw!! sslmode=disable"
	conn, err := gorm.Open("postgres", conStr)
	fmt.Println(conStr)
	if err != nil {
		fmt.Print(err)
	}
	repository.Db = conn
	return &repository
}

// Find returns a single record from DB
func (p Repository) Find(theType interface{}, where interface{}) interface{} {
	return p.Db.Model(theType).Find(theType, where)
}

// ByID returns a record for a specific id
func (p Repository) ByID(theType interface{}, id uint) interface{} {
	return p.Db.Model(theType).First(theType, id)
}

// Create a new record
func (p Repository) Create(theType interface{}) {
	p.Db.Model(theType).Create(theType)
}

// Update an existing record
func (p Repository) Update(theType interface{}) {
	p.Db.Model(theType).Update(theType)
}

// UpdateField updates a single field
func (p Repository) UpdateField(fieldName string, theType interface{}) {
	p.Db.Model(theType).Update(fieldName, theType)
}

// Delete a record
func (p Repository) Delete(theType interface{}) {
	p.Db.Model(theType).Delete(theType)
}
