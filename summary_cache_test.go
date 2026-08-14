package main

import (
	"context"
	"testing"
	"time"
)

func TestSummaryCacheEntryPersistsRevision(t *testing.T) {
	store := summaryCacheTestStore(t)
	key := summaryCacheKey{Window: "24h", Limit: 2000}
	want := summaryCacheEntry{
		data:       map[string]any{"window": "24h", "requests": float64(42)},
		cachedAt:   time.Now().Truncate(time.Second),
		durationMs: 123,
		err:        "last refresh failed",
		revision:   "u:42|q:2",
	}
	if err := store.saveSummaryCacheEntry(context.Background(), key, want); err != nil {
		t.Fatal(err)
	}
	got, ok, err := store.loadSummaryCacheEntry(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("saved summary cache entry was not loaded")
	}
	if got.revision != want.revision || got.durationMs != want.durationMs || got.err != want.err || !got.cachedAt.Equal(want.cachedAt) {
		t.Fatalf("loaded summary cache metadata = %+v, want %+v", got, want)
	}
}

func TestSummaryCacheServesRecentPriorRevision(t *testing.T) {
	cfg := defaultPluginConfig()
	cfg.SummaryCacheMaxAgeSeconds = 30
	manager := &summaryPrecomputeManager{
		cfg: cfg,
		entries: map[summaryCacheKey]summaryCacheEntry{
			{Window: "24h", Limit: 2000}: {
				data:     map[string]any{"window": "24h"},
				cachedAt: time.Now().Add(-5 * time.Second),
				revision: "old",
			},
		},
	}
	data, ok := manager.cachedEntry(summaryCacheKey{Window: "24h", Limit: 2000}, cfg, "new")
	if !ok {
		t.Fatal("recent cache should remain available while a newer revision is pending")
	}
	info, ok := data["precompute"].(summaryPrecomputeInfo)
	if !ok {
		t.Fatalf("missing precompute metadata: %#v", data["precompute"])
	}
	if info.Stale || info.Reason != "revision_pending" {
		t.Fatalf("precompute metadata = %+v, want a fresh revision_pending cache hit", info)
	}
	if data["store_revision"] != "new" || data["cache_revision"] != "old" {
		t.Fatalf("revision metadata = store:%v cache:%v", data["store_revision"], data["cache_revision"])
	}
}

func TestSummaryCacheMarksExpiredPriorRevisionStale(t *testing.T) {
	cfg := defaultPluginConfig()
	cfg.SummaryCacheMaxAgeSeconds = 30
	manager := &summaryPrecomputeManager{
		cfg: cfg,
		entries: map[summaryCacheKey]summaryCacheEntry{
			{Window: "24h", Limit: 2000}: {
				data:     map[string]any{"window": "24h"},
				cachedAt: time.Now().Add(-45 * time.Second),
				revision: "old",
			},
		},
	}
	data, ok := manager.cachedEntry(summaryCacheKey{Window: "24h", Limit: 2000}, cfg, "new")
	if !ok {
		t.Fatal("stale cache should remain available for asynchronous refresh")
	}
	info := data["precompute"].(summaryPrecomputeInfo)
	if !info.Stale || info.Reason != "revision_stale" {
		t.Fatalf("precompute metadata = %+v, want revision_stale", info)
	}
}

func TestSummaryCacheForceLookupMarksFreshEntryStale(t *testing.T) {
	cfg := defaultPluginConfig()
	cfg.SummaryCacheMaxAgeSeconds = 30
	key := summaryCacheKey{Window: "24h", Limit: 2000}
	manager := &summaryPrecomputeManager{
		cfg: cfg,
		entries: map[summaryCacheKey]summaryCacheEntry{
			key: {
				data:     map[string]any{"window": "24h"},
				cachedAt: time.Now(),
				revision: "same",
			},
		},
	}
	data, ok := manager.cachedAny(context.Background(), nil, key, cfg, "same")
	if !ok {
		t.Fatal("force refresh should return the current cache while refreshing asynchronously")
	}
	info := data["precompute"].(summaryPrecomputeInfo)
	if !info.Stale || info.Reason != "age_stale" {
		t.Fatalf("force refresh metadata = %+v, want stale cache metadata", info)
	}
}

