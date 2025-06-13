package usecase

import (
	"hackathon/model"
	"time"
	"github.com/oklog/ulid/v2" 
	"hackathon/dao"
)

type PostUsecase interface {
	FindAllPosts() ([]model.Post, error)
	CreatePost(post *model.Post) (*model.Post, error)
	FindById(id string) (*model.Post, error)
	FindAllPostsByUserId(uid string) ([]model.Post, error)
}

type postUsecase struct {
	postDAO dao.PostDao // PostDaoに依存
}

func NewPostUsecase(pd dao.PostDao) PostUsecase {
	return &postUsecase{postDAO: pd}
}

// 全ての投稿を取得するビジネスロジック
func (u *postUsecase) FindAllPosts() ([]model.Post, error) {
	// 対応するDAOのメソッドを呼び出します
	return u.postDAO.FindAll()
}

func (u *postUsecase) FindAllPostsByUserId(uid string) ([]model.Post, error) {
	return u.postDAO.FindAllByUserId(uid)
}

func (u *postUsecase) CreatePost(post *model.Post) (*model.Post, error) {
	post.Id = ulid.Make().String()
	post.CreatedAt = time.Now()
	if err := u.postDAO.Create(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (u *postUsecase) FindById(id string) (*model.Post, error) {
	return u.postDAO.FindById(id)
}