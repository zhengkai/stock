// Package app 控制所有后台活动
package app

import "project/pb"

var (
	theApp = &App{
		pool: &pool{
			m: make(map[string]*pb.AppStock, 32),
		},
	}

	GetPB = theApp.pool.GetPB
)

type App struct {
	Config *pb.Config
	pool   *pool
}
