package main

import (
	"crypto/md5"
	"errors"
	"io"
	"strings"
)

// ErrNoAvatarURL is Error when an instance of Avatar can not return avatar's URL
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

// Avatar is hold AvatarURL
type Avatar interface {
	// GetAvatarURL is avataer's URL
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar is empty struct
type AuthAvatar struct{}

// UseAuthAvatar is user auth
var UseAuthAvatar AuthAvatar

// GetAvatarURL is method which return avatarURL
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar is empty struct
type GravatarAvatar struct{}

// UseGravatar is GravatarAvatar's instance
var UseGravatar GravatarAvatar

// GetAvatarURL is method which return gravatarURL
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(useridStr))
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// FileSystemAvatar is empty struct
type FileSystemAvatar struct{}

// UseFileSystemAvatar is FileSystemAvatar's instance
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL is method which return imagepath
func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
