package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Answer struct {
	Username  string
	Question1 *string
	Question2 *string
	Question3 *string
	Question4 *string
	Question5 *string
	Question6 *string
}

type Answer2 struct {
	Username string
	Answers  []*string
}

func (c *DBClient) SaveAnswer(ctx context.Context, userId, answer string, questionNum int) error {
	// Determine the column name based on the question number
	columnName := fmt.Sprintf("question%d", questionNum)

	// Check if an answer already exists for the user and question
	var existingAnswer *string
	sqlCheck := fmt.Sprintf("SELECT %s FROM answers WHERE user_id = $1", columnName)
	err := c.conn.QueryRow(context.Background(), sqlCheck, userId).Scan(&existingAnswer)

	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("failed to check if answer exists: %v", err)
	}

	if err == pgx.ErrNoRows {
		// Insert a new answer if none exists
		sqlInsert := fmt.Sprintf("INSERT INTO answers (user_id, %s) VALUES ($1, $2)", columnName)
		_, err = c.conn.Exec(context.Background(), sqlInsert, userId, answer)
		if err != nil {
			return fmt.Errorf("failed to insert answer: %v", err)
		}
	} else {
		// Update the existing answer
		sqlUpdate := fmt.Sprintf("UPDATE answers SET %s = $1 WHERE user_id = $2", columnName)
		_, err = c.conn.Exec(context.Background(), sqlUpdate, answer, userId)
		if err != nil {
			return fmt.Errorf("failed to update answer: %v", err)
		}
	}
	return nil
}

func (c *DBClient) GetAllAnswers(ctx context.Context) ([]Answer2, error) {
	sql := `
		SELECT u.username, a.question0, a.question1, a.question2, a.question3, a.question4, a.question5
		FROM answers a
		JOIN users u ON a.user_id = u.user_id
	`
	rows, err := c.conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and collect the data
	var answers2 []Answer2
	for rows.Next() {
		var answer Answer
		err := rows.Scan(&answer.Username, &answer.Question1, &answer.Question2, &answer.Question3, &answer.Question4, &answer.Question5, &answer.Question6)
		if err != nil {
			return nil, err
		}

		// This is curry code to convert Answer to Answer2
		answer2 := Answer2{
			Username: answer.Username,
			Answers:  []*string{answer.Question1, answer.Question2, answer.Question3, answer.Question4, answer.Question5, answer.Question6},
		}
		answers2 = append(answers2, answer2)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return answers2, nil
}

func (c *DBClient) GetAnswer(ctx context.Context, userId string, questionNum int) (string, error) {
	columnName := fmt.Sprintf("question%d", questionNum)
	sql := fmt.Sprintf("SELECT answer FROM answers WHERE user_id = $1 AND %s = $2", columnName)
	var answer string
	err := c.conn.QueryRow(ctx, sql, userId, questionNum).Scan(&answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}
