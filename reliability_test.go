package main

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestReliabilityWindowContracts(t *testing.T) {
	generatedAt := time.Date(2026, time.January, 15, 12, 1, 2, 0, time.UTC)
	for _, window := range []string{"24h", "today", "7d", "30d", "all"} {
		plan := reliabilityWindowFor(window, generatedAt)
		overview := buildReliabilityOverview(plan)
		if len(overview.Buckets) != reliabilityBuckets || overview.Rows != 7 || overview.Columns != 96 || overview.Total != 672 {
			t.Fatalf("%s: invalid grid contract: %+v", window, overview)
		}
		if overview.Buckets[0].Start != plan.start || overview.Buckets[reliabilityBuckets-1].End != plan.end {
			t.Fatalf("%s: buckets do not cover the complete window", window)
		}
		for i := 1; i < len(overview.Buckets); i++ {
			if overview.Buckets[i-1].End != overview.Buckets[i].Start {
				t.Fatalf("%s: buckets %d and %d are not contiguous", window, i-1, i)
			}
		}
	}
	if got := reliabilityWindowFor("30d", generatedAt).effective; got != "7d" {
		t.Fatalf("30d effective window = %q, want 7d", got)
	}
	if got := reliabilityWindowFor("all", generatedAt).effective; got != "7d" {
		t.Fatalf("all effective window = %q, want 7d", got)
	}
}

func TestReliabilityWindowTruncatesToSeconds(t *testing.T) {
	generatedAt := time.Date(2026, time.January, 15, 12, 1, 2, 987654321, time.FixedZone("test", 8*60*60))
	for _, window := range []string{"24h", "today", "7d", "30d", "all"} {
		plan := reliabilityWindowFor(window, generatedAt)
		if plan.start.Nanosecond() != 0 || plan.end.Nanosecond() != 0 || plan.observed.Nanosecond() != 0 {
			t.Fatalf("%s: plan boundaries must be whole seconds: %+v", window, plan)
		}
		for i, bucket := range buildReliabilityOverview(plan).Buckets {
			if bucket.Start.Nanosecond() != 0 || bucket.End.Nanosecond() != 0 {
				t.Fatalf("%s: bucket %d boundaries must be whole seconds: %+v", window, i, bucket)
			}
		}
	}

	db := reliabilityTestDB(t)
	plan := reliabilityWindowFor("24h", generatedAt)
	insertReliabilityEvent(t, db, plan.observed.Unix()-1, "codex", 0, 200)
	insertReliabilityEvent(t, db, plan.observed.Unix(), "codex", 0, 200)
	overview, err := queryReliability(context.Background(), db, plan)
	if err != nil {
		t.Fatal(err)
	}
	if overview.Success != 1 {
		t.Fatalf("SQL observed boundary must match the declared whole-second boundary: %+v", overview)
	}
}

func TestReliabilityLongWindowCeilingAtExactQuarterHour(t *testing.T) {
	generatedAt := time.Date(2026, time.January, 15, 12, 15, 0, 0, time.FixedZone("test", 8*60*60))
	for _, window := range []string{"7d", "30d", "all"} {
		plan := reliabilityWindowFor(window, generatedAt)
		if !plan.end.Equal(generatedAt) {
			t.Fatalf("%s: long-window end = %s, want exact quarter-hour %s", window, plan.end, generatedAt)
		}
		if plan.effective != "7d" {
			t.Fatalf("%s: effective window = %q, want 7d", window, plan.effective)
		}
	}
}

func TestReliabilityTodayIncludesFutureBuckets(t *testing.T) {
	location := time.FixedZone("test", 8*60*60)
	generatedAt := time.Date(2026, time.January, 15, 12, 1, 2, 0, location)
	plan := reliabilityWindowFor("today", generatedAt)
	if !plan.start.Equal(time.Date(2026, time.January, 15, 0, 0, 0, 0, location)) || !plan.end.Equal(plan.start.AddDate(0, 0, 1)) {
		t.Fatalf("today must be one complete local calendar day: %+v", plan)
	}
	overview := buildReliabilityOverview(plan)
	future := 0
	for _, bucket := range overview.Buckets {
		if !bucket.Start.Before(generatedAt) {
			future++
			if bucket.Rate != -1 {
				t.Fatal("future bucket must begin empty")
			}
		}
	}
	if future == 0 {
		t.Fatal("today should retain future buckets")
	}
}

