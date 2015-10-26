package repositories

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	THREAD_COLLECTION string = "threats-entry"
)

type Thread struct {
	ID        bson.ObjectId   `json:"id,omitempty" bson:"_id"`
	TargetId  bson.ObjectId   `json:"target_id,omitempty" bson:"target_id"`
	Entry     string          `json:"entry,omitempty" bson:"entry"`
	History   []ThreadHistory `json:"history,omitempty" bson:"entry"`
	CreatedAt string          `json:"created_at,omitempty" bson:"created_at"` // rfc3339 date
	UpdatedAt string          `json:"updated_at,omitempty" bson:"updated_at"` // rfc3339 date
}

type ThreadHistory struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id"`
	TargetId  bson.ObjectId `json:"target_id,omitempty" bson:"target_id"`
	Entry     string        `json:"entry,omitempty" bson:"entry"`
	UpdatedAt string        `json:"updated_at,omitempty" bson:"updated_at"` // rfc3339 date
}

type ThreadInterface interface {
	Create(quest Thread) (bson.ObjectId, error)
	Fetch(limit ThreadFetchConfig) ([]Thread, error)
	FetchTargetID(id bson.ObjectId, conf ThreadFetchConfig) (Thread, error)
	FetchID(id bson.ObjectId) (Thread, error)
	Update(id bson.ObjectId, quest Thread) error
	Delete(id bson.ObjectId) error
}

type ThreadFetchConfig struct {
	Limit  int
	Offset int
	Search string
}

type ThreadRepository struct {
	MongoDb *mgo.Database
}

func (tr ThreadRepository) Create(d Thread) (bson.ObjectId, error) {
	// validation
	mc := tr.MongoDb.C(THREAD_COLLECTION)
	// insert mongo id
	d.ID = bson.NewObjectIdWithTime(time.Now())
	mc.Insert(d)

	return d.ID, nil
}

func (tr ThreadRepository) Fetch(conf ThreadFetchConfig) ([]Thread, error) {
	mc := tr.MongoDb.C(THREAD_COLLECTION)

	ts := []Thread{}

	err := mc.Find(bson.M{}).Limit(conf.Limit).Skip(conf.Offset).All(ts)
	if err != nil {
		return ts, err
	}

	return ts, nil
}

func (tr ThreadRepository) FetchTargetID(tid bson.ObjectId, conf ThreadFetchConfig) ([]Thread, error) {
	mc := tr.MongoDb.C(THREAD_COLLECTION)

	ts := []Thread{}

	err := mc.Find(bson.M{"target_id": tid}).Limit(conf.Limit).Skip(conf.Offset).All(ts)
	if err != nil {
		return ts, err
	}

	return ts, nil
}

func (tr ThreadRepository) FetchID(id bson.ObjectId) (Thread, error) {
	mc := tr.MongoDb.C(THREAD_COLLECTION)

	d := Thread{}
	err := mc.FindId(id).One(d)

	if err != nil {
		return d, err
	}

	return d, nil
}

func (tr ThreadRepository) Update(id bson.ObjectId, t Thread) error {
	mc := tr.MongoDb.C(THREAD_COLLECTION)

	err := mc.UpdateId(id, t)
	if err != nil {
		return err
	}

	return nil
}

func (tr ThreadRepository) Delete(id bson.ObjectId) error {
	mc := tr.MongoDb.C(THREAD_COLLECTION)

	err := mc.RemoveId(id)
	if err != nil {
		return err
	}

	return nil
}
