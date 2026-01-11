package repositories

import (
	"database/sql"
	"net/http"
	"runners-postgresql/models"
)

type ResultsRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx // Se usa para las operaciones de actualización que requieren ejecutarse dentro de una transacción
}

func NewResultsRepository(dbHAndler *sql.DB) *ResultsRepository {
	return &ResultsRepository{
		dbHandler: dbHAndler,
	}
}

func (rr ResultsRepository) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	query := `
		INSERT INTO results(runner_id, race_result, location, position, year)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	// ejecutamos la query dentro de una transaccion (estamos cambiando datos)
	rows, err := rr.transaction.Query(query, result.RunnerID, result.RaceResult, result.Location, result.Position, result.Year)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	// aseguramos que se cierre el cursor
	defer rows.Close()

	var resultId string
	// iteramos sobre el cursor
	for rows.Next() {
		// capturamos los datos recuperados con el cursor
		err := rows.Scan(&resultId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	// si hubiera un error en la iteración del cursor
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: rows.Err().Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         resultId,
		RunnerID:   result.RunnerID,
		RaceResult: result.RaceResult,
		Location:   result.Location,
		Position:   result.Position,
		Year:       result.Year,
	}, nil
}

func (rr ResultsRepository) DeleteResult(resultId string) (*models.Result, *models.ResponseError) {
	query := `
		DELETE FROM results
		WHERE id = $1
		RETURNING runner_id, race_result, year`

	rows, err := rr.transaction.Query(query, resultId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var runnerId, raceResult string
	var year int
	for rows.Next() {
		err := rows.Scan(&runnerId, &raceResult, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: rows.Err().Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         resultId,
		RunnerID:   runnerId,
		RaceResult: raceResult,
		Year:       year,
	}, nil
}

func (rr ResultsRepository) GetAllRunnersResults(runnerId string) ([]*models.Result, *models.ResponseError) {
	query := `
	SELECT id, race_result, location, position, year
	FROM results
	WHERE runner_id = $1`

	// ejecutamos la query (consulta)
	rows, err := rr.dbHandler.Query(query, runnerId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	// aseguramos que se cierre el cursor
	defer rows.Close()

	results := make([]*models.Result, 0)
	var id, raceResult, location string
	var position, year int

	// iteramos sobre el cursor
	for rows.Next() {
		// capturamos los datos recuperados con el cursor
		err := rows.Scan(&id, &raceResult, &location, &position, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		result := &models.Result{
			ID:         id,
			RunnerID:   runnerId,
			RaceResult: raceResult,
			Location:   location,
			Position:   position,
			Year:       year,
		}

		results = append(results, result)
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: rows.Err().Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return results, nil
}

func (rr ResultsRepository) GetPersonalBestResults(runnerId string) (string, *models.ResponseError) {
	query := `
	SELECT MIN(race_result)
	FROM results
	WHERE runner_id = $1`

	rows, err := rr.dbHandler.Query(query, runnerId)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var raceResult string

	for rows.Next() {
		err := rows.Scan(&raceResult)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: rows.Err().Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return raceResult, nil
}

func (rr ResultsRepository) GetSeasonBestResults(runnerId string, year int) (string, *models.ResponseError) {
	query := `
	SELECT MIN(race_result)
	FROM results
	WHERE runner_id = $1 AND year = $2`

	rows, err := rr.dbHandler.Query(query, runnerId, year)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var raceResult string

	for rows.Next() {
		err := rows.Scan(&raceResult)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: rows.Err().Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return raceResult, nil
}
