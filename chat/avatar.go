package main

import "errors"

// ErrNoAvatarURL is Error when an instance of Avatar can not return avatar's URL
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

type Avatar interface {
	// GetAvatarURL is avataer's URL
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
