package model

import "errors"

var ErrDateBusy = errors.New("данное время уже занято другим событием")
