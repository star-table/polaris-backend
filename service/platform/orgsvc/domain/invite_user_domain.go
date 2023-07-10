package domain

import (
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/polaris-backend/common/core/util/rand"
	"strconv"
)

func GenInviteCode(currentUserId int64, sourcePlatform string) string {
	return rand.RandomInviteCode(uuid.NewUuid() + strconv.FormatInt(currentUserId, 10) + sourcePlatform)
}
