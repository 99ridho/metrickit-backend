# MetricKit Backend

## Purpose

To be able to consume MetricKit payload that supplied by `MXMetricManager` via `MXMetricManagerSubscriber`, the payload must be sent to other service via HTTP, and show the payload to own dashboard to get insight about our app metrics.

## Why

As we're experimenting using `MXSignpost` to help us to gather execution duration, memory usage and CPU usage at some piece of important code (currently we implement at existing `TimeMeasurer`), we need that metrics to each View Controller which that load time is measured by `TimeMeasurer`.

So, to be able to get insight about that metrics data, we need to display the metrics with Grafana, thus this small service can help us to do so.

## References

[WWDC 2019 - Improving Battery Life and Performance](https://developer.apple.com/videos/play/wwdc2019/417/)
[MetricKit Internals. Insights into your iOS app performance.](https://appspector.com/blog/metrickit)
[NSHipster - MetricKit](https://nshipster.com/metrickit/)