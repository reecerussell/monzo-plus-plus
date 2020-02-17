package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func main() {
	for true {
		process()
	}
}

func process() {
	tokens, err := getAccessTokens()
	if err != nil {
		panic(err)
	}

	if len(tokens) > 0 {
		t := tokens[0].Expires
		if (t.Sub(time.Now().UTC()).Hours()) >= 1 {
			time.Sleep(1 * time.Hour)
		}
	} else {
		time.Sleep(5 * time.Minute)
		return
	}

	wg := &sync.WaitGroup{}

	for _, t := range tokens {
		go func(ac *AccessToken) {
			wg.Add(1)
			defer wg.Done()

			err := refreshToken(ac.UserID, ac.RefreshToken)
			if err != nil {
				log.Printf("failed to refresh token for user: %s: %v\n", ac.UserID, err)
			}
		}(t)
	}

	wg.Wait()
}

func open() {
	var err error

	if db == nil {
		db, err = sql.Open("mysql", os.Getenv("CONN_STRING"))
		if err != nil {
			panic(err)
		}
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
}

func getAccessTokens() ([]*AccessToken, error) {
	open()

	query := `SELECT 
					user_id, refresh_token, expires, token_type
				FROM
					user_tokens
				ORDER BY expires DESC;`

	ctx := context.Background()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var tokens []*AccessToken

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ac AccessToken

		err = rows.Scan(&ac.UserID, &ac.RefreshToken, &ac.Expires, &ac.TokenType)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, &ac)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func refreshToken(userID, refreshToken string) error {
	log.Printf("REFRESHING ACCESS TOKEN\n\tUser: %s\n\tToken: %s\n", userID, refreshToken)

	// get access token
	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("client_id", os.Getenv("MONZO_CLIENT_ID"))
	body.Set("client_secret", os.Getenv("MONZO_CLIENT_SECRET"))
	body.Set("refresh_token", refreshToken)

	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.PostForm("https://api.monzo.com/oauth2/token", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var message Error
		_ = json.NewDecoder(resp.Body).Decode(&message)

		return fmt.Errorf("failed to request refresh token: %d: %s", resp.StatusCode, message.Message)
	}

	var data AccessToken
	_ = json.NewDecoder(resp.Body).Decode(&data)

	query := `UPDATE user_tokens 
				SET 
					access_token = ?,
					refresh_token = ?,
					expires = ?,
					token_type = ?
				WHERE
					user_id = ?;`
	args := []interface{}{
		data.AccessToken,
		data.RefreshToken,
		time.Now().UTC().Add(time.Second * time.Duration(data.ExpiresIn)),
		data.TokenType,
		userID,
	}

	open()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	return nil
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`

	Expires time.Time `json:"-"`
}

type Error struct {
	Message string `json:"message"`
}
