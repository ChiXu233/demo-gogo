package redis

import (
	"demo-gogo/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/rs/xid"
	"time"

	"github.com/go-redis/redis"
	"github.com/wonderivan/logger"
)

const (
	ProcessStart = "start"
	Processing   = "doing"
	ProcessEnd   = "done"
	ProcessAbort = "abort"

	ResultSuccess = "success"
	ResultFailed  = "failed"
	ResultCancel  = "cancel"
)

var (
	RedisClient *redis.Client
	NilError    = redis.Nil
	LockError   = fmt.Errorf("获取锁失败")
	LockTimeOut = fmt.Errorf("获取锁超时")
	UnLockError = fmt.Errorf("释放锁失败")

	ProcessOrder = map[string]int{
		ProcessStart: 0,
		Processing:   1,
		ProcessEnd:   2,
	}
)

func InitRedis() error {
	redisConf := config.Conf.Redis
	redisOptions := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password,
		DB:       0,
	}
	RedisClient = redis.NewClient(&redisOptions)
	_, err := RedisClient.Ping().Result()
	if err != nil {
		return err
	}
	//RedisClient.LPush(config.Conf.Redis.MaTeachNerfTaskQueueKey, nil)
	//RedisClient.HMSet(config.Conf.Redis.MaTeachProgressKey, nil)
	return nil
}

type ProgressStruct struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Total     int64       `json:"total"`
	Processed int64       `json:"processed"`
	Ratio     float64     `json:"ratio"`
	Result    string      `json:"result"`
	Params    interface{} `json:"params"`
}

// LockWithTimeout 获取分布式锁
func LockWithTimeout(lockKey string, timeout time.Duration, lockTime time.Duration) (string, error) {
	uuid := xid.New().String()
	end := time.Now().Add(timeout).Unix()
	for time.Now().Unix() <= end {
		if result := RedisClient.SetNX(lockKey, uuid, lockTime).Val(); result {
			return uuid, nil
		} else {
			return "", LockError
		}
	}
	return "", LockTimeOut
}

// UnLock 释放分布式锁
func UnLock(lockKey string, uuid string) error {
	script := `if redis.call('get', KEYS[1]) == ARGV[1] then
					redis.call("del", KEYS[1])
					return 1
				else 
					return -1
				end`
	result, err := RedisClient.Eval(script, []string{lockKey}, uuid).Int()
	if err != nil || result == -1 {
		logger.Error("UnLock 释放锁失败。lockKey:[%#v] uuid:[%#v]", lockKey, uuid)
		return UnLockError
	}
	return nil
}

func ReadProgressFromRedis(redisKey, key string) (*ProgressStruct, error) {
	//读取redis
	matchKey := key
	var process ProgressStruct
	matchProgress := RedisClient.HMGet(redisKey, matchKey).Val()
	if matchProgress != nil {
		if len(matchProgress) != 0 && matchProgress[0] != nil {
			err := json.Unmarshal([]byte(matchProgress[0].(string)), &process)
			if err != nil {
				logger.Error("查询进度失败", err)
				return nil, err
			}
			return &process, nil
		} else {
			logger.Debug("redisProgress = []")
			return nil, nil
		}
	} else {
		logger.Debug("redisProgress is nil")
		return nil, nil
	}
}

func RedisPush(key string, content interface{}, method string) error {
	var err error
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if method == "L" {
		err = RedisClient.LPush(key, jsonContent).Err()
		return err
	}

	if method == "R" {
		err = RedisClient.RPush(key, jsonContent).Err()
		return err
	}
	return errors.New("unsupported method")
}

func InitRedisProgress(total int64, redisKey string, contentKey string, params interface{}) error {
	var err error
	progressMap := make(map[string]interface{})
	progressContent := ProgressStruct{
		Status:    ProcessStart,
		Message:   "开始计算",
		Total:     total,
		Processed: 0,
		Params:    params,
	}
	progressContentJson, err := json.Marshal(progressContent)
	if err != nil {
		return err
	}
	progressMap[contentKey] = progressContentJson
	if err = RedisClient.HMSet(redisKey, progressMap).Err(); err != nil {
		logger.Error(err)
	}
	return nil
}

func ResetRedisProgress(redisKey string, contentKey string) error {
	var err error
	progressMap := make(map[string]interface{})
	progressContent := ProgressStruct{
		Status:    ProcessStart,
		Message:   "",
		Total:     100,
		Processed: 0,
	}
	progressContentJson, err := json.Marshal(progressContent)
	if err != nil {
		return err
	}
	progressMap[contentKey] = progressContentJson
	if err = RedisClient.HMSet(redisKey, progressMap).Err(); err != nil {
		logger.Error(err)
	}
	return nil
}

// UpdateRedisProgress 更新Redis中的进度。保留param
func UpdateRedisProgress(redisKey string, contentKey string, processInfo *ProgressStruct) error {
	var err error
	progressMapString := make(map[string]string)
	progressMap := make(map[string]interface{})
	if err = RedisClient.HGetAll(redisKey).Err(); err != nil {
		return err
	}
	progressMapString = RedisClient.HGetAll(redisKey).Val()
	for k, v := range progressMapString {
		progressMap[k] = interface{}(v)
	}

	progressContent := ProgressStruct{}
	// progress := utils.RedisClient.HMGet(redisKey, contentKey).Val()[0]
	if v, ok := progressMap[contentKey]; ok {
		err = json.Unmarshal([]byte(v.(string)), &progressContent)
		if err != nil {
			logger.Error(progressMap[contentKey])
			logger.Error(err)
			return err
		}
		params := progressContent.Params
		err = copier.Copy(&progressContent, processInfo)
		if err != nil {
			logger.Error(err)
			return err
		}
		progressContent.Params = params
		progressContentJson, err := json.Marshal(progressContent)
		if err != nil {
			logger.Error(err)
			return err
		}
		progressMap[contentKey] = string(progressContentJson)
		if err = RedisClient.HMSet(redisKey, progressMap).Err(); err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}