func TestMarkSummaryRefreshState(t *testing.T) {
	data := map[string]any{"precompute": summaryPrecomputeInfo{Stale: true}}
	markSummaryRefreshState(data, true)
	if info := data["precompute"].(summaryPrecomputeInfo); info.Reason != "refresh_started" || !info.Stale {
		t.Fatalf("started refresh metadata = %+v", info)
	}
	markSummaryRefreshState(data, false)
	if info := data["precompute"].(summaryPrecomputeInfo); info.Reason != "refresh_in_progress" || !info.Stale {
		t.Fatalf("in-progress refresh metadata = %+v", info)
	}
}

func TestSummaryRefreshForceUpgradeIsRetained(t *testing.T) {
	key := summaryCacheKey{Window: "24h", Limit: 2000}
	manager := &summaryPrecomputeManager{
		refreshing:   map[summaryCacheKey]bool{key: true},
		refreshForce: map[summaryCacheKey]bool{},
	}
	if manager.refreshAsync(nil, defaultPluginConfig(), key, true) {
		t.Fatal("existing refresh should be upgraded instead of starting a duplicate")
	}
	manager.mu.Lock()
	force := manager.refreshForce[key]
	manager.mu.Unlock()
	if !force {
		t.Fatal("force refresh upgrade was not retained")
	}
}

func TestSummaryAsyncRefreshAdmissionDoesNotQueue(t *testing.T) {
	manager := &summaryPrecomputeManager{}
	if !manager.tryStartAsyncRefresh() {
		t.Fatal("first asynchronous refresh should be admitted")
	}
	if manager.tryStartAsyncRefresh() {
		t.Fatal("second asynchronous refresh should be skipped instead of queued")
	}
	manager.finishAsyncRefresh()
	if !manager.tryStartAsyncRefresh() {
		t.Fatal("refresh admission should reopen after completion")
	}
	manager.finishAsyncRefresh()
}

func TestSummaryCacheRejectsHardExpiredEntry(t *testing.T) {
	cfg := defaultPluginConfig()
	cfg.SummaryCacheMaxAgeSeconds = 30
	manager := &summaryPrecomputeManager{
		cfg: cfg,
		entries: map[summaryCacheKey]summaryCacheEntry{
			{Window: "24h", Limit: 2000}: {
				data:     map[string]any{"window": "24h"},
				cachedAt: time.Now().Add(-summaryCacheMaxStale(cfg) - time.Second),
				revision: "old",
			},
		},
	}
	if _, ok := manager.cachedEntry(summaryCacheKey{Window: "24h", Limit: 2000}, cfg, "new"); ok {
		t.Fatal("hard-expired summary cache must not be served")
	}
}

func TestSQLiteHeavyWorkAcquisitionHonorsContext(t *testing.T) {
	if err := acquireSQLiteHeavyWork(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer releaseSQLiteHeavyWork()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := acquireSQLiteHeavyWork(ctx); err == nil {
		t.Fatal("blocked heavy work acquisition should honor context cancellation")
	}
}

func TestSummaryCacheHardExpiryIsFixed(t *testing.T) {
	cfg := defaultPluginConfig()
	cfg.SummaryCacheMaxAgeSeconds = 3600
	if got := summaryCacheMaxStale(cfg); got != 5*time.Minute {
		t.Fatalf("hard stale limit = %s, want 5m", got)
	}
}

func summaryCacheTestStore(t *testing.T) *store {
	t.Helper()
	previousStore := globalStore
	store := &store{}
	globalStore = store
	t.Setenv("CPA_TOKEN_USAGE_DIR", t.TempDir())
	t.Cleanup(func() {
		store.close()
		globalStore = previousStore
	})
	if _, _, err := store.open(context.Background()); err != nil {
		t.Fatal(err)
	}
	return store
}
