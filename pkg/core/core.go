package core

type Sender interface {
	SendFile(filename string) bool
}
