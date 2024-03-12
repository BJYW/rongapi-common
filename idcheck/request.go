package idcheck

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"git.oschina.net/xujiang/rongapi-common/utils"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"

	"git.oschina.net/xujiang/rongapi-common/charge"
	"git.oschina.net/xujiang/rongapi-common/config"

	"github.com/golang/glog"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const sync_collection = 20

var insertChan = make(chan *Result)
var findCollectionChan = make(chan *mgo.Collection, sync_collection)

type Request struct {
	IdCode string
	Name   string
	Seqno  string
	Telno  string
}

func (r *Request) IsEmpty() bool {
	return r.Name == "" || r.IdCode == ""
}

func (r *Request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("IdCode", r.IdCode)
	enc.AddString("Name", r.Name)
	enc.AddString("Seqno", r.Seqno)
	enc.AddString("Telno", r.Telno)
	return nil
}

//BeforeRequest
func (r *Request) PreProgress() error {
	hash := md5.New()
	//虽然本处有telno做键，但固定为空
	hash.Write([]byte(fmt.Sprintf("idcode=%s&name=%s&telno=", r.IdCode, r.Name, "")))
	bts := hash.Sum(nil)
	r.Seqno = hex.EncodeToString(bts[:10])
	glog.Infoln("id check req", zap.Object("request", zapcore.ObjectMarshaler(r)))
	return nil
}

// func (r *Request) String() string {
// 	return fmt.Sprintf("idcode=%s&name=%s&seqno=%s&telno=%s", r.IdCode, r.Name, r.Seqno, r.Telno)
// }

//LoadSuccess 获取成功
func (r *Request) PostProgress(res interface{}, cost int, err error, remark map[string]charge.Remark) (interface{}, int, error, map[string]charge.Remark) {
	if result, ok := res.(*Result); ok {
		var code int
		if code, err = strconv.Atoi(result.ResultCode); err == nil {
			switch code {
			case 1000:
				//正确的时候，进行请求缓存
				result.MobilePhone = r.Telno
				result.Seqno = r.Seqno
				result.IdCode = utils.Encode(result.IdCode)
				//go result.Store(nil)
			default:
				glog.Infoln("not normal result", zap.Int("code", code))
			}
		} else {
			glog.Infoln("err in store", zap.Error(err))
		}
	} else {
		glog.Infoln("result is not *Result", zap.Reflect("result", res))
	}
	return res, cost, err, remark
}

func (r *Request) Exists() (*Result, bool) {
	for {
		c := <-findCollectionChan
		var result Result
		if err := c.Find(bson.M{"seqno": r.Seqno}).One(&result); err == nil {
			glog.Infoln("find seqno of existed", zap.String("seqno", r.Seqno), zap.Int("count", 1))
			findCollectionChan <- c
			return &result, true
		} else if err == mgo.ErrNotFound {
			glog.Infoln("not found")
			findCollectionChan <- c
			return nil, false
		} else {
			c.Database.Session.Close()
			findCollectionChan <- createCollection()
		}
	}
}

func createCollection() *mgo.Collection {
	var inter = 500 * time.Millisecond
	for {
		if newC, err := config.Config.MgoIDCheck(); err == nil {
			return newC
		} else {
			glog.Warningln("create Mongo Error", zap.Error(err))
			time.Sleep(inter)
			inter *= 2
			if inter < 30*time.Second {
				inter = 30 * time.Second
			}
		}
	}
}

func Init() {
	for index := 0; index < sync_collection; index++ {
		findCollectionChan <- createCollection()
	}
}
