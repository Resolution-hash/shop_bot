package sessions

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Resolution-hash/shop_bot/internal/card"
	db "github.com/Resolution-hash/shop_bot/internal/repository"
	user "github.com/Resolution-hash/shop_bot/internal/repository/user"
	"github.com/Resolution-hash/shop_bot/internal/services"
	"github.com/gookit/color"
	_ "github.com/lib/pq"
)

type Session struct {
	User              *user.User
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

func (sm *SessionManager) CreateSession(userInfo *user.User) *Session {
	isAdmin, err := addUserToDB(userInfo)
	if err != nil {
		color.Redln("The user has not been added to the database:", err)
	}
	if isAdmin {
		userInfo.IsAdmin = 1
	}
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
	color.Yellowln("IsAdmin:", s.User.IsAdmin)
	color.Yellowln("LastUserMessageID:", s.LastUserMessageID)
	color.Yellowln("LastBotMessageID:", s.LastBotMessageID)
	color.Yellowln("PrevStep:", s.PrevStep)
	color.Yellowln("CurrentStep:", s.CurrentStep)
	color.Yellowln("CardManager:", s.CardManager)
	color.Yellowln("CartManager:", s.CartManager)
	fmt.Print("___________________\n\n")
}

func addUserToDB(user *user.User) (bool, error) {
	db, err := db.SetupDatabase()
	if err != nil {
		color.Redln(err)
		return false, err
	}
	defer db.Close()

	svc := initUserService(db)

	isAdmin, err := svc.AddUser(*user)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func initUserService(db *sql.DB) services.UserService {
	repo := user.NewPostgresUserRepo(db)
	return *services.NewUserService(repo)
}
