package repositories

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	DUNGEON_COLLECTION string = "dungeons"
)

type Dungeon struct {
	ID            bson.ObjectId   `json:"id,omitempty" bson:"_id"`
	Name          string          `json:"name,omitempty" bson:"name"`
	Description   string          `json:"description,omitempty" bson:"description"`
	TechStack     []string        `json:"tech_stack,omitempty" bson:"tech_stack"`
	DungeonMaster []bson.ObjectId `json:"dungeon_master,omitempty" bson:"dungeon_master"`
	Guardians     []bson.ObjectId `json:"guardians,omitempty" bson:"guardians"`
	CreatedAt     string          `json:"created_at,omitempty" bson:"created_at"` // rfc3339 date
	UpdatedAt     string          `json:"updated_at,omitempty" bson:"updated_at"` // rfc3339 date
}

type DungeonInterface interface {
	Create(dungeon Dungeon) (bson.ObjectId, error)
	Fetch(limit DungeonFetchConfig) ([]Dungeon, error)
	FetchID(id bson.ObjectId) (Dungeon, error)
	Update(id bson.ObjectId, dungeon Dungeon) error
	Delete(id bson.ObjectId) error
}

type DungeonFetchConfig struct {
	Limit  int
	Offset int
	Search string
}

type DungeonRepository struct {
	MongoDb *mgo.Database
}

func (dr DungeonRepository) Create(d Dungeon) (bson.ObjectId, error) {
	// validation
	mc := dr.MongoDb.C(DUNGEON_COLLECTION)
	// insert mongo id
	d.ID = bson.NewObjectIdWithTime(time.Now())
	mc.Insert(d)

	return d.ID, nil
}

func (dr DungeonRepository) Fetch(conf DungeonFetchConfig) ([]Dungeon, error) {
	mc := dr.MongoDb.C(DUNGEON_COLLECTION)

	ds := []Dungeon{}

	err := mc.Find(bson.M{}).Limit(conf.Limit).Skip(conf.Offset).All(ds)
	if err != nil {
		return ds, err
	}

	return ds, nil
}

func (dr DungeonRepository) FetchID(id bson.ObjectId) (Dungeon, error) {
	mc := dr.MongoDb.C(DUNGEON_COLLECTION)

	d := Dungeon{}
	err := mc.FindId(id).One(d)

	if err != nil {
		return d, err
	}

	return d, nil
}

func (dr DungeonRepository) Update(id bson.ObjectId, d Dungeon) error {
	mc := dr.MongoDb.C(DUNGEON_COLLECTION)

	err := mc.UpdateId(id, d)
	if err != nil {
		return err
	}

	return nil
}

func (dr DungeonRepository) Delete(id bson.ObjectId) error {
	mc := dr.MongoDb.C(DUNGEON_COLLECTION)

	err := mc.RemoveId(id)
	if err != nil {
		return err
	}

	return nil
}
