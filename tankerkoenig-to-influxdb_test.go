package main

import "testing"

func Test_dummy(t *testing.T) {
	result := 3
	if result != 3 {
		t.Error("incorrect result: expect 3, got", result)
	}
}

func Test_parseLine(t *testing.T) {
	result := 3
	if result != 3 {
		t.Error("incorrect result: expect 3, got", result)
	}
}

// date,station_uuid,diesel,e5,e10,dieselchange,e5change,e10change
// 2022-02-04 08:15:07+01,e5215cb1-30d3-4480-9ea1-07381cd0a492,1.609,1.769,1.709,0,1,1
// 2022-02-04 08:15:07+01,77c1259a-27f4-45c0-8a25-3089e23e8866,1.629,1.759,1.699,1,0,0
// 2022-02-04 08:15:07+01,e1a71869-0ddf-4f91-ac1b-1a341212712b,1.629,1.779,1.719,1,0,0
// 2022-02-04 08:15:07+01,005056ba-7cb6-1ed2-bceb-88651ca7cd30,1.599,1.739,1.679,1,0,0
