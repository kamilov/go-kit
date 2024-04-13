package db

import "database/sql"

type Tx struct {
	*sql.Tx
}
