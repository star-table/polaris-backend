package mqtt

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/mqtt/emt"
	emitter "github.com/star-table/emitter-go-client"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
)

var log = logger.GetDefaultLogger()

func onError(c *emitter.Client, e emitter.Error) {
	log.Error(e.Message)
	if e.Status == 401 {
		log.Infof("mqtt root key失效 %s, 准备清理", e.Message)
		cacheErr := ClearRootKey()
		if cacheErr != nil {
			log.Error(cacheErr)
		}
	}
}

func Publish(channel string, payload interface{}) errs.SystemErrorInfo {
	//log.Infof("MQTT推送的channel %s", channel)

	key, err := GetRootKey()
	if err != nil {
		log.Error(err)
		return err
	}
	mqttErr := emt.Publish(key, channel, payload, 0, onError)
	if mqttErr != nil {
		log.Error(mqttErr)
		return errs.MQTTPublishError
	}
	return nil
}

func GetRootKey() (string, errs.SystemErrorInfo) {
	key, err := cache.Get(consts.CacheMQTTRootKey)
	if err != nil {
		log.Error(err)
		return "", errs.RedisOperateError
	}
	if key == "" {
		return GetRootNewKey()
	}
	return key, nil
}

func GetRootNewKey() (string, errs.SystemErrorInfo) {
	newKey, mqttErr := emt.GenerateKey(consts.MQTTChannelRoot, consts.MQTTDefaultRootPermissions, consts.MQTTDefaultTTL)
	if mqttErr != nil {
		log.Error(mqttErr)
		return "", errs.MQTTKeyGenError
	}

	err := cache.Set(consts.CacheMQTTRootKey, newKey)
	if err != nil {
		log.Error(err)
		return "", errs.RedisOperateError
	}
	return newKey, nil
}

func ClearRootKey() errs.SystemErrorInfo {
	_, err := cache.Del(consts.CacheMQTTRootKey)
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}
	return nil
}
