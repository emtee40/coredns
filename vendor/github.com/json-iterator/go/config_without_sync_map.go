//+build !go1.9

package jsoniter

import (
	"sync"
)

type frozenConfig struct {
	configBeforeFrozen            Config
	sortMapKeys                   bool
	indentionStep                 int
	objectFieldMustBeSimpleString bool
	onlyTaggedField               bool
	disallowUnknownFields         bool
	cacheLock                     *sync.RWMutex
	decoderCache                  map[uintptr]ValDecoder
	encoderCache                  map[uintptr]ValEncoder
	extensions                    []Extension
	streamPool                    *sync.Pool
	iteratorPool                  *sync.Pool
}

func (cfg *frozenConfig) initCache() {
	cfg.cacheLock = &sync.RWMutex{}
	cfg.decoderCache = map[uintptr]ValDecoder{}
	cfg.encoderCache = map[uintptr]ValEncoder{}
}

func (cfg *frozenConfig) addDecoderToCache(cacheKey uintptr, decoder ValDecoder) {
	cfg.cacheLock.Lock()
	cfg.decoderCache[cacheKey] = decoder
	cfg.cacheLock.Unlock()
}

func (cfg *frozenConfig) addEncoderToCache(cacheKey uintptr, encoder ValEncoder) {
	cfg.cacheLock.Lock()
	cfg.encoderCache[cacheKey] = encoder
	cfg.cacheLock.Unlock()
}

func (cfg *frozenConfig) getDecoderFromCache(cacheKey uintptr) ValDecoder {
	cfg.cacheLock.RLock()
	decoder, _ := cfg.decoderCache[cacheKey].(ValDecoder)
	cfg.cacheLock.RUnlock()
	return decoder
}

func (cfg *frozenConfig) getEncoderFromCache(cacheKey uintptr) ValEncoder {
	cfg.cacheLock.RLock()
	encoder, _ := cfg.encoderCache[cacheKey].(ValEncoder)
	cfg.cacheLock.RUnlock()
	return encoder
}

var cfgCacheLock = &sync.RWMutex{}
var cfgCache = map[Config]*frozenConfig{}

func getFrozenConfigFromCache(cfg Config) *frozenConfig {
	cfgCacheLock.RLock()
	frozenConfig := cfgCache[cfg]
	cfgCacheLock.RUnlock()
	return frozenConfig
}

func addFrozenConfigToCache(cfg Config, frozenConfig *frozenConfig) {
	cfgCacheLock.Lock()
	cfgCache[cfg] = frozenConfig
	cfgCacheLock.Unlock()
}
