package sessions

import (
	"fmt"
	"sync"

	"github.com/gookit/color"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Session struct {
	User        *UserInfo
	Keyboard    tgbotapi.InlineKeyboardMarkup
	CurrentStep string
	PrevStep    string
}

type UserInfo struct {
	UserID     int
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
		Keyboard:    keyboard,
		CurrentStep: currStep,
		PrevStep:    prevStep,
	}
}

func (sm *SessionManager) GetSession(userID int) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	session, exists := sm.sessions[userID]
	return session, exists
}

func (sm *SessionManager) PrintSessionData() {
	for _, s := range sm.sessions {
		color.Yellow.Println("UserID:", s.User.UserID)
		color.Yellow.Println("First_name:", s.User.First_name)
		color.Yellow.Println("Last_name:", s.User.Last_name)
		color.Yellow.Println("User_name:", s.User.User_name)
		color.Yellow.Println("Current step:", s.CurrentStep)
		color.Yellow.Println("Previous step:", s.CurrentStep)
		fmt.Print("___________________\n\n")
	}
}
