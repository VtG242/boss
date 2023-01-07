package models

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

// Define a Player type to hold the data for an individual player
type Player struct {
	PID       uint64
	Surname   string
	Firstname string
	Sex       string
	Birthdate time.Time
	Town      string
	Country   string
	Nickname  string
	Hash      string
	Email 		sql.NullString
}

// Define a PlayerModel type which wraps a sql.DB connection pool.
type PlayersModel struct {
	Pool *sql.DB
}

// create md5 hash from given string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// This will insert new player into the database.
func (db *PlayersModel) Insert(
	surname string,
	firstname string,
	sex string,
	birthdate time.Time,
	town string,
	country string,
	nickname string,
	hash string) (int, error) {

	// use parameterized SQL statement
	result, err := db.Pool.Exec(
		"INSERT INTO Players (surname,firstname,sex,birthdate,town,country,nickname,hash) VALUES (?,?,?,?,?,?,?,?)",
		surname, firstname, sex, birthdate, town, country, nickname, hash)
	if err != nil {
		return 0, err
	}

	// use the LastInsertId() method on the result to get the ID of new player
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// This will return a specific player
func (db *PlayersModel) Get(id int) (*Player, error) {
	// Initialize a pointer to a new zeroed player
	p := &Player{}

	err := db.Pool.QueryRow("SELECT * FROM Players WHERE pid=?", id).Scan(
		&p.PID, &p.Surname, &p.Firstname, &p.Sex, &p.Birthdate, &p.Town, &p.Country, &p.Nickname, &p.Hash, &p.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything went OK then return the player object.
	return p, nil
}

// This will return all players
func (db *PlayersModel) All() ([]*Player, error) {
	// Write the SQL statement we want to execute.
	stmt := `Select * FROM Players`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := db.Pool.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize an empty slice to hold the player structs.
	players := []*Player{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Player struct.
		p := &Player{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&p.PID, &p.Surname, &p.Firstname, &p.Sex, &p.Birthdate, &p.Town, &p.Country, &p.Nickname, &p.Hash, &p.Email)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		players = append(players, p)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Players slice.
	return players, nil
}
