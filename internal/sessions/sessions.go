package sessions

import (
	"fmt"
	"sync"

	"github.com/Resolution-hash/shop_bot/internal/card"
	"github.com/Resolution-hash/shop_bot/internal/repository"
	"github.com/gookit/color"
)

type Session struct {
	User              *repository.User
	LastUserMessageID int
	LastBotMessageID  int
	PrevStep          string
	CurrentStep       string
	CardManager       *card.CardManager
	CartManager       *card.CartManager
}

func (s *Session) UpdateStep(step string) {
	s.PrevStep = s.CurrentStep
	s.CurrentStep = step
}

type SessionManager struct {
	sessions map[int]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[int]*Session),
	}
}

func (sm *SessionManager) CreateSession(userInfo *repository.User) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[userInfo.UserID] = &Session{
		User:              userInfo,
		LastUserMessageID: 0,
		LastBotMessageID:  0,
		CardManager:       card.NewCardManager(),
		CartManager:       card.NewCartManager(),
	}
	return sm.sessions[userInfo.UserID]
}

func (sm *SessionManager) GetSession(userID int) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	session, exists := sm.sessions[userID]
	return session, exists
}

func (sm *SessionManager) PrintLogs(userID int) {
	s := sm.sessions[userID]
	fmt.Print("___________________\n\n")
	color.Yellowln("UserID:", s.User.UserID)
	color.Yellowln("First_name:", s.User.First_name)
	color.Yellowln("Last_name:", s.User.Last_name)
	color.Yellowln("User_name:", s.User.User_name)
	color.Yellowln("LastUserMessageID:", s.LastUserMessageID)
	color.Yellowln("LastBotMessageID:", s.LastBotMessageID)
	color.Yellowln("PrevStep:", s.PrevStep)
	color.Yellowln("CurrentStep:", s.CurrentStep)
	color.Yellowln("CardManager:", s.CardManager)
	color.Yellowln("CartManager:", s.CartManager)
	fmt.Print("___________________\n\n")

}
