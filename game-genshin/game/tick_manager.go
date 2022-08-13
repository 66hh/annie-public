package game

import (
	"flswld.com/logger"
	"time"
)

type TickManager struct {
	ticker      *time.Ticker
	tickCount   uint64
	gameManager *GameManager
}

func NewTickManager(gameManager *GameManager) (r *TickManager) {
	r = new(TickManager)
	r.ticker = time.NewTicker(time.Second * 1)
	r.gameManager = gameManager
	return r
}

func (t *TickManager) OnGameServerTick() {
	t.ticker.Stop()
	t.ticker = time.NewTicker(time.Second * 1)
	t.tickCount++
	now := time.Now().Unix()
	if t.tickCount%1 == 0 {
		t.onTickSecond(now)
	}
	if t.tickCount%60 == 0 {
		t.onTickMinute(now)
	}
	if t.tickCount%3600 == 0 {
		t.onTickHour(now)
	}
	if t.tickCount%3600*24 == 0 {
		t.onTickDay(now)
	}
	if t.tickCount%3600*24*7 == 0 {
		t.onTickWeek(now)
	}
}

func (t *TickManager) onTickWeek(now int64) {
	logger.LOG.Info("on tick week, time: %v", now)
}

func (t *TickManager) onTickDay(now int64) {
	logger.LOG.Info("on tick day, time: %v", now)
}

func (t *TickManager) onTickHour(now int64) {
	logger.LOG.Info("on tick hour, time: %v", now)
}

func (t *TickManager) onTickMinute(now int64) {
	logger.LOG.Info("on tick minute, time: %v", now)
	// 测试每分钟给在线玩家发放道具
	for _, world := range t.gameManager.worldManager.worldMap {
		for _, player := range world.playerMap {
			logger.LOG.Debug("add player item, uid: %v", player.PlayerID)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 201, ChangeCount: 10}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 202, ChangeCount: 10}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 223, ChangeCount: 1000}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 224, ChangeCount: 1000}}, true)
			t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 104003, ChangeCount: 1}}, true)
		}
	}
}

func (t *TickManager) onTickSecond(now int64) {
	//logger.LOG.Debug("on tick second, time: %v", now)
}
