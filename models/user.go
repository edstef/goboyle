package models

import (
	"context"

	"github.com/uptrace/bun"
)

type User struct {
	Id         string `bun:",pk,notnull,default:uuid_generate_v4()" json:"id"`
	Name       string `bun:"name,notnull" json:"name"`
	PictureURL string `bun:"picture,notnull,default:'/defaults/1'" json:"picture_url"`
	Theme      string `bun:"theme,notnull,default:'default_theme_1'" json:"theme"`
}

func (m *Models) CreateUser(name string) (*User, error) {
	user := &User{
		Name: name,
	}
	_, err := m.db.NewInsert().Model(user).Exec(context.Background())

	return user, err
}

func (m *Models) GetUserById(id string) (*User, error) {
	user := &User{}
	err := m.db.NewSelect().
		Model(user).
		Where("? = ?", bun.Ident("id"), id).
		Scan(context.Background())

	return user, err
}

type UserToken struct {
	UserId string `bun:"user,pk,notnull"`
	Token  string `bun:"token,notnull"`
}

func (m *Models) UpsertUserToken(userId string, token string) error {
	hashedToken, err := HashPassword(token)

	if err != nil {
		return err
	}

	userToken := &UserToken{
		UserId: userId,
		Token:  hashedToken,
	}

	_, err = m.db.NewInsert().
		Model(userToken).
		On("CONFLICT (\"user\") DO UPDATE").
		Set("token = EXCLUDED.token").
		Exec(context.Background())

	return err
}

func (m *Models) ValidateUserToken(userId string, token string) error {
	userToken := &UserToken{}
	err := m.db.NewSelect().
		Model(userToken).
		Where("? = ?", bun.Ident("user"), userId).
		Scan(context.Background())

	if err != nil {
		return err
	}

	return CheckPasswordHash(userToken.Token, token)
}
