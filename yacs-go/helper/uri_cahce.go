package helper

import (
	"sync"

	"github.com/iostrovok/yacs-go/yacs-go/utils"
)

type uriCacheSt struct {
	mu   sync.RWMutex
	docs map[string]interface{}
}

var uriCache *uriCacheSt

func init() {
	uriCache = &uriCacheSt{
		docs: map[string]interface{}{},
	}
}

func (uc *uriCacheSt) Add(uri string, doc interface{}) {

	mycopy, err := utils.DeepCopy(doc)
	if err != nil {
		return
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.docs[uri] = mycopy
}

func (uc *uriCacheSt) Get(uri string) (interface{}, bool) {

	uc.mu.RLock()
	defer uc.mu.RUnlock()

	out, find := uc.docs[uri]
	if !find {
		return nil, false
	}

	out, err := utils.DeepCopy(out)
	if err != nil {
		return nil, false
	}

	return out, true
}
