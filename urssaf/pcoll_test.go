package urssaf

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseDateEffet(t *testing.T) {
	toTest := time.Date(2021, 4, 21, 0, 0, 0, 0, time.UTC)
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"test en minuscule", args{"21Apr2021"}, toTest},
		{"test en majuscule", args{"21APR2021"}, toTest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDateEffet(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDateEffet() = %v, want %v", got, tt.want)
			}
		})
	}
}
