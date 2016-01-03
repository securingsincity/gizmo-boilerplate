package service

import "net/http"
import "github.com/gorilla/mux"
import "encoding/json"

func getEmptyForNow() error {
	return nil // This WILL compile
}

type Cat struct {
	Id      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	CatType string `json:"cat_type,omitempty"`
}

func (s *JSONService) ReadCats(r *http.Request) (int, interface{}, error) {
	var cats []Cat
	db, err := s.Config.MySQL.DB()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	rows, _ := db.Query("SELECT id, name, catType FROM cats")
	defer rows.Close()
	for rows.Next() {
		var cat Cat
		err := rows.Scan(&cat.Id, &cat.Name, &cat.CatType)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		cats = append(cats, cat)
	}
	err = rows.Err()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	res := map[string][]Cat{
		"data": cats,
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}

func (s *JSONService) CreateCat(r *http.Request) (int, interface{}, error) {
	var cat Cat
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cat)
	if cat.Name != "" && cat.CatType != "" {
		db, err := s.Config.MySQL.DB()
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		_, err = db.Exec(
			"INSERT INTO cats (name, catType) VALUES (?, ?)",
			cat.Name,
			cat.CatType,
		)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusCreated, map[string]string{"status": "success"}, nil
	}
	return http.StatusInternalServerError, nil, err
}

func (s *JSONService) UpdateCat(r *http.Request) (int, interface{}, error) {
	id := mux.Vars(r)["id"]
	db, err := s.Config.MySQL.DB()
	_, err = db.Exec("SELECT id FROM cats WHERE id = ?", id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	var cat Cat
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&cat)
	if cat.Name != "" && cat.CatType != "" {
		_, err = db.Exec("Update cats SET name = ?, catType = ? Where id = ?", cat.Name, cat.CatType, cat.Id)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusNoContent, map[string]Cat{
			"data": cat,
		}, nil
	}
	return http.StatusOK, map[string]string{"status": "success"}, nil
}

func (s *JSONService) DeleteCat(r *http.Request) (int, interface{}, error) {
	id := mux.Vars(r)["id"]
	db, err := s.Config.MySQL.DB()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	_, err = db.Exec("SELECT id FROM cats WHERE id = ?", id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	_, err = db.Exec("DELETE FROM cats WHERE id = ?", id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, map[string]string{"status": "success"}, nil
}

func (s *JSONService) ReadCat(r *http.Request) (int, interface{}, error) {
	id := mux.Vars(r)["id"]
	var cat Cat
	db, err := s.Config.MySQL.DB()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	err = db.QueryRow("SELECT id, name, catType FROM cats WHERE id = ?", id).Scan(&cat.Id, &cat.Name, &cat.CatType)
	res := map[string]Cat{
		"data": cat,
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}
