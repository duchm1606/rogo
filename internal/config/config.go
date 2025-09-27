package config

var Protocol = "tcp"
var Port = ":6380"
var Host = "localhost"
var MaxConnection = 100

var MaxKeyNumber int = 10
var EvictionRatio = 0.1

var EvictionPolicy string = "allkeys-random"

// var MaxWorkers = 10
// var MaxQueue = 100
// var MaxConnectionsPerIP = 10
// var MaxConnectionsPerIPPerSecond = 10
// var MaxConnectionsPerIPPerMinute = 100
// var MaxConnectionsPerIPPerHour = 1000
