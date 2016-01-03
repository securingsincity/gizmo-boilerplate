package service

import "net/http"
import "github.com/maxwellhealth/mgo/bson"

type Dog struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Name string        `json:"name,omitempty" bson:"name,omitempty"`
}

func (s *JSONService) ReadDogs(r *http.Request) (int, interface{}, error) {
	var dogs []Dog
	mongo := s.Config.Mongo.C("dogs")
	_ = mongo.Find(bson.M{}).All(&dogs)
	res := map[string][]Dog{
		"data": dogs,
	}
	return http.StatusOK, res, nil
}
