package db

import (
	"types"
	"upper.io/db.v3"
)

type Ratings string

const (
	ratingsTable = "rates"
)


func (*Ratings) AddRating(userRating *types.Rating) error {
	_, err := currentInstance.instance.
		InsertInto(ratingsTable).Values(userRating).Exec()
	return err
}


func (*Ratings) GetRating(customerId, shopId int64) (*types.Rating, error) {
	var rating *types.Rating
	res := currentInstance.instance.Collection(ratingsTable).Find(db.Cond{
		"user_id =" :customerId,
		"shop_id =" :shopId,
	})

	exists, err := res.Exists()
	if err != nil {
		return nil, err
	}

	if exists {
		res.One(&rating)
		return rating, nil
	} else {
		return nil, nil
	}
}


func (*Ratings) UpdateRating(rating *types.Rating) error {
	_, err := currentInstance.instance.Update(ratingsTable).
		Set("action = ?", rating.Action).
		Where("id = ?", rating.Id).
		Exec()

	return err
}


func (*Ratings) DeleteRating(customerId, shopId int64) error {
	return currentInstance.instance.Collection(ratingsTable).Find(db.Cond{
		"user_id =" :customerId,
		"shop_id =" :shopId,
	}).Delete()
}
