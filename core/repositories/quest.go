package repositories

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	QUEST_COLLECTION string = "quests"
)

type Quest struct {
	ID          bson.ObjectId    `json:"id,omitempty" bson:"_id"`
	QuestID     bson.ObjectId    `json:"quest_id,omitempty" bson:"quest_id"`
	Description string           `json:"description,omitempty" bson:"description"`
	Checklist   []QuestChecklist `json:"checklist,omitempty" bson:"checklist"`
	Members     []bson.ObjectId  `json:"members,omitempty" bson:"members"`
	Class       []string         `json:"class,omitempty" bson:"class"`
	StartDate   string           `json:"start_date,omitempty" bson:"start_date"` // rfc3339 date
	DueDate     string           `json:"due_date,omitempty" bson:"due_date"`     // rfc3339 date
	CreatedAt   string           `json:"created_at,omitempty" bson:"created_at"` // rfc3339 date
	UpdatedAt   string           `json:"updated_at,omitempty" bson:"updated_at"` // rfc3339 date
}

type QuestChecklist struct {
	Description string `json:"description,omitempty" bson:"description"`
	Status      string `json:"status,omitempty" bson:"status"`
}

type QuestInterface interface {
	Create(quest Quest) (bson.ObjectId, error)
	Fetch(config QuestFetchConfig) ([]Quest, error)
	FetchQuestID(did bson.ObjectId, config QuestFetchConfig) ([]Quest, error)
	FetchID(id bson.ObjectId) (Quest, error)
	Update(id bson.ObjectId, quest Quest) error
	Delete(id bson.ObjectId) error
}

type QuestFetchConfig struct {
	Limit  int
	Offset int
	Search string
}

type QuestRepository struct {
	MongoDb *mgo.Database
}

func (qr QuestRepository) Create(q Quest) (bson.ObjectId, error) {
	// validation
	mc := qr.MongoDb.C(QUEST_COLLECTION)
	// insert mongo id
	q.ID = bson.NewObjectIdWithTime(time.Now())
	mc.Insert(q)

	return q.ID, nil
}

func (qr QuestRepository) Fetch(conf QuestFetchConfig) ([]Quest, error) {
	mc := qr.MongoDb.C(QUEST_COLLECTION)

	qs := []Quest{}

	err := mc.Find(bson.M{}).Limit(conf.Limit).Skip(conf.Offset).All(qs)
	if err != nil {
		return qs, err
	}

	return qs, nil
}

func (qr QuestRepository) FetchQuestID(did bson.ObjectId, conf QuestFetchConfig) ([]Quest, error) {
	mc := qr.MongoDb.C(QUEST_COLLECTION)

	qs := []Quest{}

	err := mc.Find(bson.M{"quest_id": did}).Limit(conf.Limit).Skip(conf.Offset).All(qs)
	if err != nil {
		return qs, err
	}

	return qs, nil
}

func (qr QuestRepository) FetchID(id bson.ObjectId) (Quest, error) {
	mc := qr.MongoDb.C(QUEST_COLLECTION)

	q := Quest{}
	err := mc.FindId(id).One(q)

	if err != nil {
		return q, err
	}

	return q, nil
}

func (qr QuestRepository) Update(id bson.ObjectId, q Quest) error {
	mc := qr.MongoDb.C(QUEST_COLLECTION)

	err := mc.UpdateId(id, q)
	if err != nil {
		return err
	}

	return nil
}

func (qr QuestRepository) Delete(id bson.ObjectId) error {
	mc := qr.MongoDb.C(QUEST_COLLECTION)

	err := mc.RemoveId(id)
	if err != nil {
		return err
	}

	return nil
}
