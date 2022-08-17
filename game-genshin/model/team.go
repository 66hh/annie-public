package model

import (
	gdc "game-genshin/config"
	"game-genshin/constant"
)

type Team struct {
	Name         string   `bson:"name"`
	AvatarIdList []uint32 `bson:"avatarIdList"`
}

type AvatarEntity struct {
	AvatarEntityId uint32 `bson:"-"`
	WeaponEntityId uint32 `bson:"-"`
}

type TeamInfo struct {
	TeamList             []*Team         `bson:"teamList"`
	CurrTeamIndex        uint8           `bson:"currTeamIndex"`
	CurrAvatarIndex      uint8           `bson:"currAvatarIndex"`
	AvatarEntityList     []*AvatarEntity `bson:"-"`
	TeamEntityId         uint32          `bson:"-"`
	TeamResonances       map[uint16]bool `bson:"-"`
	TeamResonancesConfig map[int32]bool  `bson:"-"`
}

func NewTeamInfo() (r *TeamInfo) {
	r = &TeamInfo{
		TeamList: []*Team{
			{Name: "冒险", AvatarIdList: make([]uint32, 4)},
			{Name: "委托", AvatarIdList: make([]uint32, 4)},
			{Name: "秘境", AvatarIdList: make([]uint32, 4)},
			{Name: "深渊", AvatarIdList: make([]uint32, 4)},
		},
		CurrTeamIndex:    0,
		CurrAvatarIndex:  0,
		AvatarEntityList: make([]*AvatarEntity, 4),
		TeamEntityId:     0,
	}
	return r
}

func (t *TeamInfo) UpdateTeam(funcGetNextWorldEntityId func(uint16) uint32, funcCreateEntity func(uint16, map[uint32]float32, *Player) uint32, player *Player) {
	activeTeam := t.GetActiveTeam()
	// AvatarEntity
	entityIdTypeConst := constant.GetEntityIdTypeConst()
	t.AvatarEntityList = make([]*AvatarEntity, 4)
	for avatarIndex, avatarId := range activeTeam.AvatarIdList {
		if avatarId == 0 {
			break
		}
		avatarEntity := &AvatarEntity{
			AvatarEntityId: funcCreateEntity(entityIdTypeConst.AVATAR, player.AvatarMap[avatarId].FightPropMap, player),
			WeaponEntityId: funcGetNextWorldEntityId(entityIdTypeConst.WEAPON),
		}
		t.AvatarEntityList[avatarIndex] = avatarEntity
	}
	// 队伍元素共鸣
	elementTypeConst := constant.GetElementTypeConst()
	t.TeamResonances = make(map[uint16]bool)
	t.TeamResonancesConfig = make(map[int32]bool)
	teamElementTypeCountMap := make(map[uint16]uint8)
	avatarSkillDepotDataMapConfig := gdc.CONF.AvatarSkillDepotDataMap
	for _, avatarId := range activeTeam.AvatarIdList {
		skillData := avatarSkillDepotDataMapConfig[int32(avatarId)]
		if skillData != nil {
			teamElementTypeCountMap[skillData.ElementType.Value] += 1
		}
	}
	for k, v := range teamElementTypeCountMap {
		if v >= 2 {
			element := elementTypeConst.VALUE_MAP[k]
			if element.TeamResonanceId != 0 {
				t.TeamResonances[element.TeamResonanceId] = true
				t.TeamResonancesConfig[element.ConfigHash] = true
			}
		}
	}
	if len(t.TeamResonances) == 0 {
		t.TeamResonances[elementTypeConst.Default.TeamResonanceId] = true
		t.TeamResonancesConfig[int32(elementTypeConst.Default.TeamResonanceId)] = true
	}
}

func (t *TeamInfo) GetActiveTeamId() uint8 {
	return t.CurrTeamIndex + 1
}

func (t *TeamInfo) GetTeamByIndex(teamIndex uint8) *Team {
	if t.TeamList == nil {
		return nil
	}
	if teamIndex >= uint8(len(t.TeamList)) {
		return nil
	}
	activeTeam := t.TeamList[teamIndex]
	return activeTeam
}

func (t *TeamInfo) GetActiveTeam() *Team {
	return t.GetTeamByIndex(t.CurrTeamIndex)
}

func (t *TeamInfo) ClearTeamAvatar(teamIndex uint8) {
	team := t.GetTeamByIndex(teamIndex)
	if team == nil {
		return
	}
	team.AvatarIdList = make([]uint32, 4)
}

func (t *TeamInfo) AddAvatarToTeam(avatarId uint32, teamIndex uint8) {
	team := t.GetTeamByIndex(teamIndex)
	if team == nil {
		return
	}
	for i, v := range team.AvatarIdList {
		if v == 0 {
			team.AvatarIdList[i] = avatarId
			break
		}
	}
}

func (t *TeamInfo) GetActiveAvatarId() uint32 {
	activeTeam := t.GetActiveTeam()
	if activeTeam == nil {
		return 0
	}
	if t.CurrAvatarIndex >= uint8(len(activeTeam.AvatarIdList)) {
		return 0
	}
	return activeTeam.AvatarIdList[t.CurrAvatarIndex]
}

func (t *TeamInfo) GetAvatarEntityByIndex(index uint8) *AvatarEntity {
	if t.AvatarEntityList == nil {
		return nil
	}
	if index >= uint8(len(t.AvatarEntityList)) {
		return nil
	}
	return t.AvatarEntityList[index]
}

func (t *TeamInfo) GetActiveAvatarEntity() *AvatarEntity {
	return t.GetAvatarEntityByIndex(t.CurrAvatarIndex)
}
