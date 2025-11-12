package helpers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var pgxErrCode = map[string]int{ //nolint:gochecknoglobals // postgres sqlstate mapping
	"23505": EResourceAlreadyExists,       // unique violation
	"23502": EDBNotNullViolation,          // not-null constraint violation
	"23514": EDBCheckViolation,            // check constraint violation
	"23000": EDBDataIntegrityViolation,    // data integrity violation
	"22007": EDBInvalidDatetimeFormat,     // invalid datetime format
	"22003": EDBNumericValueOutOfRange,    // numeric value out of range
	"22801": EDBStringLengthExceeded,      // string length exceeded
	"22P02": EDBInvalidTextRepresentation, // e.g., invalid UUID format
}

func ConvertPgErrToAppCode(err error) int {
	if err == nil {
		return Success
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return EResourceNotFound
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return EResourceAlreadyExists
	}

	var pgxErr *pgconn.PgError
	if ok := errors.As(err, &pgxErr); ok {
		if appCode, exists := pgxErrCode[pgxErr.Code]; exists {
			return appCode
		}
	}
	return EDatabaseError
}
