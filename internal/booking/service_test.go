package booking

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	store := NewMemoryStore()
	svc := NewService(store)

	const numGoroutines = 100_000

	var (
		successes atomic.Int64
		failures  atomic.Int64
		wg        sync.WaitGroup
	)

	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			err := svc.Book(Booking{
				MovieID: "screen-1",
				SeatID:  "A1",
				UserID:  uuid.New().String(),
			})
			if err == nil {
				successes.Add(1)
			} else if errors.Is(err, ErrSeatAlreadyBooked) {
				failures.Add(1)
			}
		}()
	}

	wg.Wait()

	if got := successes.Load(); got != 1 {
		t.Errorf("expected exactly 1 success, got %d", got)
	}
	if got := failures.Load(); got != int64(numGoroutines-1) {
		t.Errorf("expected %d ErrSeatAlreadyBooked failures, got %d", numGoroutines-1, got)
	}
}