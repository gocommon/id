/**
 * 用这个生成id
 */
package id

import (
	"os"
	"strconv"
)

// Default Default
var Default *IDWorker

func init() {
	if Default == nil {
		Default, _ = NewIDWorker(0)
	}
}

// Init Init
func Init(id ...int64) error {
	var nodeID int64
	if len(id) > 0 {
		nodeID = id[0]
	}

	if nodeID == 0 {
		envs := os.Getenv("SNOWFLAKE_WORKER_ID")
		if len(envs) > 0 {
			nodeID, _ = strconv.ParseInt(envs, 10, 64)
		}
	}

	var err error
	Default, err = NewIDWorker(nodeID)
	if err != nil {
		return err
	}

	return nil
}

// ID 使用过程中时间倒退会出错
func ID() (int64, error) {
	return Default.NextID()
}

// MustID panic error
func MustID() int64 {
	id, err := Default.NextID()
	if err != nil {
		panic(err)
	}
	return id
}
