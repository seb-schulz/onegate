package jwt

type UserJwtConverter interface {
	Subject() string
}

type anonymousUser struct{}

func (u anonymousUser) Subject() string {
	return "anon"
}

var AnonymousUser = anonymousUser{}
