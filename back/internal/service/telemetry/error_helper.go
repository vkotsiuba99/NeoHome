package telemetry

import (
	"errors"

	"github.com/vkotsiuba99/NeoHome/back/internal/service"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
)

func mapStorageErr(err error) error {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return service.ErrNotFound
	case errors.Is(err, storage.ErrConflict):
		return service.ErrConflict
	default:
		return err
	}
}
