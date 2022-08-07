package model

type Item struct {
	ItemId uint32 `bson:"itemId"` // 道具id
	Count  uint32 `bson:"count"`  // 道具数量
	Guid   uint64 `bson:"-"`
}
