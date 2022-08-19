package singleton

import "sync"

type Singleton struct {
}

var (
	singleton *Singleton
	mux       sync.Mutex
	once      sync.Once
)

// 懒汉模式
func SingletonFun() *Singleton {
	if singleton == nil {
		singleton = &Singleton{}
	}
	return singleton
}

// 懒汉+锁模式
func SingletonMuxFun() *Singleton {
	mux.Lock()
	defer mux.Unlock()
	if singleton == nil {
		singleton = &Singleton{}
	}
	return singleton
}

// 懒汉+检查锁模式
func SingletonCheckMuxFun() *Singleton {
	if singleton == nil {
		mux.Lock()
		defer mux.Unlock()
		if singleton == nil {
			singleton = &Singleton{}
		}
	}
	return singleton
}

// sync.once
func SingletonOnce() *Singleton {
	once.Do(func() {
		singleton = &Singleton{}
	})
	return singleton
}
