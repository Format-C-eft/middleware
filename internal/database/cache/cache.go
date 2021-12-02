package cache

import (
	"context"
	"errors"
	"time"

	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/go-redis/redis/v8"
)

var _ Client = &clientStuct{}

// ErrorNotFound - global error not found for cache DB
var ErrorNotFound = errors.New("error not found")

// Client - interface DB
type Client interface {
	SaveSession(session *config.SessionInfo) error
	RefreshExpire(session *config.SessionInfo, refreshCookie bool) error
	GetSessionInfo(uid string) (*config.SessionInfo, error)
	GetListSession(login, uid string, activeOnly bool) (*[]SessionDescription, error)
	DropSession(uid, login string) error
	FlushBase() error
}

type clientStuct struct {
	Client  *redis.Client
	Expires struct {
		Session  time.Duration
		User     time.Duration
		Cookie1C time.Duration
	}
}

// ClientDB - global cache client
var ClientDB Client

// InitClient - init client cache DB
func InitClient(cfg *config.Config) error {

	result := clientStuct{
		Client: redis.NewClient(&redis.Options{
			Addr:         cfg.Database.CacheDB.Host,
			Password:     cfg.Database.CacheDB.Password,
			DB:           cfg.Database.CacheDB.DB,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
			MinIdleConns: 5,
		}),
		Expires: struct {
			Session  time.Duration
			User     time.Duration
			Cookie1C time.Duration
		}{
			Session:  cfg.Database.CacheDB.KeysTimeExpires.Session,
			User:     cfg.Database.CacheDB.KeysTimeExpires.User,
			Cookie1C: cfg.Database.CacheDB.KeysTimeExpires.Cookie,
		},
	}

	if _, err := result.Client.Ping(context.TODO()).Result(); err != nil {
		return err
	}

	ClientDB = &result

	return nil

}

type nameKeys struct {
	SessionUserKey   string
	SessionCookie1C  string
	Users            string
	UsersSessionList string
}

func getNamesKeys(Login, ID string) *nameKeys {
	return &nameKeys{
		SessionUserKey:   "Session:" + ID + ":UserKey",
		SessionCookie1C:  "Session:" + ID + ":Cookie1C",
		Users:            "Users:" + Login + ":Info",
		UsersSessionList: "Users:" + Login + ":SessionList",
	}
}

// SessionDescription - session description
type SessionDescription struct {
	SessionID string           `json:"sessionID"`
	IsActive  bool             `json:"isActive"`
	Info      config.AgentInfo `json:"info"`
}

// CheckActive - check active
func (session *SessionDescription) CheckActive() {
	session.IsActive = session.Info.EndsTime.After(time.Now().UTC())
}
