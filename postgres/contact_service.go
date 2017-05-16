package postgres

import (
	"github.com/jorgeastorga/contakx/contakx"
	_"database/sql"
	"github.com/jinzhu/gorm"
	"log"
)

type ContactService struct {
	AppDB *gorm.DB
}

///ContactService returns a user for a given id.
func (c *ContactService) Contact(id int) (contact contakx.Contact, err error){
	//var *contact contakx.Contact
	c.AppDB.First(&contact, id)
	return contact, err
}

func (c *ContactService) CreateContact(contact *contakx.Contact) error {
	err := c.AppDB.Create(&contakx.Contact{
	FirstName:contact.FirstName,
	LastName: contact.LastName,
	Address1: contact.Address1,
	Address2: contact.Address2,
	City: contact.City,
	State: contact.State,
	ZipCode: contact.ZipCode,
	Company: contact.Company,
	Email: contact.Email,
	Phone: contact.Phone}).Error

	if err != nil {
		log.Println("Error creating the contact: ", err.Error())
	}

	return err
}

func (c *ContactService)OpenDB(dbConnection string){
	//Database Connection Setup: PostgreSQL
	var err error
	c.AppDB, err = gorm.Open("postgres", dbConnection)

	if err != nil{
		log.Println("Failed to connect to the database", err.Error())
	} else {
		log.Println("DB Connection: connected to the database successfully")
	}

}


