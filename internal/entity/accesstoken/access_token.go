package accesstoken

type AccessToken struct {
	deviceId string
	teamId   int
	token    string
	userId   int
}

func NewAccessToken(deviceId string, teamId int, token string, userId int) *AccessToken {
	return &AccessToken{
		deviceId: deviceId,
		teamId:   teamId,
		token:    token,
		userId:   userId,
	}
}

func (entity *AccessToken) GetDeviceId() string {
	return entity.deviceId
}

func (entity *AccessToken) GetTeamId() int {
	return entity.teamId
}

func (entity *AccessToken) GetToken() string {
	return entity.token
}

func (entity *AccessToken) GetUserId() int {
	return entity.userId
}
