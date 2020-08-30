package services

import (
	"golang.org/x/crypto/bcrypt"
	"imooc-shop/datamodels"
	"imooc-shop/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository	repositories.IUserRepository
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)  {
	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}

	isOk, _ = ValidatePassword(pwd, user.HasPassword)
	if !isOk {
		return &datamodels.User{}, false
	}

	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error)  {
	pwdByte, errPwd := GeneratePassword(user.HasPassword)
	if errPwd != nil {
		return userId, errPwd
	}

	user.HasPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error)  {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword, hashed string) (isOK bool, err error)  {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, err
	}

	return true, nil
}

