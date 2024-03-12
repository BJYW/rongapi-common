package config

import (
	"log"
	"math/rand"

	"time"

	"git.oschina.net/xujiang/rongapi-common/mgopool"
	"gopkg.in/mgo.v2"
)

//Config 全局配置
var Config TConfig
var initSession *mgo.Session

const (
	ERROR_INTERNAL = iota
)

type TConfig struct {
	MongoDbURL                string `goblet:"mongodb_url,localhost" toml:"mongodb_url"`
	MongoDbName               string `goblet:"mongodb_name,apihub" toml:"mongodb_name"`
	MongoIdcheckDbCollection  string `goblet:"mongodb_idcheck_collection,idcheck_new" toml:"mongodb_idcheck_collection"`
	MongoAuthDbCollection     string `goblet:"mongodb_auth_collection,user"`
	MongoKeeperDbCollection   string `goblet:"mongodb_keeper_collection,keeper"`
	MongoEmployeeDbCollection string `goblet:"mongodb_employee_collection,employee"`
	MongoRechargeDbCollection string `goblet:"mongodb_recharge_collection,recharge"`
	MongoAPIDbCollection      string `goblet:"mongodb_api_collection,api"`
	MongoSdkDbCollection      string `goblet:"mongodb_sdk_collection,sdk"`
	MongoSftUserDbCollection  string `goblet:"mongodb_stf_user_collection,sft_user"`
	MongoLogDbCollection      string `goblet:"mongodb_log_collection,log"`
	MongoBlacklistCollection  string `goblet:"mongodb_blacklist_collection,blacklist"`
	MongoApiApplyCollection   string `goblet:"mongbdb_apply_collection,apply"`
	MongoResult               string `goblet:"mongodb_result,result"`
	InfluxDBUrl               string `goblet:"influx_url,http://192.168.1.201:8086"`
	InfluxDBUserName          string `goblet:"influx_name,"`
	InfluxDBPassword          string `goblet:"influx_pwd,"`
	InfluxDBName              string `goblet:"influx_db,rongapi"`
	InfluxDBCollection        string `goblet:"influx_collection,bill"`
	BillExcelPath             string `goblet:"billpath,../www/public/upload/"`
	BillPathReplace           string `goblet:"pathreplace,../www/public"`

	DHFKey      string `goblet:"dhf_key,mingkeshichaungvip"`
	DHFPassword string `goblet:"dhf_password,mingkeshichaung0523@dhf"`

	ShuJiaUrl string `goblet:"shujia_url,localhost" toml:"shujia_url"`

	AlertsValue float64 `goblet:"alert_value,5000"`
}

func (t *TConfig) DialMgo() (session *mgo.Session, err error) {
	if initSession == nil {
		initSession, err = mgo.Dial(t.MongoDbURL)
		if err != nil {
			log.Printf("initialize MongoOutput failed, %s for %s", err.Error(), t.MongoDbURL)
			return nil, err
		} else {
			initSession.SetMode(mgo.Monotonic, true)
		}
	}
	cloned := initSession.Clone()
	err = cloned.Ping()
	init_n := 500
	for err != nil {
		time.Sleep(time.Duration(rand.Intn(init_n)) * time.Millisecond)
		cloned = initSession.Clone()
		err = cloned.Ping()
		init_n *= 2
	}
	return cloned, nil
}

func (t *TConfig) MgoEmployee() (c *mgo.Collection, err error) {
	if s, err := t.DialMgo(); err == nil {
		return s.DB(t.MongoDbName).C(t.MongoEmployeeDbCollection), nil
	} else {
		return nil, err
	}
}

func (t *TConfig) MgoIDCheck() (c *mgo.Collection, err error) {
	if s, err := t.DialMgo(); err == nil {
		return s.DB(t.MongoDbName).C(t.MongoIdcheckDbCollection), nil
	} else {
		return nil, err
	}
}

func (t *TConfig) MgoAuth() (c *mgo.Collection, err error) {
	if s, err := t.DialMgo(); err == nil {
		return s.DB(t.MongoDbName).C(t.MongoAuthDbCollection), nil
	} else {
		return nil, err
	}
}

func (t *TConfig) MgoSDK() (c *mgo.Collection, err error) {
	if s, err := t.DialMgo(); err == nil {
		return s.DB(t.MongoDbName).C(t.MongoSdkDbCollection), nil
	} else {
		return nil, err
	}
}

func (t *TConfig) MgoAPI() (c *mgo.Collection, err error) {
	if s, err := t.DialMgo(); err == nil {
		return s.DB(t.MongoDbName).C(t.MongoAPIDbCollection), nil
	} else {
		return nil, err
	}
}

var pool *mgopool.Pool

func (t *TConfig) MgoLOG() (c *mgo.Collection, helper *mgopool.SessionHelper, err error) {
	if pool == nil {
		pool, err = mgopool.New(t.MongoDbURL)
	}
	var session *mgo.Session
	session, helper = pool.Session()
	c = session.DB(t.MongoDbName).C(t.MongoLogDbCollection)

	return
}

func (t *TConfig) MgoRecharge() (c *mgo.Collection, helper *mgopool.SessionHelper, err error) {
	if pool == nil {
		pool, err = mgopool.New(t.MongoDbURL)
	}
	var session *mgo.Session
	session, helper = pool.Session()
	c = session.DB(t.MongoDbName).C(t.MongoRechargeDbCollection)

	return
}

func (t *TConfig) MgoKeeper() (c *mgo.Collection, helper *mgopool.SessionHelper, err error) {
	if pool == nil {
		pool, err = mgopool.New(t.MongoDbURL)
	}
	var session *mgo.Session
	session, helper = pool.Session()
	c = session.DB(t.MongoDbName).C(t.MongoKeeperDbCollection)

	return
}

func (t *TConfig) MgoSftUser() (c *mgo.Collection, helper *mgopool.SessionHelper, err error) {
	if pool == nil {
		pool, err = mgopool.New(t.MongoDbURL)
	}
	var session *mgo.Session
	session, helper = pool.Session()
	c = session.DB(t.MongoDbName).C(t.MongoSftUserDbCollection)

	return
}

func (t *TConfig) MgoApply() (c *mgo.Collection, helper *mgopool.SessionHelper, err error) {
	if pool == nil {
		pool, err = mgopool.New(t.MongoDbURL)
	}
	var session *mgo.Session
	session, helper = pool.Session()
	c = session.DB(t.MongoDbName).C(t.MongoApiApplyCollection)

	return
}

func (t *TConfig) MgoResult() (c *mgo.Collection, err error) {
	if s, err := t.DialMgo(); err == nil {
		return s.DB(t.MongoDbName).C(t.MongoResult), nil
	} else {
		return nil, err
	}
}
