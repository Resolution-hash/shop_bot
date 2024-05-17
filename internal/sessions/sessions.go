package sessions

import (
	"fmt"
	"sync"

	"github.com/Resolution-hash/shop_bot/internal/card"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gookit/color"
)

type Session struct {
	User        *UserInfo
	CardManager *card.CardManager
	Keyboard    tgbotapi.InlineKeyboardMarkup
	CurrentStep string
	PrevStep    string
}

type UserInfo struct {
	UserID     int
	MessageID  int
	First_name string
	Last_name  string
	User_name  string
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

func (sm *SessionManager) CreateSession(userInfo *UserInfo, keyboard tgbotapi.InlineKeyboardMarkup, currStep string, prevStep string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[userInfo.UserID] = &Session{
		User:        userInfo,
		CardManager: card.NewCardManager(),
		Keyboard:    keyboard,
		CurrentStep: currStep,
		PrevStep:    prevStep,
	}
}

func (sm *SessionManager) UpdateSession(userID int, keyboard tgbotapi.InlineKeyboardMarkup, newStep string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session := sm.sessions[userID]
	session.Keyboard = keyboard
	if newStep != "prev" && newStep != "next" {
		session.PrevStep = session.CurrentStep
		session.CurrentStep = newStep
	}

}

func (sm *SessionManager) UpdateStep(userID int, newStep string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session := sm.sessions[userID]
	if newStep != "prev" && newStep != "next" {
		session.PrevStep = session.CurrentStep
		session.CurrentStep = newStep
	}

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
	color.Yellowln("MessageID:", s.User.MessageID)
	color.Yellowln("CardManager:", s.CardManager)
	color.Yellowln("First_name:", s.User.First_name)
	color.Yellowln("Last_name:", s.User.Last_name)
	color.Yellowln("User_name:", s.User.User_name)
	color.Yellowln("Current step:", s.CurrentStep)
	color.Yellowln("Previous step:", s.PrevStep)
	fmt.Print("___________________\n\n")

}
