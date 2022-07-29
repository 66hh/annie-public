package game

import (
	"flswld.com/logger"
	"game-genshin/dao"
	"game-genshin/entity"
	"sync"
	"time"
)

type UserManager struct {
	log           *logger.Logger
	dao           *dao.Dao
	playerMap     map[uint32]*entity.Player
	playerMapLock sync.RWMutex
}

func NewUserManager(log *logger.Logger, dao *dao.Dao) (r *UserManager) {
	r = new(UserManager)
	r.log = log
	r.dao = dao
	r.playerMap = make(map[uint32]*entity.Player)
	return r
}

func (u *UserManager) GetTargetUser(userId uint32) *entity.Player {
	u.playerMapLock.RLock()
	player, exist := u.playerMap[userId]
	u.playerMapLock.RUnlock()
	if exist {
		return player
	} else {
		return u.LoadUserFromDb(userId)
	}
}

func (u *UserManager) LoadUserFromDb(userId uint32) *entity.Player {
	player, err := u.dao.QueryPlayerByID(userId)
	if err != nil {
		u.log.Error("query player error: %v", err)
		return nil
	}
	player.DbState = entity.DbNormal
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
	return player
}

func (u *UserManager) AddUser(player *entity.Player) {
	if player == nil {
		return
	}
	u.ChangeUserDbState(player, entity.DbInsert)
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
}

func (u *UserManager) DeleteUser(player *entity.Player) {
	u.ChangeUserDbState(player, entity.DbDelete)
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
}

func (u *UserManager) UpdateUser(player *entity.Player) {
	if player == nil {
		return
	}
	u.ChangeUserDbState(player, entity.DbUpdate)
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
}

func (u *UserManager) ChangeUserDbState(player *entity.Player, state int) {
	if player == nil {
		return
	}
	switch player.DbState {
	case entity.DbInsert:
		if state == entity.DbDelete {
			player.DbState = entity.DbDelete
		}
	case entity.DbDelete:
	case entity.DbUpdate:
		if state == entity.DbDelete {
			player.DbState = entity.DbDelete
		}
	case entity.DbNormal:
		if state == entity.DbDelete {
			player.DbState = entity.DbDelete
		} else if state == entity.DbUpdate {
			player.DbState = entity.DbUpdate
		}
	}
}

func (u *UserManager) StartAutoSaveUser() {
	go func() {
		var ticker *time.Ticker = nil
		for {
			u.log.Info("auto save user start")
			playerMapTemp := make(map[uint32]*entity.Player)
			u.playerMapLock.RLock()
			for k, v := range u.playerMap {
				playerMapTemp[k] = v
			}
			u.playerMapLock.RUnlock()
			u.log.Info("copy user map finish")
			insertList := make([]*entity.Player, 0)
			deleteList := make([]uint32, 0)
			updateList := make([]*entity.Player, 0)
			for k, v := range playerMapTemp {
				switch v.DbState {
				case entity.DbInsert:
					insertList = append(insertList, v)
					playerMapTemp[k].DbState = entity.DbNormal
				case entity.DbDelete:
					deleteList = append(deleteList, v.PlayerID)
					delete(playerMapTemp, k)
				case entity.DbUpdate:
					updateList = append(updateList, v)
					playerMapTemp[k].DbState = entity.DbNormal
				case entity.DbNormal:
					continue
				}
			}
			u.log.Info("db state init finish")
			err := u.dao.InsertPlayerList(insertList)
			if err != nil {
				u.log.Error("insert player list error: %v", err)
			}
			err = u.dao.DeletePlayerList(deleteList)
			if err != nil {
				u.log.Error("delete player error: %v", err)
			}
			err = u.dao.UpdatePlayerList(updateList)
			if err != nil {
				u.log.Error("update player error: %v", err)
			}
			u.log.Info("db write finish")
			u.playerMapLock.Lock()
			u.playerMap = playerMapTemp
			u.playerMapLock.Unlock()
			u.log.Info("auto save user finish")
			ticker = time.NewTicker(time.Minute * 10)
			<-ticker.C
			ticker.Stop()
		}
	}()
}
