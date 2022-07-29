package dao

import (
	"context"
	"flswld.com/common/config"
	"flswld.com/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dao struct {
	conf   *config.Config
	log    *logger.Logger
	client *mongo.Client
	db     *mongo.Database
}

func NewDao(conf *config.Config, log *logger.Logger) (r *Dao) {
	r = new(Dao)
	r.conf = conf
	r.log = log
	clientOptions := options.Client().ApplyURI(conf.Database.Url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		r.log.Error("mongo connect error: %v", err)
		return nil
	}
	r.client = client
	r.db = client.Database("gate_genshin")
	return r
}

func (d *Dao) CloseDao() {
	err := d.client.Disconnect(context.TODO())
	if err != nil {
		d.log.Error("mongo close error: %v", err)
	}
}
