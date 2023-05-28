package models

import (
	"context"

	"github.com/uptrace/bun"
)

type Profile struct {
	Id         string `bun:",pk,notnull,type:uuid,default:uuid_generate_v4()" json:"id"`
	Name       string `bun:"name,notnull" json:"name"`
	PictureURL string `bun:"name,notnull,default:'/defaults/1'" json:"picture_url"`
	Theme      string `bun:"theme,notnull,default:'default_theme_1'" json:"theme"`
}

func (m *Models) CreateProfile(name string) (*Profile, error) {
	profile := &Profile{
		Name: name,
	}
	_, err := m.db.NewInsert().Model(profile).Exec(context.Background())

	return profile, err
}

func (m *Models) GetProfileById(id string) (*Profile, error) {
	profile := &Profile{}
	err := m.db.NewSelect().
		Model(profile).
		Where("? = ?", bun.Ident("id"), id).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return profile, nil
}
