package appError

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func FromDBError(err error) *AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return New(http.StatusNotFound, ErrNotFound, "record not found", err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return New(http.StatusConflict, ErrDuplicate, "a record with this value already exists", err)
		case "23502":
			return New(http.StatusBadRequest, ErrFieldRequired, "a required field is missing", err)
		case "23514":
			return New(http.StatusBadRequest, ErrCheckViolation, "value does not meet the required constraints", err)
		}
	}

	return New(http.StatusInternalServerError, ErrInternal, "internal server error", err)
}
