package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          uint      `gorm:"primaryKey"`
	FirstName   string    `json:"name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email" gorm:"unique" diff:"-"`
	Telephone   string    `json:"telephone" gorm:"index:idx_phone_site_id,unique"`
	State       string    `json:"state"`
	City        string    `json:"city"`
	Area        string    `json:"area"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	LastSeen    time.Time `json:"last_seen"`
	Lng         float64   `json:"lng"`
	Lat         float64   `json:"lat"`

	Experiences string `gorm:"type:json" json:"experiences" diff:"-"` // Was Posts
	Map         string `gorm:"type:json" json:"map" diff:"-"`

	Password []byte `json:"-"`
	Site     uint32 `json:"site"`
	SiteId   string `json:"site_id" gorm:"index:idx_phone_site_id,unique"` // gorm:"index:,unique"` Cant register with this
	RoleId   uint16 `json:"role_id"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleId"`
	// Favorites  []UFavorites `json:"favorites"`
	// Likes      []UFavorites `json:"likes" gorm:"foreignKey:FavoritedID;references:ID"`
	// Prices     []UPrice     `json:"prices"`
	// Hours      []UHour      `json:"hours"`
	// Attributes []UAttribute `gorm:"many2many:user_attributes;" json:"attributes"`
	// Languages  []*ULanguage `gorm:"many2many:user_languages"`
	// Updates    []*UDiff     `diff:"-"`
	// Photos     []UPhotos    `json:"photos"`
	// Videos     []UVideos    `json:"videos"`
	DOB       time.Time `json:"dob"`
	MemberAt  time.Time `json:"member_at" diff:"-"`
	CreatedAt time.Time `diff:"-"`
	UpdatedAt time.Time `diff:"-"`

	// This is only used to send status to Websockets (No DB or Diff)
	Status string `gorm:"-" diff:"-" json:"status"`
	// gorm.Model
}

type UFavorites struct {
	ID          uint      `gorm:"primaryKey" diff:"-"`
	UserID      int       `json:"user_id" diff:"-"`
	FavoritedID uint      `json:"favorited_id" diff:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProfilesPage struct {
	ID        uint       `gorm:"primaryKey"`
	FirstName string     `json:"name"`
	State     string     `json:"state"`
	City      string     `json:"city"`
	Area      string     `json:"area"`
	Title     string     `json:"title"`
	LastSeen  time.Time  `json:"last_seen"`
	Photos    []*UPhotos `json:"photos" gorm:"foreignKey:UserID;references:ID"`
	Prices    []*UPrice  `json:"prices" gorm:"foreignKey:UserID;references:ID"`
	// Attributes []UAttribute `gorm:"many2many:user_attributes;foreignKey:user_id;references:id" json:"attributes"`
}

type UPrice struct {
	ID      uint   `gorm:"primaryKey;" diff:"-"`
	UserID  uint   `json:"user_id" diff:"-"`
	SiteID  string `json:"site_id" diff:"-"`
	Price   int64  `json:"price"`
	Time    string `json:"time"`
	Minutes int64  `json:"minutes" diff:"-" gorm:"-"`
	// Users   []*User `gorm:"many2many:user_prices;"`
	// gorm.Model
}

type UHour struct {
	ID     uint   `gorm:"primaryKey;" diff:"-"`
	UserID uint   `json:"user_id" diff:"-"`
	SiteID string `json:"site_id" diff:"-"`
	Days   string `json:"day"`
	Hours  string `json:"hour"`
	From   int    `json:"from" diff:"-" gorm:"-"`
	To     int    `json:"to" diff:"-" gorm:"-"`
	// Users []*User `gorm:"many2many:user_hours;"`
	// gorm.Model
}

type UPhotos struct {
	ID        uint      `gorm:"primaryKey" diff:"-"`
	UserID    uint      `json:"user_id" diff:"-"`
	SiteID    string    `json:"site_id" diff:"-"`
	Path      string    `json:"path" diff:"-" gorm:"index:idx_path_filename,unique"`
	Url       string    `json:"url" diff:"-"`
	FileName  string    `json:"filename" diff:"-" gorm:"index:idx_path_filename,unique"`
	SitePath  string    `json:"site_path" gorm:"index:,unique;"`
	Active    int8      `json:"active"`
	LinkAt    time.Time `diff:"-" json:"link_at"`
	CreatedAt time.Time `diff:"-"`
	// UpdatedAt time.Time `diff:"-"`
	// gorm.Model
}

type UVideos struct {
	ID        uint      `gorm:"primaryKey" diff:"-"`
	UserID    uint      `json:"user_id" diff:"-"`
	SiteID    string    `json:"site_id" gorm:"index:idx_siteid_sitepath,unique" diff:"-"`
	Path      string    `json:"path" diff:"-" gorm:"index:idx_path_filename,unique"`
	Url       string    `json:"url" diff:"-"`
	FileName  string    `json:"filename" diff:"-" gorm:"index:idx_path_filename,unique"`
	SitePath  string    `json:"site_path" gorm:"index:idx_siteid_sitepath,unique"`
	Active    int8      `json:"active"`
	LinkAt    time.Time `diff:"-" json:"link_at"`
	CreatedAt time.Time `diff:"-"`
	// UpdatedAt time.Time `diff:"-"`
	// gorm.Model
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) Name() string {
	return user.FirstName + " " + user.LastName
}

func (user *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&User{}).Count(&total)

	return total
}

func (user *User) Take(db *gorm.DB, limit int, offset int) interface{} {
	var users []User

	db.Preload("Role").Offset(offset).Limit(limit).Find(&users)

	return users
}
