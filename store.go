package main

import (
	"sync"
	"time"
)

var urlStore = make(map[string]ShortURL) //in memory storage of short_url generated
var cacheMutex sync.RWMutex              //mutex var for inmemory storage

var clickCounter = make(map[string]int) //counter to count number of times short_url is fetched
var clickCounterMutex sync.RWMutex      //mutex var for counter map

var requestLog = make(map[string][]time.Time) //in memory storage to check rate limit of short_url fetch count
var rateLimitMutex sync.RWMutex               // mutex var for rate limit var

var eventClick = make(chan ClickEvent, 100) //channel to receive short_url fetch count
