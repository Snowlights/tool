package vprometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"unicode/utf8"
)

func makePrometheusLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}

// LabelValues is a type alias that provides validation on its With method.
// Metrics may include it as a member to help them satisfy With semantics and
// save some code duplication.
type LabelValues []string

// With validates the input, and returns a new aggregate labelValues.
func (lvs LabelValues) With(labelValues ...string) LabelValues {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}
	return append(lvs, labelValues...)
}

// Check check label valid
func (lvs LabelValues) Check() error {
	for i := 1; i < len(lvs); i += 2 {
		if !utf8.ValidString(lvs[i]) {
			return fmt.Errorf("label %s: value %q is not valid UTF-8", lvs[i-1], lvs[i])
		}
	}
	return nil
}
