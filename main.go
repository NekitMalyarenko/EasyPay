package src

import (
	"db"
	"net/http"
	"io/ioutil"
	"log"
	"my_errors"
	"encoding/json"
	"types"
	"errors"
	myHandlers "handlers"
)


type handler struct {
	method                 func (input map[string]interface{}) (string, error)
	requirements           []string
	requiresAuthentication bool
	requiresUser           bool
}


var handlers map[string]*handler

var rawMethodHandlers map[string]func(input map[string]interface{}, rawData []byte)string


func main() {
	initHandlers()
	initRawHandlers()

	db.GetInstance()

	http.HandleFunc("/", Handle)

	//http.HandleFunc("/raw", RawHandle)
	//http.ListenAndServe(":8080", nil)

	//services.SendSMS("test", "380967519036")
}


func Handle(w http.ResponseWriter, r *http.Request) {
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(my_errors.Errors[my_errors.CantReadBody].Error()))
		return
	}
	defer r.Body.Close()

	var parsedInput map[string]interface{}
	err = json.Unmarshal(input, &parsedInput)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(my_errors.Errors[my_errors.JsonUnmarshalError].Error()))
		return
	}

	log.Println("New connection with input data:", parsedInput)

	if parsedInput["method_name"] != nil {
		methodName := parsedInput["method_name"].(string)
		handler := handlers[methodName]
		var (
			methodData map[string]interface{}
			user       *types.User
		)

		if handler != nil {
			methodData = parsedInput["method_data"].(map[string]interface{})

			if handler.requiresAuthentication {
				handler.requirements = append(handler.requirements, "phone_number")
				handler.requirements = append(handler.requirements, "password")
			} else if handler.requiresUser {
				handler.requirements = append(handler.requirements, "phone_number")
			}

			if len(handler.requirements) != 0 {
				if !handler.checkRequirements(methodData) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(my_errors.Errors[my_errors.ArgumentsError].Error()))
					return
				}
			}

			if handler.requiresAuthentication {
				user = authenticateUser(methodData)
				if user == nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(my_errors.Errors[my_errors.AuthenticationError].Error()))
					return
				} else if handler.requiresUser {
					methodData["user"] = user
				}
			}

			if handler.requiresUser && !handler.requiresAuthentication {
				user = getUser(methodData)
				if user != nil {
					methodData["user"] = user
				}
			}

			response, err := handler.method(methodData)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(my_errors.Errors[my_errors.NoSuchMethodError].Error()))
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(my_errors.Errors[my_errors.MethodNameIsEmpty].Error()))
		return
	}

	//w.WriteHeader(200)
}


func RawHandle(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()

	rawFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Println(len(rawFile))

	var parsedMethodData map[string]interface{}
	err = json.Unmarshal([]byte(r.FormValue("method_data")), &parsedMethodData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("json parsing error"))
		return
	}

	if parsedMethodData["method_name"] != nil {
		methodName := parsedMethodData["method_name"].(string)
		var methodData = parsedMethodData["method_data"].(map[string]interface{})

		if rawMethodHandlers[methodName] != nil {
			response := rawMethodHandlers[methodName](methodData, rawFile)
			log.Println("response:", response)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		} else {
			log.Println(errors.New("no such method error"))
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("no such method error"))
			return
		}
	} else {
		log.Println(errors.New("method_name is empty error"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("method_name is empty error"))
		return
	}

}


func (handler *handler) checkRequirements(inputData map[string]interface{}) bool {
	for _, requirement := range handler.requirements {
		if inputData[requirement] == nil {
			return false
		}
	}

	return true
}


func authenticateUser(inputData map[string]interface{}) (*types.User) {
	customer, err := db.GetInstance().Customers.GetCustomer(inputData["phone_number"].(string))
	if err != nil {
		log.Println(err)
		return nil
	}

	if customer != nil && customer.Password == inputData["password"].(string) {
		return &types.User{Customer : customer, Seller : nil}
	} else {
		seller, err := db.GetInstance().Sellers.GetSeller(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return nil
		}

		if seller != nil && seller.Password == inputData["password"].(string) {
			return &types.User{Customer : nil, Seller : seller}
		} else {
			return nil
		}
	}

}


func getUser(inputData map[string]interface{}) *types.User {
	customer, err := db.GetInstance().Customers.GetCustomer(inputData["phone_number"].(string))
	if err != nil {
		log.Println(err)
		return nil
	}

	if customer != nil  {
		return &types.User{Customer : customer, Seller : nil}
	} else {
		seller, err := db.GetInstance().Sellers.GetSeller(inputData["phone_number"].(string))
		if err != nil {
			log.Println(err)
			return nil
		}

		if seller != nil  {
			return &types.User{Customer : nil, Seller : seller}
		} else {
			return nil
		}
	}
}


func initHandlers() {
	handlers = make(map[string]*handler)

	handlers["userLogin"] = &handler{
		method : myHandlers.UserLogin,
		requirements : []string{},
		requiresAuthentication : true,
		requiresUser : true,
	}

	handlers["startVerification"] = &handler{
		method : myHandlers.StartVerification,
		requirements : []string{},
		requiresUser : true,
	}
	handlers["verifyPhone"] = &handler{
		method : myHandlers.VerifyPhone,
		requirements : []string{"verification_code"},
		requiresUser : true,
	}

	handlers["customerRegister"] = &handler{
		method : myHandlers.CustomerRegister,
		requirements : []string{"first_name", "last_name", "password",
			"phone_number"},
	}

	handlers["sellerRegister"] = &handler{
		method : myHandlers.SellerRegister,
		requirements : []string{"first_name", "last_name", "password",
			"phone_number", "description"},
	}

	handlers["addTransaction"] = &handler{
		method : myHandlers.AddTransaction,
		requirements : []string{"shop_id", "transaction"},
		requiresAuthentication : true,
		requiresUser : true,
	}

	/*
	methodHandlers["getTransactionById"] = handlers.GetTransactionById
	methodHandlers["getNewestTransactions"] = handlers.GetNewestTransactions
	methodHandlers["getOldestTransactions"] = handlers.GetOldestTransactions*/
}


func initRawHandlers() {
	rawMethodHandlers = make(map[string]func(input map[string]interface{}, rawData []byte)string)
	//rawMethodHandlers["testUpload"] = handlers.TestPhoto
}