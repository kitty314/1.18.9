package profile

import (
	"github.com/kitty314/1.18.9/common/atomic"
)

// StoreSelected is a global switch for storing selected proxy to cache
var StoreSelected = atomic.NewBool(true)
