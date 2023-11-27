package reminder

import "go-reminder-bot/pkg/util"

type GetListParams struct {
	util.GetListParams
	CreatedBy *string
}
