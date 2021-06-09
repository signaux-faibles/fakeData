package common

import (
	"fmt"
	"github.com/golang-collections/collections/set"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_RandEffectif(t *testing.T) {
	var input string
	var actual string
	var err error
	input = "52�047"
	actual, err = RandInt(input)
	if len(actual) == 5 && err == nil {
		t.Logf("[OK] randomize effectif %s has correct size, and is %s", input, actual)
	} else {
		t.Errorf("randomize effectif %s should have length of 5, and has length %d (%s),error is %s", input, len(actual), actual, err)
	}

	input = "52 047"
	actual, err = RandInt(input)
	if len(actual) == 5 && err == nil {
		t.Logf("[OK] randomize effectif %s has correct size, and is %s", input, actual)
	} else {
		t.Errorf("randomize effectif %s should have length of 5, and has length %d (%s),error is %s", input, len(actual), actual, err)
	}

	input = "toto va 17 fois à la plage"
	actual, err = RandInt(input)
	if input == actual && err != nil {
		t.Logf("[OK] effectif '%s' has non digit chars, so there's no randomization -> '%s', error is %s", input, actual, err)
	} else {
		t.Errorf("randomize effectif '%s' should not be modified and is '%s',error is %s", input, actual, err)
	}

	input = ""
	actual, err = RandInt(input)
	if len(actual) == 0 && err == nil {
		t.Logf("[OK] randomize effectif '%s' is empty", input)
	} else {
		t.Errorf("randomize effectif '%s' should be empty and is '%s, error is %s", input, actual, err)
	}
}

func Test_FalsifyNumber(t *testing.T) {
	var result1, result2 string
	var err error
	var generated = set.New()
	var falsified = set.New()

	for i := 0; i < 100000; i++ {
		input := RandStringBytesRmndr(10)
		result1, err = FalsifyNumber(input)
		if !generated.Has(input) {
			generated.Insert(input)
			if falsified.Has(result1) {
				t.Errorf("when falsify '%s', falsified results (here '%s') already exist", input, result1)
			}
			falsified.Insert(result1)
		}
		if err != nil {
			t.Errorf("error when falsify '%s', error is '%s'", input, err)
		}
		if len(result1) != len(input) {
			t.Errorf("when falsify '%s', result (here '%s') must have same size", input, result1)
		}
		if input == result1 {
			t.Errorf("when falsify '%s', result (here '%s') should be different", input, result1)
		}

		result2, err = FalsifyNumber(input)
		if err != nil {
			t.Errorf("error when falsify '%s', error is '%s'", input, err)
		}
		if result1 != result2 {
			t.Errorf("when falsify '%s', results (here '%s' and '%s') should be equals", input, result1, result2)
		}
	}
}

func Test_randDateAround(t *testing.T) {
	now := time.Now()
	type args struct {
		input time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"test now", args{now}, now},
		{"test one day ago", args{now.AddDate(0, 0, -1)}, now},
		{"test one month ago", args{now.AddDate(0, -1, 0)}, now},
		{"test one year ago", args{now.AddDate(-1, 0, 0)}, now},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randDateAround(tt.args.input); !got.Before(tt.want) {
				t.Errorf("randDateAround(%v) = %v, should be before %v", tt.args, got, tt.want)
			}
		})
	}
}

func Test_randDateAround_inFuture(t *testing.T) {
	now := time.Now()
	type args struct {
		input time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"test in one year", args{now.AddDate(1, 0, 0)}, now},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randDateAround(tt.args.input); !got.After(tt.want) {
				t.Errorf("randDateAround(%v) = %v, should be before %v", tt.args, got, tt.want)
			}
		})
	}
}

func Test_randDateAround_withClosure(t *testing.T) {
	now := time.Now()
	type args struct {
		input time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{"test now", args{now}, now},
		{"test one day ago", args{now.AddDate(0, 0, -1)}, now},
		{"test one month ago", args{now.AddDate(0, -1, 0)}, now},
		{"test one year ago", args{now.AddDate(-1, 0, 0)}, now},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randDateAround(tt.args.input)
			assert.Condition(t,
				func() bool {
					return got.Before(tt.want)
				},
				fmt.Sprintf("Result -> %v, should be before %v", got, tt.want))
		})
	}
}

func Test_RandDateAroundAsString(t *testing.T) {
	type args struct {
		value  string
		layout string
	}
	tests := []struct {
		name string
		args args
	}{
		{"pcoll format", args{"22MAR2021", "02Jan2006"}},
		{"empty date", args{"", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RandDate(tt.args.layout, tt.args.value)
			assert.Nil(t, err, "error when parsing %s with layout as %s", tt.args.value, tt.args.layout)
			parsed, err := time.Parse(tt.args.layout, result)
			assert.Nil(t, err, "error when parsing result %s with layout as %s", parsed, tt.args.layout)
		})
	}
}
