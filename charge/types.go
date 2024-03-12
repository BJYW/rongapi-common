package charge

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/souriki/ali_mns"

	"github.com/influxdata/influxdb/client/v2"
	"gopkg.in/mgo.v2/bson"
)

//用户可以快速调用的阈值，该值应该和系统输出能力相匹配
const Threshold = 100

type Charge struct {
	UserID        bson.ObjectId
	TransactionID string
	Unit          int
	Stat          Stat
}

func (c *Charge) WriteToInfluxDB(cli client.Client) error {

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "idtag",
	})

	// Create a point and add to batch
	tags := map[string]string{"userid": c.UserID.Hex(), "api": c.Stat.API, "source": c.Stat.Source, "remote": c.Stat.RemoteIp}
	for k, v := range c.Stat.APIRemark {
		if !v.JustForLog {
			tags[k] = fmt.Sprintf("%s", v.Value)
		}
	}

	diff := c.Stat.RecvAt.Sub(c.Stat.SendAt)

	fields := map[string]interface{}{
		"cost":  c.Unit,
		"delay": diff.Seconds(),
	}
	pt, err := client.NewPoint("bill", tags, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)

	// Write the batch
	return cli.Write(bp)
}

type RemoveAPIMessage struct {
	APPID string
	API   string
	//该客户不需要等待返回的窗口
}

type Stat struct {
	StartAt    time.Time
	SendAt     time.Time
	RecvAt     time.Time
	CallBackAt time.Time
	API        string
	Source     string
	RemoteIp   string
	Request    interface{}
	APIRemark  map[string]Remark
}

type Remark struct {
	Value      interface{}
	JustForLog bool
}

//Push 向阿里推送相关的消息
func Push(obj interface{}, queue ali_mns.AliMNSQueue) (err error) {
	var bts []byte
	if bts, err = json.Marshal(obj); err == nil {
		msg := ali_mns.MessageSendRequest{
			MessageBody:  string(bts),
			DelaySeconds: 0,
			Priority:     1}
		_, err = queue.SendMessage(msg)
	}
	return err
}
