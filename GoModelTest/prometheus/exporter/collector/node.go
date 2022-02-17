package collector

import (
	"runtime"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/host"
)

var (
	hostname string
)

type NodeCollector struct {
	threadsDesc *prometheus.Desc //Gauge
	mutex       sync.Mutex
}

func NewNodeCollector() prometheus.Collector {
	host, _ := host.Info()
	hostname = host.Hostname
	return &NodeCollector{}
}

// Describe returns all descriptions of the collector.
//实现采集器Describe接口
func (n *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.threadsDesc
}

// Collect returns the current state of all metrics of the collector.
//实现采集器Collect接口,真正采集动作
func (n *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	num, _ := runtime.ThreadCreateProfile(nil)
	ch <- prometheus.MustNewConstMetric(n.threadsDesc, prometheus.GaugeValue, float64(num))
	n.mutex.Unlock()
}
