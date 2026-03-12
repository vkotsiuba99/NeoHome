package auth

import (
	"errors"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func mapStorageErr(err error, notFoundErr ...error) error {
	notFoundMappedErr := service.ErrNotFound
	if len(notFoundErr) > 0 && notFoundErr[0] != nil {
		notFoundMappedErr = notFoundErr[0]
	}

	switch {
	case errors.Is(err, storage.ErrNotFound):
		return notFoundMappedErr
	case errors.Is(err, storage.ErrConflict):
		return service.ErrConflict
	default:
		return err
	}
}
