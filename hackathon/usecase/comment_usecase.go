package usecase

import (
	"hackathon/dao"
	"hackathon/model"
	"time"
	"github.com/oklog/ulid/v2"
	"log"
)

type CommentUsecase interface {
	CreateComment(comment *model.Comment) (*model.Comment, error)
	GetCommentsByPostId(postId string) ([]model.Comment, error)
}

type commentUsecase struct {
	commentDao dao.CommentDao
}

func NewCommentUsecase(commentDao dao.CommentDao) CommentUsecase {
	return &commentUsecase{commentDao: commentDao}
}

func (c *commentUsecase) CreateComment(comment *model.Comment) (*model.Comment, error) {
	comment.Id = ulid.Make().String()
	comment.CreatedAt = time.Now()
	if err := c.commentDao.CreateComment(comment); err != nil {
		log.Printf("ERROR: Failed to create comment: %v", err)
		return nil, err
	}
	return comment, nil
}

func (c *commentUsecase) GetCommentsByPostId(postId string) ([]model.Comment, error) {
	comments, err := c.commentDao.GetCommentsByPostId(postId)
	if err != nil {
		log.Printf("ERROR: Failed to get comments by post id: %v", err)
		return nil, err
	}
	return comments, nil
}