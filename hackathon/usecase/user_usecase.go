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
	GetUserByFirebaseUID(firebaseUID string) (*model.User, error)
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

func (uc *userUsecase) GetUserByFirebaseUID(firebaseUID string) (*model.User, error) {
	if firebaseUID == "" {
		return nil, fmt.Errorf("firebaseUID is empty")
	}
	user, err := uc.userDao.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) SearchUsers(query string) ([]model.User, error) {
	return u.userDao.SearchByName(query)
}