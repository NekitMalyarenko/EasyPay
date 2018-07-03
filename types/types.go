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

type Verification struct {
	Id               int64  `db:"id,omitempty"`
	PhoneNumber      string `db:"phone_number"`
	VerificationCode int64  `db:"verification_code"`
	IsVerified       bool   `db:"is_verified"`
}
