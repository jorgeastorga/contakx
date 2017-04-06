package main

import (
	"github.com/jinzhu/gorm"
	"log"
)

type Contact struct{
	gorm.Model //
	//Id int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City string `json:"city"`
	State string `json:"state"`
	ZipCode int `json:"zip"`
	Company string `json:"company"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

//TODO Move these to a contact_controller file per Anu's recommendation. This logic should be separate.
/**
 * Function used to create  a contact in the db
 */
func (c *Contact)create() (err error){

	err = AppDB.Create(&Contact{
		FirstName:c.FirstName,
		LastName: c.LastName,
		Address1: c.Address1,
		Address2: c.Address2,
		City: c.City,
		State: c.State,
		ZipCode: c.ZipCode,
		Company: c.Company,
		Email: c.Email,
		Phone: c.Phone}).Error

	if err != nil{
		log.Println("Error Creating the Contact: ", err.Error())
	}

	return err
}

/**
* Function used to retrieve a Contact from the db
*
*/
func retrieve(id int) (contact Contact, err error) {

	AppDB.First(&contact, id)
	return contact, err
}


/**
* Function used to update a Contact from the db
*
*
*/
//TODO: update the update method to use Gorm
func (c *Contact) update() (err error){
	/*_, err = Db. Exec("update contacts set firstname = $2, lastname = $3, address1 = $3 where id = $1",
		c.Id, c.FirstName, c.LastName, c.Address1)*/
	return err
}

/**
* Function used to delete a Contact from the db
*
 */
//TODO: update the delete method to use Gorm
func (c *Contact)delete() (err error){
	//_, err = Db. Exec("delete from contacts where id = $1", c.Id)
	return
}
