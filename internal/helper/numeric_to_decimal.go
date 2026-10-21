package helper

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func NumericToDecimal(n pgtype.Numeric) (decimal.Decimal, error) {
	if !n.Valid {
		return decimal.Decimal{}, errors.New("invalid numeric")
	}

	if n.NaN || n.InfinityModifier != pgtype.Finite {
		return decimal.Decimal{}, errors.New("numeric is not a finite decimal")
	}

	if n.Int == nil {
		return decimal.Zero, nil
	}

	return decimal.NewFromBigInt(n.Int, n.Exp), nil
}
