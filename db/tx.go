package db

import "database/sql"

// Tx transaction control struct
type Tx struct {
	*sql.Tx
}
