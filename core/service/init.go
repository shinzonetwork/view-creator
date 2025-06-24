package service

import (
	"strconv"
	"time"

	"github.com/shinzonetwork/view-creator/core/models"
	"github.com/shinzonetwork/view-creator/core/store"
)

func InitView(name string, s store.ViewStore) (models.View, error) {
	return s.Create(name, strconv.FormatInt(time.Now().Unix(), 10))
}

func InspectView(name string, s store.ViewStore) (models.View, error) {
	return s.Load(name)
}

func DeleteView(name string, s store.ViewStore) error {
	return s.Delete(name)
}
