package mysql

import (
	"Echo/models"
	"fmt"
	"strings"
)

func CreatePost(p *models.Post) (err error) {
	err = DB.Create(p).Error
	return
}

func GetPostById(pid int64) (*models.Post, error) {
	p := new(models.Post)
	err := DB.Where("post_id = ?", pid).First(p).Error
	return p, err
}

// GetPostListByIDs 根据给定的 ID 列表查询帖子，并严格保持给定 ID 的顺序
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	if len(ids) == 0 {
		return nil, nil
	}
	orderStr := fmt.Sprintf("FIELD(post_id, %s)", strings.Join(ids, ","))

	// 执行查询
	err = DB.Where("post_id IN (?)", ids).
		Order(orderStr).
		Find(&postList).Error

	return postList, err
}
