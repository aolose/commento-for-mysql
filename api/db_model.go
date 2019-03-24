package main

import (
	"os"
	"time"
)
import "github.com/jinzhu/gorm"

func initDb() {
	if os.Getenv("DATABASE_TYPE") == "mysql" {
		db = db.Set("gorm:table_options", "CHARSET=r")
	}
	db.AutoMigrate(&Owners{})
	db.AutoMigrate(&Migrations{})
	db.AutoMigrate(&OwnerSessions{})
	db.AutoMigrate(&OwnerConfirmHexes{})
	db.AutoMigrate(&OwnerResetHexes{})
	db.AutoMigrate(&Domains{})
	db.AutoMigrate(&Moderators{})
	db.AutoMigrate(&Commenters{})
	db.AutoMigrate(&CommenterSessions{})
	db.AutoMigrate(&Comments{})
	db.AutoMigrate(&Votes{})
	db.AutoMigrate(&Views{})
	db.AutoMigrate(&Pages{})
	db.AutoMigrate(&Config{})
	db.AutoMigrate(&Exports{})
	db.AutoMigrate(&Emails{})

	var count int32
	db.Model(&Comments{}).Where("state = ?", "approved").Group("domain, path").Count(&count)
	logger.Infof("SET Page's comment count %d", count)
	db.Model(&Pages{}).Update("comment_count", count)
	db.Model(&Config{}).Update("version", "v1.6.0")
}

type Migrations struct {
	Filename string `gorm:"unique;not null"`
}

type Owners struct {
	OwnerHex       string    `gorm:"unique;not null;primary_key"`
	Email          string    `gorm:"unique;not null"`
	Name           string    `gorm:"unique;not null"`
	PasswordHash   string    `gorm:"not null"`
	ConfirmedEmail string    `gorm:"default:false;not null"`
	JoinDate       time.Time `gorm:"not null"`
}
type OwnerSessions struct {
	OwnerToken string    `gorm:"unique;not null;primary_key"`
	OwnerHex   string    `gorm:"not null"`
	LoginDate  time.Time `gorm:"not null"`
}

type OwnerConfirmHexes struct {
	ConfirmHex string `gorm:"unique;not null;primary_key"`
	OwnerHex   string `gorm:"not null"`
	SendDate   string `gorm:"not null"`
}

type OwnerResetHexes struct {
	ResetHex string `gorm:"unique;not null;primary_key"`
	OwnerHex string `gorm:"not null"`
	SendDate string `gorm:"not null"`
}

type Domains struct {
	Domain                  string    `gorm:"unique;not null;primary_key"`
	OwnerHex                string    `gorm:"not null"`
	Name                    string    `gorm:"not null"`
	CreationDate            time.Time `gorm:"not null"`
	State                   string    `gorm:"not null;default:'unfrozen'"`
	ImportedComments        string    `gorm:"not null;default:false"`
	AutoSpamFilter          bool      `gorm:"not null;default:true"`
	RequireModeration       bool      `gorm:"not null;default:false"`
	RequireIdentification   bool      `gorm:"not null;default:true"`
	ViewsThisMonth          int32     `gorm:"not null;default:0"`
	ModerateAllAnonymous    bool      `gorm:"default:true"`
	EmailNotificationPolicy string    `gorm:"default:'pending-moderation'"`
}

type Moderators struct {
	Domain  string    `gorm:"not null;primary_key"`
	Email   string    `gorm:"not null;primary_key"`
	AddDate time.Time `gorm:"not null"`
}

