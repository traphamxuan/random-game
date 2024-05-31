package authentication

import (
	"context"
	"game-random-api/internal/api/dto"
	"game-random-api/internal/orm"
	"game-random-api/package/logger"

	"github.com/traphamxuan/gobs"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	log      *logger.Logger
	orm      *orm.Orm
	jwtToken *JwtToken
}

var _ gobs.IService = (*Authentication)(nil)

func (a *Authentication) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&orm.Orm{},
		&JwtToken{},
	}
	onSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		a.log = dependencies[0].(*logger.Logger)
		a.orm = dependencies[1].(*orm.Orm)
		a.jwtToken = dependencies[2].(*JwtToken)
		return nil
	}
	sb.OnSetup = &onSetup
	return nil
}

func (a *Authentication) Signin(ctx context.Context, creds dto.Credentials) (*dto.RespToken, error) {
	// Get the expected password from our in memory map
	user, err := a.orm.User.GetOne(ctx, orm.QueryUserFilter{Email: &creds.Email})
	if nil != err || user == nil {
		a.log.Error(err.Error())
		return nil, err
	}
	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Pass)); nil != err {
		return nil, err
	}

	userId := user.ID
	tokenStr, expired, err := a.jwtToken.ComposeToken(userId)
	if nil != err {
		a.log.Error(err.Error())
		return nil, err
	}

	return &dto.RespToken{
		Token:       tokenStr,
		UserID:      userId,
		AccessToken: "",
		ExpiresAt:   expired,
	}, nil
}

func (a *Authentication) RefreshToken(ctx context.Context, tokenStr string) (*dto.RespToken, error) {
	tokenStr, expired, err := a.jwtToken.RefreshToken(tokenStr)
	if nil != err {
		return nil, err
	}
	return &dto.RespToken{
		Token:     tokenStr,
		ExpiresAt: expired,
	}, nil
}
