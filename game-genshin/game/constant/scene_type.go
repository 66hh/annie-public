package constant

type SceneType struct {
	SCENE_NONE       uint16
	SCENE_WORLD      uint16
	SCENE_DUNGEON    uint16
	SCENE_ROOM       uint16
	SCENE_HOME_WORLD uint16
	SCENE_HOME_ROOM  uint16
	SCENE_ACTIVITY   uint16
}

func GetSceneTypeConst() (r *SceneType) {
	r = new(SceneType)
	r.SCENE_NONE = 0
	r.SCENE_WORLD = 1
	r.SCENE_DUNGEON = 2
	r.SCENE_ROOM = 3
	r.SCENE_HOME_WORLD = 4
	r.SCENE_HOME_ROOM = 5
	r.SCENE_ACTIVITY = 6
	return r
}
