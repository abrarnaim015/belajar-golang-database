package belajargolangdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestExecSql(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES('naim3', 'Naim3')"
	_, err := db.ExecContext(ctx, script)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert new Customer")
}

func TestExecSql2(t *testing.T)  {
	wg := sync.WaitGroup {}
	db := GetConnection()
	defer db.Close()

	script := []string {
		"INSERT INTO customer (id, name, email, balance, rating, birth_date, married) VALUES ('naim', 'Naim', 'naim@gmail.com', 100000, 5.0, '1999-9-9', false)",
		"INSERT INTO customer (id, name, email, balance, rating, birth_date, married) VALUES ('abrar', 'Abrar', 'Abrar@gmail.com', 100000, 5.0, '1999-9-9', false)",
	}

	ctx := context.Background()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(script); i++ {
			_, err := db.ExecContext(ctx, script[i])
			if err != nil {
				panic(err)
			}
		}
	}()

	wg.Wait()
	fmt.Println("Success Insert new Customer")
}

func TestQuerySql(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, name string
		)
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("Id :", id)
		fmt.Println("Name :", name)
	}
}

func TestQuerySqlComplex(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, name string
			email sql.NullString
			balance int32
			rating float64
			birthDate sql.NullTime
			createdAt time.Time
			married bool
		)
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("==================")
		fmt.Println("Id :", id)
		fmt.Println("Name :", name)
		// fmt.Println("email :", email.String)
		if email.Valid {
			fmt.Println("email :", email.String)
		}
		fmt.Println("Balance :", balance)
		fmt.Println("Rating :", rating)
		// fmt.Println("BirthDate :", birthDate.Time)
		if birthDate.Valid {
			fmt.Println("BirthDate :", birthDate.Time)
		}
		fmt.Println("Married :", married)
		fmt.Println("CreatedAt :", createdAt)
	}
	fmt.Println(runtime.NumGoroutine())
}

func TestAddUser(t *testing.T)  {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	script := "INSERT INTO USER(username, PASSWORD) VALUES ('admin', 'admin')"

	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert new User")
}

func TestSqlInjection(t *testing.T)  {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	username := "admin'; #"
	password := "salah" // padahal passwordnya salah tp sukses login karna strnya di manipulasi
	// username := "admin"
	// password := "admin"

	script := "SELECT username FROM user WHERE username = '"+username+
	"' AND password = '"+password+"' LIMIT 1"

	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var (
			username string
		)
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Success login", username)
	} else {
		fmt.Println("Gagal Login")
	}
}

func TestSqlInjectionSafe(t *testing.T)  {
	db := GetConnection()
	defer db.Close()
	ctx := context.Background()

	username := "admin'; #"
	password := "salah"

	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"

	rows, err := db.QueryContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var (
			username string
		)
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Success login", username)
	} else {
		fmt.Println("Gagal Login")
	}
}

func TestExecSqlParameter(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "pio"
	password := "pio"

	script := "INSERT INTO user(username, password) VALUES(?, ?)"
	_, err := db.ExecContext(ctx, script, username, password)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert new User")
}

func TestAutoIncrement(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "naim@gmail.com"
	comment := "Test comment"

	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, script, email, comment)

	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert new Comment with id", insertId)
}

func TestPrepareStatement(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "naim" + strconv.Itoa(i) + "@gamil.com"
		comment := "Ini comment ke " + strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment id ", id)
	}
}

func TestTransaction(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	// do Transaction
	for i := 0; i < 10; i++ {
		email := "naim" + strconv.Itoa(i) + "@gamil.com"
		comment := "Ini comment ke " + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment id ", id)
	}

	// err = tx.Commit()
	err = tx.Rollback()
	if err != nil {
		panic(err)
	}
}

func TestNumGoroutine(t *testing.T)  {
	fmt.Println(runtime.NumGoroutine())
}