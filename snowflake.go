package id

import (
	"errors"
	"sync"
	"time"
)

var (
	// CEpoch CEpoch 可以根据项目时间更新
	CEpoch int64 = 1560407323000
)

const (

	// CWorkerIDBits 最大支持16
	CWorkerIDBits = 5 // Num of WorkerID Bits
	// CSenquenceBits CSenquenceBits
	CSenquenceBits = 12 // Num of Sequence Bits

	// CWorkerIDShift CWorkerIDShift
	CWorkerIDShift = 12
	// CTimeStampShift CTimeStampShift
	CTimeStampShift = 17
)

// IDWorker Struct
type IDWorker struct {
	workerID      int64
	lastTimeStamp int64
	sequence      int64
	sequenceMask  int64
	maxWorkerID   int64
	lock          *sync.Mutex
}

// NewIDWorker Func: Generate NewIDWorker with Given workeriD
func NewIDWorker(workeriD int64) (iw *IDWorker, err error) {
	iw = new(IDWorker)

	iw.maxWorkerID = getMaxWorkerID()

	if workeriD > iw.maxWorkerID || workeriD < 0 {
		return nil, errors.New("workerid not fit")
	}
	iw.workerID = workeriD
	iw.lastTimeStamp = -1
	iw.sequence = 0
	iw.sequenceMask = getSequenceMask()
	iw.lock = new(sync.Mutex)
	return iw, nil
}

func getMaxWorkerID() int64 {
	return -1 ^ -1<<CWorkerIDBits
}

func getSequenceMask() int64 {
	return -1 ^ -1<<CSenquenceBits
}

func (iw *IDWorker) timeGen() int64 {
	return int64(time.Now().UnixNano() / 1000000)
}

func (iw *IDWorker) timeReGen(last int64) int64 {
	ts := int64(time.Now().UnixNano() / 1000000)
	for {
		if ts <= last {
			ts = iw.timeGen()
		} else {
			break
		}
	}
	return ts
}

// NextID Func: Generate next id
func (iw *IDWorker) NextID() (ts int64, err error) {
	iw.lock.Lock()
	defer iw.lock.Unlock()
	ts = iw.timeGen()
	if ts == iw.lastTimeStamp {
		iw.sequence = (iw.sequence + 1) & iw.sequenceMask
		if iw.sequence == 0 {
			ts = iw.timeReGen(ts)
		}
	} else {
		iw.sequence = 0
	}

	// 时间倒退了，可以生成重复id
	if ts < iw.lastTimeStamp {
		err = errors.New("Clock moved backwards, Refuse gen id")
		return 0, err
	}
	iw.lastTimeStamp = ts
	ts = ((ts - CEpoch) << CTimeStampShift) | (iw.workerID << CWorkerIDShift) | iw.sequence
	return ts, nil
}
