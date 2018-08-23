package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL is Error when an instance of Avatar can not return avatar's URL
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

// Avatar is hold AvatarURL
type Avatar interface {
	// GetAvatarURL is avataer's URL
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars is Avatars' slice
type TryAvatars []Avatar

// GetAvatarURL is method which check all avatar
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar is empty struct
type AuthAvatar struct{}

// UseAuthAvatar is user auth
var UseAuthAvatar AuthAvatar

// GetAvatarURL is method which return avatarURL
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar is empty struct
type GravatarAvatar struct{}

// UseGravatar is GravatarAvatar's instance
var UseGravatar GravatarAvatar

// GetAvatarURL is method which return gravatarURL
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar is empty struct
type FileSystemAvatar struct{}

// UseFileSystemAvatar is FileSystemAvatar's instance
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL is method which return imagepath
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
