package sUsers

import "github.com/remrafvil/Auriga_API/internal/repositories/rUsers"

type Service interface {
	UserInfo(id uint) (msUser, error)
	UserList() []msUser
	UserAdd(UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (msUser, error)
	UserUpdate(id uint, UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (msUser, error)
	UserDel(id uint) error
	UserLogin(UserName string, Email string, Password string) (msUser, error)
}

type service struct {
	repository rUsers.Repository
}

func New(repository rUsers.Repository) Service {
	return &service{
		repository: repository,
	}
}

type msUser struct {
	UserName   string
	Password   string
	Email      string
	FirstName  string
	LastName   string
	Company    string
	Address    string
	City       string
	State      string
	PostalCode uint
	Role       uint
}

func (s *service) UserInfo(id uint) (msUser, error) {
	a := msUser{}
	return a, nil
}

func (s *service) UserList() []msUser {
	a := []msUser{}
	return a
}

func (s *service) UserAdd(UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (msUser, error) {
	a := msUser{}
	return a, nil
}

func (s *service) UserUpdate(id uint, UserName string, Email string, Password string, FirstName string, LastName string, Company string, Address string, City string, State string, PostalCode uint, Role uint) (msUser, error) {
	a := msUser{}
	return a, nil
}

func (s *service) UserDel(id uint) error {
	return nil
}
func (s *service) UserLogin(UserName string, Email string, Password string) (msUser, error) {
	a := msUser{}
	return a, nil
}
