package idcheck

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"git.oschina.net/xujiang/rongapi-common/config"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Result struct {
	Id_ bson.ObjectId `bson:"_id"`
	// :业务流水号
	Seqno string `json:"seqno"`
	// : 姓名
	Name string `json:"name" xml:"nameIdInfos>nameIdInfo>outputXm"`
	// : 身份证号码
	IdCode string `json:"idcode" xml:"nameIdInfos>nameIdInfo>outputZjhm"`
	// : 核查结果（1000一致；1002-库中无此号；1001-不一致； 2003-身份号码不合规则；2004-姓名不合规则；9901-系统异常）
	ResultCode string `json:"resultCode" xml:"nameIdInfos>nameIdInfo>code"`
	// :结果描述
	ResultMsg string `json:"resultMsg" xml:"nameIdInfos>nameIdInfo>message"`
	// :性别（男性；女性）
	Gender string `json:"gender" xml:"nameIdInfos>nameIdInfo>xb"`
	// :生日
	Birthday string `json:"birthday" xml:"nameIdInfos>nameIdInfo>csrq"`
	// :民族
	Nationality string `json:"nationality" xml:"nameIdInfos>nameIdInfo>mz"`
	// :所属省市县区
	Ssssxq string `json:"ssssxq" xml:"nameIdInfos>nameIdInfo>ssssxq"`

	Photo interface{} `json:"photo" xml:"nameIdInfos>nameIdInfo>xp"`
	// :住址
	Address string `json:"address" xml:"nameIdInfos>nameIdInfo>zz"`
	// :识别码
	Oid string `json:"oid"`
	// :手机号
	MobilePhone interface{} `json:"mobilephone"`
	// :坐过火车？
	TrainChecked bool `json:"trainchecked"`
}

func (r *Result) CheckSeqNo() {
	if r.Seqno == "" {
		hash := md5.New()
		//虽然本处有telno做键，但固定为空
		hash.Write([]byte(fmt.Sprintf("idcode=%s&name=%s&telno=", r.IdCode, r.Name, "")))
		bts := hash.Sum(nil)
		r.Seqno = hex.EncodeToString(bts[:10])
	}
	if !r.Id_.Valid() {
		r.Id_ = bson.NewObjectId()
	}
}

func (r *Result) Store(onSuccess func()) {
	r.CheckSeqNo()
	for {
		var inter = 500 * time.Millisecond
		var err error
		var c *mgo.Collection
		if c, err := config.Config.MgoIDCheck(); err == nil {
			defer c.Database.Session.Close()
			iter := c.Find(bson.M{"seqno": r.Seqno}).Iter()
			var hasOld = 0

			old := new(Result)
			firstOld := new(Result)
			for iter.Next(old) {
				hasOld++
				if hasOld == 1 {
					old.merge(r)
					firstOld = old
				} else {
					firstOld.merge(old)
					c.RemoveId(old.Id_)
				}
				old = nil
			}
			if hasOld > 0 {
				err = c.UpdateId(firstOld.Id_, firstOld)
			} else if hasOld == 0 {
				err = c.Insert(r)
			}
		}
		if err != nil {
			time.Sleep(inter)
			inter *= 2
			if inter > 10*time.Second {
				inter = 10 * time.Second
			}
			c.Database.Session.Close()
			glog.Infoln("Upsert Mongo Error", err, "unstored result")
		} else {
			if onSuccess != nil {
				onSuccess()
			}
			return
		}
	}
}

func (oldChecked *Result) merge(newChecked *Result) {
	oldChecked.Photo = merge(oldChecked.Photo, newChecked.Photo)
	oldChecked.MobilePhone = merge(oldChecked.MobilePhone, newChecked.MobilePhone)
}

//
func merge(oldChecked interface{}, newChecked interface{}) interface{} {
	//将old merge 到new里面
	switch told := oldChecked.(type) {
	//如果old里面是string
	case string:
		if told != "" {
			switch tnew := newChecked.(type) {
			//如果新的是string
			case string: //ok
				if tnew != "" {
					if told != tnew {
						newChecked = []interface{}{told, tnew}
					} else {
						newChecked = tnew
					}
				} else {
					newChecked = told
				}
				//如果新的是数组
			case []interface{}:
				//将旧的转换成数组
				tmp := []interface{}{told}
				for _, singleNew := range tnew {
					if singleNew != told && singleNew != "" {
						tmp = append(tmp, singleNew)
					}
				}
				newChecked = tmp
			}
		}
		//如果old里面是数组
	case []interface{}:
		switch tnew := newChecked.(type) {
		case string:
			if tnew != "" {
				var shouldStore = true
				for _, singleOld := range told {
					shouldStore = shouldStore && (tnew != singleOld)
				}
				if shouldStore {
					newChecked = append(told, tnew)
				}
				return newChecked
			}
			newChecked = told
		case []interface{}:
			var tmp = told
			for _, singleNew := range tnew {
				if sSingleNew, ok := singleNew.(string); ok {
					var shouldStore = true
					for _, singleOld := range told {
						if sSingleOld, ok := singleOld.(string); ok {
							shouldStore = shouldStore && (sSingleNew != sSingleOld && sSingleNew != "")
						}
					}
					if shouldStore {
						tmp = append(tmp, sSingleNew)
					}
				}
			}
			newChecked = tmp
		}
	}
	return newChecked
}
