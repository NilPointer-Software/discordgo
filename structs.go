// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains all structures for the discordgo package.  These
// may be moved about later into separate files but I find it easier to have
// them all located together.

package discordgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// A Session represents a connection to the Discord API.
type Session struct {
	sync.RWMutex

	// General configurable settings.

	// Authentication token for this session
	Token string
	MFA   bool

	// Debug for printing JSON request/responses
	Debug    bool // Deprecated, will be removed.
	LogLevel int

	// Should the session reconnect the websocket on errors.
	ShouldReconnectOnError bool

	// Should the session request compressed websocket data.
	Compress bool

	// Sharding
	ShardID    int
	ShardCount int

	// Should state tracking be enabled.
	// State tracking is the best way for getting the the users
	// active guilds and the members of the guilds.
	StateEnabled bool

	// Whether or not to call event handlers synchronously.
	// e.g false = launch event handlers in their own goroutines.
	SyncEvents bool

	// Exposed but should not be modified by User.

	// Whether the Data Websocket is ready
	DataReady bool // NOTE: Maye be deprecated soon

	// Max number of REST API retries
	MaxRestRetries int

	// Status stores the currect status of the websocket connection
	// this is being tested, may stay, may go away.
	status int32

	// Whether the Voice Websocket is ready
	VoiceReady bool // NOTE: Deprecated.

	// Whether the UDP Connection is ready
	UDPReady bool // NOTE: Deprecated

	// Stores a mapping of guild id's to VoiceConnections
	VoiceConnections map[string]*VoiceConnection

	// Managed state object, updated internally with events when
	// StateEnabled is true.
	State *State

	// The http client used for REST requests
	Client *http.Client

	// The user agent used for REST APIs
	UserAgent string

	// Stores the last HeartbeatAck that was recieved (in UTC)
	LastHeartbeatAck time.Time

	// Stores the last Heartbeat sent (in UTC)
	LastHeartbeatSent time.Time

	// used to deal with rate limits
	Ratelimiter *RateLimiter

	// Event handlers
	handlersMu   sync.RWMutex
	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance

	// The websocket connection.
	wsConn *websocket.Conn

	// When nil, the session is not listening.
	listening chan interface{}

	// sequence tracks the current gateway api websocket sequence number
	sequence *int64

	// stores sessions current Discord Gateway
	gateway string

	// stores session ID of current Gateway connection
	sessionID string

	// used to make sure gateway websocket writes do not happen concurrently
	wsMutex sync.Mutex

	// used to send intents to the gateway
	Intents int
}

// UserConnection is a Connection returned from the UserConnections endpoint
type UserConnection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked"`
	Integrations []*Integration `json:"integrations"`
	Verified     bool           `json:"verified"`
	FriendSync   bool           `json:"friend_sync"`
	ShowActivity bool           `json:"show_activity"`
	Visibility   int            `json:"visibility"` // 0 = None, 1 = Everyone
}

