package rUsers

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	UserList() []MrUser
	UserInfo(id uint) (MrUser, error)
	UserAdd(UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (MrUser, error)
	UserUpdate(id uint, UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (MrUser, error)
	UserDel(id uint) error
	UserLogin(UserName string, Email string, Password string) (MrUser, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

type MrUser struct {
	gorm.Model
	UserName   string `gorm:"type:varchar(20);	not null;		index:username" 	validate:"required"`
	Password   string `gorm:"type:char(32);		not null;" 							validate:"required"`
	Email      string `gorm:"type:varchar(100);	DEFAULT '';"`
	FirstName  string `gorm:"type:varchar(100);	DEFAULT '';"`
	LastName   string `gorm:"type:varchar(100);	DEFAULT '';"`
	Company    string `gorm:"type:varchar(100);	DEFAULT '';"`
	Address    string `gorm:"type:varchar(100);	DEFAULT '';"`
	City       string `gorm:"type:varchar(100);	DEFAULT '';"`
	State      string `gorm:"type:varchar(100);	DEFAULT '';"`
	PostalCode uint
	Role       uint
}

func (m *repository) UserList() []MrUser {
	var data = []MrUser{}
	err := m.db.Find(&data).Error
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

func (m *repository) UserInfo(id uint) (MrUser, error) {
	var data MrUser

	if err := m.db.Where("id  = ? ", id).Find(&data).Error; err != nil {
		return MrUser{}, errors.New("No existe registro con es id")
	}

	return data, nil
}

func (m *repository) UserAdd(UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (MrUser, error) {
	var data MrUser
	data.UserName = UserName
	data.Email = Email
	data.Password = GetMD5Hash(Password)
	data.FirstName = FirstName
	data.LastName = LastName
	data.Company = Company
	data.Address = Address
	data.City = City
	data.State = State
	data.PostalCode = PostalCode
	data.Role = Role

	if err := m.db.Create(&data).Error; err != nil {
		return MrUser{}, err
	} else {

		return data, nil
	}

}

func (m *repository) UserUpdate(id uint, UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (MrUser, error) {
	var data MrUser

	if err := m.db.Where("id = ? ", id).First(&data).Error; err != nil {
		return MrUser{}, errors.New("No hay registro con ese id")
	}
	data.UserName = UserName
	data.Email = Email
	data.Password = Password
	data.FirstName = FirstName
	data.LastName = LastName
	data.Company = Company
	data.Address = Address
	data.City = City
	data.State = State
	data.PostalCode = PostalCode
	data.Role = Role

	if err := m.db.Save(&data).Error; err != nil {
		return MrUser{}, errors.New("no se pudo actualizar")
	}
	return data, nil

}

func (m *repository) UserDel(id uint) error { // Borra uno
	var data MrUser
	if err := m.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}

/**************************************************/

func (m *repository) UserLogin(UserName string, Email string, Password string) (MrUser, error) {
	var data MrUser
	data.UserName = UserName
	data.Email = Email
	data.Password = Password

	if UserName == "" || Password == "" {
		return MrUser{}, errors.New("Cuenta o contraseña está vacía")
	}
	md5_password := GetMD5Hash(data.Password)
	if err := m.db.Where("user_name = ? AND password = ?", data.UserName, md5_password).First(&data).Error; err != nil {
		return MrUser{}, errors.New("Cuenta o contraseña incorrecta")
	}
	return data, nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

//////         INITIALITATION             ////////////

func DB_InitUsers(c *gorm.DB) { //
	c.Create(&MrUser{UserName: "Jesus", Email: "jr.cxpn@gmail.com", Password: GetMD5Hash("1234"), FirstName: "Jesus", LastName: "del Rio", Company: "Coexpan", Address: "Avda. Madrid, 72", City: "Alcalá de Henares", State: "Madrid - ESPAÑA", PostalCode: 28802, Role: 10})
	c.Create(&MrUser{UserName: "Jesus", Email: "jdelrio@coexpan.com", Password: GetMD5Hash("1234"), FirstName: "Jesus", LastName: "del Rio", Company: "Coexpan", Address: "Avda. Madrid, 72", City: "Alcalá de Henares", State: "Madrid - ESPAÑA", PostalCode: 28802, Role: 8})
}
