package args

import (
	"math/big"
	"testing"
	"time"
)

func TestBitCountToMaxDecimalCount(t *testing.T) {
	bitStart := 1
	v := big.NewInt(1)
	v.Lsh(v, uint(bitStart-1))

	for bitCount := bitStart; bitCount < 10000; bitCount++ {
		decimal := v.String()
		want := len(decimal)
		if got := BitCountToMaxDecimalCount(bitCount); got != want {
			t.Errorf("BitCountToMaxDecimalCount() = %v, want %v", got, want)
		}
		v.Lsh(v, 1)
	}
}

func TestDurationFixedLen(t *testing.T) {
	type args struct {
		d        time.Duration
		fixedLen int
	}
	tests := []struct {
		name string
		args args
	}{
		{"0µs", args{0, 0}},
		{"0.00µs", args{0, 7}},
		{"0.50ms", args{500*time.Microsecond + 1, 6}},
		{"1.00ms", args{1 * time.Millisecond, 6}},
		{"1.000s", args{1 * time.Second, 6}},
		{"60.00s", args{1 * time.Minute, 6}},
		{"60m00s", args{1 * time.Hour, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationFixedLen(tt.args.d, tt.args.fixedLen); got != tt.name {
				t.Errorf("DurationFixedLen() = %v, want %v", got, tt.name)
			}
		})
	}
}

func Test_fmtAboveSeconds(t *testing.T) {
	type args struct {
		d        time.Duration
		fixedLen int
	}
	tests := []struct {
		name string
		args args
	}{
		{"0s", args{0, 0}},
		{"1s", args{1 * time.Second, 0}},
		{"0m", args{10 * time.Second, 0}},
		{"10s", args{10 * time.Second, 3}},
		{"2m", args{100 * time.Second, 0}},
		{"2m", args{100 * time.Second, 4}},
		{"1m40s", args{100 * time.Second, 5}},
		{"1m40s", args{100 * time.Second, 100}},

		{"17m", args{1000 * time.Second, 0}},
		{"16m40s", args{1000 * time.Second, 6}},
		{"16m40s", args{1000 * time.Second, 100}},
		{"167m", args{10000 * time.Second, 0}},
		{"166m40s", args{10000 * time.Second, 7}},
		{"166m40s", args{10000 * time.Second, 100}},
		{"28h", args{100000 * time.Second, 0}},
		{"27h46m", args{100000 * time.Second, 6}},
		{"27h46m", args{100000 * time.Second, 100}},
		{"278h", args{1000000 * time.Second, 0}},
		{"277h46m", args{1000000 * time.Second, 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmtAboveSeconds(tt.args.d, tt.args.fixedLen); got != tt.name {
				t.Errorf("fmtAboveSeconds() = %v, want %v", got, tt.name)
			}
		})
	}
}

//
func Test_fmtPortions(t *testing.T) {
	type args struct {
		valueLo  uint64
		valueHi  uint64
		fixedLen int
	}
	tests := []struct {
		name string
		args args
	}{
		{"0min", args{0, 0, 0}},
		{"0min", args{0, 0, 100}},
		{"1hr", args{31, 0, 0}},
		{"2hr", args{31, 1, 0}},
		{"3hr", args{31, 2, 1}},
		{"8hr", args{31, 7, 7}},
		{"8hr59min", args{59, 8, 8}},
		{"9hr59min", args{59, 9, 9}},
		{"80hr", args{29, 80, 8}},
		{"90hr29min", args{29, 90, 9}},
		{"100hr", args{29, 100, 9}},
		{"100hr29min", args{29, 100, 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmtPortions(tt.args.valueLo, "min", tt.args.valueHi, "hr", tt.args.fixedLen); got != tt.name {
				t.Errorf("fmtPortions() = %v, want %v", got, tt.name)
			}
		})
	}
}

type args struct {
	d time.Duration
}

