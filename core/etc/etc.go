package etc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	validEmailTest = regexp.MustCompile(`^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$`)
)

func Duperr(err error) bool {
	return err != nil && strings.Contains(err.Error(),
		"duplicate key value violates unique constraint",
	)
}

func SanitizeEmail(email string) (string, error) {
	email = strings.ToLower(email)

	if !validEmailTest.MatchString(email) {
		return "", errors.New("invalid email")
	}

	return email, nil
}

func NewTLSConfig(tlsCertPath, tlsKeyPath string) (*tls.Config, error) {
	var tlsErr error = nil

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	tlsConfig.Certificates = make([]tls.Certificate, 1)
	tlsConfig.Certificates[0], tlsErr = tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)

	return tlsConfig, tlsErr

}

// AppendRange adds a LIMIT and OFFSET clause to an existing SQL string, also
// adding values to params and incrementing i, the param counter.
func AppendRange(q string, params []interface{}, i int, from, to uint64) (newQ string, newParams []interface{}, newI int) {

	// Default to old values
	newQ = q
	newI = i
	newParams = append(newParams, params...)

	// If from/to are not empty, then use them for range
	if !(from == 0 && to == 0) {
		newQ = q + fmt.Sprintf(` LIMIT $%d OFFSET $%d `, i+1, i+2)

		newParams = append(newParams, to-from+1, from)

		newI = i + 2
	}

	return newQ, newParams, newI
}

// DBHealthCheck is an HTTP handler that can check that an API is
// up and able to speak to a Postgres DB through sqlx
type DBHealthCheck struct {
	db *sqlx.DB
}

// NewDBHealthCheck returns a DBHealthCheck object that uses this
// db connection
func NewDBHealthCheck(db *sqlx.DB) *DBHealthCheck {
	return &DBHealthCheck{
		db: db,
	}
}

// ServeHTTP tries to run a simple query on the db. On success, it sends
// HTTP 200 with body "OK" to the http.ResponseWriter. On error, it sends
// HTTP 500 with body "ERROR"
func (h *DBHealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var t time.Time

	row := h.db.QueryRow(`SELECT CURRENT_TIMESTAMP`)
	err := row.Scan(&t)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "ERROR\n")
	} else {
		fmt.Fprint(w, "OK\n")
	}
}
