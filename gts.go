package Warp10Exporter

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Datapoint are datapoint in a GTS
type Datapoint struct {
	Timestamp time.Time
	Value     interface{}
}

// Datapoints are a slice of Datapoint
type Datapoints []Datapoint

// GTS are GeoTimeSeries
type GTS struct {
	Classname  string
	Labels     map[string]string
	Datapoints Datapoints
}

// Labels type, used to pass to `WithLabels`.
type Labels map[string]string

// NewGTS is creating a new GTS with a name and a value
func NewGTS(classname string) *GTS {

	gts := &GTS{
		Classname: classname,
		Labels:    make(map[string]string),
	}
	return gts
}

// AddDatapoint is adding a datapoint to a GTS
func (gts *GTS) AddDatapoint(ts time.Time, value interface{}) *GTS {

	gts.Datapoints = append(gts.Datapoints, Datapoint{Timestamp: ts, Value: value})
	return gts
}

// WithLabels is adding Labels
func (gts *GTS) WithLabels(labels Labels) *GTS {
	gts.Labels = labels
	return gts
}

// PrintValue is Printing the value
// It's supporting string
func (dp *Datapoint) PrintValue() string {
	switch v := dp.Value.(type) {
	case bool:
		return fmt.Sprintf("%v", strings.ToUpper(strconv.FormatBool(v)[0:1]))
	case string:
		return fmt.Sprintf("\"%v\"", v)
	}
	return fmt.Sprintf("%v", dp.Value)
}

// AddLabel is pushing a new label to the GTS
func (gts *GTS) AddLabel(key string, value string) *GTS {
	gts.Labels[key] = value
	return gts
}

// Print is printing the Warp10 Input Format
// it respects the following format:
// TS// NAME{LABELS} VALUE
func (gts *GTS) Print(b *bytes.Buffer) {

	for i, dp := range gts.Datapoints {
		ts := dp.Timestamp.Unix() * 1000 * 1000
		b.WriteString(fmt.Sprintf("%d// %s{%s} %v", ts, gts.Classname, gts.getLabels(), dp.PrintValue()))
		if i != len(gts.Datapoints)-1 {
			b.WriteString("\n")
		}
	}
}

// getLabels format the map into the right form
func (gts *GTS) getLabels() string {

	var s string
	for key, value := range gts.Labels {

		s = s + key + "=" + value + ","
	}
	// Removing last comma
	s = strings.TrimSuffix(s, ",")
	return s
}

// Push is pushing a single GTS to a warp10 endpoint
func (gts *GTS) Push(warp10Endpoint string, warp10Token string) int {
	var b bytes.Buffer
	gts.Print(&b)
	return pushGTS(&b, warp10Endpoint, warp10Token)
}
