package helper

import "github.com/sirupsen/logrus"

// WrapCloser helper to wrap closer
func WrapCloser(close func() error) {
	if err := close(); err != nil {
		logrus.Error(err)
	}
}
