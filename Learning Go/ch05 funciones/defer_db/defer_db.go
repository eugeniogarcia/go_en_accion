package defer_db

import (
	"context"
	"database/sql"
)

func DoSomeInserts(ctx context.Context, db *sql.DB, value1, value2 string) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
		if err != nil {
			tx.Rollback()
		}
	}() //En defer incluimos una llamada a una funcion, anónima o no, que se ejecutará al salir del ámbito de la función que la contiene. Las variables que se capturan, se capturan en el momento de hacer el defer, no en el momento de ejecutar la función (al salir del ámbito)

	_, err = tx.ExecContext(ctx, "INSERT INTO FOO (val) values $1", value1)
	if err != nil {
		return err
	}
	// use tx to do more database inserts here
	return nil
}
