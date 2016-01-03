package service

import "net/http"
import "github.com/gorilla/mux"
import "encoding/json"
import "strconv"

type Cat struct {
	Id      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	CatType string `json:"cat_type,omitempty" gorm:"column:catType"`
}

func (s *JSONService) ReadCats(r *http.Request) (int, interface{}, error) {
	var cats []Cat
	db := s.Config.DB
	db.Find(&cats)
	res := map[string][]Cat{
		"data": cats,
	}
	return http.StatusOK, res, nil
}

func (s *JSONService) CreateCat(r *http.Request) (int, interface{}, error) {
	var cat Cat
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cat)
	if cat.Name != "" && cat.CatType != "" {
		db := s.Config.DB
		db.Create(&cat)
		return http.StatusCreated, map[string]string{"status": "success"}, nil
	}
	return http.StatusInternalServerError, nil, err
}

func (s *JSONService) UpdateCat(r *http.Request) (int, interface{}, error) {
	var cat Cat
	db := s.Config.DB
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cat)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	if cat.Name != "" && cat.CatType != "" {
		cat.Id, _ = strconv.ParseInt(mux.Vars(r)["id"], 0, 64)
		db.Save(&cat)
		return http.StatusNoContent, map[string]Cat{
			"data": cat,
		}, nil
	}
	return http.StatusOK, map[string]string{"status": "success"}, nil
}

func (s *JSONService) DeleteCat(r *http.Request) (int, interface{}, error) {
	var cat Cat
	id := mux.Vars(r)["id"]
	db := s.Config.DB
	db.First(&cat, id)
	db.Delete(&cat)
	return http.StatusOK, map[string]string{"status": "success"}, nil
}

func (s *JSONService) ReadCat(r *http.Request) (int, interface{}, error) {
	var cat Cat
	id := mux.Vars(r)["id"]
	db := s.Config.DB
	db.First(&cat, id)
	res := map[string]Cat{
		"data": cat,
	}
	return http.StatusOK, res, nil
}
