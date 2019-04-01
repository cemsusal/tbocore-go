package core

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Repository database instance
type Repository struct {
	Db *gorm.DB
}

// Init initiates Postgre DB instance
func (r *Repository) Init(conString string) {
	conn, err := gorm.Open("postgres", conString)
	if err != nil {
		fmt.Print(err)
	}
	r.Db = conn
	// defer p.Db.Close()
}

// NewRepository initiates a new repository instance
func NewRepository(config *Config) *Repository {
	repository := Repository{}
	conStr := "host=" + config.Database.Host + " port=" + config.Database.Port + " user=" + config.Database.User + " dbname=" + config.Database.DbName + " password=" + config.Database.Password + " sslmode=" + config.Database.SslMode
	conn, err := gorm.Open("postgres", conStr)
	if err != nil {
		fmt.Print(err)
	}
	repository.Db = conn
	return &repository
}

// Find returns a single record from DB
func (r Repository) Find(entity interface{}, where interface{}) interface{} {
	t := reflect.TypeOf(entity)
	return r.Db.Model(t).Find(entity, where)
}

// ByID returns a record for a specific id
func (r Repository) ByID(entity interface{}, id uint) (interface{}, error) {
	t := reflect.TypeOf(entity)
	retval := r.Db.Model(t).First(entity, id)
	return retval.Value, retval.Error
}

// Create a new record
func (r Repository) Create(entity interface{}) (interface{}, error) {
	t := reflect.TypeOf(entity)
	d := r.Db.Model(t).Create(entity)
	return d.Value, d.Error
}

// Update an existing record
func (r Repository) Update(entity interface{}) (interface{}, error) {
	t := reflect.TypeOf(entity)
	d := r.Db.Model(t).Update(entity)
	return d.Value, d.Error
}

// UpdateField updates a single field
func (r Repository) UpdateField(fieldName string, entity interface{}) (interface{}, error) {
	t := reflect.TypeOf(entity)
	d := r.Db.Model(t).Update(fieldName, entity)
	return d.Value, d.Error
}

// Delete a record
func (r Repository) Delete(entity interface{}) (interface{}, error) {
	t := reflect.TypeOf(entity)
	d := r.Db.Model(t).Delete(entity)
	return d.Value, d.Error
}
