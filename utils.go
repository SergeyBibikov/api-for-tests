package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
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
	if len(pass) < 8 {
		return &PasswordValidationError{"The password must be at least 8 characters long"}
	}
	if strings.ToLower(pass) == pass {
		return &PasswordValidationError{"The password must contain uppercase and lowercase letters"}
	}
	return nil
}

func validateEmail(email string) error {
	re := regexp.MustCompile(`^([a-zA-Z0-9\.\-\+_]+)@([a-zA-Z0-9\.\-_]+)\.([a-zA-Z]{2,5})$`)

	const message = "The email has an invalid format"
	if isFound := re.Find([]byte(email)); isFound == nil {
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

func getUserRole(token string) string {
	t := strings.Split(token, "_")
	return t[0]
}
