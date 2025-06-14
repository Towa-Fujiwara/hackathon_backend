package usecase

import (
	"fmt"
	"hackathon/dao"
	"hackathon/model"
)

// UserUsecase はユーザー関連のビジネスロジックを担当します。
// 依存するDAOをフィールドとして持ちます。
type userUsecase struct {
	userDao dao.UserDao
}


type UserUsecase interface {
	RegisterUser(userId, firebaseUID, name, bio, iconURL string) (*model.User, error)
	SearchUserExist(userId string) (*model.User, error)
	SearchUsers(query string) ([]model.User, error)
}

// NewUserUsecase はUserUsecaseの新しいインスタンスを生成します。
func NewUserUsecase(userDao dao.UserDao) *userUsecase {
	return &userUsecase{userDao: userDao}
}


func (uc *userUsecase) RegisterUser(userId, firebaseUID, name, bio, iconURL string) (*model.User, error) {
	
	user, err := model.NewUser(userId, firebaseUID, name, bio, iconURL)
	if err != nil {
		return nil, err
	}

	if err := uc.userDao.Create(user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}

func (uc *userUsecase) SearchUserExist(userId string) (*model.User, error) {
	if userId == "" {
		return nil, fmt.Errorf("userId is empty")
	}

	// 構造体が持つuserDaoのメソッドを呼び出します。
	user, err := uc.userDao.FindById(userId)
	if err != nil {
		return nil, err 
	}
	return user, nil
}

func (u *userUsecase) SearchUsers(query string) ([]model.User, error) {
	return u.userDao.SearchByName(query)
}