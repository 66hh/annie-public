package game

import (
	"flswld.com/logger"
	"game-genshin/dao"
	"game-genshin/model"
	"sync"
	"time"
)

type UserManager struct {
	log           *logger.Logger
	dao           *dao.Dao
	playerMap     map[uint32]*model.Player
	playerMapLock sync.RWMutex
}

func NewUserManager(log *logger.Logger, dao *dao.Dao) (r *UserManager) {
	r = new(UserManager)
	r.log = log
	r.dao = dao
	r.playerMap = make(map[uint32]*model.Player)
	return r
}

func (u *UserManager) GetTargetUser(userId uint32) *model.Player {
	u.playerMapLock.RLock()
	player, exist := u.playerMap[userId]
	u.playerMapLock.RUnlock()
	if exist {
		return player
	} else {
		return u.LoadUserFromDb(userId)
	}
}

func (u *UserManager) LoadUserFromDb(userId uint32) *model.Player {
	player, err := u.dao.QueryPlayerByID(userId)
	if err != nil {
		u.log.Error("query player error: %v", err)
		return nil
	}
	player.DbState = model.DbNormal
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
	return player
}

func (u *UserManager) AddUser(player *model.Player) {
	if player == nil {
		return
	}
	u.ChangeUserDbState(player, model.DbInsert)
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
}

func (u *UserManager) DeleteUser(player *model.Player) {
	u.ChangeUserDbState(player, model.DbDelete)
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
}

func (u *UserManager) UpdateUser(player *model.Player) {
	if player == nil {
		return
	}
	u.ChangeUserDbState(player, model.DbUpdate)
	u.playerMapLock.Lock()
	u.playerMap[player.PlayerID] = player
	u.playerMapLock.Unlock()
}

func (u *UserManager) ChangeUserDbState(player *model.Player, state int) {
	if player == nil {
		return
	}
	switch player.DbState {
	case model.DbInsert:
		if state == model.DbDelete {
			player.DbState = model.DbDelete
		}
	case model.DbDelete:
	case model.DbUpdate:
		if state == model.DbDelete {
			player.DbState = model.DbDelete
		}
	case model.DbNormal:
		if state == model.DbDelete {
			player.DbState = model.DbDelete
		} else if state == model.DbUpdate {
			player.DbState = model.DbUpdate
		}
	}
}

func (u *UserManager) StartAutoSaveUser() {
	go func() {
		var ticker *time.Ticker = nil
		for {
			u.log.Info("auto save user start")
			playerMapTemp := make(map[uint32]*model.Player)
			u.playerMapLock.RLock()
			for k, v := range u.playerMap {
				playerMapTemp[k] = v
			}
			u.playerMapLock.RUnlock()
			u.log.Info("copy user map finish")
			insertList := make([]*model.Player, 0)
			deleteList := make([]uint32, 0)
			updateList := make([]*model.Player, 0)
			for k, v := range playerMapTemp {
				switch v.DbState {
				case model.DbInsert:
					insertList = append(insertList, v)
					playerMapTemp[k].DbState = model.DbNormal
				case model.DbDelete:
					deleteList = append(deleteList, v.PlayerID)
					delete(playerMapTemp, k)
				case model.DbUpdate:
					updateList = append(updateList, v)
					playerMapTemp[k].DbState = model.DbNormal
				case model.DbNormal:
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
			ticker = time.NewTicker(time.Minute * 1)
			<-ticker.C
			ticker.Stop()
		}
	}()
}
