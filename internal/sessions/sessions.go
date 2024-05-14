package sessions

import (
	"fmt"
	"sync"

	"github.com/gookit/color"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Session struct {
	User        *UserInfo
	Chat        *ChatInfo
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

type ChatInfo struct {
	ChatID    int64
	MessageID int
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

func (sm *SessionManager) CreateSession(userInfo *UserInfo, chatInfo *ChatInfo, keyboard tgbotapi.InlineKeyboardMarkup, currStep string, prevStep string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[userInfo.UserID] = &Session{
		User:        userInfo,
		Chat:        chatInfo,
		Keyboard:    keyboard,
		CurrentStep: currStep,
		PrevStep:    prevStep,
	}
}

func (sm *SessionManager) UpdateSession(userID int, keyboard tgbotapi.InlineKeyboardMarkup, currStep string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session := sm.sessions[userID]
	session.Keyboard = keyboard
	session.CurrentStep = currStep
}

func (sm *SessionManager) UpdatePrevStep(userID int, prevStep string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session := sm.sessions[userID]
	session.PrevStep = prevStep
}

func (sm *SessionManager) GetSession(userID int) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	session, exists := sm.sessions[userID]
	return session, exists
}

func (sm *SessionManager) PrintSession() {
	for _, s := range sm.sessions {
		fmt.Print("___________________\n\n")
		color.Yellow.Println("UserID:", s.User.UserID)
		color.Yellow.Println("First_name:", s.User.First_name)
		color.Yellow.Println("Last_name:", s.User.Last_name)
		color.Yellow.Println("User_name:", s.User.User_name)
		color.Yellow.Println("ChatID:", s.Chat.ChatID)
		color.Yellow.Println("MessageID:", s.Chat.MessageID)
		color.Yellow.Println("Current step:", s.CurrentStep)
		color.Yellow.Println("Previous step:", s.PrevStep)
		fmt.Print("___________________\n\n")
	}
}

func (sm *SessionManager) PrintSessionByID(userID int) {
	s := sm.sessions[userID]
	fmt.Print("___________________\n\n")
	color.Yellow.Println("UserID:", s.User.UserID)
	color.Yellow.Println("First_name:", s.User.First_name)
	color.Yellow.Println("Last_name:", s.User.Last_name)
	color.Yellow.Println("User_name:", s.User.User_name)
	color.Yellow.Println("ChatID:", s.Chat.ChatID)
	color.Yellow.Println("MessageID:", s.Chat.MessageID)
	color.Yellow.Println("Current step:", s.CurrentStep)
	color.Yellow.Println("Previous step:", s.PrevStep)
	fmt.Print("___________________\n\n")
}
