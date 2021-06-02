package common

import "testing"

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
