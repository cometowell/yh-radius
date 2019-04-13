package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Provider interface {
	CreateSession(string, string) ISession
	DestroySession(string) error
	ReadSession(string) ISession
	UpdateSession(string, string, interface{}) error
	GetActiveSessions() []ISession
	SessionGC(int64)
}

// 使用内存存储session
// 需实现Provider接口
type MemoryProvider struct {
	Lock     sync.RWMutex
	Sessions map[string]*list.Element
	SesList  *list.List
}

func (r *MemoryProvider) GetActiveSessions() (s []ISession) {
	s = make([]ISession, r.SesList.Len())
	sesList := r.SesList
	for e := sesList.Front(); e != nil; e = e.Next() {
		session := e.Value.(*Session)
		s = append(s, session)
	}
	return
}

func (r *MemoryProvider) CreateSession(sid, host string) (s ISession) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	session := Session{Id: sid, CreateTime: time.Now().Unix(), LastAccessTime: time.Now().Unix(), Host: host}
	session.Attributes = make(map[string]interface{})
	s = &session
	e := r.SesList.PushFront(s)
	r.Sessions[sid] = e
	return
}

func (r *MemoryProvider) DestroySession(sid string) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	if e, ok := r.Sessions[sid]; ok {
		delete(r.Sessions, sid)
		r.SesList.Remove(e)
	}
	return nil
}

func (r *MemoryProvider) ReadSession(sid string) ISession {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	if e, ok := r.Sessions[sid]; ok {
		se, ok := e.Value.(*Session)
		if !ok {
			return nil
		}
		se.LastAccessTime = time.Now().Unix()
		r.SesList.MoveToFront(e)
		return se
	}
	return nil
}

func (r *MemoryProvider) UpdateSession(sid, key string, value interface{}) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	if e, ok := r.Sessions[sid]; ok {
		se, ok := e.Value.(*Session)
		if !ok {
			return fmt.Errorf("session id is not existed")
		}
		se.LastAccessTime = time.Now().Unix()
		se.SetAttr(key, value)
		r.SesList.MoveToFront(e)
	}
	return nil
}

func (r *MemoryProvider) SessionGC(timeout int64) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	for {
		ele := r.SesList.Back()
		if ele == nil {
			break
		}

		session, _ := ele.Value.(*Session)
		if session.LastAccessTime+timeout <= time.Now().Unix() {
			r.SesList.Remove(ele)
			delete(r.Sessions, session.Id)
		} else {
			break
		}
	}
}

// 使用redis管理session
// 需实现Provider接口
type RedisProvider struct {
	MaxLifeTime int64
}
