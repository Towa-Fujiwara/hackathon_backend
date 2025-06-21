package dao

import (
	"database/sql"
	"hackathon/model"
)

type FollowUserDao interface {
	FollowUser(follow *model.Follow) error
	UnfollowUser(userId, followUserId string) error
	IsFollowing(userId, followUserId string) (bool, error)
	GetFollowers(userId string) ([]model.Follow, error)
	GetFollowing(userId string) ([]model.Follow, error)
}

type followDao struct {
	db *sql.DB
}

func NewFollowUserDao(db *sql.DB) FollowUserDao {
	return &followDao{db: db}
}

func (f *followDao) FollowUser(follow *model.Follow) error {
	query := "INSERT INTO follows (id, userId, followUserId, createdAt) VALUES (?, ?, ?, ?)"
	_, err := f.db.Exec(query, follow.Id, follow.UserId, follow.FollowUserId, follow.CreatedAt)
	return err
}

func (f *followDao) UnfollowUser(userId, followUserId string) error {
	query := "DELETE FROM follows WHERE userId = ? AND followUserId = ?"
	_, err := f.db.Exec(query, userId, followUserId)
	return err
}

func (f *followDao) IsFollowing(userId, followUserId string) (bool, error) {
	query := "SELECT COUNT(*) FROM follows WHERE userId = ? AND followUserId = ?"
	var count int
	err := f.db.QueryRow(query, userId, followUserId).Scan(&count)
	return count > 0, err
}

func (f *followDao) GetFollowers(userId string) ([]model.Follow, error) {
	query := "SELECT * FROM follows WHERE followUserId = ?"
	rows, err := f.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []model.Follow
	for rows.Next() {
		var follow model.Follow
		err := rows.Scan(&follow.Id, &follow.UserId, &follow.FollowUserId, &follow.CreatedAt)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follow)
	}
	return followers, nil
}

func (f *followDao) GetFollowing(userId string) ([]model.Follow, error) {
	query := "SELECT * FROM follows WHERE userId = ?"
	rows, err := f.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []model.Follow
	for rows.Next() {
		var follow model.Follow
		err := rows.Scan(&follow.Id, &follow.UserId, &follow.FollowUserId, &follow.CreatedAt)
		if err != nil {
			return nil, err
		}
		following = append(following, follow)
	}
	return following, nil
}