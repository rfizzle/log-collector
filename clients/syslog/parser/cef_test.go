package parser

import "testing"

var cefMessage1 = "0|illusive|illusive|3.1.128.1719|illusive:heartbeat|Heartbeat|0|dvc=10.118.182.162 rt=1600239263565 cat=illusive:SYS"
var cefMessage2 = "CEF:0|Cool Vendor|Cool Product|1.0|FLAKY_EVENT|Something flaky happened.|3|requestClientApplication=Go-http-client/1.1 sourceAddress=127.0.0.1"
var cefMessage3 = `0|illusive|illusive|3.1.128.1719|illusive:audit|Audit|5|msg=theuser@domain.local logged out {User role \\= ROLE_ADMIN; Source address \\= 10.120.10.152}  dvc=10.105.33.50 rt=1600239250955 duser=theuser@domain.local cat=illusive:info outcome=SUCCESS`
var cefMessage4 = `illusive|illusive|3.1.128.1719|illusive:heartbeat|Heartbeat|0|dvc=10.118.182.162 rt=1600239263565 cat=illusive:SYS`

func TestParseCef(t *testing.T) {
	expectedVersion := "0"
	expectedDeviceVendor := "illusive"
	expectedDeviceEventClassId := "illusive:heartbeat"
	expectedName := "Heartbeat"

	if _, err := ParseCef(cefMessage1); err != nil {
		t.Fatalf("failed to parse CEF message: %v", err)
	}

	cefObj, err := cefStringToObject(cefMessage1)

	if err != nil {
		t.Errorf("failed to parse CEF message: %v", err)
	}

	if cefObj == nil {
		t.Errorf("null CEF message returned")
	}

	if cefObj.Version != expectedVersion {
		t.Errorf("cefStringToObject(message1).Version got %s; expected %s", cefObj.Version, expectedVersion)
	}

	if cefObj.DeviceVendor != expectedDeviceVendor {
		t.Errorf("cefStringToObject(message1).DeviceVendor got %s; expected %s", cefObj.DeviceVendor, expectedDeviceVendor)
	}

	if cefObj.Name != expectedName {
		t.Errorf("cefStringToObject(message1).Name got %s; expected %s", cefObj.Name, expectedName)
	}

	if cefObj.DeviceEventClassId != expectedDeviceEventClassId {
		t.Errorf("cefStringToObject(message1).DeviceEventClassId got %s; expected %s", cefObj.DeviceEventClassId, expectedDeviceEventClassId)
	}

}

func TestParseCef2(t *testing.T) {
	expectedVersion := "0"
	expectedDeviceVendor := "Cool Vendor"
	expectedDeviceProduct := "Cool Product"
	expectedDeviceEventClassId := "FLAKY_EVENT"

	if _, err := ParseCef(cefMessage2); err != nil {
		t.Errorf("failed to parse CEF message: %v", err)
	}

	cefObj, err := cefStringToObject(cefMessage2)

	if err != nil {
		t.Errorf("failed to parse CEF message: %v", err)
	}

	if cefObj == nil {
		t.Errorf("null CEF message returned")
	}

	if cefObj.Version != expectedVersion {
		t.Errorf("cefStringToObject(message2).Version got %s; expected %s", cefObj.Version, expectedVersion)
	}

	if cefObj.DeviceVendor != expectedDeviceVendor {
		t.Errorf("cefStringToObject(message2).DeviceVendor got %s; expected %s", cefObj.DeviceVendor, expectedDeviceVendor)
	}

	if cefObj.DeviceProduct != expectedDeviceProduct {
		t.Errorf("cefStringToObject(message2).DeviceProduct got %s; expected %s", cefObj.DeviceProduct, expectedDeviceProduct)
	}

	if cefObj.DeviceEventClassId != expectedDeviceEventClassId {
		t.Errorf("cefStringToObject(message2).DeviceEventClassId got %s; expected %s", cefObj.DeviceEventClassId, expectedDeviceEventClassId)
	}

}

func TestParseCef3(t *testing.T) {
	cefExpectedKeyValuePair1 := []string{"msg", `theuser@domain.local logged out {User role \\= ROLE_ADMIN; Source address \\= 10.120.10.152}`}
	cefExpectedKeyValuePair2 := []string{"duser", "theuser@domain.local"}
	cefExpectedKeyValuePair3 := []string{"outcome", "SUCCESS"}
	cefExpectedKeyValues := [][]string{cefExpectedKeyValuePair1, cefExpectedKeyValuePair2, cefExpectedKeyValuePair3}

	if _, err := ParseCef(cefMessage3); err != nil {
		t.Errorf("failed to parse CEF message: %v", err)
	}

	cefObj, err := cefStringToObject(cefMessage3)

	if err != nil {
		t.Errorf("failed to parse CEF message: %v", err)
	}

	if cefObj == nil {
		t.Errorf("null CEF message returned")
	}

	for _, v := range cefExpectedKeyValues {
		if cefObj.Extensions[v[0]] != v[1] {
			t.Errorf(`cefStringToObject(message3).Extensions["%s"] got %s; expected %s`, v[0], cefObj.Extensions[v[0]], v[1])
		}
	}

}

func TestParseCef4(t *testing.T) {
	if _, err := ParseCef(cefMessage4); err == nil {
		t.Errorf("failed to error on invalid CEF message: %v", err)
	}
}