var metricTests = []struct {
	name                string
	args                args
	want                time.Duration
	want1, want2, want3 string
}{
	/*
		NB! "µs" require 3 bytes, not 2 bytes, hence the difference with "ms" on formatting
	*/
	{"0.000µs", args{0}, time.Microsecond, "µs", "0µs", "0.0µs"},
	{"0.001µs", args{1 * time.Nanosecond}, time.Microsecond, "µs", "0µs", "0.0µs"},
	{"0.999µs", args{1*time.Microsecond - 1}, time.Microsecond, "µs", "1µs", "0.9µs"},
	{"1.000µs", args{1 * time.Microsecond}, time.Microsecond, "µs", "1µs", "1.0µs"},
	{"99.00µs", args{99 * time.Microsecond}, time.Microsecond, "µs", "0ms", "99µs"},
	{"99.87µs", args{99876 * time.Nanosecond}, time.Microsecond, "µs", "0ms", "100µs"},
	{"499.9µs", args{500*time.Microsecond - 1}, time.Microsecond, "µs", "0ms", "500µs"},
	{"500.0µs", args{500 * time.Microsecond}, time.Microsecond, "µs", "0ms", "500µs"},
	{"0.5000ms", args{500*time.Microsecond + 1}, time.Millisecond, "ms", "1ms", "0.5ms"},
	{"1.0000ms", args{1 * time.Millisecond}, time.Millisecond, "ms", "1ms", "1.0ms"},
	{"99.000ms", args{99 * time.Millisecond}, time.Millisecond, "ms", "0s", "99ms"},
	{"99.876ms", args{99876 * time.Microsecond}, time.Millisecond, "ms", "0s", "100ms"},
	{"499.99ms", args{500*time.Millisecond - time.Microsecond}, time.Millisecond, "ms", "0s", "500ms"},
	{"499.99ms", args{500*time.Millisecond - 1}, time.Millisecond, "ms", "0s", "500ms"},
	{"500.00ms", args{500 * time.Millisecond}, time.Millisecond, "ms", "0s", "500ms"},
	{"0.50000s", args{500*time.Millisecond + 1}, time.Second, "s", "1s", "0.50s"},
	{"1.00000s", args{1 * time.Second}, time.Second, "s", "1s", "1.00s"},
	{"30.0000s", args{30 * time.Second}, time.Second, "s", "0m", "30.0s"},
	{"59.0000s", args{59 * time.Second}, time.Second, "s", "1m", "59.0s"},
	{"300.000s", args{300 * time.Second}, time.Second, "s", "5m", "300s"},
	{"599.000s", args{10*time.Minute - time.Second}, time.Second, "s", "10m", "599s"},
	{"599.999s", args{10*time.Minute - 1}, time.Second, "s", "10m", "600s"},
	{"10m00s", args{10 * time.Minute}, time.Minute, "m", "", ""},
	{"10000m", args{10 * time.Minute}, time.Minute, "m", "", ""},
}

func Test_metricDurationBase(t *testing.T) {
	for _, tt := range metricTests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := metricDurationBase(tt.args.d)
			if got != tt.want {
				t.Errorf("metricDurationBase() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("metricDurationBase() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_fmtMetric(t *testing.T) {
	for _, tt := range metricTests {
		if tt.want == time.Minute {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			if got3 := fmtMetric(tt.args.d, tt.want, tt.want1, 8); got3 != tt.name {
				t.Errorf("fmtMetric() = %v, want %v", got3, tt.name)
			}
		})
	}
}

func Test_fmtMetricTooShort(t *testing.T) {
	for _, tt := range metricTests {
		if tt.want == time.Minute {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			if got3 := fmtMetric(tt.args.d, tt.want, tt.want1, 2); got3 != tt.want2 {
				t.Errorf("fmtMetric() = %v, want %v", got3, tt.want2)
			}
		})
	}
}

func Test_fmtMetricTooShort2(t *testing.T) {
	for _, tt := range metricTests {
		if tt.want == time.Minute {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			expectedLen := 5
			if tt.want == time.Microsecond {
				/*
					NB! "µs" require 3 bytes, not 2 bytes, hence the difference with "ms" on formatting
				*/
				expectedLen++
			}
			if got3 := fmtMetric(tt.args.d, tt.want, tt.want1, expectedLen); got3 != tt.want3 {
				t.Errorf("fmtMetric() = %v, want %v", got3, tt.want3)
			}
		})
	}
}
