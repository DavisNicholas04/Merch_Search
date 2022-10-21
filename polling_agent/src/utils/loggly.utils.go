package utils

import "github.com/jamespearly/loggly"

func InstantiateClient(tag string) *loggly.ClientType {
	return loggly.New(tag)
}
