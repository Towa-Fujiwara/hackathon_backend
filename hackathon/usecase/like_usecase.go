package usecase

import (
	"hackathon/dao"
	"hackathon/model"
	"log"
	"time"
	"github.com/oklog/ulid/v2" 
)

type PostLikeUsecase interface {
	ToggleLike(userId string, postId string) (*model.Like, error)
}

type postLikeUsecase struct {
	postLikeDao dao.LikeDao
}

func NewPostLikeUsecase(postLikeDao dao.LikeDao) PostLikeUsecase {
	return &postLikeUsecase{postLikeDao: postLikeDao}
}

func (uc *postLikeUsecase) ToggleLike(userId string, postId string) (*model.Like, error) {
	existingLike, err := uc.postLikeDao.FindById(userId, postId)
	if err != nil {
		log.Printf("ERROR: Failed to find like by user and post: %v", err)
		return nil, err
	}
	if existingLike != nil {
		log.Printf("INFO: Like already exists. Deleting like ID: %s", existingLike.Id)
		if err := uc.postLikeDao.Delete(existingLike.Id); err != nil {
			log.Printf("ERROR: Failed to delete like: %v", err)
			return nil, err
		}
		return nil, nil
	} else {
		like := &model.Like{
			Id:        ulid.Make().String(),
			UserId:    userId,
			PostId:    postId,
			CreatedAt: time.Now(),          
		}
		if err := uc.postLikeDao.Create(like); err != nil {
			return nil, err
		}
		return like, nil
	}
}