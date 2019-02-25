package monitor

import (
	"log"
)

// RecentChangeHandler handles recent changes as they arrive
type RecentChangeHandler struct {
	changesQueue []RecentChange
	logger       *log.Logger
	fetcher      RevisionFetcher
}

// NewRecentChangeHandler creates a handler for recent changes
func NewRecentChangeHandler(fetcher RevisionFetcher, logger *log.Logger) RecentChangeHandler {
	return RecentChangeHandler{
		fetcher: fetcher,
		logger:  logger,
	}
}

func (m RecentChangeHandler) Handle(rc RecentChange) {
	// s, _ := json.MarshalIndent(rc, "", "\t")
	// m.logger.Println(string(s))

	if rc.LogAction == LogActionDelete {
		m.logger.Println("action delete noted")
		// RemovePageRevisionsFromQueue(rc.Title)
	}

	if rc.Revision.New == nil {
		m.logger.Println("")
		return
	}

	new := *rc.Revision.New
	m.fetcher.Queue(new, m.handleFetchResponse)
}

func (m RecentChangeHandler) handleFetchResponse(body []byte, err error) {
	if err != nil {
		m.logger.Println("Fetcher received error: %+v", err)
		return
	}

	m.logger.Println("handleFetchResponse success")
	// m.logger.Println(string(body))
}

// RemovePageRevisionsFromQueue removes the given page revisions from the queue
func (m RecentChangeHandler) RemovePageRevisionsFromQueue(title string) {
	// changes := []RecentChange
	// for i, r := range m.changesQueue {
	// 	if r.Title != title {
	// 		changes := append(changes, r)
	// 	}
	// }
	// m.changesQueue = changes
}
