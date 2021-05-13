package memory

import (
	"Cookie_Session/data/common"

	//"Cookie_Session/data"
	"container/list"
	"sync"
	"time"
)

//实现IProvider接口
/*type IProvider interface {
	SessionInit(sid string)(ISession, error)
	SessionRead(sid string)(ISession, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}*/
type Provider struct{
	lock sync.Mutex
	sessions map[string]*list.Element
	list *list.List
}

func NewProvider()*Provider{
	return &Provider{
		lock: sync.Mutex{},
		sessions: make(map[string]*list.Element),
		list: list.New(),
	}
}

func (pder *Provider) SessionInit(sid string)(common.ISession, error){
	pder.lock.Lock()
	pder.lock.Unlock()

	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *Provider)SessionRead(sid string)(common.ISession, error){
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (pder *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *Provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Front()
		if element == nil {
			break
		}

		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) <
			time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}


