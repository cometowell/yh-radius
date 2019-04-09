package main

import (
	"container/list"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"strings"
	"sync"
	"time"
)

// 全局Provide映射关系
var providers = make(map[string]Provider)
var GlobalSessionManager *SessionManager

func initWeb() {
	provider := &MemoryProvider{SesList: list.New()}
	providers["memory"] = provider
	provider.Sessions = make(map[string]*list.Element)
	gsm, err := CreateSessionManager(SessionName, "memory", int64(config["web.session.timeout"].(float64)))
	if err != nil {
		panic(err)
	}
	GlobalSessionManager = gsm
	go GlobalSessionManager.Gc()
}

// session interface
type ISession interface {
	GetAttr(string) interface{}
	SetAttr(string, interface{})
	DelAttr(string)
	SessionId() string
}

type Session struct {
	Id             string
	CreateTime     int64
	LastAccessTime int64
	Attributes     map[string]interface{}
	Host           string
}

func (r *Session) GetCreateTime() time.Time {
	return time.Unix(r.CreateTime, 0)
}

func (r *Session) GetHost() string {
	return r.Host
}

func (r *Session) GetAttr(key string) interface{} {
	val, ok := r.Attributes[key]
	if !ok {
		val = nil
	}
	return val
}

func (r *Session) SetAttr(key string, value interface{}) {
	r.Attributes[key] = value
}

func (r *Session) DelAttr(key string) {
	delete(r.Attributes, key)
}

func (r *Session) SessionId() string {
	return r.Id
}

// ********************** session manager ***************************
type SessionManager struct {
	tokenName   string // cookie or token name
	Lock        sync.RWMutex
	Provider    Provider
	MaxLifeTime int64 // session timeout
}

func CreateSessionManager(tokenName string, providerName string, maxLifeTime int64) (*SessionManager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("unkown provide be used %s, please init it", providerName)
	}
	return &SessionManager{tokenName: tokenName, Provider: provider, MaxLifeTime: maxLifeTime}, nil
}

func (*SessionManager) genSessionId() string {
	uuidVal := uuid.NewV4()
	sessionId := uuidVal.String()
	return strings.Replace(sessionId, "-", "", -1)
}

func (mgr *SessionManager) CreateSession(c *gin.Context) (session ISession) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	sessionId := mgr.genSessionId()
	session = mgr.Provider.CreateSession(sessionId, c.Request.Host)
	return
}

func (mgr *SessionManager) DestroySession(c *gin.Context) error {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	token := c.GetHeader(mgr.tokenName)
	mgr.Provider.DestroySession(token)
	return nil
}

const MaxGcInterval = 5 * 60

func (mgr *SessionManager) Gc() {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()
	logger.Info("web service session timeout scheduled task begins execution")
	mgr.Provider.SessionGC(mgr.MaxLifeTime)
	time.AfterFunc(time.Duration(int64(time.Second)*MaxGcInterval), func() {
		mgr.Gc()
	})
}

func (mgr *SessionManager) GetActiveSessions() []ISession {
	return mgr.Provider.GetActiveSessions()
}

func (mgr *SessionManager) GetSession(sid string) ISession {
	return mgr.Provider.ReadSession(sid)
}

func (mgr *SessionManager) GetSessionByGinContext(c *gin.Context) ISession {
	accessToken := c.GetHeader(SessionName)
	return mgr.Provider.ReadSession(accessToken)
}