type Commenters struct {
	CommenterHex string    `gorm:"unique;not null;primary_key"`
	Email        string    `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Link         string    `gorm:"not null"`
	Photo        string    `gorm:"not null"`
	Provider     string    `gorm:"not null"`
	JoinDate     time.Time `gorm:"not null"`
	State        string    `gorm:"not null;default:'ok'"`
	PasswordHash string    `gorm:"not null;default:''"`
}

type CommenterSessions struct {
	CommenterToken string    `gorm:"unique;not null;primary_key"`
	CommenterHex   string    `gorm:"not null;default:'none'"`
	CreationDate   time.Time `gorm:"not null"`
}

type Comments struct {
	CommentHex   string    `gorm:"unique;not null;primary_key"`
	Domain       string    `gorm:"not null"`
	Path         string    `gorm:"not null"`
	CommenterHex string    `gorm:"not null"`
	Markdown     string    `gorm:"not null"`
	Html         string    `gorm:"not null"`
	ParentHex    string    `gorm:"not null"`
	Score        int32     `gorm:"not null;default:0"`
	State        string    `gorm:"not null;default:'unapproved'"`
	CreationDate time.Time `gorm:"not null"`
}

func (u *Comments) AfterCreate(tx *gorm.DB) (err error) {
	if u.Domain != "" {
		db.Model(Pages{}).Where("domain = ? AND path = ?", u.Domain, u.Path).
			Update("comment_count", gorm.Expr("comment_count + ?", 1))

	}
	return
}
func (u *Comments) AfterDelete(tx *gorm.DB) (err error) {
	if u.CommentHex != "" {
		db.Delete(Comments{}, "parent_hex = ?", u.CommentHex)
	}
	return
}

type Votes struct {
	CommentHex   string    `gorm:"unique;not null"`
	CommenterHex string    `gorm:"unique;not null"`
	Direction    int       `gorm:"not null"`
	VoteDate     time.Time `gorm:"not null"`
}

func (u *Votes) AfterCreate(tx *gorm.DB) (err error) {
	if u.CommentHex != "" {
		db.Model(Comments{}).Where("comment_hex = ?", u.CommentHex).
			Update("score", gorm.Expr("score + ?", u.Direction))
	}
	return
}
func (u *Votes) beforeUpdate(tx *gorm.DB) (err error) {
	if u.CommentHex != "" {
		db.Model(Comments{}).Where("comment_hex = ?", u.CommentHex).
			Update("score", gorm.Expr("score - ?", u.Direction))
	}
	return
}
func (u *Votes) AfterUpdate(tx *gorm.DB) (err error) {
	if u.CommentHex != "" {
		db.Model(Comments{}).Where("comment_hex = ?", u.CommentHex).
			Update("score", gorm.Expr("score + ?", u.Direction))

	}
	return
}

type Views struct {
	Domain       string    `gorm:"not null;index"`
	CommenterHex string    `gorm:"not null"`
	ViewDate     time.Time `gorm:"not null"`
}

func (u *Views) AfterCreate(tx *gorm.DB) (err error) {
	if u.Domain != "" {
		db.Model(Domains{}).Where("domain = ?", u.Domain).
			Update("viewsThisMonth", gorm.Expr("viewsThisMonth + ?", 1))
	}
	return
}

type Pages struct {
	Domain           string `gorm:"not null;unique_index"`
	Path             string `gorm:"not null"`
	IsLocked         bool   `gorm:"not null;default:false"`
	CommentCount     int32  `gorm:"not null;default:0"`
	StickyCommentHex string `gorm:"not null;default:'none'"`
	Title            string `gorm:"default:''"`
}

type Config struct {
	Version string `gorm:"not null"`
}

type Exports struct {
	ExportHex    string    `gorm:"unique;not null;primary_key"`
	BinData      []byte    `gorm:"not null"`
	Domain       string    `gorm:"not null"`
	CreationDate time.Time `gorm:"not null"`
}

//
type Emails struct {
	Email                      string    `gorm:"unique;not null;primary_key"`
	UnsubscribeSecretHex       string    `gorm:"not null;index"`
	LastEmailNotificationDate  time.Time `gorm:"not null"`
	PendingEmails              int       `gorm:"not null;default:0"`
	SendReplyNotifications     bool      `gorm:"not null;default:false"`
	SendModeratorNotifications bool      `gorm:"not null;default:true"`
}
