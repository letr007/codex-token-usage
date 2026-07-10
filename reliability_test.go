package main

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestReliabilityWindowContracts(t *testing.T) {
	generatedAt := time.Date(2026, time.January, 15, 12, 1, 2, 0, time.UTC)
	for _, test := range []struct {
		window        string
		buckets       int
		rows, columns int
		bucketWidth   time.Duration
	}{
		{"24h", 288, 6, 48, 5 * time.Minute},
		{"today", 288, 6, 48, 5 * time.Minute},
		{"7d", 336, 7, 48, 30 * time.Minute},
		{"30d", 336, 7, 48, 30 * time.Minute},
		{"all", 336, 7, 48, 30 * time.Minute},
	} {
		plan := reliabilityWindowFor(test.window, generatedAt)
		overview := buildReliabilityOverview(plan)
		if len(overview.Buckets) != test.buckets || overview.Rows != test.rows || overview.Columns != test.columns || overview.Total != test.buckets {
			t.Fatalf("%s: invalid grid contract: %+v", test.window, overview)
		}
		if overview.Buckets[0].Start != plan.start || overview.Buckets[len(overview.Buckets)-1].End != plan.end {
			t.Fatalf("%s: buckets do not cover the complete window", test.window)
		}
		for i := 1; i < len(overview.Buckets); i++ {
			if overview.Buckets[i-1].End != overview.Buckets[i].Start {
				t.Fatalf("%s: buckets %d and %d are not contiguous", test.window, i-1, i)
			}
		}
		for i, bucket := range overview.Buckets {
			if bucket.End.Sub(bucket.Start) != test.bucketWidth {
				t.Fatalf("%s: bucket %d width = %s, want %s", test.window, i, bucket.End.Sub(bucket.Start), test.bucketWidth)
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

func TestReliabilityTodayUsesDSTSafeCalendarDay(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	generatedAt := time.Date(2026, time.March, 8, 12, 0, 0, 0, location)
	plan := reliabilityWindowFor("today", generatedAt)
	if !plan.end.Equal(plan.start.AddDate(0, 0, 1)) {
		t.Fatalf("today end = %s, want next local midnight after %s", plan.end, plan.start)
	}
	if plan.end.Sub(plan.start) != 23*time.Hour {
		t.Fatalf("DST day duration = %s, want 23h", plan.end.Sub(plan.start))
	}
	overview := buildReliabilityOverview(plan)
	if len(overview.Buckets) != reliability24hBuckets || overview.Buckets[0].Start != plan.start || overview.Buckets[len(overview.Buckets)-1].End != plan.end {
		t.Fatalf("DST grid does not cover the complete local day: %+v", overview)
	}
	for i := 1; i < len(overview.Buckets); i++ {
		if overview.Buckets[i-1].End != overview.Buckets[i].Start {
			t.Fatalf("DST buckets %d and %d are not contiguous", i-1, i)
		}
	}
}

func TestReliabilityTodayFallBackGrid(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	plan := reliabilityWindowFor("today", time.Date(2026, time.November, 1, 12, 0, 0, 0, location))
	if plan.end.Sub(plan.start) != 25*time.Hour {
		t.Fatalf("fall-back day duration = %s, want 25h", plan.end.Sub(plan.start))
	}
	overview := buildReliabilityOverview(plan)
	if len(overview.Buckets) != reliability24hBuckets || overview.Buckets[0].Start != plan.start || overview.Buckets[len(overview.Buckets)-1].End != plan.end {
		t.Fatalf("fall-back grid does not cover the complete local day: %+v", overview)
	}
	for i, bucket := range overview.Buckets {
		width := bucket.End.Sub(bucket.Start)
		if !bucket.Start.Before(bucket.End) || width < 4*time.Minute+45*time.Second || width > 5*time.Minute+15*time.Second {
			t.Fatalf("fall-back bucket %d has invalid bounds: %+v", i, bucket)
		}
		if i > 0 && overview.Buckets[i-1].End != bucket.Start {
			t.Fatalf("fall-back buckets %d and %d are not contiguous", i-1, i)
		}
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

func TestReliabilityLongWindowLocalHalfHourCeiling(t *testing.T) {
	location, err := time.LoadLocation("Asia/Kathmandu")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name        string
		generatedAt time.Time
		end         time.Time
	}{
		{"exact boundary", time.Date(2026, time.January, 15, 12, 30, 0, 0, location), time.Date(2026, time.January, 15, 12, 30, 0, 0, location)},
		{"between boundaries", time.Date(2026, time.January, 15, 12, 10, 0, 0, location), time.Date(2026, time.January, 15, 12, 30, 0, 0, location)},
		{"seconds past boundary", time.Date(2026, time.January, 15, 12, 30, 1, 0, location), time.Date(2026, time.January, 15, 13, 0, 0, 0, location)},
	} {
		for _, window := range []string{"7d", "30d", "all"} {
			plan := reliabilityWindowFor(window, test.generatedAt)
			if !plan.end.Equal(test.end) {
				t.Fatalf("%s %s: long-window end = %s, want %s", test.name, window, plan.end, test.end)
			}
			if plan.end.Location() != location || plan.end.Minute() != test.end.Minute() || (plan.end.Minute() != 0 && plan.end.Minute() != 30) {
				t.Fatalf("%s %s: long-window end is not on a local half-hour boundary: %s", test.name, window, plan.end)
			}
			if plan.end.Sub(plan.start) != 7*24*time.Hour {
				t.Fatalf("%s %s: long-window duration = %s, want 7d", test.name, window, plan.end.Sub(plan.start))
			}
			if plan.effective != "7d" {
				t.Fatalf("%s %s: effective window = %q, want 7d", test.name, window, plan.effective)
			}
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
	insertReliabilityEvent(t, db, start+1799, "codex", 0, 200)           // bucket 0
	insertReliabilityEvent(t, db, start+1800, "codex", 1, 429)           // bucket 1
	insertReliabilityEvent(t, db, start+1800, "other", 0, 200)           // wrong scope
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

func TestQueryReliability24hFiveMinuteBucketBoundaries(t *testing.T) {
	db := reliabilityTestDB(t)
	plan := reliabilityWindowFor("24h", time.Date(2026, time.January, 15, 12, 1, 2, 0, time.UTC))
	boundary := plan.start.Unix() + int64(5*time.Minute/time.Second)
	insertReliabilityEvent(t, db, boundary-1, "codex", 0, 200)
	insertReliabilityEvent(t, db, boundary, "codex", 0, 200)
	result, err := queryReliability(context.Background(), db, plan)
	if err != nil {
		t.Fatal(err)
	}
	if result.Buckets[0].Success != 1 || result.Buckets[1].Success != 1 {
		t.Fatalf("five-minute boundary assigned to wrong buckets: before=%+v after=%+v", result.Buckets[0], result.Buckets[1])
	}
}

func TestQueryReliabilityDSTTodayBucketBoundaries(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name        string
		generatedAt time.Time
	}{
		{"spring forward", time.Date(2026, time.March, 8, 12, 0, 0, 0, location)},
		{"fall back", time.Date(2026, time.November, 1, 12, 0, 0, 0, location)},
	} {
		t.Run(test.name, func(t *testing.T) {
			db := reliabilityTestDB(t)
			plan := reliabilityWindowFor("today", test.generatedAt)
			overview := buildReliabilityOverview(plan)
			boundary := overview.Buckets[0].End.Unix()
			insertReliabilityEvent(t, db, boundary-1, "codex", 0, 200)
			insertReliabilityEvent(t, db, boundary, "codex", 0, 200)

			result, err := queryReliability(context.Background(), db, plan)
			if err != nil {
				t.Fatal(err)
			}
			if result.Buckets[0].Success != 1 || result.Buckets[1].Success != 1 {
				t.Fatalf("DST bucket boundary assigned incorrectly: before=%+v after=%+v", result.Buckets[0], result.Buckets[1])
			}
		})
	}
}

func TestQueryReliabilityEmptyGrid(t *testing.T) {
	db := reliabilityTestDB(t)
	plan := reliabilityWindowFor("24h", time.Date(2026, time.January, 15, 12, 1, 2, 0, time.UTC))
	overview, err := queryReliability(context.Background(), db, plan)
	if err != nil {
		t.Fatal(err)
	}
	if overview.Rate != -1 || overview.Success != 0 || overview.Failure != 0 || len(overview.Buckets) != reliability24hBuckets {
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
