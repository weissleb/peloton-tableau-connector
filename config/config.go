package config

import "time"

const LogLevel = "INFO"

var Protocol = "http"
var Host = "localhost"
var Port = "8889"
var RequireAuth = true
var UseWorkoutCache = false
var CacheExpireSeconds = 30

var Banner = `
    ____       __      __            
   / __ \___  / /___  / /_____  ____
  / /_/ / _ \/ / __ \/ __/ __ \/ __ \
 / ____/  __/ / /_/ / /_/ /_/ / / / /
/_/    \___/_/\____/\__/\____/_/ /_/ 
`

var PeloPageLimit = 100
var PeloAllPages = true

var OutputFileName = "data_" + time.Now().UTC().Format("20060102150405") + ".csv"