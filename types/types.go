package types

type Customer struct {
	Id          int64  `db:"id,omitempty"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
	Image       string `db:"image"`
}

type Seller struct {
	Id          int64  `db:"id,omitempty"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Description string `db:"description"`
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
	ShopId      int64  `db:"shop_id"`
	Image       string `db:"image"`
}


type User struct {
	*Customer
	*Seller
}

type Verification struct {
	Id               int64  `db:"id,omitempty"`
	PhoneNumber      string `db:"phone_number"`
	VerificationCode int64  `db:"verification_code"`
	IsVerified       bool   `db:"is_verified"`
	StartTime        string `db:"start_time"`
}

type Transaction struct {
	Id               int64  `db:"id,omitempty"`
	UserId           int64  `db:"user_id"           json:"user_id"`
	ShopId           int64  `db:"shop_id"           json:"shop_id"`
	Products         string `db:"products"          json:"products"`
	TotalPrice       int    `db:"total_price"       json:"total_price"`
	Date             string `db:"date"              json:"date"`
	VerificationCode int64  `db:"verification_code" json:"verification_code"`
}

type Shop struct {
	Id          int64  `db:"id,omitempty"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Image 		string `db:"image"`
	Likes       int	   `db:"likes"`
	Dislikes    int	   `db:"dislikes"`
	CardNumber  int64  `db:"card_number"`
	Email       string `db:"email"`
	RowSellers  string `db:"sellers"`
	RowProducts string `db:"products"`
}

type Product struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

type Rating struct {
	Id     int64  `db:"id,omitempty"`
	UserId int64  `db:"user_id"`
	ShopId int64  `db:"shop_id"`
	Action string `db:"action"`
}

type PaymentData struct {
	Id int64 `json:"id"`
	Transaction
}