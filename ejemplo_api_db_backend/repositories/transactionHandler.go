package repositories

import (
	"context"
	"database/sql"
)

func BeginTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	// creamos un contexto para la transacción
	ctx := context.Background()
	// iniciamos la transacción
	transaction, err := resultsRepository.dbHandler.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	// asignamos la transacción a ambos repositorios
	runnersRepository.transaction = transaction
	resultsRepository.transaction = transaction

	return nil
}

func RollbackTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	// toma la transacción de uno de los repositorios (los dos tienen la misma transacción así que da igual cual utilicemos)
	transaction := runnersRepository.transaction
	// limpiamos la transacción en ambos repositorios
	runnersRepository.transaction = nil
	resultsRepository.transaction = nil
	// hacemos el rollback
	return transaction.Rollback()
}

func CommitTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	// toma la transacción de uno de los repositorios (los dos tienen la misma transacción así que da igual cual utilicemos)
	transaction := runnersRepository.transaction

	// limpiamos la transacción en ambos repositorios
	runnersRepository.transaction = nil
	resultsRepository.transaction = nil

	// hacemos el commit
	return transaction.Commit()
}
