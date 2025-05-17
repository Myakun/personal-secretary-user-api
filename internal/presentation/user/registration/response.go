package registration

//import entityUser "github.com/Myakun/personal-secretary-user-api/internal/entity/user"

type ErrorResponse struct {
	Err string `json:"error"`
}

type RegisterUserResult struct {
	ErrorResponse   *ErrorResponse
	Success         bool
	SuccessResponse *SuccessResponse
}

type SuccessResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	//	User         *entityUser.UserDTO `json:"user"`
}

func newErrorResponse(err string) *ErrorResponse {
	return &ErrorResponse{
		Err: err,
	}
}
