package discordgo

import (
	"encoding/json"
)

// This file contains all the possible structs that can be
// handled by AddHandler/EventHandler.
// DO NOT ADD ANYTHING BUT EVENT HANDLER STRUCTS TO THIS FILE.
//go:generate go run tools/cmd/eventhandlers/main.go

// Connect is the data for a Connect event.
// This is a synthetic event and is not dispatched by Discord.
type Connect struct{}

// Disconnect is the data for a Disconnect event.
// This is a synthetic event and is not dispatched by Discord.
type Disconnect struct{}

// RateLimit is the data for a RateLimit event.
// This is a synthetic event and is not dispatched by Discord.
type RateLimit struct {
	*TooManyRequests
	URL string
}

// Event provides a basic initial struct for all websocket events.
type Event struct {
	Operation int             `json:"op"`
	Sequence  int64           `json:"s"`
	Type      string          `json:"t"`
	RawData   json.RawMessage `json:"d"`
	// Struct contains one of the other types in this file.
	Struct interface{} `json:"-"`
}

// A Ready stores all data for the websocket READY event.
type Ready struct {
	Version         int        `json:"v"`
	User            *User      `json:"user"`
	PrivateChannels []*Channel `json:"private_channels"`
	Guilds          []*Guild   `json:"guilds"`
	SessionID       string     `json:"session_id"`
	Shard           []int      `json:"shard"`

	// Undocumented fields
	ReadState         []*ReadState         `json:"read_state"`
	Settings          *Settings            `json:"user_settings"`
	UserGuildSettings []*UserGuildSettings `json:"user_guild_settings"`
	Relationships     []*Relationship      `json:"relationships"`
	Presences         []*Presence          `json:"presences"`
	Notes             map[string]string    `json:"notes"`
}

// ChannelCreate is the data for a ChannelCreate event.
type ChannelCreate struct {
	*Channel
}

// ChannelUpdate is the data for a ChannelUpdate event.
type ChannelUpdate struct {
	*Channel
}

// ChannelDelete is the data for a ChannelDelete event.
type ChannelDelete struct {
	*Channel
}

// ChannelPinsUpdate stores data for a ChannelPinsUpdate event.
type ChannelPinsUpdate struct {
	GuildID          string `json:"guild_id,omitempty"`
	ChannelID        string `json:"channel_id"`
	LastPinTimestamp string `json:"last_pin_timestamp"`
}

// GuildCreate is the data for a GuildCreate event.
type GuildCreate struct {
	*Guild
}

// GuildUpdate is the data for a GuildUpdate event.
type GuildUpdate struct {
	*Guild
}

// GuildDelete is the data for a GuildDelete event.
type GuildDelete struct {
	*Guild
}

// GuildBanAdd is the data for a GuildBanAdd event.
type GuildBanAdd struct {
	GuildID string `json:"guild_id"`
	User    *User  `json:"user"`
}

// GuildBanRemove is the data for a GuildBanRemove event.
type GuildBanRemove struct {
	GuildID string `json:"guild_id"`
	User    *User  `json:"user"`
}

// GuildMemberAdd is the data for a GuildMemberAdd event.
type GuildMemberAdd struct {
	*Member
}

// GuildMemberUpdate is the data for a GuildMemberUpdate event.
type GuildMemberUpdate struct {
	*Member
}

// GuildMemberRemove is the data for a GuildMemberRemove event.
type GuildMemberRemove struct {
	*Member
}

// GuildRoleCreate is the data for a GuildRoleCreate event.
type GuildRoleCreate struct {
	*GuildRole
}

// GuildRoleUpdate is the data for a GuildRoleUpdate event.
type GuildRoleUpdate struct {
	*GuildRole
}

// A GuildRoleDelete is the data for a GuildRoleDelete event.
type GuildRoleDelete struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
}

// A GuildEmojisUpdate is the data for a guild emoji update event.
type GuildEmojisUpdate struct {
	GuildID string   `json:"guild_id"`
	Emojis  []*Emoji `json:"emojis"`
}

// A GuildMembersChunk is the data for a GuildMembersChunk event.
type GuildMembersChunk struct {
	GuildID    string      `json:"guild_id"`
	Members    []*Member   `json:"members"`
	ChunkIndex int         `json:"chunk_index"`
	ChunkCount int         `json:"chunk_count"`
	NotFound   []string    `json:"not_found"`
	Presences  []*Presence `json:"presences"`
	Nonce      string      `json:"nonce"`
}

// GuildIntegrationsUpdate is the data for a GuildIntegrationsUpdate event.
type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

// MessageAck is the data for a MessageAck event.
type MessageAck struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
}

// MessageCreate is the data for a MessageCreate event.
type MessageCreate struct {
	*Message
}

// MessageUpdate is the data for a MessageUpdate event.
type MessageUpdate struct {
	*Message
	// BeforeUpdate will be nil if the Message was not previously cached in the state cache.
	BeforeUpdate *Message `json:"-"`
}

