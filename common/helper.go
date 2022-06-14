package common

import "go.mongodb.org/mongo-driver/bson/primitive"

const TokenClaimsKey = "token-claims"

func Int64P(v int64) *int64 {
	return &v
}

func BoolP(v bool) *bool {
	return &v
}

func ObjectIDP(v interface{}) *primitive.ObjectID {
	if id, ok := v.(primitive.ObjectID); ok {
		return &id
	}
	return nil
}
