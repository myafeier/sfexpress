package sf_express

import "testing"

func TestQueryRouteInfo(t *testing.T) {
	req := new(RouteRequestBody)
	req.MethodType = 1
	req.TrackingNumber = "1000100000000172"
	req.TrackingType = 2

	result, err := QueryRouteInfo(req)
	if err != nil {
		t.Fatalf("%s", err.Error())
		return
	}
	for _, v := range result {
		t.Logf("result: %+v", v)
	}

}
