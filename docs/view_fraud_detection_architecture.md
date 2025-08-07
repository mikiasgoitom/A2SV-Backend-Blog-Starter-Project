# View Fraud Detection Architecture

## 1. Introduction

In a content-driven platform, authentic engagement metrics are crucial for authors, readers, and the platform's credibility. The view count is a primary indicator of a blog post's reach and popularity. To maintain the integrity of this metric, a robust system is in place to detect and prevent artificial inflation of view counts from bots or malicious users.

This document outlines the multi-layered architecture designed to ensure that view counts are a genuine reflection of human interest.

## 2. Multi-Layered Defense Strategy

Our system employs a sophisticated, multi-layered approach to validate each view request before it is counted. A request must pass through all layers of defense to be considered legitimate.

### Layer 1: Basic Bot Filtering

- **Mechanism**: User-Agent Analysis
- **Description**: The first line of defense is a blocklist that filters out requests from known bots, crawlers, and spiders. The `User-Agent` string of every request is inspected. If it contains common bot signatures (e.g., "bot", "crawler", "spider"), the view is immediately discarded. This is a low-cost, high-efficiency check that eliminates a significant portion of non-human traffic.

### Layer 2: Unique View Cooldown

- **Mechanism**: Time-based Debouncing
- **Description**: To prevent a single user from repeatedly refreshing a page to inflate views, we enforce a "cooldown" period. The system checks if the same user (identified by their User ID if logged in, or their IP address if anonymous) has already viewed the _same blog post_ within the last 24 hours. If so, the subsequent views are not counted until the cooldown period expires. This is managed by the `HasViewedRecently` check.

### Layer 3: Advanced Behavioral Analysis

This is the most sophisticated layer, designed to catch more advanced fraudulent patterns that simple debouncing would miss.

- **Mechanism 1: IP Velocity Check**
  - **Description**: This check prevents a single entity (like a script running on a server) from rapidly viewing many _different_ blog posts. The system tracks the number of views originating from a single IP address across the entire platform within a short time window (e.g., 5 minutes). If this number exceeds a configured threshold (`maxIpVelocity`), all subsequent views from that IP are temporarily blocked and flagged as suspicious.
- **Mechanism 2: User-IP Rotation Check**
  - **Description**: This check targets scenarios where a malicious actor uses a single user account but cycles through multiple IP addresses (e.g., using a proxy network) to appear as different users. The system tracks how many unique IP addresses are associated with a single logged-in user account over a medium time window (e.g., 60 minutes). If this exceeds a threshold (`maxUserIPs`), it indicates suspicious activity, and the view is rejected.

## 3. Scenario Walkthrough: A Bot's Failed Attempt

Let's walk through how the system would handle a malicious script attempting to inflate the view count of a blog post.

**The Attacker:** A script running on a server, designed to send 100 view requests to `blog-post-123` in one minute.

1. **Request 1:**

   - **User-Agent**: `curl/7.68.0` (a common script user-agent).
   - **IP Address**: `198.51.100.10`
   - **Layer 1 (Bot Filter)**: The User-Agent is not on the explicit bot blocklist. **PASSES**.
   - **Layer 2 (Cooldown)**: The system checks if IP `198.51.100.10` has viewed `blog-post-123` recently. It has not. **PASSES**.
   - **Layer 3 (Behavioral)**: The system checks the IP velocity. This is the first view from this IP in the last 5 minutes. **PASSES**.
   - **Result**: The view is **COUNTED**. The system records that `blog-post-123` was viewed by IP `198.51.100.10`.

2. **Request 2 (5 seconds later):**

   - **User-Agent**: `curl/7.68.0`
   - **IP Address**: `198.51.100.10`
   - **Layer 1 (Bot Filter)**: **PASSES**.
   - **Layer 2 (Cooldown)**: The system checks if IP `198.51.100.10` has viewed `blog-post-123` recently. It has. **FAILS**.
   - **Result**: The view is **REJECTED**. The process stops here.

3. **Request 11 (The script now tries to view a different post, `blog-post-456`):**
   - **User-Agent**: `curl/7.68.0`
   - **IP Address**: `198.51.100.10`
   - **Layer 1 (Bot Filter)**: **PASSES**.
   - **Layer 2 (Cooldown)**: The system checks if IP `198.51.100.10` has viewed `blog-post-456` recently. It has not. **PASSES**.
   - **Layer 3 (Behavioral - IP Velocity)**: The system checks how many views have come from IP `198.51.100.10` in the last 5 minutes. It finds 10 previous views (to 10 different articles). The threshold `maxIpVelocity` is 10. This request is the 11th. **FAILS**.
   - **Result**: The view is **REJECTED** due to high IP velocity. A warning is logged.

## 4. Technical Component Breakdown

The logic for this entire process is orchestrated within the `TrackBlogView` method in the `BlogUseCaseImpl`.

- `TrackBlogView(ctx, blogID, userID, ipAddress, userAgent)`: This is the entry point. It receives all necessary information from the HTTP handler.

- **Sequence of Operations:**

  1.  It first calls the internal `isBot(userAgent)` helper for the **Layer 1** check.
  2.  It then calls `blogRepo.HasViewedRecently(...)` for the **Layer 2** check.
  3.  For **Layer 3**, it performs two repository calls:
      - `blogRepo.GetRecentViewsByIP(ctx, ipAddress, shortWindow)`
      - `blogRepo.GetRecentViewsByUser(ctx, userID, mediumWindow)`
  4.  It analyzes the results of these calls against the configured thresholds (`maxIpVelocity`, `maxUserIPs`).
  5.  If all checks pass, it finalizes the process by calling `blogRepo.IncrementViewCount(ctx, blogID)` and `blogRepo.RecordView(...)`.

- **Repository Dependencies (`IBlogRepository`):**
  - `HasViewedRecently(...)`: Checks for a view on a specific blog from a user/IP.
  - `GetRecentViewsByIP(...)`: Fetches all views from an IP in a given timeframe.
  - `GetRecentViewsByUser(...)`: Fetches all views from a user in a given timeframe.
  - `RecordView(...)`: Inserts a new view record into the `blog_views` collection, which has a TTL index to automatically purge old records.
  - `IncrementViewCount(...)`: Atomically increments the `viewCount` field on the main `blogs` collection.

This structured, multi-layered approach ensures that our view count metric remains a reliable and trustworthy indicator of genuine reader engagement.
