--- a/util/flowcontrol/throttle.go
+++ b/util/flowcontrol/throttle.go
@@
-import (
-    "golang.org/x/time/rate"
-)
+import (
+    "golang.org/x/time/rate"
+    "k8s.io/klog/v2"
+)
@@
-func NewTokenBucketRateLimiter(qps float32, burst int) RateLimiter {
-    // existing logic
-}
+func NewTokenBucketRateLimiter(qps float32, burst int) RateLimiter {
+    // Sanitize burst to avoid indefinite blocking when QPS > 0.
+    if qps > 0 && burst <= 0 {
+        klog.Warningf("RateLimiter burst must be positive. Adjusting burst from %d to 1.", burst)
+        burst = 1
+    }
+    // existing logic
+}