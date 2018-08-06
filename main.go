package main

import (
	"db"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"types"
	myHandlers "handlers"
	"my_errors"
	"local_storage"
	"strings"
)


const (
	liqpayServerAddress = "54.229.105.178"
)


type handler struct {
	method func (input map[string]interface{}) (string, error)
	requirements           []string
	requiresAuthentication bool
	requiresUser           bool
}

type fileHandler struct {
	method func (input map[string]interface{}, rawFile []byte) (string, error)
	requirements           []string
	requiresAuthentication bool
	requiresUser           bool
}


var handlers map[string]*handler

var fileHandlers  map[string]*fileHandler


func main() {
	//localhost:8080/api?data=%7B%22method_data%22%3A%7B%22password%22%3A%221234%22%2C%22phone_number%22%3A%22380967519036%22%2C%22products%22%3A%5B%7B%22id%22%3A1%2C%22quantity%22%3A2%7D%2C%7B%22id%22%3A2%2C%22quantity%22%3A4%7D%5D%2C%22shop_id%22%3A23%7D%2C%22method_name%22%3A%22pay%22%7D
	initHandlers()
	initRawHandlers()

	db.GetInstance()
	defer db.GetInstance().CloseConnection()

	local_storage.GetInstance()

	http.HandleFunc("/api", MainHandler)
	http.HandleFunc("/files", FileHandler)
	http.HandleFunc("/callback", LiqpayCallbackHandler)
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("<h1>Hello World!</h1>"))
	})
	http.ListenAndServe(":8080", nil)

	/*u, err := url.Parse("localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	parsedData, _ := json.Marshal(map[string]interface{}{
		"method_data" : map[string]interface{}{
			"shop_id" : 23,
			"products" : []map[string]interface{}{
				{
					"id" : 1,
					"quantity" : 2,
				},
				{
					"id" : 2,
					"quantity" : 4,
				},
			},

			"phone_number" : "380967519036",
			"password" : "1234",
		},
		"method_name" : "pay",
	})
	q.Set("data", string(parsedData))
	u.RawQuery = q.Encode()
	fmt.Println(u)*/
	//services.SendSMS("test", "380967519036")
}


func LiqpayCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ip := []byte(r.RemoteAddr)[:strings.Index(r.RemoteAddr, ":")]

	if r.Method == "POST" && string(ip) == liqpayServerAddress {
		myHandlers.CheckAllPayments()
	} else {
		log.Println("DENIED", string(ip))
	}
}


func MainHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Println("METHOD:", r.Method, "from:", r.RemoteAddr)

	var (
		input []byte
		err   error
	)

	if r.Method == "POST" {
		input, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(my_errors.Errors[my_errors.CantReadBody].Error()))
			return
		}
	} else if r.Method == "GET" {
		input = []byte(r.URL.Query().Get("data"))
	}

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
			user *types.User
		)

		if handler != nil {
			if parsedInput["method_data"] != nil {
				methodData = parsedInput["method_data"].(map[string]interface{})
			} else {
				methodData = make(map[string]interface{})
			}

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
					log.Println("Not enough arguments")
					return
				}
			}

			if handler.requiresAuthentication {
				user = authenticateUser(methodData)
				if user == nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(my_errors.Errors[my_errors.AuthenticationError].Error()))
					log.Println("Authentication Error")
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
				log.Println("error:", err.Error())
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
			log.Println("response:", response)
			return

		} else {
			log.Println(my_errors.Errors[my_errors.NoSuchMethodError].Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(my_errors.Errors[my_errors.NoSuchMethodError].Error()))
			return
		}

	} else {
		log.Println(my_errors.Errors[my_errors.MethodNameIsEmpty].Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(my_errors.Errors[my_errors.MethodNameIsEmpty].Error()))
		return
	}
}


func FileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if len(r.FormValue("method_data")) != 0 {
		var parsedInput map[string]interface{}
		err := json.Unmarshal([]byte(r.FormValue("method_data")), &parsedInput)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(my_errors.Errors[my_errors.JsonUnmarshalError].Error()))
			return
		}
		log.Println("New connection with input data:", parsedInput)

		file, headers, err := r.FormFile("file")
		if err != nil {
			log.Println("1", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Println("headers:", headers.Header)

		rawFile, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("2", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if parsedInput["method_name"] != nil {
			methodName := parsedInput["method_name"].(string)
			handler := fileHandlers[methodName]
			var (
				methodData map[string]interface{}
				user *types.User
			)

			if handler != nil {
				if parsedInput["method_data"] != nil {
					methodData = parsedInput["method_data"].(map[string]interface{})
				} else {
					methodData = make(map[string]interface{})
				}

				if parsedInput["method_data"] != nil {
					methodData = parsedInput["method_data"].(map[string]interface{})
				} else {
					methodData = make(map[string]interface{})
				}

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
					}
				}

				if handler.requiresAuthentication {
					user = authenticateUser(methodData)
					if user == nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(my_errors.Errors[my_errors.AuthenticationError].Error()))
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

				response, err := handler.method(methodData, rawFile)
				if err != nil {
					log.Println("error:", err.Error())
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				log.Println("response:", response)
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

	} else {
		log.Println(my_errors.Errors[my_errors.ArgumentsError].Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(my_errors.Errors[my_errors.ArgumentsError].Error()))
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


func (handler *fileHandler) checkRequirements(inputData map[string]interface{}) bool {
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

	handlers["healthCheck"] = &handler{
		method : myHandlers.HealthCheck,
	}

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
	handlers["getSeller"] = &handler{
		method : myHandlers.GetSeller,
		requirements : []string{"seller_id"},
	}

	handlers["addTransaction"] = &handler{
		method : myHandlers.AddTransaction,
		requirements : []string{"shop_id", "transaction"},
		requiresAuthentication : true,
		requiresUser : true,
	}
	handlers["getOldestTransactions"] = &handler{
		method : myHandlers.GetOldestTransactions,
		requirements : []string{"transaction_id", "number_of_transactions"},
		requiresAuthentication : true,
		requiresUser : true,
	}
	handlers["getNewestTransactions"] = &handler{
		method : myHandlers.GetNewestTransactions,
		requirements : []string{"last_transaction_id"},
		requiresAuthentication : true,
		requiresUser : true,
	}
	handlers["getTransactionById"] = &handler{
		method : myHandlers.GetTransactionById,
		requirements : []string{"transaction_id"},
		requiresAuthentication: true,
	}
	handlers["getProductImage"] = &handler{
		method : myHandlers.GetTransactionById,
		requirements : []string{"transaction_id", "product_id"},
		requiresAuthentication: true,
	}

	handlers["shopRegister"] = &handler{
		method : myHandlers.ShopRegister,
		requirements : []string{"shop"},
		requiresAuthentication : true,
		requiresUser : true,
	}
	handlers["shopAddSeller"] = &handler{
		method : myHandlers.ShopAddSeller,
		requirements : []string{"seller"},
		requiresAuthentication : true,
		requiresUser : true,
	}
	handlers["shopAddProducts"] = &handler{
		method : myHandlers.ShopAddProducts,
		requirements : []string{"products"},
		requiresAuthentication : true,
		requiresUser : true,
	}
	handlers["getShop"] = &handler{
		method : myHandlers.GetShop,
		requirements : []string{"shop_id"},
	}
	handlers["getShopCardNumber"] = &handler{
		method : myHandlers.GetShopCardNumber,
		requirements : []string{"shop_id"},
		requiresAuthentication : true,
	}
	handlers["getShopProducts"] = &handler{
		method : myHandlers.GetShopProducts,
		requirements : []string{"shop_id"},
	}

	handlers["addLike"] = &handler{
		method: myHandlers.AddLike,
		requirements: []string{"shop_id"},
		requiresAuthentication: true,
		requiresUser: true,
	}
	handlers["addDislike"] = &handler{
		method: myHandlers.AddDislike,
		requirements: []string{"shop_id"},
		requiresAuthentication: true,
		requiresUser: true,
	}
	handlers["getUserRating"] = &handler{
		method: myHandlers.GetRating,
		requirements: []string{"shop_id"},
		requiresAuthentication: true,
		requiresUser: true,
	}

	handlers["pay"] = &handler{
		method: myHandlers.Pay,
		requirements: []string{"shop_id", "products"},
		requiresAuthentication: true,
		requiresUser : true,
	}
}


func initRawHandlers() {
	fileHandlers = make(map[string]*fileHandler)

	fileHandlers["addCustomerImage"] = &fileHandler{
		method : myHandlers.AddCustomerImage,
		requiresAuthentication : true,
		requiresUser : true,
	}
}