func TestQueryReliabilityBoundariesAndScope(t *testing.T) {
	db := reliabilityTestDB(t)
	generatedAt := time.Date(2026, time.January, 15, 12, 1, 2, 0, time.UTC)
	plan := reliabilityWindowFor("7d", generatedAt)
	start := plan.start.Unix()
	insertReliabilityEvent(t, db, start-1, "codex", 0, 200)              // before range
	insertReliabilityEvent(t, db, start+899, "codex", 0, 200)            // bucket 0
	insertReliabilityEvent(t, db, start+900, "codex", 1, 429)            // bucket 1
	insertReliabilityEvent(t, db, start+900, "other", 0, 200)            // wrong scope
	insertReliabilityEvent(t, db, plan.observed.Unix(), "codex", 0, 200) // observed end is exclusive

	overview, err := queryReliability(context.Background(), db, plan)
	if err != nil {
		t.Fatal(err)
	}
	if overview.Success != 1 || overview.Failure != 1 || overview.RateLimited != 1 || overview.Rate != 0.5 {
		t.Fatalf("unexpected aggregate: %+v", overview)
	}
	if bucket := overview.Buckets[0]; bucket.Success != 1 || bucket.Failure != 0 || bucket.Rate != 1 {
		t.Fatalf("unexpected first bucket: %+v", bucket)
	}
	if bucket := overview.Buckets[1]; bucket.Success != 0 || bucket.Failure != 1 || bucket.RateLimited != 1 || bucket.Rate != 0 {
		t.Fatalf("unexpected second bucket: %+v", bucket)
	}
	if overview.Buckets[2].Rate != -1 {
		t.Fatal("empty buckets must not affect the aggregate rate")
	}
}

func TestQueryReliability24hNonIntegralBucketBoundaries(t *testing.T) {
	db := reliabilityTestDB(t)
	plan := reliabilityWindowFor("24h", time.Date(2026, time.January, 15, 12, 1, 2, 0, time.UTC))
	overview := buildReliabilityOverview(plan)

	indexesBySpan := make(map[time.Duration]int)
	for i := 0; i < len(overview.Buckets)-1; i++ {
		span := overview.Buckets[i].End.Sub(overview.Buckets[i].Start)
		if (span == 128*time.Second || span == 129*time.Second) && indexesBySpan[span] == 0 {
			indexesBySpan[span] = i + 1
		}
		if len(indexesBySpan) == 2 {
			break
		}
	}
	if len(indexesBySpan) != 2 {
		t.Fatal("24h grid must contain both 128-second and 129-second buckets")
	}
	boundaryIndexes := []int{indexesBySpan[128*time.Second] - 1, indexesBySpan[129*time.Second] - 1}
	for _, index := range boundaryIndexes {
		boundary := overview.Buckets[index].End.Unix()
		insertReliabilityEvent(t, db, boundary-1, "codex", 0, 200)
		insertReliabilityEvent(t, db, boundary, "codex", 0, 200)
	}
	result, err := queryReliability(context.Background(), db, plan)
	if err != nil {
		t.Fatal(err)
	}
	for _, index := range boundaryIndexes {
		if result.Buckets[index].Success != 1 || result.Buckets[index+1].Success != 1 {
			t.Fatalf("boundary after bucket %d assigned to wrong buckets: before=%+v after=%+v", index, result.Buckets[index], result.Buckets[index+1])
		}
	}
}

func TestQueryReliabilityEmptyGrid(t *testing.T) {
	db := reliabilityTestDB(t)
	plan := reliabilityWindowFor("24h", time.Date(2026, time.January, 15, 12, 1, 2, 0, time.UTC))
	overview, err := queryReliability(context.Background(), db, plan)
	if err != nil {
		t.Fatal(err)
	}
	if overview.Rate != -1 || overview.Success != 0 || overview.Failure != 0 || len(overview.Buckets) != 672 {
		t.Fatalf("unexpected empty overview: %+v", overview)
	}
	for i, bucket := range overview.Buckets {
		if bucket.Rate != -1 {
			t.Fatalf("empty bucket %d rate = %v, want -1", i, bucket.Rate)
		}
	}
}

func reliabilityTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatal(err)
	}
	return db
}

func insertReliabilityEvent(t *testing.T, db *sql.DB, requestedAt int64, provider string, failed, statusCode int) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO usage_events (requested_at, provider, failed, status_code) VALUES (?, ?, ?, ?)`, requestedAt, provider, failed, statusCode); err != nil {
		t.Fatal(err)
	}
}
