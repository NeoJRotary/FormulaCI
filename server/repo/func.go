package repo

import (
	"strings"

	D "github.com/NeoJRotary/describe-go"
	"github.com/google/uuid"
)

// UUID ...
func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// ResolveSrc ...
func ResolveSrc(src string) (hub string, user string, name string, err error) {
	s := D.String(src)
	if !s.HasSuffix(".git") {
		return "", "", "", D.NewErr("invalid git src")
	}

	if s.HasPrefix("https://") {
		ls := s.TrimPrefix("https://").TrimSuffix(".git").Split("/").Get()
		if len(ls) != 3 {
			err = D.NewErr("invalid length of src. Get", len(ls), "should be", 3)
		} else {
			hub = ls[0]
			user = ls[1]
			name = ls[2]
		}
	} else if s.HasPrefix("git@") {
		hub = s.RangeBetween("@", ":").Get()
		user = s.RangeBetween(":", "/").Get()
		name = s.RangeBetween("/", ".").Get()
	} else {
		err = D.NewErr("invalid git src")
	}

	return hub, user, name, err
}
