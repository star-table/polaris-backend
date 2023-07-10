package domain

import (
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/dao"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/po"
	pkgErr "github.com/pkg/errors"
	"upper.io/db.v3"
)

func InsertMqFailMsgToDB(msg model.MqMessage, msgType int, orgId int64) (int64, errs.SystemErrorInfo) {
	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableMessage)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	msgJson := json.ToJsonIgnoreError(msg)

	msgPo := po.PpmTakMessage{
		Id:        id,
		OrgId:     orgId,
		Topic:     msg.Topic,
		Type:      msgType, //待修改
		ProjectId: 0,
		IssueId:   0,
		TrendsId:  0,
		Info:      msgJson,
		Content:   msg.Body,
		Status:    consts.MessageStatusWait,
		IsDelete:  consts.AppIsNoDelete,
	}

	dbErr := dao.InsertMessage(msgPo)
	if dbErr != nil {
		log.Errorf("消息入表失败，消息内容%s, 错误信息%v", msgJson, dbErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return id, nil
}

func GetMessageBoList(page uint, size uint, cond db.Cond) (*[]bo.MessageBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectMessageByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "id desc",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.MessageBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetMessageBo(id int64) (*bo.MessageBo, errs.SystemErrorInfo) {
	po, err := dao.SelectMessageById(id)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.MessageBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func UpdateMsgStatus(id int64, newStatus int) errs.SystemErrorInfo {
	_, err := dao.UpdateMessageById(id, mysql.Upd{
		consts.TcStatus: newStatus,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

// InsertFailMsgToDB 将处理错误的消息放入数据库中，以便后续的排查和处理。
func InsertFailMsgToDB(msgBody string, msgType int, orgId int64) (int64, error) {
	id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableMessage)
	if err != nil {
		// 使用了 err wrap 就不用再记录，只需在外层 api 调用的地方取出错误栈即可。
		// log.Error(err)
		return 0, pkgErr.Wrap(err, "idfacade error")
	}

	msgPo := po.PpmTakMessage{
		Id:        id,
		OrgId:     orgId,
		Topic:     "write_some_failed_msg",
		Type:      msgType, //待修改
		ProjectId: 0,
		IssueId:   0,
		TrendsId:  0,
		Info:      "",
		Content:   msgBody,
		Status:    consts.MessageStatusFail,
		IsDelete:  consts.AppIsNoDelete,
	}

	dbErr := dao.InsertMessage(msgPo)
	if dbErr != nil {
		// 使用了 err wrap 就不用再记录，只需在外层 api 调用的地方取出错误栈即可。
		// log.Errorf("[InsertFailMsgToDB] 消息入表失败，错误信息%v", dbErr)
		return 0, pkgErr.Wrap(dbErr, "insert error")
	}
	return id, nil
}
