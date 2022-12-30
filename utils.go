package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type PasswordValidationError struct {
	message string
}

func (p *PasswordValidationError) Error() string {
	return p.message
}

type EmailValidationError struct {
	message string
}

func (e *EmailValidationError) Error() string {
	return e.message
}

func validatePassword(pass string) error {

	mesShort := "The password must be at least 8 characters long"
	mesInvalidCharSet :=
		"The password must contain uppercase, lowercase letters and at least one number"

	if len(pass) < 8 {
		return &PasswordValidationError{mesShort}
	}

	m, _ := regexp.Match("[0-9]", []byte(pass))
	if strings.ToLower(pass) == pass || !m {
		return &PasswordValidationError{mesInvalidCharSet}
	}

	return nil
}

func validateEmail(email string) error {
	prefixIsInvalid := func(_email string) bool {
		switch _email[0] {
		case '.', '_', '-':
			return true
		}
		prefix := strings.Split(_email, "@")[0]
		prLastChar := prefix[len(prefix)-1]
		if hasDouble, _ := regexp.MatchString(`\.\.|__|--`, prefix); hasDouble {
			return true
		}
		switch prLastChar {
		case '.', '_', '-':
			return true
		}
		return false
	}
	domainIsInvalid := func(_email string) bool {
		return false
	}
	const message = "The email has an invalid format"
	re := `^([a-zA-Z0-9\.\-\+_]+)@([a-zA-Z0-9_]+)(\.[a-zA-Z]{2,5})+$`
	matched, _ := regexp.MatchString(re, email)
	switch {
	case !matched,
		prefixIsInvalid(email),
		domainIsInvalid(email):
		return &EmailValidationError{message}
	}

	return nil
}

type QueryParams struct {
	name, conf, div, year string
}

func getTeamsQueryParams(c *gin.Context) QueryParams {
	return QueryParams{
		name: c.Query("name"),
		conf: c.Query("conference"),
		div:  c.Query("division"),
		year: c.Query("est_year"),
	}
}

func validateTeamsQueryParams(c *gin.Context) error {
	q := getTeamsQueryParams(c)
	if q.name != "" && (q.conf != "" || q.div != "" || q.year != "") {
		return errors.New("if name filter is present, other filters are not allowed")
	}
	return nil
}

func getTeamsSQLQuery(c *gin.Context) string {
	qp := getTeamsQueryParams(c)

	var query string

	if qp.name != "" {
		query = fmt.Sprintf("select * from teams where name='%s'", strings.Trim(qp.name, "\""))
	} else {
		query = "select * from teams"
		var filters []string
		if qp.conf != "" {
			filters = append(filters, fmt.Sprintf("conference = '%s'", strings.Trim(qp.conf, "\"")))
		}
		if qp.div != "" {
			filters = append(filters, fmt.Sprintf("division = '%s'", strings.Trim(qp.div, "\"")))
		}
		if qp.year != "" {
			filters = append(filters, fmt.Sprintf("est_year = %s", strings.Trim(qp.year, "\"")))
		}
		if len(filters) > 0 {
			query += " where "
			filters := strings.Join(filters, " and ")
			query += filters
		}
	}
	fmt.Printf("Teams sql query:\n[ %s ]\n", query)
	return query
}

type Token struct {
	Username, Role string
}

func parseToken(tok string) (Token, error) {
	tokenParts := strings.Split(tok, "_")
	if len(tokenParts) != 3 {
		return Token{}, errors.New("incorrect token format. Proper format: role_token_username")
	}
	return Token{Username: tokenParts[2], Role: tokenParts[0]}, nil
}

// func verifyToken(token string) error {
// 	ctx := context.TODO()

// 	conn, err := getDbConnection()
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close(ctx)

// 	tok, err := parseToken(token)
// 	if err != nil {
// 		return err
// 	}

// 	var role string
// 	row := conn.QueryRow(
// 		ctx,
// 		"Select r.name from users u "+
// 			"join roles r on u.roleid = r.id "+
// 			"where username=$1",
// 		tok.Username)
// 	row.Scan(&role)

// 	if role == "" {
// 		return errors.New("invalid username")
// 	}
// 	if role != tok.Role {
// 		return errors.New("incorrect user role")
// 	}
// 	return nil
// }

func getDbConnection() (*pgx.Conn, error) {
	connUrl := fmt.Sprintf("postgres://postgres:%s@localhost:5432/postgres", os.Getenv("DBPass"))
	conn, err := pgx.Connect(context.TODO(), connUrl)
	if err != nil {
		return nil, fmt.Errorf("Error establishing a connections: \n%v", err)
	}
	return conn, nil
}
