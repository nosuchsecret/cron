package cron

type JobMeta struct {
	Id string
	Force int //force run thing
	Data interface{}
}