// MessageDelete is the data for a MessageDelete event.
type MessageDelete struct {
	*Message
}

// MessageReactionAdd is the data for a MessageReactionAdd event.
type MessageReactionAdd struct {
	*MessageReaction
}

// MessageReactionRemove is the data for a MessageReactionRemove event.
type MessageReactionRemove struct {
	*MessageReaction
}

// MessageReactionRemoveAll is the data for a MessageReactionRemoveAll event.
type MessageReactionRemoveAll struct {
	*MessageReaction
}

// MessageReactionRemoveEmoji is the data for a MessageReactionRemoveEmoji event.
type MessageReactionRemoveEmoji struct {
	*MessageReaction
}

// PresencesReplace is the data for a PresencesReplace event.
type PresencesReplace []*Presence

// PresenceUpdate is the data for a PresenceUpdate event.
type PresenceUpdate struct {
	*Presence
}

// Resumed is the data for a Resumed event.
type Resumed struct {
	Trace []string `json:"_trace"`
}

// RelationshipAdd is the data for a RelationshipAdd event.
type RelationshipAdd struct {
	*Relationship
}

// RelationshipRemove is the data for a RelationshipRemove event.
type RelationshipRemove struct {
	*Relationship
}

// TypingStart is the data for a TypingStart event.
type TypingStart struct {
	ChannelID string  `json:"channel_id"`
	GuildID   string  `json:"guild_id,omitempty"`
	UserID    string  `json:"user_id"`
	Timestamp int     `json:"timestamp"`
	Member    *Member `json:"member"`
}

// UserUpdate is the data for a UserUpdate event.
type UserUpdate struct {
	*User
}

// UserSettingsUpdate is the data for a UserSettingsUpdate event.
type UserSettingsUpdate map[string]interface{}

// UserGuildSettingsUpdate is the data for a UserGuildSettingsUpdate event.
type UserGuildSettingsUpdate struct {
	*UserGuildSettings
}

// UserNoteUpdate is the data for a UserNoteUpdate event.
type UserNoteUpdate struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}

// VoiceServerUpdate is the data for a VoiceServerUpdate event.
type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

// VoiceStateUpdate is the data for a VoiceStateUpdate event.
type VoiceStateUpdate struct {
	*VoiceState
}

// MessageDeleteBulk is the data for a MessageDeleteBulk event
type MessageDeleteBulk struct {
	Messages  []string `json:"ids"`
	ChannelID string   `json:"channel_id"`
	GuildID   string   `json:"guild_id"`
}

// WebhooksUpdate is the data for a WebhooksUpdate event
type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// InviteCreate is the data for InviteCreate event
type InviteCreate struct {
	ChannelID  string    `json:"channel_id"`
	Code       string    `json:"code"`
	CreatedAt  Timestamp `json:"created_at"`
	GuildID    string    `json:"guild_id"`
	Inviter    *User     `json:"inviter"`
	MaxAge     int       `json:"max_age"`
	MaxUses    int       `json:"max_uses"`
	Target     *User     `json:"target_user"`
	TargetType int       `json:"target_user_type"`
	Temporary  bool      `json:"temporary"`
	Uses       int       `json:"uses"`
}

// InviteDelete is the data for InviteDelete event
type InviteDelete struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Code      string `json:"code"`
}

type ApplicationCommandCreate struct {
	*ApplicationCommand
	GuildID *string `json:"guild_id"`
}

type ApplicationCommandUpdate struct {
	*ApplicationCommandCreate
}

type ApplicationCommandDelete struct {
	*ApplicationCommandCreate
}

type IntegrationCreate struct {
	*Integration
	GuildID string `json:"guild_id"`
}

type IntegrationUpdate struct {
	*IntegrationCreate
}

type IntegrationDelete struct {
	*IntegrationCreate
	ID            string  `json:"id"`
	ApplicationID *string `json:"application_id"`
}

type ThreadCreate struct {
	*Channel
}

type ThreadUpdate struct {
	*Channel
}

type ThreadDelete struct {
	ID       string      `json:"id"`
	GuildID  string      `json:"guild_id"`
	ParentID string      `json:"parent_id"`
	Type     ChannelType `json:"type"`
}

type ThreadListSync struct {
	GuildID    string         `json:"guild_id"`
	ChannelIDs *[]string      `json:"channel_ids"`
	Threads    []Channel      `json:"threads"`
	Members    []ThreadMember `json:"members"`
}

type ThreadMemberUpdate struct {
	*ThreadMember
}

type ThreadMembersUpdate struct {
	ID       string      `json:"id"`
	GuildID  string      `json:"guild_id"`
	MemberCount int `json:"member_count"`
	AddedMembers *[]ThreadMember `json:"added_members"`
	RemovedMemberIDs *[]string `json:"removed_member_ids"`
}
