package repository

import (
	"context"
	"fmt"
	"testing"

	belajargolangdatabase "github.com/abrarnaim015/belajar-golang-database"
	"github.com/abrarnaim015/belajar-golang-database/entity"
	_ "github.com/go-sql-driver/mysql"
)

var commentRepository = NewCommentRepository(belajargolangdatabase.GetConnection())

func TestCommentInsert(t *testing.T)  {
	ctx := context.Background()
	comment := entity.Comment {
		Email: "repository@test.com",
		Comment: "Test Repository",
	}

	result, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func TestCommentFindById(t *testing.T)  {
	ctx := context.Background()

	comment, err := commentRepository.FindById(ctx, 55)
	if err != nil {
		panic(err)
	}

	fmt.Println(comment)
}

func TestCommentFindAll(t *testing.T)  {
	ctx := context.Background()

	comments, err := commentRepository.FindAll(ctx)
	if err != nil {
		panic(err)
	}

	for _, comment := range comments {
		fmt.Println(comment)
	}
}

func TestCommentDelete(t *testing.T)  {
	ctx := context.Background()
	id := 5

	comment, err := commentRepository.Delete(ctx, int32(id))
	if err != nil {
		panic(err)
	}

	fmt.Println(comment)
}

func TestCOmmentUpdateById(t *testing.T)  {
	ctx := context.Background()
	id := 9
	comment := entity.Comment {
		Email: "update@test.com",
		Comment: "Update Test",
	}

	result, err := commentRepository.UpdateById(ctx, int32(id), comment)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}