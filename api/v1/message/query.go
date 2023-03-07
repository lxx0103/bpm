package message

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type messageQuery struct {
	conn *sqlx.DB
}

func NewMessageQuery(connection *sqlx.DB) *messageQuery {
	return &messageQuery{
		conn: connection,
	}
}

func (r *messageQuery) GetMessageByID(id int64) (*Message, error) {
	var message Message
	err := r.conn.Get(&message, "SELECT * FROM messages WHERE id = ? AND status > 0", id)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
func (r *messageQuery) GetMessageCount(filter MessageFilter) (int, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	var count int
	err := r.conn.Get(&count, `
		SELECT count(1) as count
		FROM messages
		WHERE `+strings.Join(where, " AND "), args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *messageQuery) GetMessageList(filter MessageFilter) (*[]Message, error) {
	where, args := []string{"status > 0"}, []interface{}{}
	if v := filter.Name; v != "" {
		where, args = append(where, "name like ?"), append(args, "%"+v+"%")
	}
	args = append(args, filter.PageId*filter.PageSize-filter.PageSize)
	args = append(args, filter.PageSize)
	var messages []Message
	err := r.conn.Select(&messages, `
		SELECT *
		FROM messages
		WHERE `+strings.Join(where, " AND ")+`
		LIMIT ?, ?
	`, args...)
	if err != nil {
		return nil, err
	}
	return &messages, nil
}

// func (r *messageQuery) GetUserByPosition(positionID int64) (*[]string, error) {
// 	var openIDs []string
// 	err := r.conn.Select(&openIDs, `
// 		SELECT identifier
// 		FROM users
// 		WHERE position_id = ?
// 		AND status = 1
// 	`, positionID)
// 	return &openIDs, err
// }

func (r *messageQuery) GetUserByIDAndProject(userID, projectID int64) (string, error) {
	var openID string
	err := r.conn.Get(&openID, `
		SELECT identifier
		FROM users
		WHERE id = ?
		AND status = 1
		AND id IN (
			SELECT  user_id from project_members where project_id  = ? and status > 0
		)
		`, userID, projectID)
	return openID, err
}

func (r *messageQuery) GetUserByPositionAndProject(positionID, projectID int64) (*[]string, error) {
	var openIDs []string
	err := r.conn.Select(&openIDs, `
		SELECT identifier
		FROM users
		WHERE position_id = ?
		AND status = 1
		AND id IN (
			SELECT  user_id from project_members where project_id  = ? and status > 0
		)
	`, positionID, projectID)
	return &openIDs, err
}

func (r *messageQuery) GetOtherMemberByProject(projectID, userID int64) ([]string, error) {
	var openID []string
	err := r.conn.Select(&openID, `
		SELECT identifier
		FROM users
		WHERE status = 1
		AND id IN (
			SELECT  user_id from project_members where project_id  = ? AND user_id != ? and status > 0
		)
		`, projectID, userID)
	return openID, err
}
