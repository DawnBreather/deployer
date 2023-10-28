package deployment_test

import (
	. "github.com/DawnBreather/go-commons/app/deployer/commons/models/deployment"
	"reflect"
	"testing"
)

func TestParseHttpStatusCodeRanges(t *testing.T) {
	t.Parallel()
	rangeString := make([]string, 5)
	rangeString[0] = "200-399,401"
	rangeString[1] = "301,401-404,305,200-399"
	rangeString[2] = "300"
	rangeString[3] = ""
	rangeString[4] = "19203"

	//optionDefault := "200-399"

	wantRange := make([]ValidStatusCodes, 5)
	wantRange[0] = ValidStatusCodes{
		AtomicRanges: []ValidStatusCodesAtomicRange{
			{
				Min: 200,
				Max: 399,
			},
		},
		AtomicValues: []int{
			401,
		},
	}

	wantRange[1] = ValidStatusCodes{
		AtomicRanges: []ValidStatusCodesAtomicRange{
			{
				Min: 401,
				Max: 404,
			},
			{
				Min: 200,
				Max: 399,
			},
		},
		AtomicValues: []int{
			301, 305,
		},
	}

	wantRange[2] = ValidStatusCodes{
		AtomicValues: []int{
			300,
		},
	}

	wantRange[3] = ValidStatusCodes{
		AtomicRanges: []ValidStatusCodesAtomicRange{
			{
				Min: 200,
				Max: 399,
			},
		},
	}

	wantRange[4] = ValidStatusCodes{
		AtomicRanges: []ValidStatusCodesAtomicRange{
			{
				Min: 200,
				Max: 399,
			},
		},
	}

	gotRange := [5]ValidStatusCodes{}
	gotRange[0] = ValidStatusCodes{}
	gotRange[1] = ValidStatusCodes{}
	gotRange[2] = ValidStatusCodes{}
	gotRange[3] = ValidStatusCodes{}
	gotRange[4] = ValidStatusCodes{}

	//ParseHttpStatusCode(rangeString[0], &gotRange[0], "", -1)
	//ParseHttpStatusCode(rangeString[1], &gotRange[1], "", -1)
	//ParseHttpStatusCode(rangeString[2], &gotRange[2], "", -1)
	//ParseHttpStatusCode(rangeString[3], &gotRange[3], "", -1)
	//ParseHttpStatusCode(rangeString[4], &gotRange[4], "", -1)

	for index, theRangeString := range rangeString {
		ParseHttpStatusCode(theRangeString, &gotRange[index], "", -1)
	}

	for index, theGotRange := range gotRange {
		if !reflect.DeepEqual(theGotRange, wantRange[index]) {
			t.Errorf("initial: { %s } | wanted { %s }, but gotProbeResult { %s }", rangeString[index], wantRange[index].ToString(), theGotRange.ToString())
		}
	}
	//if !reflect.DeepEqual(gotRange[0], wantRange[0]) {
	//	t.Errorf("initial: { %s } | wanted { %s }, but gotProbeResult { %s }", range1String, wantRange1.ToString(), gotRange1.ToString())
	//}
	//
	//if !reflect.DeepEqual(gotRange[1], wantRange[1]) {
	//	t.Errorf("initial: { %s } | wanted { %s }, but gotProbeResult { %s }", rangeString[1], wantRange[1].ToString(), gotRange[1].ToString())
	//}
	//
	//if !reflect.DeepEqual(gotRange[2], wantRange[2]) {
	//	t.Errorf("initial: { %s } | wanted { %s }, but gotProbeResult { %s }", rangeString[2], wantRange[2].ToString(), gotRange[2].ToString())
	//}
	//
	//if !reflect.DeepEqual(gotRange[3], wantRange[3]) {
	//	t.Errorf("initial: { %s } | wanted { %s }, but gotProbeResult { %s }", rangeString[3], wantRange[3].ToString(), gotRange[3].ToString())
	//}
	//
	//if !reflect.DeepEqual(gotRange[4], wantRange[4]) {
	//	t.Errorf("initial: { %s } | wanted { %s }, but gotProbeResult { %s }", rangeString[4], wantRange[4].ToString(), gotRange[4].ToString())
	//}

	probes := [5][]int{}
	probes[0] = []int{202, 401, 300, 400}
	probes[1] = []int{200, 305, 199, 201, 300, 301}
	probes[2] = []int{300, 200}
	probes[3] = []int{-1, 199, 200, 399, 400}
	probes[4] = []int{1, 19203, 200, 399, 400}

	wantProbeResult := [5][]bool{}
	wantProbeResult[0] = []bool{true, true, true, false}
	wantProbeResult[1] = []bool{true, true, false, true, true, true}
	wantProbeResult[2] = []bool{true, false}
	wantProbeResult[3] = []bool{false, false, true, true, false}
	wantProbeResult[4] = []bool{false, false, true, true, false}

	gotProbeResult := [5][]bool{}
	gotProbeResult[0] = []bool{}
	gotProbeResult[1] = []bool{}
	gotProbeResult[2] = []bool{}
	gotProbeResult[3] = []bool{}
	gotProbeResult[4] = []bool{}

	for index, theProbes := range probes {
		for _, probe := range theProbes {
			gotProbeResult[index] = append(gotProbeResult[index], gotRange[index].IsValid(probe))
		}
	}

	for index, probeResult := range gotProbeResult {
		if !reflect.DeepEqual(probeResult, wantProbeResult[index]) {
			t.Errorf("index: { %d } | range: { %v } | initial: { %v } | wanted { %v }, but gotProbeResult { %v }", index, gotRange[index].ToString(), probes[index], wantProbeResult[index], probeResult)
		}
	}

}
