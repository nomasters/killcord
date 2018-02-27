// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

// NotificationToMap converts a Notification into a map[string]interface{}
func NotificationToMap(notif *gnmi.Notification) (map[string]interface{}, error) {
	m := make(map[string]interface{}, 1)
	m["timestamp"] = notif.Timestamp
	m["path"] = StrPath(notif.Prefix)
	if len(notif.Update) != 0 {
		updates := make(map[string]interface{}, len(notif.Update))
		var err error
		for _, update := range notif.Update {
			updates[StrPath(update.Path)] = StrUpdateVal(update)
			if err != nil {
				return nil, err
			}
		}
		m["updates"] = updates
	}
	if len(notif.Delete) != 0 {
		deletes := make([]string, len(notif.Delete))
		for i, del := range notif.Delete {
			deletes[i] = StrPath(del)
		}
		m["deletes"] = deletes
	}
	return m, nil
}
