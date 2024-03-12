package mgopool

import (
	"math/rand"
	"sync"
	"time"

	"git.oschina.net/xujiang/rongapi-common/mgopool/cmap_string_session"

	"go.uber.org/zap"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

const (
	defaultPoolSize = 20
)

var initSessions = cmap_string_session.New()

func New(url string, size ...int) (pool *Pool, err error) {
	var initSession *mgo.Session
	var ok bool
	if initSession, ok = initSessions.Get(url); !ok {
		initSession, err = mgo.Dial(url)
		if err != nil {
			glog.Errorln("initialize MongoOutput failed", err)
			return nil, err
		} else {
			initSession.SetMode(mgo.Monotonic, true)
		}
		initSessions.Set(url, initSession)
	}

	var poolSize = defaultPoolSize
	if len(size) > 0 && size[0] > 0 {
		poolSize = size[0]
	}

	pool = new(Pool)
	pool.initSession = initSession
	pool.recircleChan = make(chan *mgo.Session, poolSize)
	pool.errChan = make(chan *mgo.Session, poolSize)
	go pool.handleError()
	return pool, nil
}

type Pool struct {
	db           string
	collection   string
	recircleChan chan *mgo.Session
	errChan      chan *mgo.Session
	initSession  *mgo.Session
	createdSize  int
	lock         sync.RWMutex
}

func (p *Pool) Session() (session *mgo.Session, helper *SessionHelper) {
	session, helper = p.session()
	err := session.Ping()
	init_n := 500
	for err != nil {
		glog.Errorln("get session from pool happen error,retry...", zap.Error(err))
		time.Sleep(time.Duration(rand.Intn(init_n)) * time.Millisecond)
		helper.Close(err)
		session, helper = p.session()
		err = session.Ping()
		init_n *= 2
	}
	return
}

func (p *Pool) session() (session *mgo.Session, helper *SessionHelper) {

	if len(p.recircleChan) == 0 && cap(p.recircleChan) > p.createdSize {
		//should created
		p.lock.Lock()
		if cap(p.recircleChan) > p.createdSize {
			p.createdSize++
		} else {
			//有人比我更快的获得了创建的权利
			goto wait
		}
		p.lock.Unlock()

		session = p.initSession.Clone()
		err := session.Ping()
		init_n := 500
		for err != nil {
			time.Sleep(time.Duration(rand.Intn(init_n)) * time.Millisecond)
			session = p.initSession.Clone()
			err = session.Ping()
			init_n *= 2
		}
		helper = &SessionHelper{session, p.recircleChan, p.errChan}
		return
	}
wait:
	session = <-p.recircleChan
	helper = &SessionHelper{session, p.recircleChan, p.errChan}
	return
}

func (p *Pool) handleError() {
	for {
		select {
		case c := <-p.errChan:
			c.Close()
			p.createdSize--
		}
	}
}

type SessionHelper struct {
	session      *mgo.Session
	recircleChan chan *mgo.Session
	errChan      chan *mgo.Session
}

func (s *SessionHelper) Close(err error) {
	if err == nil {
		s.recircleChan <- s.session
	} else {
		glog.Errorln("recycle mongo by error", zap.Error(err))
		s.errChan <- s.session
	}
}
