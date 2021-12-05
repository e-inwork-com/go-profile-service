package data

import (
	"database/sql"
	"errors"

	dataUser "github.com/e-inwork-com/golang-user-microservice/pkg/data"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Profiles	ProfileModel
	Addresses	AddressModel
	Users 		dataUser.UserModel
}

func InitModels(db *sql.DB) Models {
	return Models{
		Profiles:	ProfileModel{DB: db},
		Addresses:	AddressModel{DB: db},
		Users: 		dataUser.UserModel{DB: db},
	}
}

