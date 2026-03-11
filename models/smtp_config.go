package models

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	FromName string
	From     string
	Subject  string
}
