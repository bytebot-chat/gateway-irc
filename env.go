package main

import (
	"flag"
	"os"
	"strings"
)

// This file contains environnement variables parsing related methods,
// for configuration purpose.

// parseEnv parses configuration environnement variables.
func parseEnv() {
	//FIXME this is very ugly
	if !isFlagSet("channels") {
		parseChannelsFromEnv()
	}

	if !isFlagSet("server") {
		*serv = parseStringFromEnv("BYTEBOT_SERVER", "localhost:6667")
	}

	if !isFlagSet("redis") {
		*redisAddr = parseStringFromEnv("BYTEBOT_REDIS", "localhost:6379")
	}

	if !isFlagSet("nick") {
		*nick = parseStringFromEnv("BYTEBOT_NICK", "bytebot")
	}

	if !isFlagSet("id") {
		*id = parseStringFromEnv("BYTEBOT_ID", "irc")
	}

	if !isFlagSet("inbound") {
		*inbound = parseStringFromEnv("BYTEBOT_INBOUND", "irc-inbound")
	}

	if !isFlagSet("outbound") {
		*outbound = parseStringFromEnv("BYTEBOT_OUTBOUND", *id)
	}

	if !isFlagSet("tls") {
		_, set := os.LookupEnv("BYTEBOT_TLS")
		if set {
			*tls = true
		}
	}

}

// Parses channels and sets them
func parseChannelsFromEnv() {
	val, set := os.LookupEnv("BYTEBOT_CHANNELS")
	if set {
		channels = strings.Split(val, ",")
	}
}

// Parses a string from an env variable and returns it.
func parseStringFromEnv(varName, defaultVal string) string {
	val, set := os.LookupEnv(varName)
	if set {
		return val
	}
	return defaultVal
}

// This is used to check if a flag was set
// Must be called after flag.Parse()
func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