// Integration stores integration information
type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           bool               `json:"enabled"`
	Syncing           bool               `json:"syncing"`
	RoleID            string             `json:"role_id"`
	EnableEmoticons   bool               `json:"enable_emoticons"`
	ExpireBehavior    int                `json:"expire_behavior"` // 0 = remove role, 1 = kick
	ExpireGracePeriod int                `json:"expire_grace_period"`
	User              *User              `json:"user"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          Timestamp          `json:"synced_at"`
}

// IntegrationAccount is integration account information
// sent by the UserConnections endpoint
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// A VoiceRegion stores data for a specific voice region server.
type VoiceRegion struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	VIP        bool   `json:"vip"`
	Optimal    bool   `json:"optimal"`
	Deprecated bool   `json:"deprecated"`
	Custom     bool   `json:"custom"`
}

// A VoiceICE stores data for voice ICE servers.
type VoiceICE struct {
	TTL     string       `json:"ttl"`
	Servers []*ICEServer `json:"servers"`
}

// A ICEServer stores data for a specific voice ICE server.
type ICEServer struct {
	URL        string `json:"url"`
	Username   string `json:"username"`
	Credential string `json:"credential"`
}

// A Invite stores all data related to a specific Discord Guild or Channel invite.
type Invite struct {
	Code       string   `json:"code"`
	Guild      *Guild   `json:"guild"`
	Channel    *Channel `json:"channel"`
	Inviter    *User    `json:"inviter"`
	Target     *User    `json:"target_user"`
	TargetType int      `json:"target_user_type"`

	// will only be filled when using InviteWithCounts
	ApproximatePresenceCount int `json:"approximate_presence_count"`
	ApproximateMemberCount   int `json:"approximate_member_count"`

	// Metadata
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt Timestamp `json:"created_at"`
}

// ChannelType is the type of a Channel
type ChannelType int

// Block contains known ChannelType values
const (
	ChannelTypeGuildText ChannelType = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildNews
	ChannelTypeGuildStore
)

// A Channel holds all data related to an individual Discord channel.
type Channel struct {
	// The ID of the channel.
	ID string `json:"id"`

	// The type of the channel.
	Type ChannelType `json:"type"`

	// The ID of the guild to which the channel belongs, if it is in a guild.
	// Else, this ID is empty (e.g. DM channels).
	GuildID string `json:"guild_id"`

	// The position of the channel, used for sorting in client.
	Position int `json:"position"`

	// A list of permission overwrites present for the channel.
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`

	// The name of the channel.
	Name string `json:"name"`

	// The topic of the channel.
	Topic string `json:"topic"`

	// Whether the channel is marked as NSFW.
	NSFW bool `json:"nsfw"`

	// The ID of the last message sent in the channel. This is not
	// guaranteed to be an ID of a valid message.
	LastMessageID string `json:"last_message_id"`

	// The bitrate of the channel, if it is a voice channel.
	Bitrate int `json:"bitrate"`

	// The user limit of the voice channel.
	UserLimit int `json:"user_limit"`

	// Amount of seconds a user has to wait before sending another message (0-21600)
	// bots, as well as users with the permission manage_messages or manage_channel, are unaffected
	RateLimitPerUser int `json:"rate_limit_per_user"`

	// The recipients of the channel. This is only populated in DM channels.
	Recipients []*User `json:"recipients"`

	// Icon of the group DM channel.
	Icon string `json:"icon"`

	// ID of the DM channel creator
	OwnerID string `json:"owner_id"`

	// Application ID of DM group creator if it was bot-created
	ApplicationID string `json:"application_id"`

	// The ID of the parent channel, if the channel is under a category
	ParentID string `json:"parent_id"`

	// The timestamp of the last pinned message in the channel.
	// Empty if the channel has no pinned messages.
	LastPinTimestamp *Timestamp `json:"last_pin_timestamp"`

	// Non-Discord property
	// The messages in the channel. This is only present in state-cached channels,
	// and State.MaxMessageCount must be non-zero.
	Messages []*Message `json:"-"`
}

// Mention returns a string which mentions the channel
func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

// A ChannelEdit holds Channel Field data for a channel edit.
type ChannelEdit struct {
	Name                 string                 `json:"name,omitempty"`
	Type                 ChannelType            `json:"type,omitempty"` // Only on Guild with "NEWS" features and conversion between text and news
	Position             int                    `json:"position,omitempty"`
	Topic                string                 `json:"topic,omitempty"`
	NSFW                 *bool                  `json:"nsfw,omitempty"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
	Bitrate              int                    `json:"bitrate,omitempty"`
	UserLimit            int                    `json:"user_limit,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
}

// FollowChannel structure holds information from ChannelFollow function
type FollowChannel struct {
	ChannelID string `json:"channel_id"`
	WebHookID string `json:"webhook_id"`
}

// A PermissionOverwrite holds permission overwrite data for a Channel
type PermissionOverwrite struct {
	ID    string                  `json:"id"`
	Type  PermissionOverwriteType `json:"type"`
	Deny  int                     `json:"deny,string"`
	Allow int                     `json:"allow,string"`
}

type PermissionOverwriteType int

const (
	PermissionOverwriteTypeRole   PermissionOverwriteType = 0
	PermissionOverwriteTypeMember PermissionOverwriteType = 1
)

