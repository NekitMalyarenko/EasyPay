package types

type User struct {
	Id          int64  `db:"id,omitempty"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	//PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	Image       string `db:"image"`
}
