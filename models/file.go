package models

// File describes user model structure
type File struct {
	ID                  int        `db:"id" json:"id" goqu:"skipinsert"`
	Fullname            string     `db:"fullname" json:"fullname"`
	Name                string     `db:"name" json:"name"`
	Ext                 string     `db:"ext" json:"ext"`
	Path                string     `db:"path" json:"path"`
	Size                int64      `db:"size" json:"size"`
	Created             string     `db:"created" json:"created"`
	Updated             string     `db:"updated" json:"updated"`
}