// Emoji struct holds data related to Emoji's
type Emoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	User          *User    `json:"user"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
	Available     bool     `json:"available"`
}

// MessageFormat returns a correctly formatted Emoji for use in Message content and embeds
func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

// APIName returns an correctly formatted API name for use in the MessageReactions endpoints.
func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

// VerificationLevel type definition
type VerificationLevel int

// Constants for VerificationLevel levels from 0 to 4 inclusive
const (
	VerificationLevelNone VerificationLevel = iota
	VerificationLevelLow
	VerificationLevelMedium
	VerificationLevelHigh
	VerificationLevelVeryHigh
)

// ExplicitContentFilterLevel type definition
type ExplicitContentFilterLevel int

// Constants for ExplicitContentFilterLevel levels from 0 to 2 inclusive
const (
	ExplicitContentFilterDisabled ExplicitContentFilterLevel = iota
	ExplicitContentFilterMembersWithoutRoles
	ExplicitContentFilterAllMembers
)

// MfaLevel type definition
type MfaLevel int

// Constants for MfaLevel levels from 0 to 1 inclusive
const (
	MfaLevelNone MfaLevel = iota
	MfaLevelElevated
)

// PremiumTier type definition
type PremiumTier int

// Constants for PremiumTier levels from 0 to 3 inclusive
const (
	PremiumTierNone PremiumTier = iota
	PremiumTier1
	PremiumTier2
	PremiumTier3
)

// A Guild holds all data related to a specific Discord Guild.  Guilds are also
// sometimes referred to as Servers in the Discord client.
type Guild struct {
	// The ID of the guild.
	ID string `json:"id"`

	// The name of the guild. (2–100 characters)
	Name string `json:"name"`

	// The hash of the guild's icon. Use Session.GuildIcon
	// to retrieve the icon itself.
	Icon string `json:"icon"`

	// The hash of the guild's splash.
	Splash string `json:"splash"`

	// The Hash of the guild's discovery splash. Only available for guilds with "DISCOVERABLE" feature
	DiscoverySplash string `json:"discovery_splash"`

	// The user ID of the owner of the guild.
	OwnerID string `json:"owner_id"`

	// The voice region of the guild.
	Region string `json:"region"`

	// The ID of the AFK voice channel.
	AfkChannelID string `json:"afk_channel_id"`

	// The timeout, in seconds, before a user is considered AFK in voice.
	AfkTimeout int `json:"afk_timeout"`

	// Whether or not the Server Widget is enabled
	WidgetEnabled bool `json:"widget_enabled"`

	// The Channel ID for the Server Widget
	WidgetChannelID string `json:"widget_channel_id"`

	// The verification level required for the guild.
	VerificationLevel VerificationLevel `json:"verification_level"`

	// The default message notification setting for the guild.
	// 0 == all messages, 1 == mentions only.
	DefaultMessageNotifications int `json:"default_message_notifications"`

	// The explicit content filter level
	ExplicitContentFilter ExplicitContentFilterLevel `json:"explicit_content_filter"`

	// A list of roles in the guild.
	Roles []*Role `json:"roles"`

	// A list of the custom emojis present in the guild.
	Emojis []*Emoji `json:"emojis"`

	// The list of enabled guild features
	Features []string `json:"features"`

	// Required MFA level for the guild
	MfaLevel MfaLevel `json:"mfa_level"`

	// Application ID of the Guild Creator if it is bot-created
	ApplicationID string `json:"application_id"`

	// The Channel ID to which system messages are sent (eg join and leave messages)
	SystemChannelID string `json:"system_channel_id"`

	// System Channel Flags
	SystemChannelFlags int `json:"system_channel_flags"`

	// Channel ID of the channel where guilds with "PUBLIC" feature can display rules and/or guidelines
	RulesChannelID string `json:"rules_channel_id"`

	// The time at which the current user joined the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	JoinedAt Timestamp `json:"joined_at"`

	// Whether the guild is considered large. This is
	// determined by a member threshold in the identify packet,
	// and is currently hard-coded at 250 members in the library.
	Large bool `json:"large"`

	// Whether this guild is currently unavailable (most likely due to outage).
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Unavailable bool `json:"unavailable"`

	// The number of members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	MemberCount int `json:"member_count"`

	// A list of voice states for the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	VoiceStates []*VoiceState `json:"voice_states"`

	// A list of the members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Members []*Member `json:"members"`

	// A list of channels in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Channels []*Channel `json:"channels"`

	// A list of partial presence objects for members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Presences []*Presence `json:"presences"`

	// The maximum number of presences for the guild (if nil, assume the default value, currently of 25000)
	MaxPresences *int `json:"max_presences"`

	// The maximum number of member fot the guild
	MaxMembers int `json:"max_members"`

	// the vanity url code for the guild
	VanityURLCode string `json:"vanity_url_code"`

	// the description for the guild
	Description string `json:"description"`

	// The hash of the guild's banner
	Banner string `json:"banner"`

	// The premium tier of the guild
	PremiumTier PremiumTier `json:"premium_tier"`

	// The total number of users currently boosting this server
	PremiumSubscriptionCount int `json:"premium_subscription_count"`

	// Preferred Local for guild with "PUBLIC" feature. Defaults to "en_US"
	PreferredLocale string `json:"preferred_locale"`

	// Channel ID of the channel where guilds with "PUBLIC" feature receive notices from Discord
	PublicUpdatesChannelID string `json:"public_updates_channel_id"`

	// The maximum amount of users in video channel
	MaxVideoChannelUsers int `json:"max_video_channel_users"`

	// Approximate member count
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate presence count
	ApproximatePresenceCount int `json:"approximate_presence_count"`
}

// IconURL returns a URL to the guild's icon.
func (g *Guild) IconURL() string {
	if g.Icon == "" {
		return ""
	}

	if strings.HasPrefix(g.Icon, "a_") {
		return EndpointGuildIconAnimated(g.ID, g.Icon)
	}

	return EndpointGuildIcon(g.ID, g.Icon)
}

// A UserGuild holds a brief version of a Guild
type UserGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int    `json:"permissions,string"`
}

// A GuildParams stores all the data needed to update discord guild settings
type GuildParams struct {
	Name                        string                     `json:"name,omitempty"`
	Region                      string                     `json:"region,omitempty"`
	VerificationLevel           *VerificationLevel         `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                        `json:"default_message_notifications,omitempty"` // TODO: Separate type?
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
	AfkChannelID                string                     `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                        `json:"afk_timeout,omitempty"`
	Icon                        string                     `json:"icon,omitempty"`
	OwnerID                     string                     `json:"owner_id,omitempty"`
	Splash                      string                     `json:"splash,omitempty"`
	Banner                      string                     `json:"banner,omitempty"`
	SystemChannelID             string                     `json:"system_channel_id,omitempty"`
	RulesChannelID              string                     `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      string                     `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             string                     `json:"preferred_locale,omitempty"`
}

// A short summary of a discoverable guild available to anyone
type GuildPreview struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Icon                     string   `json:"icon"`
	Splash                   string   `json:"splash"`
	DiscoverySplash          string   `json:"discovery_splash"`
	Emojis                   []*Emoji `json:"emojis"`
	Features                 []string `json:"features"`
	ApproximateMemberCount   int      `json:"approximate_member_count"`
	ApproximatePresenceCount int      `json:"approximate_presence_count"`
	Description              string   `json:"description"`
}

// A Role stores information about Discord guild member roles.
type Role struct {
	// The ID of the role.
	ID string `json:"id"`

	// The name of the role.
	Name string `json:"name"`

	// The hex color of this role.
	Color int `json:"color"`

	// Whether this role is hoisted (shows up separately in member list).
	Hoist bool `json:"hoist"`

	// The position of this role in the guild's role hierarchy.
	Position int `json:"position"`

	// The permissions of the role on the guild (doesn't include channel overrides).
	// This is a combination of bit masks; the presence of a certain permission can
	// be checked by performing a bitwise AND between this int and the permission.
	Permissions int `json:"permissions,string"`

	// Whether this role is managed by an integration, and
	// thus cannot be manually added to, or taken from, members.
	Managed bool `json:"managed"`

	// Whether this role is mentionable.
	Mentionable bool `json:"mentionable"`
}

// Mention returns a string which mentions the role
func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%s>", r.ID)
}

// Roles are a collection of Role
type Roles []*Role

func (r Roles) Len() int {
	return len(r)
}

func (r Roles) Less(i, j int) bool {
	return r[i].Position > r[j].Position
}

func (r Roles) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// A VoiceState stores the voice states of Guilds
type VoiceState struct {
	GuildID    string  `json:"guild_id"`
	ChannelID  string  `json:"channel_id"`
	UserID     string  `json:"user_id"`
	Member     *Member `json:"member"`
	SessionID  string  `json:"session_id"`
	Deaf       bool    `json:"deaf"`
	Mute       bool    `json:"mute"`
	SelfDeaf   bool    `json:"self_deaf"`
	SelfMute   bool    `json:"self_mute"`
	SelfStream bool    `json:"self_stream"`
	SelfVideo  bool    `json:"self_video"`
	Suppress   bool    `json:"suppress"`
}

// A Presence stores the online, offline, or idle and game status of Guild members.
type Presence struct {
	User         *User        `json:"user"`
	GuildID      string       `json:"guild_id"`
	Status       Status       `json:"status"`
	Activities   []*Activity  `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}

