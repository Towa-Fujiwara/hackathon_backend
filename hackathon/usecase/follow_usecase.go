package usecase

import (
	"hackathon/model"
	"hackathon/dao"
	"log"
)

type FollowUserUsecase interface {
	FollowUser(follow *model.Follow) (*model.Follow, error)
	GetFollowers(userId string) ([]model.Follow, error)
	GetFollowing(userId string) ([]model.Follow, error)
	IsFollowing(userId, followUserId string) (bool, error)
}

type followUserUsecase struct {
	followUserDao dao.FollowUserDao
}

func NewFollowUserUsecase(followUserDao dao.FollowUserDao) FollowUserUsecase {
	return &followUserUsecase{followUserDao: followUserDao}
}

func (u *followUserUsecase) FollowUser(follow *model.Follow) (*model.Follow, error) {
	if err := u.followUserDao.FollowUser(follow); err != nil {
		log.Printf("ERROR: Failed to follow user: %v", err)
		return nil, err
	}
	return follow, nil
}
func (u *followUserUsecase) GetFollowers(userId string) ([]model.Follow, error) {
	followers, err := u.followUserDao.GetFollowers(userId)
	if err != nil {
		log.Printf("ERROR: Failed to get followers: %v", err)
		return nil, err	
	}
	return followers, nil
}
func (u *followUserUsecase) GetFollowing(userId string) ([]model.Follow, error) {
	following, err := u.followUserDao.GetFollowing(userId)
	if err != nil {
		log.Printf("ERROR: Failed to get following: %v", err)
		return nil, err
	}
	return following, nil
}
func (u *followUserUsecase) IsFollowing(userId, followUserId string) (bool, error) {
	isFollowing, err := u.followUserDao.IsFollowing(userId, followUserId)
	if err != nil {
		log.Printf("ERROR: Failed to check if following: %v", err)
		return false, err
	}
	return isFollowing, nil
}
