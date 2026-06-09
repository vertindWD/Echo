package mysql

import "Echo/models"

func UpsertVote(userID, postID int64, direction int8) error {
	return DB.Exec(
		"INSERT INTO vote (user_id, post_id, direction) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE direction = VALUES(direction)",
		userID, postID, direction,
	).Error
}

func DeleteVote(userID, postID int64) error {
	return DB.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Vote{}).Error
}