// ActivityType is the type of "activity" (see ActivityType* consts) in the Activity struct
type ActivityType int

// Valid ActivityType values
const (
	ActivityTypeGame ActivityType = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
	ActivityTypeCustom
	ActivityTypeCompeting
)

// A Activity struct holds the name of the "playing .." activity for a user
type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	URL           string       `json:"url,omitempty"`
	CreatedAt     int          `json:"created_at"`
	Timestamps    Timestamps   `json:"timestamps,omitempty"`
	ApplicationID string       `json:"application_id,omitempty"`
	Details       string       `json:"details,omitempty"`
	State         string       `json:"state,omitempty"`
	Emoji         GatewayEmoji `json:"emoji,omitempty"`
	Party         Party        `json:"party,omitempty"`
	Assets        Assets       `json:"assets,omitempty"`
	Secrets       Secrets      `json:"secrets,omitempty"`
	Instance      bool         `json:"instance,omitempty"`
	Flags         int          `json:"flags,omitempty"`
}

type GatewayEmoji struct {
	ID       *string `json:"id"`
	Name     *string `json:"name"`
	Animated *bool   `json:"animated,omitempty"`
}

// A ClientStatus struct holds information about the users platform
type ClientStatus struct {
	Desktop *Status `json:"desktop"`
	Mobile  *Status `json:"mobile"`
	Web     *Status `json:"web"`
}

