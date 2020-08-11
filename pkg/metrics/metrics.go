package metrics

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/IBM/prometheus-metrics-server-simulator/pkg/config"
)

var generator *Generator

type Generator struct {
	config     *config.GeneratorConfig
	r          *prometheus.Registry
	metrics    []Metric
	metricsMap map[string]Metric
}
type Metric interface {
	run()
	valueMode() config.ValueMode
	setValue(v float64)
}
type CounterMetric struct {
	config *config.CounterConfig
	metric prometheus.Counter
	mode   config.ValueMode
	dataCh chan float64
}

func (m *CounterMetric) valueMode() config.ValueMode {
	return m.mode

}
func (m *CounterMetric) setValue(v float64) {
	m.dataCh <- v

}
func (m *CounterMetric) run() {
	if m.mode == config.ValueModeAuto {
		go func() {
			for {
				m.metric.Add(float64(rand.Intn(10)))
				time.Sleep(time.Second)
			}
		}()
	} else {
		go func() {
			for {
				select {
				case data := <-m.dataCh:
					m.metric.Add(data)
				default:
					time.Sleep(time.Millisecond * 50)
				}
			}
		}()

	}

}

type GaugeMetric struct {
	config *config.GaugeConfig
	metric prometheus.Gauge
	currV  float64
	mode   config.ValueMode
	dataCh chan float64
}

func (m *GaugeMetric) valueMode() config.ValueMode {
	return m.mode

}
func (m *GaugeMetric) setValue(v float64) {
	m.dataCh <- v

}
func (m *GaugeMetric) run() {
	if m.mode == config.ValueModeAuto {
		go func() {
			for {
				delta := float64(rand.Intn(20) - 10)
				if m.currV+delta > float64(m.config.Range.Upper) || m.currV+delta < float64(m.config.Range.Lower) {
					delta = 0.0
				}
				m.currV = m.currV + delta
				m.metric.Set(m.currV)
				time.Sleep(time.Second)
			}
		}()
	} else {
		go func() {
			for {
				select {
				case data := <-m.dataCh:
					m.metric.Set(data)
				default:
					time.Sleep(time.Millisecond * 50)
				}
			}
		}()
	}

}
func NewGenerator(file string, r *prometheus.Registry) (*Generator, error) {
	config, err := config.ParseConfig(file)
	generator = &Generator{config: config,
		r:          r,
		metrics:    []Metric{},
		metricsMap: make(map[string]Metric),
	}
	if err != nil {
		return nil, err
	}
	return generator, nil
}
func (g *Generator) Start() {
	g.addMetrics()
	g.runMetrics()
}

func (g *Generator) SetValue(mname string, labels map[string]string, v float64) error {
	m, ok := g.metricsMap[g.generateMapkey(mname, labels)]
	if !ok {
		return fmt.Errorf(fmt.Sprintf("metricdoes not exist. name: %v. labels: %v"), mname, labels)
	}
	if m.valueMode() == config.ValueModeAuto {
		return fmt.Errorf("metric value is auto generated and not changeable")
	}
	m.setValue(v)
	return nil

}
func (g *Generator) generateMapkey(mname string, labels map[string]string) string {
	key := mname

	for k, v := range labels {
		key = key + ":" + k + "," + v

	}
	return key

}
func (g *Generator) addMetrics() {
	if g.config.Counters != nil {
		for _, counter := range g.config.Counters {
			for i := 0; i < counter.Number; i++ {

				labels := make(map[string]string)
				for _, lsetting := range counter.Labels {
					v := selectRandomLabelValue(lsetting.ValueSet)
					labels[lsetting.Name] = v
				}
				opts := prometheus.CounterOpts{}
				if counter.Number == 1 {
					opts.Name = counter.Prefix
				} else {
					opts.Name = fmt.Sprint(counter.Prefix, "_", i)
				}
				opts.ConstLabels = labels
				pmetric := prometheus.NewCounter(opts)
				cmetric := CounterMetric{config: counter,
					metric: pmetric,
					mode:   counter.ValueMode,
					dataCh: make(chan float64)}
				g.r.MustRegister(pmetric)
				g.metrics = append(g.metrics, &cmetric)
				g.metricsMap[g.generateMapkey(opts.Name, opts.ConstLabels)] = &cmetric

			}

		}

	}
	if g.config.Gauges != nil {
		for _, gauge := range g.config.Gauges {
			for i := 0; i < gauge.Number; i++ {
				labels := make(map[string]string)
				for _, lsetting := range gauge.Labels {
					v := selectRandomLabelValue(lsetting.ValueSet)
					labels[lsetting.Name] = v
				}
				opts := prometheus.GaugeOpts{}
				opts.Name = fmt.Sprint(gauge.Prefix, "_", i)
				opts.ConstLabels = labels
				pmetric := prometheus.NewGauge(opts)
				g.r.MustRegister(pmetric)
				cmetric := GaugeMetric{config: gauge,
					metric: pmetric,
					currV:  float64(gauge.Range.Lower),
					mode:   gauge.ValueMode,
					dataCh: make(chan float64)}
				g.metrics = append(g.metrics, &cmetric)
				g.metricsMap[g.generateMapkey(opts.Name, opts.ConstLabels)] = &cmetric

			}

		}

	}

}

func (g *Generator) runMetrics() {
	for _, m := range g.metrics {
		m.run()
	}
}
func selectRandomLabelValue(labelValues []string) string {
	if len(labelValues) == 0 {
		return ""
	}
	i := rand.Intn(len(labelValues))
	return labelValues[i]

}
