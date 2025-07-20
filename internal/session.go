package internal

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

type SessionData struct {
	UserID int
	Role   string
}

var sessions = make(map[string]SessionData)

func CreateSession(userID int, role string) (string, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)
	sessions[token] = SessionData{
		UserID: userID,
		Role:   role,
	}
	log.Printf("Session created for UserID: %d with token: %s", userID, token)
	return token, nil
}

func GetSessionData(token string) (SessionData, bool) {
	session, found := sessions[token]
	return session, found
}

func DeleteSession(token string) {
	delete(sessions, token)
	log.Printf("Session deleted for token: %s", token)
}
