package bourses

import (
	db "ToDaMoon/database"
)

type database struct {
	DB map[string]db.DBM
}
