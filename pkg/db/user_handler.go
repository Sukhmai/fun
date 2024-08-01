package db

import "context"

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func (c DBClient) SaveUser(ctx context.Context, user User) error {
	sql := "INSERT INTO users (user_id, username) VALUES ($1, $2)"
	_, err := c.conn.Exec(ctx, sql, user.UserID, user.Username)
	if err != nil {
		return err
	}
	return nil
}
