package internal

import (
	"changeme/internal/define"
)

func (a *App) GetHomeInfo() define.M {

	a.Info.Video = a.Vm.TotalInfo

	return define.M{
		"code": 200,
		"data": a.Info,
	}
}
