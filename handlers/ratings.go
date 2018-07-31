package handlers

import (
	"types"
	"db"
	"my_errors"
	"log"
	"encoding/json"
)


const (
	actionLike    = "like"
	actionDislike = "dislike"
)


func AddLike(inputData map[string]interface{}) (string, error) {
	customer := inputData["user"].(*types.User).Customer
	shopId := int64(inputData["shop_id"].(float64))

	dbRating, err := db.GetInstance().Ratings.GetRating(customer.Id, shopId)
	if err != nil {
		log.Println(err)
		return my_errors.GetError(my_errors.DBError)
	}

	if dbRating != nil {

		switch dbRating.Action {

			case actionLike:
				err = db.GetInstance().Shops.AddLikes(-1, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}

				err = db.GetInstance().Ratings.DeleteRating(dbRating.UserId, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}
			break

			case actionDislike:
				err = db.GetInstance().Shops.AddLikes(1, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}

				err = db.GetInstance().Shops.AddDislikes(-1, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}

				dbRating.Action = actionLike
				err = db.GetInstance().UpdateRating(dbRating)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}
			break
		}

	} else {
		log.Println("3")
		err = addRatingToDB(&types.Rating{
			Action: actionLike,
			UserId: customer.Id,
			ShopId: shopId,
		})
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		}
	}

	return my_errors.SuccessfullyOperation()
}


func AddDislike(inputData map[string]interface{}) (string, error) {
	customer := inputData["user"].(*types.User).Customer
	shopId := int64(inputData["shop_id"].(float64))

	dbRating, err := db.GetInstance().Ratings.GetRating(customer.Id, shopId)
	if err != nil {
		log.Println(err)
		return my_errors.GetError(my_errors.DBError)
	}

	if dbRating != nil {

		switch dbRating.Action {

			case actionDislike:
				err = db.GetInstance().Shops.AddDislikes(-1, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}

				err = db.GetInstance().DeleteRating(dbRating.UserId, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}
				break

			case actionLike:
				err = db.GetInstance().Shops.AddDislikes(1, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}

				err = db.GetInstance().Shops.AddLikes(-1, dbRating.ShopId)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}

				dbRating.Action = actionDislike
				err = db.GetInstance().UpdateRating(dbRating)
				if err != nil {
					log.Println(err)
					return my_errors.GetError(my_errors.DBError)
				}
				break
		}

	} else {
		err = addRatingToDB(&types.Rating{
			Action: actionDislike,
			UserId: customer.Id,
			ShopId: shopId,
		})
		if err != nil {
			return my_errors.GetError(my_errors.DBError)
		}
	}

	return my_errors.SuccessfullyOperation()
}


func GetRating(inputData map[string]interface{}) (string, error) {
	user := inputData["user"].(*types.User).Customer
	shopId := int64(inputData["shop_id"].(float64))

	rating, err := db.GetInstance().Ratings.GetRating(user.Id, shopId)
	if err != nil {
		log.Println(err)
		return my_errors.GetError(my_errors.DBError)
	}

	var action string
	if rating != nil {
		action = rating.Action
	} else {
		action = "none"
	}

	response, _ := json.Marshal(map[string]interface{}{
		"status" : "ok",
		"action" : action,
	})

	return string(response), nil
}


func addRatingToDB(rating *types.Rating) error {
	err := db.GetInstance().Ratings.AddRating(rating)
	if err != nil {
		return err
	}

	if rating.Action == actionLike {
		err = db.GetInstance().Shops.AddLikes(1, rating.ShopId)
	} else {
		err = db.GetInstance().Shops.AddDislikes(1, rating.ShopId)
	}

	return err
}