package game

import (
	"flswld.com/logger"
	"time"
)

type TickManager struct {
	gameManager *GameManager
}

func NewTickManager(gameManager *GameManager) (r *TickManager) {
	r = new(TickManager)
	r.gameManager = gameManager
	return r
}

func (t *TickManager) Start() {
	logger.LOG.Info("start game server tick")
	go func() {
		var ticker *time.Ticker = nil
		for {
			ticker = time.NewTicker(time.Minute * 1)
			<-ticker.C
			ticker.Stop()
			logger.LOG.Info("tick run")
			for _, world := range t.gameManager.worldManager.worldMap {
				for _, player := range world.playerMap {
					// TODO 线程不安全
					logger.LOG.Info("add player item, uid: %v", player.PlayerID)
					t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 201, ChangeCount: 10}}, true)
					t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 202, ChangeCount: 10}}, true)
					t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 223, ChangeCount: 10}}, true)
					t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 224, ChangeCount: 10}}, true)
					t.gameManager.AddUserItem(player.PlayerID, []*UserItem{{ItemId: 104003, ChangeCount: 1}}, true)
				}
			}
		}
	}()
}
