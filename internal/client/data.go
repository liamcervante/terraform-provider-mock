package client

type CreateResource func() Data

type Data interface {
	GetId() string
	SetId(string)
}
