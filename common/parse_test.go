package common

import "testing"

func Test_selectOnlyDigits(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with not cutable space", args{"52�047"}, "52047"},
		{"with 0 after comma", args{"520,"}, "520"},
		{"with 1 after comma", args{"520,4"}, "520,4"},
		{"with 2 after comma", args{"520,47"}, "520,47"},
		{"with 3 after comma", args{"520,478"}, "520,478"},
		{"with dot", args{"520.478"}, "520.478"},
		{"with blank", args{"52 047"}, "52047"},
		{"with negative int", args{"-52 047"}, "-52047"},
		{"with negative float", args{"-52,047"}, "-52,047"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectOnlyDigits(tt.args.input); got != tt.want {
				t.Errorf("selectOnlyDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFloat(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"with comma", args{"520,54"}, 520.54},
		{"with dot", args{"520.54"}, 520.54},
		{"without comma or dot", args{"520"}, 520},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := parseFloat(tt.args.input); got != tt.want {
				if err != nil {
					t.Errorf("error when parseFloat(%v) -> err is %v", tt.args.input, err)
				}
				t.Errorf("parseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_floatToStringWith2Digits(t *testing.T) {
	type args struct {
		input float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with 2 digits", args{520.54}, "520,54"},
		{"with 3 digits", args{520.545}, "520,54"},
		{"with 0 digit", args{520}, "520,00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := floatToStringWith2Digits(tt.args.input); got != tt.want {
				t.Errorf("floatToStringWith2Digits() = %v, want %v", got, tt.want)
			}
		})
	}
}
