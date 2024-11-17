package modelx

import (
	"github.com/eddieowens/opts"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

type ModelOpts struct {
	CacheConf cache.CacheConf
	CacheOpts []cache.Option
}

func (opts ModelOpts) DefaultOptions() ModelOpts {
	return ModelOpts{}
}

func WithCacheConf(cacheConf cache.CacheConf) opts.Opt[ModelOpts] {
	return func(o *ModelOpts) {
		o.CacheConf = cacheConf
	}
}

func WithCacheOpts(cacheOpts ...cache.Option) opts.Opt[ModelOpts] {
	return func(o *ModelOpts) {
		o.CacheOpts = cacheOpts
	}
}
