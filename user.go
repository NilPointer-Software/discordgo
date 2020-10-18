package discordgo

import "strings"

// A User stores all data for an individual Discord user.
type User struct {
	// The ID of the user.
	ID string `json:"id"`

	// The user's username.
	Username string `json:"username"`

	// The discriminator of the user (4 numbers after name).
	Discriminator string `json:"discriminator"`

	// The hash of the user's avatar. Use Session.UserAvatar
	// to retrieve the avatar itself.
	Avatar string `json:"avatar"`

	// Whether the user is a bot.
	Bot bool `json:"bot"`

	// Whether the user is a system user.
	System bool `json:"system"`

	// Whether the user has multi-factor authentication enabled.
	MFAEnabled bool `json:"mfa_enabled"`

	// The user's chosen language option.
	Locale string `json:"locale"`

	// Whether the user's email is verified.
	Verified bool `json:"verified"`

	// The email of the user. This is only present when
	// the application possesses the email scope for the user.
	Email string `json:"email"`

	// User Account Flags
	Flags int `json:"flags"`

	// Type of Nitro Subscriptions
	PremiumType PremiumType `json:"premium_type"`

	// Public User Account Flags
	PublicFlags int `json:"public_flags"`

	// The token of the user. This is only present for
	// the user represented by the current session.
	Token string `json:"token"`
}

type PremiumType int

const (
	PremiumTypeNone PremiumType = iota
	PremiumTypeNitroClassic
	PremiumTypeNitro
)

// String returns a unique identifier of the form username#discriminator
func (u *User) String() string {
	return u.Username + "#" + u.Discriminator
}

// Mention return a string which mentions the user
func (u *User) Mention() string {
	return "<@" + u.ID + ">"
}

// AvatarURL returns a URL to the user's avatar.
//    size:    The size of the user's avatar as a power of two
//             if size is an empty string, no size parameter will
//             be added to the URL.
func (u *User) AvatarURL(size string) string {
	var URL string
	if u.Avatar == "" {
		URL = EndpointDefaultUserAvatar(u.Discriminator)
	} else if strings.HasPrefix(u.Avatar, "a_") {
		URL = EndpointUserAvatarAnimated(u.ID, u.Avatar)
	} else {
		URL = EndpointUserAvatar(u.ID, u.Avatar)
	}

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}
