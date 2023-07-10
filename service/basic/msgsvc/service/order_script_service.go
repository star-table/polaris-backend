package service

import (
	"strconv"
	"time"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/msgvo"
	"github.com/star-table/polaris-backend/common/model/vo/ordervo"
	"github.com/star-table/polaris-backend/facade/orderfacade"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/dao"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/domain"
	"upper.io/db.v3"
)

// FixAddOrderForFeiShu 修复消费失败的飞书方的订单回调
func FixAddOrderForFeiShu(input msgvo.FixAddOrderForFeiShuReqData) errs.SystemErrorInfo {
	cond := db.Cond{
		consts.TcTopic:    "topic_feishu_callback_order_gray",
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcVersion:  1,
	}
	if input.StartTime != nil && input.EndTime != nil {
		cond["create_time"] = db.Gte(input.StartTime)
		cond[" create_time"] = db.Lte(input.EndTime)
	}
	for i := 0; i < 5000; i++ {
		list, _, err := domain.GetMessageBoList(1, 10000, cond)
		if err != nil {
			if err == db.ErrNoMoreRows {
				break
			}
			log.Error(err)
			return err
		}
		if len(*list) < 1 {
			break
		}
		cond["id"] = db.Gt((*list)[0].Id)
		newIdArr := make([]int64, 0)
		for _, oneMsg := range *list {
			if oneMsg.Content == nil || len(*oneMsg.Content) < 2 {
				continue
			}
			contentObj := msgvo.OrderContent{}
			if err := json.FromJson(*oneMsg.Content, &contentObj); err != nil {
				log.Errorf("[FixAddOrderForFeiShu] from json err: %s", err)
				continue
			}
			contentDataObj := msgvo.OrderContentData{}
			if err := json.FromJson(contentObj.Data, &contentDataObj); err != nil {
				log.Errorf("[FixAddOrderForFeiShu] from json err: %s", err)
				continue
			}
			eventData := contentDataObj.Event
			orderFsBo := bo.OrderFsBo{
				OrderId:       eventData.OrderId,
				PricePlanId:   eventData.PricePlanId,
				PricePlanType: eventData.PricePlanType,
				Seats:         eventData.Seats,
				BuyCount:      eventData.BuyCount,
				Status:        consts.FsOrderStatusNormal,
				BuyType:       eventData.BuyType,
				SrcOrderId:    eventData.SrcOrderId,
				DstOrderId:    "",
				OrderPayPrice: eventData.OrderPayPrice,
				TenantKey:     eventData.TenantKey,
			}
			if eventData.PayTime != "" {
				payTime, err := strconv.ParseInt(eventData.PayTime, 10, 64)
				if err != nil {
					log.Error(err)
				}
				orderFsBo.PaidTime = time.Unix(payTime, 0)
			}
			if eventData.CreateTime != "" {
				createTime, err := strconv.ParseInt(eventData.CreateTime, 10, 64)
				if err != nil {
					log.Error(err)
				}
				orderFsBo.CreateTime = time.Unix(createTime, 0)
			}
			resp := orderfacade.AddFsOrder(ordervo.AddFsOrderReq{
				Data: orderFsBo,
			})
			if resp.Failure() {
				log.Error(resp.Error())
				continue
			} else {
				// 收集新数据的 id，修改成功后，将 version 标识改为 2，表示处理过
				newId := resp.Void.ID
				if newId > 0 {
					newIdArr = append(newIdArr, newId)
				}
			}
		}
		// 修改 version
		if len(newIdArr) > 0 {
			if _, oriErr := dao.UpdateMessageByCond(db.Cond{
				consts.TcId: db.In(newIdArr),
			}, mysql.Upd{
				consts.TcVersion: 2,
			}); oriErr != nil {
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
			}
		}
	}

	return nil
}
