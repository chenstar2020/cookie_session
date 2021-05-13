package main

import (
	"Cookie_Session/data"
	"Cookie_Session/data/common"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	lock        sync.Mutex
	cookieName  string           //private cookieName
	provider    common.IProvider //
	maxLifeTime int64            //max life time
}

func NewManager(provideName, cookieName string, maxLifeTime int64)(*Manager, error){
	provider, ok := data.Provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgottenimport?)", provideName)
	}

	manager := &Manager{lock: sync.Mutex{}, cookieName: cookieName, provider: provider, maxLifeTime: maxLifeTime}

	go manager.GC()

	return manager, nil
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request)(session common.ISession){
	manager.lock.Lock()
	defer manager.lock.Unlock()

	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid),
			Path:"/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	}else{
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}else{
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		cookie := http.Cookie{Name:manager.cookieName, Path:"/",
			HttpOnly:true, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC(){
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func(){manager.GC()})
}