package main

import "database/sql"

type worker struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	WorkerType workerType `json:"worker_type"`
}

type workerType struct {
	ID             int    `json:"id"`
	WorkerTypeName string `json:"worker_type_name"`
}

func (wt *workerType) getWorkTypes(db *sql.DB) ([]workerType, error) {
	rows, err := db.Query("SELECT id, workertypename FROM workertypes")
	if err != nil {
		return nil, err
	}

	var workerTypes []workerType
	defer rows.Close()

	for rows.Next() {
		var wt workerType
		if err := rows.Scan(&wt.ID, &wt.WorkerTypeName); err != nil {
			return nil, err
		}
		workerTypes = append(workerTypes, wt)
	}
	return workerTypes, nil
}

func (w *worker) getWorkers(db *sql.DB) ([]worker, error) {
	rows, err := db.Query("SELECT w.id, w.name, wt.id AS wt_id, wt.workertypename AS worker_type_name FROM workers AS w INNER JOIN workertypes wt ON w.workertype = wt.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workers := []worker{}
	for rows.Next() {
		var w worker
		if err := rows.Scan(&w.ID, &w.Name, &w.WorkerType.ID, &w.WorkerType.WorkerTypeName); err != nil {
			return nil, err
		}
		workers = append(workers, w)
	}
	return workers, nil
}
