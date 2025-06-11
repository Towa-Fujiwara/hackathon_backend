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
	RegisterUser(id, name, password, displayName, bio, iconURL string, age int) (*model.User, error)
	SearchUserExist(id string) (*model.User, error)
	SearchUsers(query string) ([]model.User, error)
}

// NewUserUsecase はUserUsecaseの新しいインスタンスを生成します。
func NewUserUsecase(userDao dao.UserDao) *userUsecase {
	return &userUsecase{userDao: userDao}
}

// RegisterUser はユーザーを登録するビジネスロジックです。
// 引数を、model.UserProfileの定義に合わせて修正しました。
func (uc *userUsecase) RegisterUser(id, name, password, displayName, bio, iconURL string, age int) (*model.User, error) {
	
	user, err := model.NewUser(id, name, password, displayName, bio, iconURL, age)
	if err != nil {
		return nil, err
	}
	// 構造体が持つuserDaoのメソッドを呼び出します。
	if err := uc.userDao.Create(user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}

// SearchUser はユーザーをIDで検索するビジネスロジックです。
// 不要なコードを削除し、ロジックを修正しました。
func (uc *userUsecase) SearchUserExist(id string) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}

	// 構造体が持つuserDaoのメソッドを呼び出します。
	user, err := uc.userDao.FindById(id)
	if err != nil {
		return nil, err // エラーをそのまま返す
	}
	return user, nil
}

func (u *userUsecase) SearchUsers(query string) ([]model.User, error) {
	return u.userDao.SearchByName(query)
}