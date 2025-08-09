package metrics

import "github.com/prometheus/client_golang/prometheus"


var (
    BlogCacheDetailHits = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "blog",
        Subsystem: "cache",
        Name:      "detail_hits_total",
        Help:      "Total cache hits for blog detail (by slug)",
    })
    BlogCacheDetailMiss = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "blog",
        Subsystem: "cache",
        Name:      "detail_miss_total",
        Help:      "Total cache misses for blog detail (by slug)",
    })
    BlogCacheListHits = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "blog",
        Subsystem: "cache",
        Name:      "list_hits_total",
        Help:      "Total cache hits for blogs list",
    })
    BlogCacheListMiss = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "blog",
        Subsystem: "cache",
        Name:      "list_miss_total",
        Help:      "Total cache misses for blogs list",
    })

    // Total duration (seconds) spent serving cache hits and misses
    BlogCacheHitDuration = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "blog",
        Subsystem: "cache",
        Name:      "hit_duration_seconds",
        Help:      "Total time spent serving cache hits (seconds)",
    })
    BlogCacheMissDuration = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "blog",
        Subsystem: "cache",
        Name:      "miss_duration_seconds",
        Help:      "Total time spent serving cache misses (seconds)",
    })
)


func init() {
    prometheus.MustRegister(
        BlogCacheDetailHits,
        BlogCacheDetailMiss,
        BlogCacheListHits,
        BlogCacheListMiss,
        BlogCacheHitDuration,
        BlogCacheMissDuration,
    )
}


func IncDetailHit()  { BlogCacheDetailHits.Inc() }
func IncDetailMiss() { BlogCacheDetailMiss.Inc() }
func IncListHit()    { BlogCacheListHits.Inc() }
func IncListMiss()   { BlogCacheListMiss.Inc() }

// Add duration (in seconds) to the total hit/miss duration counters
func AddHitDuration(seconds float64)  { BlogCacheHitDuration.Add(seconds) }
func AddMissDuration(seconds float64) { BlogCacheMissDuration.Add(seconds) }
