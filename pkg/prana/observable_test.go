package prana

import (
	"sync"
	"testing"
	"time"
)

func TestObservable_PanicRecovery(t *testing.T) {
	obs := NewObservable(0)
	
	panicWatcher := func(old, new int) {
		panic("watcher panic!")
	}
	
	var normalCalled bool
	var mu sync.Mutex
	
	// Register panic watcher first
	obs.Watch(panicWatcher)
	
	// Register normal watcher second
	obs.Watch(func(old, new int) {
		mu.Lock()
		normalCalled = true
		mu.Unlock()
	})
	
	// Update should not panic and should call all watchers
	obs.Set(42)
	
	// Give a bit of time for watchers to complete
	time.Sleep(10 * time.Millisecond)
	
	// Normal watcher should have been called despite panic in first watcher
	mu.Lock()
	if !normalCalled {
		t.Error("Normal watcher was not called after panic in another watcher")
	}
	mu.Unlock()
}

func TestObservable_GetSet(t *testing.T) {
	obs := NewObservable(10)
	
	if obs.Get() != 10 {
		t.Errorf("Get() = %v, want 10", obs.Get())
	}
	
	var gotOld, gotNew int
	obs.Watch(func(old, new int) {
		gotOld = old
		gotNew = new
	})
	
	obs.Set(20)
	
	time.Sleep(10 * time.Millisecond)
	
	if gotOld != 10 {
		t.Errorf("Watcher received old = %v, want 10", gotOld)
	}
	if gotNew != 20 {
		t.Errorf("Watcher received new = %v, want 20", gotNew)
	}
	if obs.Get() != 20 {
		t.Errorf("Get() = %v, want 20", obs.Get())
	}
}

func TestObservable_Update(t *testing.T) {
	obs := NewObservable(5)
	
	obs.Update(func(v int) int {
		return v * 2
	})
	
	if obs.Get() != 10 {
		t.Errorf("After Update, Get() = %v, want 10", obs.Get())
	}
}

func TestObservable_Watch(t *testing.T) {
	obs := NewObservable(0)
	
	called := 0
	unwatch := obs.Watch(func(old, new int) {
		called++
	})
	
	obs.Set(1)
	time.Sleep(10 * time.Millisecond)
	
	if called != 1 {
		t.Errorf("Watcher called %d times, want 1", called)
	}
	
	// Unwatch
	unwatch()
	
	obs.Set(2)
	time.Sleep(10 * time.Millisecond)
	
	if called != 1 {
		t.Errorf("Watcher called %d times after unwatch, want 1", called)
	}
}

func TestObservable_WatcherCount(t *testing.T) {
	obs := NewObservable(0)
	
	if obs.WatcherCount() != 0 {
		t.Errorf("Initial WatcherCount() = %v, want 0", obs.WatcherCount())
	}
	
	unwatch1 := obs.Watch(func(old, new int) {})
	if obs.WatcherCount() != 1 {
		t.Errorf("WatcherCount() = %v, want 1", obs.WatcherCount())
	}
	
	unwatch2 := obs.Watch(func(old, new int) {})
	if obs.WatcherCount() != 2 {
		t.Errorf("WatcherCount() = %v, want 2", obs.WatcherCount())
	}
	
	unwatch1()
	if obs.WatcherCount() != 1 {
		t.Errorf("After unwatch, WatcherCount() = %v, want 1", obs.WatcherCount())
	}
	
	unwatch2()
	if obs.WatcherCount() != 0 {
		t.Errorf("After all unwatch, WatcherCount() = %v, want 0", obs.WatcherCount())
	}
}

func TestObservable_ClearWatchers(t *testing.T) {
	obs := NewObservable(0)
	
	obs.Watch(func(old, new int) {})
	obs.Watch(func(old, new int) {})
	
	if obs.WatcherCount() != 2 {
		t.Errorf("WatcherCount() = %v, want 2", obs.WatcherCount())
	}
	
	obs.ClearWatchers()
	
	if obs.WatcherCount() != 0 {
		t.Errorf("After ClearWatchers, WatcherCount() = %v, want 0", obs.WatcherCount())
	}
}

func TestObservable_ConcurrentAccess(t *testing.T) {
	obs := NewObservable(0)
	
	var wg sync.WaitGroup
	
	// Multiple goroutines setting values
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			obs.Set(val)
		}(i)
	}
	
	// Multiple goroutines reading values
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = obs.Get()
		}()
	}
	
	// Multiple goroutines adding/removing watchers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unwatch := obs.Watch(func(old, new int) {})
			time.Sleep(1 * time.Millisecond)
			unwatch()
		}()
	}
	
	wg.Wait()
	
	// Should not panic or deadlock
}

func BenchmarkObservable_Set(b *testing.B) {
	obs := NewObservable(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obs.Set(i)
	}
}

func BenchmarkObservable_Get(b *testing.B) {
	obs := NewObservable(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = obs.Get()
	}
}

func BenchmarkObservable_SetWithWatchers(b *testing.B) {
	obs := NewObservable(0)
	for i := 0; i < 10; i++ {
		obs.Watch(func(old, new int) {})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obs.Set(i)
	}
}
