package action

import (
	"net/http"

	"github.com/lbryio/notifica/app/store"

	"github.com/lbryio/notifica/app/types"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

// Root Handler is the default handler
func Root(r *http.Request) api.Response {
	if r.URL.Path == "/" {
		return api.Response{Data: "Welcome to Rick Reports!"}
	}
	return api.Response{Status: http.StatusNotFound, Error: errors.Err("404 Not Found")}
}

// Test only used when in dev mode
func Test(r *http.Request) api.Response {
	reply := &types.Notification{
		Type:       &types.Notification_Comment{Comment: &types.Comment{}},
		Name:       "my notification",
		Title:      "This is my title",
		Text:       "please reply to my comment asap",
		EmailData:  nil,
		ClaimData:  nil,
		DeviceData: nil,
		AppData:    nil,
	}
	reply.GetType()
	bucketHash := "5uJsckhx7art7i4tRdVkHqzqkYfUP8AK" //crypto.RandString(32)
	s := store.Retrieve(bucketHash)
	err := s.Add(reply)
	if err != nil {
		return api.Response{Error: errors.Err(err)}
	}

	return api.Response{Data: "ok"}
}

// Status is the default handler
func Status(r *http.Request) api.Response {
	bucketHash := "5uJsckhx7art7i4tRdVkHqzqkYfUP8AK" //crypto.RandString(32)
	s := store.Retrieve(bucketHash)
	var err error

	notification, err := s.Next()
	if err != nil {
		return api.Response{Error: errors.Err(err)}
	}
	notifications := []*types.Notification{notification}
	for notification != nil {
		notification, err = s.Next()
		if err != nil {
			return api.Response{Error: errors.Err(err)}
		}
		notifications = append(notifications, notification)
	}
	return api.Response{Data: notifications}
}
