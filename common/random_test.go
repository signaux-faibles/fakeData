package common

import (
	"github.com/golang-collections/collections/set"
	"testing"
)

func Test_RandEffectif(t *testing.T) {
	var input string
	var actual string
	var err error
	input = "52�047"
	actual, err = RandEffectif(input)
	if len(actual) == 5 && err == nil {
		t.Logf("[OK] randomize effectif %s has correct size, and is %s", input, actual)
	} else {
		t.Errorf("randomize effectif %s should have length of 5, and has length %d (%s),error is %s", input, len(actual), actual, err)
	}

	input = "52 047"
	actual, err = RandEffectif(input)
	if len(actual) == 5 && err == nil {
		t.Logf("[OK] randomize effectif %s has correct size, and is %s", input, actual)
	} else {
		t.Errorf("randomize effectif %s should have length of 5, and has length %d (%s),error is %s", input, len(actual), actual, err)
	}

	input = "toto va 17 fois à la plage"
	actual, err = RandEffectif(input)
	if input == actual && err != nil {
		t.Logf("[OK] effectif '%s' has non digit chars, so there's no randomization -> '%s', error is %s", input, actual, err)
	} else {
		t.Errorf("randomize effectif '%s' should not be modified and is '%s',error is %s", input, actual, err)
	}

	input = ""
	actual, err = RandEffectif(input)
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
		input := RandStringBytesRmndr(i%6 + 9)
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

		//anotherInput := RandStringBytesRmndr(i%20)
		//if anotherInput != input {
		//	anotherResult, err = FalsifyNumber(anotherInput)
		//	if err != nil {
		//		t.Errorf("error when falsify '%s', error is '%s'", input, err)
		//	}
		//	if result1 == anotherResult {
		//		t.Errorf("when falsify '%s' and '%s', results (here '%s') should be differents", input, anotherInput, anotherResult)
		//	}
		//}
	}

}
