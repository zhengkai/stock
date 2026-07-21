package app

import (
	"fmt"
	"project/pb"
	"project/util"
	"strings"
)

type priceCheckMeta struct {
	code      string
	name      string
	price     int64
	ts        uint32
	up        priceCheck
	down      priceCheck
	alertSkip int
}

type priceCheck struct {
	name      string
	alert     pb.PriceAlert
	file      *util.File
	isUp      bool
	threshold int64
	limit     int64
	limitName string // 涨停/跌停
	skip      bool
}

func (pm *priceCheckMeta) loadAlert(pc *priceCheck, isUp bool, threshold uint64) {

	direction := `up`
	pc.name = `上涨`
	if !isUp {
		direction = `down`
		pc.name = `下跌`
	}
	pc.isUp = isUp
	pc.threshold = int64(threshold)
	pc.file = util.NewFileF(`alert/price-%s-%s.json`, pm.code, direction)

	pm._loadAlert(pc)
	if pc.skip {
		pm.alertSkip++
	}
}

func (pm *priceCheckMeta) _loadAlert(pc *priceCheck) {

	if pc.threshold == 0 {
		pc.skip = true
		pm.msg(pc.name, `未设阈值`)
		return
	}

	pc.file.ReadJSON(&pc.alert)
	if pc.alert.GetExpire() > pm.ts {
		pc.skip = true
		pm.msg(pc.name, `等待`, pc.alert.GetExpire()-pm.ts)
		return
	}
}

func (pm *priceCheckMeta) do(pc *priceCheck, limit int64, limitName string) {

	if pc.skip {
		return
	}
	pc.limit = limit
	pc.limitName = limitName

	msg, sleep, ok := pm.check(pc)
	if ok {
		util.WebPushAll(`股票价格提醒`, msg)
	}

	if sleep == 0 {
		sleep = 600
	}
	if msg == `` {
		msg = pm.msg(pc.name, `nothing`)
	}

	d := &pc.alert
	d.SetMsg(msg)
	d.SetTs(pm.ts)
	d.SetExpire(pm.ts + sleep)

	pc.file.WriteJSON(d)

	last := util.NewFile(`last-alert.json`)
	last.WriteJSON(d)
}

func (pm *priceCheckMeta) check(pc *priceCheck) (msg string, sleep uint32, ok bool) {

	var di int64 = 1
	if !pc.isUp {
		di = -1
	}

	inLimit := (pc.limit - pc.threshold) * di
	if inLimit < 0 {
		msg = pm.msgF(`%s价 %s 未达到%s价 %s，今天放弃`, pc.limitName, util.Price(pc.limit), pc.name, util.Price(pc.threshold))
		sleep = 7200
		return
	}

	remain := (pc.threshold - pm.price) * di
	if remain <= 0 {
		sleep = 43200 // 半天
		msg = pm.msgF(`当前价格 %s 已达到%s价 %s `, util.Price(pm.price), pc.name, util.Price(pc.threshold))
		ok = true
		return
	}

	remainRate := float64(remain) / float64(pc.threshold) * 100
	if remain < 10 || remainRate < 2 {
		sleep = 50
		return
	}
	if remainRate < 3 {
		sleep = 100
		return
	}
	if remainRate < 5 {
		sleep = 300
		return
	}
	return
}

func (pm *priceCheckMeta) msg(msg ...any) string {
	args := []any{fmt.Sprintf(`%s(%s)`, pm.name, pm.code)}
	args = append(args, msg...)

	s := strings.TrimSuffix(fmt.Sprintln(args...), "\n")
	fmt.Println(s)
	return s
}

func (pm *priceCheckMeta) msgF(format string, a ...any) string {
	return pm.msg(fmt.Sprintf(format, a...))
}

func (a *App) checkAlert() {

	fmt.Println()
	fmt.Println(`alert list`, len(a.Config.GetAlert()))
	for _, al := range a.Config.GetAlert() {
		a.checkAlertOne(al)
		// if life.Sleep(sleep) != nil {
		// fmt.Println(`break`)
		// break
		// }
	}
}

func (a *App) checkAlertOne(al *pb.Alert) (msg string, sleep int, ok bool) {

	pc := &priceCheckMeta{
		code: al.GetCode(),
		name: al.GetName(),
		ts:   util.TS(),
	}
	pc.loadAlert(&pc.up, true, al.GetMax())
	pc.loadAlert(&pc.down, false, al.GetMin())
	if pc.alertSkip > 2 {
		return
	}

	q := Stock(al.GetCode())
	if q == nil || q.GetPrice() < 1 {
		return
	}
	pc.price = q.GetPrice()

	pc.do(&pc.up, q.GetLimitUp(), `涨停`)
	pc.do(&pc.down, q.GetLimitDown(), `跌停`)

	return
}
