package serializer

import (
	"encoding/json"

	"github.com/Malijoe/djanGo-unchained/models"
	"github.com/hashicorp/go-multierror"
)

type Serializer struct {
	Meta
	models.Model
}

func NewModelSerializer(model models.Model, meta Meta) *Serializer {
	initMeta(&meta)
	return &Serializer{
		Model: model,
		Meta:  meta,
	}
}

func (s *Serializer) PerformModifications(t ModFieldTime) error {
	var errs error
	if modifiers, ok := s.Modifiers[t]; ok {
		for _, m := range modifiers {
			if err := ManageModelFields(s.Model, m.Fields, m.Modifier); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs
}

func (s *Serializer) MarshalJSON() ([]byte, error) {
	// make sure the Model is initialized
	s.Init()
	// perform any pre-write modifications
	if err := s.PerformModifications(PreRead); err != nil {
		return nil, err
	}

	// marshal json data
	data, err := json.Marshal(s.Model)
	if err != nil {
		return nil, err
	}

	// perform any post-write modifications
	if err = s.PerformModifications(PostRead); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Serializer) UnmarshalJSON(data []byte) error {
	// make sure the Model is initialized
	s.Init()
	if err := s.PerformModifications(PreWrite); err != nil {
		return err
	}

	if err := json.Unmarshal(data, s.Model); err != nil {
		return err
	}

	return s.PerformModifications(PostWrite)
}
