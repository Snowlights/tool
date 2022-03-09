package vcollector

//Nerve SD configurations allow retrieving scrape targets from [AirBnB's Nerve]
// (https://github.com/airbnb/nerve) which are stored in Zookeeper.

// Serverset SD configurations allow retrieving scrape targets from
// [Serversets] (https://github.com/twitter/finagle/tree/master/finagle-serversets)
// which are stored in Zookeeper. Serversets are commonly used by Finagle and Aurora.

// Zookeeper only supports these two kinds of structured data,
// but we implement service registration and discovery by ourselves, so we don't use ZK for service index statistics

const DefaultMetricPath = "/metric"

func NewCollector() {

}
