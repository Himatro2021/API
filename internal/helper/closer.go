package helper

import "github.com/sirupsen/logrus"

// WrapCloser helper to wrap closer
func WrapCloser(close func() error) {
	if err := close(); err != nil {
		logrus.Error(err)
	}
}

// LogIfErr will log with error level if received non nil error
func LogIfErr(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

// PanicIfErr will call logrus.Panic if received non nil error
func PanicIfErr(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}
