package caches

import (
	"github.com/noble-gase/ne/helper"
	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

// Discard 丢弃数据，不缓存
const Discard = helper.NilError("caches: discarded")
