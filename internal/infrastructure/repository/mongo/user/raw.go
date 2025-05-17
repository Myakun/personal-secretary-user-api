package user

import "github.com/Myakun/personal-secretary-user-api/internal/domain/user"

type rawUser struct {
	Email    string `bson:"email"`
	Id       string `bson:"_id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

func (r *rawUser) ToUser() *user.User {
	return user.NewUserFromStorage(r.Email, r.Id, r.Name, r.Password)
}
