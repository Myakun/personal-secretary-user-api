package user

type User struct {
	email      string
	id         string
	isInserted bool
	name       string
	password   string
}

type UserDTO struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	Name  string `json:"name"`
}

func NewUser(email string, id string, name string, password string) *User {
	return &User{
		email:      email,
		id:         id,
		isInserted: false,
		name:       name,
		password:   password,
	}
}

func (entity *User) GetEmail() string {
	return entity.email
}

func (entity *User) GetId() string {
	return entity.id
}

func (entity *User) GetName() string {
	return entity.name
}

func (entity *User) GetPassword() string {
	return entity.password
}

func (entity *User) IsInserted() bool {
	return entity.isInserted
}

func (entity *User) setIsInserted(isInserted bool) {
	entity.isInserted = isInserted
}

func (entity *User) SetPassword(password string) {
	entity.password = password
}
