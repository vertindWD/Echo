package snowflakeID

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.ParseInLocation("2006-01-02", startTime, time.Local)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

func GenID() int64 {
	if node == nil {
		panic("snowflake node not initialized")
	}
	return node.Generate().Int64()
}
