package memory

import (
	"container/list"
	"time"
)

var pder = &Provider{list: list.New()}
//实现session接口
/*type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}*/
type SessionStore struct {
	sid string
	timeAccessed time.Time
	value map[interface{}]interface{}
}

func (st *SessionStore)Set(key, value interface{}) error{
	pder.SessionUpdate(st.sid)
	st.value[key] = value
	return nil
}

func (st *SessionStore)Get(key interface{}) interface{}{
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (st *SessionStore)Delete(key interface{}) error {
	pder.SessionUpdate(st.sid)
	delete(st.value, key)
	return nil
}

func (st *SessionStore)SessionID()string{
	return st.sid
}

