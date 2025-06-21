package usecase

import (
	"context"
	"fmt"
	"hackathon/dao"
	"hackathon/model"
	"log"
	"time"

	"github.com/google/uuid"
)

type FollowUserUsecase interface {
	FollowUser(ctx context.Context, currentFirebaseUID string, targetUserId string) error
	UnfollowUser(ctx context.Context, currentFirebaseUID string, targetUserId string) error
	GetFollowers(ctx context.Context, targetUserId string) ([]model.User, error)
	GetFollowing(ctx context.Context, targetUserId string) ([]model.User, error)
	IsFollowing(ctx context.Context, currentFirebaseUID string, targetUserId string) (bool, error)
	GetUserFollowCounts(ctx context.Context, targetUserId string) (int, int, error)
}

type followUserUsecase struct {
	followUserDao dao.FollowUserDao
	userDao       dao.UserDao
}

func NewFollowUserUsecase(followUserDao dao.FollowUserDao, userDao dao.UserDao) FollowUserUsecase {
	return &followUserUsecase{
		followUserDao: followUserDao,
		userDao:       userDao,
	}
}

func (u *followUserUsecase) getCurrentUserIdFromFirebaseUID(ctx context.Context, firebaseUID string) (string, error) {
	userProfile, err := u.userDao.FindByFirebaseUID(firebaseUID)
	if err != nil {
		log.Printf("ERROR: Failed to get current user profile by FirebaseUID (%s): %v", firebaseUID, err)
		return "", fmt.Errorf("現在のユーザーのプロフィール取得に失敗しました: %w", err)
	}
	if userProfile == nil {
		return "", fmt.Errorf("user not found for firebaseUID: %s", firebaseUID)
	}
	return userProfile.UserId, nil
}

func (u *followUserUsecase) FollowUser(ctx context.Context, currentFirebaseUID string, targetUserId string) error {
	currentUserId, err := u.getCurrentUserIdFromFirebaseUID(ctx, currentFirebaseUID)
	if err != nil {
		return err
	}

	if currentUserId == targetUserId {
		return fmt.Errorf("自分自身をフォローすることはできません")
	}

	isFollowing, err := u.followUserDao.IsFollowing(currentUserId, targetUserId)
	if err != nil {
		return fmt.Errorf("フォロー状態の確認に失敗しました: %w", err)
	}
	if isFollowing {
		return fmt.Errorf("既にフォローしています")
	}

	follow := &model.Follow{
		Id:           uuid.New().String(),
		UserId:       currentUserId,
		FollowUserId: targetUserId,
		CreatedAt:    time.Now(),
	}

	if err := u.followUserDao.FollowUser(follow); err != nil {
		log.Printf("ERROR: Failed to follow user: %v", err)
		return err
	}
	return nil
}

func (u *followUserUsecase) UnfollowUser(ctx context.Context, currentFirebaseUID string, targetUserId string) error {
	currentUserId, err := u.getCurrentUserIdFromFirebaseUID(ctx, currentFirebaseUID)
	if err != nil {
		return err
	}

	isFollowing, err := u.followUserDao.IsFollowing(currentUserId, targetUserId)
	if err != nil {
		return fmt.Errorf("フォロー状態の確認に失敗しました: %w", err)
	}
	if !isFollowing {
		return fmt.Errorf("フォローしていません")
	}

	if err := u.followUserDao.UnfollowUser(currentUserId, targetUserId); err != nil {
		log.Printf("ERROR: Failed to unfollow user: %v", err)
		return err
	}
	return nil
}

func (u *followUserUsecase) GetFollowers(ctx context.Context, targetUserId string) ([]model.User, error) {
	followers, err := u.followUserDao.GetFollowers(targetUserId)
	if err != nil {
		log.Printf("ERROR: Failed to get followers from DAO: %v", err)
		return nil, err
	}

	var followerUsers []model.User
	for _, follow := range followers {
		userProfile, err := u.userDao.FindById(follow.UserId)
		if err != nil {
			log.Printf("Warning: Failed to get follower user profile for UserId %s: %v", follow.UserId, err)
			continue
		}
		if userProfile != nil {
			followerUsers = append(followerUsers, *userProfile)
		}
	}
	return followerUsers, nil
}

func (u *followUserUsecase) GetFollowing(ctx context.Context, targetUserId string) ([]model.User, error) {
	following, err := u.followUserDao.GetFollowing(targetUserId)
	if err != nil {
		log.Printf("ERROR: Failed to get following from DAO: %v", err)
		return nil, err
	}

	var followingUsers []model.User
	for _, follow := range following {
		userProfile, err := u.userDao.FindById(follow.FollowUserId)
		if err != nil {
			log.Printf("Warning: Failed to get following user profile for FollowUserId %s: %v", follow.FollowUserId, err)
			continue
		}
		if userProfile != nil {
			followingUsers = append(followingUsers, *userProfile)
		}
	}
	return followingUsers, nil
}

func (u *followUserUsecase) IsFollowing(ctx context.Context, currentFirebaseUID string, targetUserId string) (bool, error) {
	currentUserId, err := u.getCurrentUserIdFromFirebaseUID(ctx, currentFirebaseUID)
	if err != nil {
		// ユーザーが見つからない場合はフォローしていない
		return false, nil
	}
	isFollowing, err := u.followUserDao.IsFollowing(currentUserId, targetUserId)
	if err != nil {
		log.Printf("ERROR: Failed to check if following: %v", err)
		return false, err
	}
	return isFollowing, nil
}

func (u *followUserUsecase) GetUserFollowCounts(ctx context.Context, targetUserId string) (int, int, error) {
	followers, err := u.followUserDao.GetFollowers(targetUserId)
	if err != nil {
		return 0, 0, fmt.Errorf("フォロワー数の取得に失敗しました: %w", err)
	}
	following, err := u.followUserDao.GetFollowing(targetUserId)
	if err != nil {
		return 0, 0, fmt.Errorf("フォロー数の取得に失敗しました: %w", err)
	}
	return len(following), len(followers), nil
}
