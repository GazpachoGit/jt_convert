package model

import "time"

type Model struct {
	CreationDate time.Time
	JTFileName   string
	PMIs         []PMI
}

type PMI struct {
	Label      string
	Type       string
	Attributes map[string]string
}
