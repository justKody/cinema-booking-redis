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

// func TestBook_SameSeatTwice_SecondFails(t *testing.T) {
// 	svc := NewService(NewMemoryStore())

// 	if err := svc.Book(Booking{MovieID: "screen-1", SeatID: "A1", UserID: "user-1"}); err != nil {
// 		t.Fatalf("first booking: unexpected error: %v", err)
// 	}

// 	err := svc.Book(Booking{MovieID: "screen-1", SeatID: "A1", UserID: "user-2"})
// 	if !errors.Is(err, ErrSeatAlreadyBooked) {
// 		t.Errorf("second booking: expected ErrSeatAlreadyBooked, got %v", err)
// 	}
// }

// func TestBook_DifferentSeats_BothSucceed(t *testing.T) {
// 	svc := NewService(NewMemoryStore())

// 	for _, seat := range []string{"A1", "A2"} {
// 		if err := svc.Book(Booking{MovieID: "screen-1", SeatID: seat, UserID: "user-1"}); err != nil {
// 			t.Errorf("booking seat %s: unexpected error: %v", seat, err)
// 		}
// 	}
// }

// func TestBook_SameSeatDifferentMovies_BothSucceed(t *testing.T) {
// 	svc := NewService(NewMemoryStore())

// 	for _, movie := range []string{"screen-1", "screen-2"} {
// 		if err := svc.Book(Booking{MovieID: movie, SeatID: "A1", UserID: "user-1"}); err != nil {
// 			t.Errorf("booking A1 for %s: unexpected error: %v", movie, err)
// 		}
// 	}
// }

// func TestListBooking_FiltersByMovie(t *testing.T) {
// 	store := NewMemoryStore()
// 	svc := NewService(store)

// 	bookings := []Booking{
// 		{MovieID: "screen-1", SeatID: "A1", UserID: "user-1"},
// 		{MovieID: "screen-1", SeatID: "A2", UserID: "user-2"},
// 		{MovieID: "screen-2", SeatID: "B1", UserID: "user-3"},
// 	}
// 	for _, b := range bookings {
// 		if err := svc.Book(b); err != nil {
// 			t.Fatalf("booking %s/%s: unexpected error: %v", b.MovieID, b.SeatID, err)
// 		}
// 	}

// 	got := store.ListBooking("screen-1")
// 	if len(got) != 2 {
// 		t.Fatalf("expected 2 bookings for screen-1, got %d", len(got))
// 	}

// 	seats := map[string]bool{}
// 	for _, b := range got {
// 		if b.MovieID != "screen-1" {
// 			t.Errorf("got booking for movie %q, want screen-1", b.MovieID)
// 		}
// 		seats[b.SeatID] = true
// 	}
// 	if !seats["A1"] || !seats["A2"] {
// 		t.Errorf("expected seats A1 and A2, got %v", seats)
// 	}
// }

// func TestListBooking_UnknownMovie_ReturnsEmpty(t *testing.T) {
// 	store := NewMemoryStore()

// 	if err := store.Book(Booking{MovieID: "screen-1", SeatID: "A1", UserID: "user-1"}); err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if got := store.ListBooking("screen-99"); len(got) != 0 {
// 		t.Errorf("expected no bookings, got %d", len(got))
// 	}
// }
