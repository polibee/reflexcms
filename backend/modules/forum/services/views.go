package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"reflexcms/backend/app/facades"
)

const (
	pendingViewsKey = "views:topic:pending"
	viewBufferTTL   = 24 * time.Hour
)

func viewKey(topicID uint64) string {
	return "views:topic:" + strconv.FormatUint(topicID, 10)
}

func counterGet(key string) int64 {
	s := strings.TrimSpace(fmt.Sprint(facades.Cache().Get(key)))
	if s == "" || s == "<nil>" {
		return 0
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func counterPut(key string, n int64) error {
	return facades.Cache().Put(key, n, viewBufferTTL)
}

// IncrView buffers one view for the topic and registers it in the flush set.
func IncrView(topicID uint64) error {
	key := viewKey(topicID)
	if err := counterPut(key, counterGet(key)+1); err != nil {
		return err
	}

	raw := fmt.Sprint(facades.Cache().Get(pendingViewsKey))
	seen := map[string]bool{}
	for _, part := range strings.Fields(raw) {
		seen[part] = true
	}
	seen[strconv.FormatUint(topicID, 10)] = true

	parts := make([]string, 0, len(seen))
	for part := range seen {
		parts = append(parts, part)
	}
	return facades.Cache().Put(pendingViewsKey, strings.Join(parts, "\n"), viewBufferTTL)
}

// FlushViews drains buffered deltas into topics.view_count. Failures keep
// their buffer for the next round; succeeded ones are forgotten.
func FlushViews() error {
	parts := strings.Fields(fmt.Sprint(facades.Cache().Get(pendingViewsKey)))
	var remaining []string

	for _, part := range parts {
		id, err := strconv.ParseUint(part, 10, 64)
		if err != nil || id == 0 {
			continue
		}

		delta := counterGet(viewKey(id))
		if delta > 0 {
			if _, err := query().Exec(
				"UPDATE topics SET view_count = view_count + ? WHERE id = ?",
				delta, id,
			); err != nil {
				facades.Log().Warning("forum: flush views for topic " + part + ": " + err.Error())
				remaining = append(remaining, part)
				continue
			}
		}
		facades.Cache().Forget(viewKey(id))
	}

	if len(remaining) == 0 {
		facades.Cache().Forget(pendingViewsKey)
		return nil
	}
	return facades.Cache().Put(pendingViewsKey, strings.Join(remaining, "\n"), viewBufferTTL)
}
