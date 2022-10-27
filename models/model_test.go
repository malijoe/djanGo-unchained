package models

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/malijoe/djanGo-unchained/db"
)

var (
	cache  = make(map[uint]testType)
	lastId uint
	repo   = db.NewMockRepository[testType](
		db.WithGetFn(func(id uint) (*testType, error) {
			t, ok := cache[id]
			if !ok {
				return nil, fmt.Errorf("no item found with id %d", id)
			}
			return &t, nil
		}),
		db.WithSaveFn(func(t *testType) (*testType, error) {
			lastId++
			t.ID = lastId
			cache[t.ID] = *t
			return t, nil
		}),
		db.WithUpdateFn(func(t *testType) (*testType, error) {
			if t.ID == 0 {
				return nil, errors.New("missing id")
			}
			_, ok := cache[t.ID]
			if !ok {
				return nil, fmt.Errorf("no item found with id %d", t.ID)
			}
			cache[t.ID] = *t
			return t, nil
		}),
		db.WithDeleteFn[testType](func(id uint) error {
			_, ok := cache[id]
			if !ok {
				return fmt.Errorf("no item found with id %d", id)
			}
			delete(cache, id)
			return nil
		}),
	)
)

type testType struct {
	ID     uint
	Name   string
	Number int
}

func (t testType) Objects() DataModel[testType] {
	return &testModel{
		DataModel: NewDataModel(repo, &t),
	}
}

type testModel struct {
	DataModel[testType]
}

func TestDataModel(t *testing.T) {
	t1 := testType{
		Name:   "one",
		Number: 1,
	}
	dm := t1.Objects()
	if dm.Instance() == nil {
		t.Error("data model instance is nil")
		return
	}
	n, err := dm.Create(context.Background())
	if err != nil {
		t.Error("error saving data model", err)
		return
	}
	if n.ID != 1 {
		t.Errorf("unexpected id of first inserted model: %d", n.ID)
	}
	if !repo.SaveInvoked {
		t.Error("repo's save method was not called")
		return
	}

	t2 := testType{
		ID:     1,
		Name:   "two",
		Number: 1,
	}
	dm = t2.Objects()
	if dm.Instance() == nil {
		t.Error("data model instance is nil")
		return
	}
	n, err = dm.Update(context.Background())
	if err != nil {
		t.Error("error updating data model", err)
		return
	}
	if n.ID != 1 {
		t.Errorf("unexpected id of updated model: %d", n.ID)
	}
	if n.Name != "two" {
		t.Errorf("unexpected name value of updated model. expected: %s; got: %s", "two", n.Name)
	}
	if n.Number != 1 {
		t.Errorf("unexpected number value of update model. expected: %d; got %d", 1, n.Number)
	}
	if !repo.UpdateInvoked {
		t.Error("repo's update method was not called")
	}

	t3 := testType{}
	dm = t3.Objects()
	n, err = dm.Get(context.Background(), 1)
	if err != nil {
		t.Error("error getting data model", err)
		return
	}
	if !repo.GetInvoked {
		t.Error("repo's get method was not called")
	}
	if n.ID != 1 {
		t.Errorf("unexpected id value of retrieved model. expected: %d; got %d", 1, n.ID)
	}
	if n.Name != "two" {
		t.Errorf("unexpected name value of updated model. expected: %s; got: %s", "two", n.Name)
	}
	if n.Number != 1 {
		t.Errorf("unexpected number value of update model. expected: %d; got %d", 1, n.Number)
	}
	err = dm.Delete(context.Background(), 1)
	if err != nil {
		t.Error("error deleting data model", err)
	}
	if !repo.DeleteInvoked {
		t.Error("repo's delete method was not called")
	}
}
