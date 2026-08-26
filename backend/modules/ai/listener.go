package ai

import (
	"fmt"

	"github.com/goravel/framework/contracts/event"

	"reflexcms/backend/app/facades"
)

// SummaryListener is the AI-module placeholder (plan §O8 / M2). It fires
// whenever an article is published; the real implementation will enqueue a
// model call via safefetch and write back articles.ai_summary/ai_keywords.
// Keeping it inert-by-design means the pipeline exists without any provider
// configured.
type SummaryListener struct{}

func NewSummaryListener() *SummaryListener { return &SummaryListener{} }

func (l *SummaryListener) Signature() string { return "ai.article-summary-placeholder" }

// Queue keeps the listener off any queue: with no provider it must stay a
// cheap, synchronous no-op.
func (l *SummaryListener) Queue(args ...any) event.Queue {
	return event.Queue{Enable: false}
}

func (l *SummaryListener) Handle(args ...any) error {
	var articleID uint64
	if len(args) > 0 {
		switch v := args[0].(type) {
		case uint64:
			articleID = v
		case int:
			articleID = uint64(v)
		}
	}
	facades.Log().Info(fmt.Sprintf(
		"[ai-placeholder] summary generation queued for article #%d "+
			"(no provider configured; wire modules/ai to enable)", articleID))
	return nil
}