// A TimeStamps struct contains start and end times used in the rich presence "playing .." Activity
type Timestamps struct {
	StartTimestamp int64 `json:"start,omitempty"`
	EndTimestamp   int64 `json:"end,omitempty"`
}

// UnmarshalJSON unmarshals JSON into TimeStamps struct
func (t *Timestamps) UnmarshalJSON(b []byte) error {
	temp := struct {
		End   float64 `json:"end,omitempty"`
		Start float64 `json:"start,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	t.EndTimestamp = int64(temp.End)
	t.StartTimestamp = int64(temp.Start)
	return nil
}

// An Assets struct contains assets and labels used in the rich presence "playing .." Activity
type Assets struct {
	LargeImageID string `json:"large_image,omitempty"`
	LargeText    string `json:"large_text,omitempty"`
	SmallImageID string `json:"small_image,omitempty"`
	SmallText    string `json:"small_text,omitempty"`
}

// An Party struct contains the id and size of a party used in rich presence Activity
type Party struct {
	ID   string `json:"id,omitempty"`
	Size []int  `json:"size,omitempty"`
}

// An Secrets struct contains secrets for rich presence Activity
type Secrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

// A Member stores user information for Guild members. A guild
// member represents a certain user's presence in a guild.
type Member struct {
	// The underlying user on which the member is based.
	User *User `json:"user"`

	// The nickname of the member, if they have one.
	Nick string `json:"nick"`

	// A list of IDs of the roles which are possessed by the member.
	Roles []string `json:"roles"`

	// The time at which the member joined the guild, in ISO8601.
	JoinedAt Timestamp `json:"joined_at"`

	// When the user used their Nitro boost on the server
	PremiumSince Timestamp `json:"premium_since"`

	// Whether the member is deafened at a guild level.
	Deaf bool `json:"deaf"`

	// Whether the member is muted at a guild level.
	Mute bool `json:"mute"`

	// The guild ID on which the member exists.
	// This property is not provided by Discord. I only left it here to not break the state.go
	GuildID string `json:"guild_id"`
}

// Mention creates a member mention
func (m *Member) Mention() string {
	return "<@!" + m.User.ID + ">"
}

// A Settings stores data for a specific users Discord client settings.
type Settings struct {
	RenderEmbeds           bool               `json:"render_embeds"`
	InlineEmbedMedia       bool               `json:"inline_embed_media"`
	InlineAttachmentMedia  bool               `json:"inline_attachment_media"`
	EnableTtsCommand       bool               `json:"enable_tts_command"`
	MessageDisplayCompact  bool               `json:"message_display_compact"`
	ShowCurrentGame        bool               `json:"show_current_game"`
	ConvertEmoticons       bool               `json:"convert_emoticons"`
	Locale                 string             `json:"locale"`
	Theme                  string             `json:"theme"`
	GuildPositions         []string           `json:"guild_positions"`
	RestrictedGuilds       []string           `json:"restricted_guilds"`
	FriendSourceFlags      *FriendSourceFlags `json:"friend_source_flags"`
	Status                 Status             `json:"status"`
	DetectPlatformAccounts bool               `json:"detect_platform_accounts"`
	DeveloperMode          bool               `json:"developer_mode"`
}

// Status type definition
type Status string

// Constants for Status with the different current available status
const (
	StatusOnline       Status = "online"
	StatusIdle         Status = "idle"
	StatusDoNotDisturb Status = "dnd"
	StatusInvisible    Status = "invisible"
	StatusOffline      Status = "offline"
)

// FriendSourceFlags stores ... TODO :)
type FriendSourceFlags struct {
	All           bool `json:"all"`
	MutualGuilds  bool `json:"mutual_guilds"`
	MutualFriends bool `json:"mutual_friends"`
}

// A Relationship between the logged in user and Relationship.User
type Relationship struct {
	User *User  `json:"user"`
	Type int    `json:"type"` // 1 = friend, 2 = blocked, 3 = incoming friend req, 4 = sent friend req
	ID   string `json:"id"`
}

// A TooManyRequests struct holds information received from Discord
// when receiving a HTTP 429 response.
type TooManyRequests struct {
	Bucket     string  `json:"bucket"`
	Message    string  `json:"message"`
	RetryAfter float64 `json:"retry_after"`
}

// A ReadState stores data on the read state of channels.
type ReadState struct {
	MentionCount  int    `json:"mention_count"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

// An Ack is used to ack messages
type Ack struct {
	Token string `json:"token"`
}

// A GuildRole stores data for guild roles.
type GuildRole struct {
	GuildID string `json:"guild_id"`
	Role    *Role  `json:"role"`
}

// A GuildBan stores data for a guild ban.
type GuildBan struct {
	Reason string `json:"reason"`
	User   *User  `json:"user"`
}

// A GuildWidget stores data for a guild widget.
type GuildWidget struct {
	Enabled   bool   `json:"enabled"`
	ChannelID string `json:"channel_id"`
}

// A GuildAuditLog stores data for a guild audit log.
type GuildAuditLog struct {
	Webhooks        []*Webhook       `json:"webhooks,omitempty"`
	Users           []*User          `json:"users,omitempty"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
	Integrations    []*Integration   `json:"integrations"`
}

type AuditLogEntry struct {
	TargetID   string             `json:"target_id"`
	Changes    []Change           `json:"changes,omitempty"`
	UserID     string             `json:"user_id"`
	ID         string             `json:"id"`
	ActionType AuditLogActionType `json:"action_type"`
	Options    AuditOptions       `json:"options,omitempty"`
	Reason     string             `json:"reason"`
}

type Change struct {
	NewValue interface{} `json:"new_value"`
	OldValue interface{} `json:"old_value"`
	Key      string      `json:"key"`
}

type AuditOptions struct {
	DeleteMembersDay string `json:"delete_member_days"`
	MembersRemoved   string `json:"members_removed"`
	ChannelID        string `json:"channel_id"`
	MessageID        string `json:"message_id"`
	Count            string `json:"count"`
	ID               string `json:"id"`
	Type             string `json:"type"` // 0 = role, 1 = member
	RoleName         string `json:"role_name"`
}

type AuditLogActionType int

// Block contains Discord Audit Log Action Types
const (
	AuditLogActionGuildUpdate AuditLogActionType = 1

	AuditLogActionChannelCreate          AuditLogActionType = 10
	AuditLogActionChannelUpdate          AuditLogActionType = 11
	AuditLogActionChannelDelete          AuditLogActionType = 12
	AuditLogActionChannelOverwriteCreate AuditLogActionType = 13
	AuditLogActionChannelOverwriteUpdate AuditLogActionType = 14
	AuditLogActionChannelOverwriteDelete AuditLogActionType = 15

	AuditLogActionMemberKick       AuditLogActionType = 20
	AuditLogActionMemberPrune      AuditLogActionType = 21
	AuditLogActionMemberBanAdd     AuditLogActionType = 22
	AuditLogActionMemberBanRemove  AuditLogActionType = 23
	AuditLogActionMemberUpdate     AuditLogActionType = 24
	AuditLogActionMemberRoleUpdate AuditLogActionType = 25
	AuditLogActionMemberMove       AuditLogActionType = 26
	AuditLogActionMemberDisconnect AuditLogActionType = 27

	AuditLogActionBotAdd AuditLogActionType = 28

	AuditLogActionRoleCreate AuditLogActionType = 30
	AuditLogActionRoleUpdate AuditLogActionType = 31
	AuditLogActionRoleDelete AuditLogActionType = 32

	AuditLogActionInviteCreate AuditLogActionType = 40
	AuditLogActionInviteUpdate AuditLogActionType = 41
	AuditLogActionInviteDelete AuditLogActionType = 42

	AuditLogActionWebhookCreate AuditLogActionType = 50
	AuditLogActionWebhookUpdate AuditLogActionType = 51
	AuditLogActionWebhookDelete AuditLogActionType = 52

	AuditLogActionEmojiCreate AuditLogActionType = 60
	AuditLogActionEmojiUpdate AuditLogActionType = 61
	AuditLogActionEmojiDelete AuditLogActionType = 62

	AuditLogActionMessageDelete     AuditLogActionType = 72
	AuditLogActionMessageBulkDelete AuditLogActionType = 73
	AuditLogActionMessagePin        AuditLogActionType = 74
	AuditLogActionMessageUnpin      AuditLogActionType = 75

	AuditLogActionIntegrationCreate AuditLogActionType = 80
	AuditLogActionIntegrationUpdate AuditLogActionType = 81
	AuditLogActionIntegrationDelete AuditLogActionType = 82
)

// A UserGuildSettingsChannelOverride stores data for a channel override for a users guild settings.
type UserGuildSettingsChannelOverride struct {
	Muted                bool   `json:"muted"`
	MessageNotifications int    `json:"message_notifications"`
	ChannelID            string `json:"channel_id"`
}

// A UserGuildSettings stores data for a users guild settings.
type UserGuildSettings struct {
	SupressEveryone      bool                                `json:"suppress_everyone"`
	Muted                bool                                `json:"muted"`
	MobilePush           bool                                `json:"mobile_push"`
	MessageNotifications int                                 `json:"message_notifications"`
	GuildID              string                              `json:"guild_id"`
	ChannelOverrides     []*UserGuildSettingsChannelOverride `json:"channel_overrides"`
}

// A UserGuildSettingsEdit stores data for editing UserGuildSettings
type UserGuildSettingsEdit struct {
	SupressEveryone      bool                                         `json:"suppress_everyone"`
	Muted                bool                                         `json:"muted"`
	MobilePush           bool                                         `json:"mobile_push"`
	MessageNotifications int                                          `json:"message_notifications"`
	ChannelOverrides     map[string]*UserGuildSettingsChannelOverride `json:"channel_overrides"`
}

// An APIErrorMessage is an api error message returned from discord
type APIErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Webhook stores the data for a webhook.
type Webhook struct {
	ID            string `json:"id"`
	Type          int    `json:"type"`
	GuildID       string `json:"guild_id"`
	ChannelID     string `json:"channel_id"`
	User          *User  `json:"user"`
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	Token         string `json:"token"`
	ApplicationID string `json:"application_id"`
}

// WebhookParams is a struct for webhook params, used in the WebhookExecute command.
type WebhookParams struct {
	Content         string          `json:"content,omitempty"`
	Username        string          `json:"username,omitempty"`
	AvatarURL       string          `json:"avatar_url,omitempty"`
	TTS             bool            `json:"tts,omitempty"`
	File            string          `json:"file,omitempty"`
	Embeds          []*MessageEmbed `json:"embeds,omitempty"`
	AllowedMentions AllowMention    `json:"allowed_mentions"`
}

type AllowMention struct {
	Parse []string `json:"parse"`
	Roles []string `json:"roles"`
	Users []string `json:"users"`
}

// MessageReaction stores the data for a message reaction.
type MessageReaction struct {
	UserID    string  `json:"user_id"`
	ChannelID string  `json:"channel_id"`
	MessageID string  `json:"message_id"`
	GuildID   string  `json:"guild_id,omitempty"`
	Member    *Member `json:"member"`
	Emoji     *Emoji  `json:"emoji"`
}

// GatewayBotResponse stores the data for the gateway/bot response
type GatewayBotResponse struct {
	URL               string            `json:"url"`
	Shards            int               `json:"shards"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

type SessionStartLimit struct {
	Total      int   `json:"total"`
	Remaining  int   `json:"remaining"`
	ResetAfter int64 `json:"reset_after"`
}

type ApplicationCommand struct {
	ID            string                      `json:"id,omitempty"`
	ApplicationID string                      `json:"application_id,omitempty"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Options       *[]ApplicationCommandOption `json:"options,omitempty"`
}

type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType      `json:"type"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Default     bool                              `json:"default,omitempty"`
	Required    bool                              `json:"required,omitempty"`
	Choices     *[]ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options     *[]ApplicationCommandOption       `json:"options,omitempty"`
}

type ApplicationCommandOptionType int

type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`
}

type Interaction struct {
	ID        string                             `json:"id"`
	Type      InteractionType                    `json:"type"`
	Data      *ApplicationCommandInteractionData `json:"data"`
	GuildID   string                             `json:"guild_id"`
	ChannelID string                             `json:"channel_id"`
	Member    *Member                            `json:"member"`
	User      *User                              `json:"user"`
	Token     string                             `json:"token"`
	Version   int                                `json:"version"`
}

type InteractionType int

const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
)

type ApplicationCommandInteractionData struct {
	ID      string                                     `json:"id"`
	Name    string                                     `json:"name"`
	Options *[]ApplicationCommandInteractionDataOption `json:"options"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Value   interface{}                                `json:"value"` // =/ Oh god, this will be VERY FUNNY to handle ~MBSA
	Options *[]ApplicationCommandInteractionDataOption `json:"options"`
}

type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}

type InteractionResponseType int

const (
	InteractionResponseTypePong InteractionResponseType     = 1
	InteractionResponseTypeChannelMessageWithSource         = 4
	InteractionResponseTypeDeferredChannelMessageWithSource = 5
)

type InteractionApplicationCommandCallbackData struct {
	TTS             bool            `json:"tts,omitempty"`
	Content         string          `json:"content"`
	Embeds          *[]MessageEmbed `json:"embeds,omitempty"`
	AllowedMentions *AllowMention   `json:"allowed_mentions,omitempty"` // TODO: Fix AllowMention in the fork
	Flags           int             `json:"flags,omitempty"`
}

const (
	OptionTypeSubCommand ApplicationCommandOptionType = iota + 1
	OptionTypeSubCommandGroup
	OptionTypeString
	OptionTypeInteger
	OptionTypeBoolean
	OptionTypeUser
	OptionTypeChannel
	OptionTypeRole
)

// Constants for the different bit offsets of text channel permissions
const (
	PermissionReadMessages = 1 << (iota + 10)
	PermissionSendMessages
	PermissionSendTTSMessages
	PermissionManageMessages
	PermissionEmbedLinks
	PermissionAttachFiles
	PermissionReadMessageHistory
	PermissionMentionEveryone
	PermissionUseExternalEmojis
)

// Constants for the different bit offsets of voice permissions
const (
	PermissionVoiceConnect = 1 << (iota + 20)
	PermissionVoiceSpeak
	PermissionVoiceMuteMembers
	PermissionVoiceDeafenMembers
	PermissionVoiceMoveMembers
	PermissionVoiceUseVAD
	PermissionVoicePrioritySpeaker = 1 << (iota + 2)
)

// Constants for general management.
const (
	PermissionChangeNickname = 1 << (iota + 26)
	PermissionManageNicknames
	PermissionManageRoles
	PermissionManageWebhooks
	PermissionManageEmojis
)

// Constants for the different bit offsets of general permissions
const (
	PermissionCreateInstantInvite = 1 << iota
	PermissionKickMembers
	PermissionBanMembers
	PermissionAdministrator
	PermissionManageChannels
	PermissionManageServer
	PermissionAddReactions
	PermissionViewAuditLogs

	PermissionAllText = PermissionReadMessages |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone
	PermissionAllVoice = PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker
	PermissionAllChannel = PermissionAllText |
		PermissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	PermissionAll = PermissionAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis
)

// Block contains Discord JSON Error Response codes
const (
	ErrCodeUnknownAccount     = 10001
	ErrCodeUnknownApplication = 10002
	ErrCodeUnknownChannel     = 10003
	ErrCodeUnknownGuild       = 10004
	ErrCodeUnknownIntegration = 10005
	ErrCodeUnknownInvite      = 10006
	ErrCodeUnknownMember      = 10007
	ErrCodeUnknownMessage     = 10008
	ErrCodeUnknownOverwrite   = 10009
	ErrCodeUnknownProvider    = 10010
	ErrCodeUnknownRole        = 10011
	ErrCodeUnknownToken       = 10012
	ErrCodeUnknownUser        = 10013
	ErrCodeUnknownEmoji       = 10014
	ErrCodeUnknownWebhook     = 10015

	ErrCodeBotsCannotUseEndpoint  = 20001
	ErrCodeOnlyBotsCanUseEndpoint = 20002

	ErrCodeMaximumGuildsReached     = 30001
	ErrCodeMaximumFriendsReached    = 30002
	ErrCodeMaximumPinsReached       = 30003
	ErrCodeMaximumGuildRolesReached = 30005
	ErrCodeTooManyReactions         = 30010

	ErrCodeUnauthorized = 40001

	ErrCodeMissingAccess                             = 50001
	ErrCodeInvalidAccountType                        = 50002
	ErrCodeCannotExecuteActionOnDMChannel            = 50003
	ErrCodeEmbedDisabled                             = 50004
	ErrCodeCannotEditFromAnotherUser                 = 50005
	ErrCodeCannotSendEmptyMessage                    = 50006
	ErrCodeCannotSendMessagesToThisUser              = 50007
	ErrCodeCannotSendMessagesInVoiceChannel          = 50008
	ErrCodeChannelVerificationLevelTooHigh           = 50009
	ErrCodeOAuth2ApplicationDoesNotHaveBot           = 50010
	ErrCodeOAuth2ApplicationLimitReached             = 50011
	ErrCodeInvalidOAuthState                         = 50012
	ErrCodeMissingPermissions                        = 50013
	ErrCodeInvalidAuthenticationToken                = 50014
	ErrCodeNoteTooLong                               = 50015
	ErrCodeTooFewOrTooManyMessagesToDelete           = 50016
	ErrCodeCanOnlyPinMessageToOriginatingChannel     = 50019
	ErrCodeCannotExecuteActionOnSystemMessage        = 50021
	ErrCodeMessageProvidedTooOldForBulkDelete        = 50034
	ErrCodeInvalidFormBody                           = 50035
	ErrCodeInviteAcceptedToGuildApplicationsBotNotIn = 50036

	ErrCodeReactionBlocked = 90001
)
