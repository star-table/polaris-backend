package service

//func GetMQTTChannelKey(orgId int64, userId int64, input vo.GetMQTTChannelKeyReq) (*vo.GetMQTTChannelKeyResp, errs.SystemErrorInfo) {
//	mqttConfig := config.GetMQTTConfig()
//	if mqttConfig == nil {
//		log.Error("[GetMQTTChannelKey] mqtt config is nil")
//		return nil, errs.MQTTMissingConfigError
//	}
//
//	channelType := input.ChannelType
//	channel := ""
//
//	switch channelType {
//	case consts.MQTTChannelTypeProject:
//		if input.ProjectID == nil {
//			log.Errorf("[GetMQTTChannelKey] ChannelType: %d projectId is nil, orgId: %d, userId: %d, input: %v",
//				channelType, orgId, userId, input)
//			return nil, errs.ProjectNotExist
//		}
//		channel = util.GetMQTTProjectChannel(orgId, *input.ProjectID)
//	case consts.MQTTChannelTypeOrg:
//		channel = util.GetMQTTOrgRootChannel(orgId)
//	case consts.MQTTChannelTypeUser:
//		channel = util.GetMQTTUserChannel(orgId, userId)
//	case consts.MQTTChannelTypeApp:
//		if input.AppID == nil {
//			log.Errorf("[GetMQTTChannelKey] ChannelType: %d appId is nil, orgId: %d, userId: %d, input: %v",
//				channelType, orgId, userId, input)
//			return nil, errs.AppNotExist
//		}
//		channel = util.GetMQTTAppChannel(orgId, *input.AppID)
//	default:
//		return nil, errs.IllegalityMQTTChannelType
//	}
//
//	// check cache
//	cacheKey := fmt.Sprintf("%v_%v_%v", channel, consts.MQTTDefaultPermissions, consts.MQTTDefaultKeyTTLOneDay)
//	key, err := cache.Get(cacheKey)
//	if key == "" || err != nil {
//		key, err = emt.GenerateKey(channel, consts.MQTTDefaultPermissions, consts.MQTTDefaultKeyTTLOneDay)
//		if err != nil {
//			log.Errorf("[GetMQTTChannelKey] GenerateKey channel :%v, permissions: %v, ttl: %v, failed: %v",
//				channel, consts.MQTTDefaultPermissions, consts.MQTTDefaultKeyTTLOneDay, err)
//			return nil, errs.MQTTKeyGenError
//		}
//		log.Infof("[GetMQTTChannelKey] cache miss %v, gen new key: %v channel :%v, permissions: %v, ttl: %v",
//			cacheKey, key, channel, consts.MQTTDefaultPermissions, consts.MQTTDefaultKeyTTLOneDay)
//		// set cache
//		cache.SetEx(cacheKey, key, consts.MQTTDefaultKeyTTLOneDay/2)
//	} else {
//		log.Infof("[GetMQTTChannelKey] cache got: %v:%v, channel :%v, permissions: %v, ttl: %v",
//			cacheKey, key, channel, consts.MQTTDefaultPermissions, consts.MQTTDefaultKeyTTLOneDay)
//	}
//
//	port := (*int)(nil)
//	if mqttConfig.Port > 0 {
//		port = &mqttConfig.Port
//	}
//	return &vo.GetMQTTChannelKeyResp{
//		Key:     key,
//		Host:    mqttConfig.Host,
//		Port:    port,
//		Channel: channel,
//		Address: mqttConfig.Address,
//	}, nil
//}

//func MqttPublish(req *msgvo.MqttPublishReqVo) errs.SystemErrorInfo {
//	log.Infof("[MqttPublish] channel: %v, data: %q", req.Data.Channel, req.Data.Body)
//	err := mqtt.Publish(req.Data.Channel, util.GzipEncoding(req.Data.Body))
//	if err != nil {
//		log.Errorf("[MqttPublish] failed: %v, channel: %v, msg: %q",
//			err, req.Data.Channel, req.Data.Body)
//		return err
//	}
//	return nil
//}
