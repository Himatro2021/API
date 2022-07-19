package model

import "github.com/lib/pq"

// GroupMember :nodoc:
type GroupMember struct {
	ID      int64         `json:"id" gorm:"primaryKey"`
	Name    string        `json:"name" gorm:"unique"`
	Members pq.Int64Array `json:"members"`
}
