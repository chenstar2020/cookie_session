package common

type IProvider interface {
	SessionInit(sid string)(ISession, error)
	SessionRead(sid string)(ISession, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}


type ISession interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}