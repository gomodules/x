package time

import (
	"testing"
	"time"
)

func TestAdjustForWeekend(t *testing.T) {
	type args struct {
		now time.Time
		adj WeekendAdjustment
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Monday-NoChange",
			args: args{
				now: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
				adj: NoChange,
			},
			want: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Monday-Before",
			args: args{
				now: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
				adj: Before,
			},
			want: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Monday-After",
			args: args{
				now: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
				adj: After,
			},
			want: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Friday-NoChange",
			args: args{
				now: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
				adj: NoChange,
			},
			want: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Friday-Before",
			args: args{
				now: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
				adj: Before,
			},
			want: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Friday-After",
			args: args{
				now: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
				adj: After,
			},
			want: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Saturday-NoChange",
			args: args{
				now: time.Date(2021, 12, 4, 10, 0, 0, 0, time.UTC),
				adj: NoChange,
			},
			want: time.Date(2021, 12, 4, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Saturday-Before",
			args: args{
				now: time.Date(2021, 12, 4, 10, 0, 0, 0, time.UTC),
				adj: Before,
			},
			want: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Saturday-After",
			args: args{
				now: time.Date(2021, 12, 4, 10, 0, 0, 0, time.UTC),
				adj: After,
			},
			want: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Sunday-NoChange",
			args: args{
				now: time.Date(2021, 12, 5, 10, 0, 0, 0, time.UTC),
				adj: NoChange,
			},
			want: time.Date(2021, 12, 5, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Sunday-Before",
			args: args{
				now: time.Date(2021, 12, 5, 10, 0, 0, 0, time.UTC),
				adj: Before,
			},
			want: time.Date(2021, 12, 3, 10, 0, 0, 0, time.UTC),
		},
		{
			name: "Sunday-After",
			args: args{
				now: time.Date(2021, 12, 5, 10, 0, 0, 0, time.UTC),
				adj: After,
			},
			want: time.Date(2021, 12, 6, 10, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AdjustForWeekend(tt.args.now, tt.args.adj); !got.Equal(tt.want) {
				t.Errorf("AdjustForWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}
