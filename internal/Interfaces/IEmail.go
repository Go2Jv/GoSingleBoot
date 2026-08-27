package Interfaces

type IEmail interface {
	Send(who string) error
}
