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
	GetUserByUserId(userId string) (*model.User, error)
	UpdateUser(firebaseUID, name, bio, iconURL string) (*model.User, error)
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

func (uc *userUsecase) GetUserByUserId(userId string) (*model.User, error) {
	if userId == "" {
		return nil, fmt.Errorf("userId is empty")
	}
	user, err := uc.userDao.FindById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) SearchUsers(query string) ([]model.User, error) {
	return u.userDao.SearchByName(query)
}

func (u *userUsecase) UpdateUser(firebaseUID, name, bio, iconURL string) (*model.User, error) {
	if firebaseUID == "" {
		return nil, fmt.Errorf("firebaseUID is empty")
	}

	// 既存のユーザーを取得
	existingUser, err := u.userDao.FindByFirebaseUID(firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	// 更新するフィールドのみを設定
	if name != "" {
		existingUser.Name = name
	}
	if bio != "" {
		existingUser.Bio = bio
	}
	if iconURL != "" {
		existingUser.IconUrl = iconURL
	}

	// データベースを更新
	if err := u.userDao.Update(existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}