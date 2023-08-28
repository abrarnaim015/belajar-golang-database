package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/abrarnaim015/belajar-golang-database/entity"
)

type commentRepositoryImpl struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}

func (repository *commentRepositoryImpl) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	script := "INSERT INTO comments(email, comment) VALUES (?, ?)"
	result, err := repository.DB.ExecContext(ctx, script, comment.Email, comment.Comment)
	if err != nil {
		return comment, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}

	comment.Id = int32(id)
	return comment, nil
}

func (repository *commentRepositoryImpl) FindById(ctx context.Context, id int32) (entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments WHERE id=? LIMIT 1"
	rows, err := repository.DB.QueryContext(ctx, script, id)
	
	comment := entity.Comment{}
	if err != nil {
		return comment, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		if err != nil {
			return comment, err
		}
		return comment, nil
	} else {
		return comment, errors.New("Id " + strconv.Itoa(int(id)) + " Not Found")
	}
}

func (repository *commentRepositoryImpl) FindAll(ctx context.Context) ([]entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments LIMIT 5"
	rows, err := repository.DB.QueryContext(ctx, script)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment

	for rows.Next() {
		comment := entity.Comment {}
		rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		comments = append(comments, comment)
	}

	return comments, nil
}

func (repository *commentRepositoryImpl) UpdateById(ctx context.Context, id int32, comment entity.Comment) (entity.Comment, error) {
	script := "UPDATE comments SET email=?, comment=? WHERE id=?"
	result, err := repository.DB.ExecContext(ctx, script, comment.Email, comment.Comment, id)

	newComment := entity.Comment{}

	if err != nil {
		return newComment, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return newComment, err
	}

	if rowsAffected == 0 {
		return comment, errors.New("Id " + strconv.Itoa(int(id)) + " Not Found")
	}

	comment.Id = id
	return comment, nil
}

func (repository *commentRepositoryImpl) Delete(ctx context.Context, id int32) (bool, error) {
	script := "DELETE FROM comments WHERE id=?"
	result, err := repository.DB.ExecContext(ctx, script, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil
	}
	
	return true, nil
}