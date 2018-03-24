package models

import (
	"testing"
)

func TestHasRecordTypeName(t *testing.T) {
	x := &RecordConfig{
		Type: "A",
	}
	x.SetLabel("@", "example.tld")
	dc := DomainConfig{}
	if dc.HasRecordTypeName("A", "@") {
		t.Errorf("%v: expected (%v) got (%v)\n", dc.Records, false, true)
	}
	dc.Records = append(dc.Records, x)
	if !dc.HasRecordTypeName("A", "@") {
		t.Errorf("%v: expected (%v) got (%v)\n", dc.Records, true, false)
	}
	if dc.HasRecordTypeName("AAAA", "@") {
		t.Errorf("%v: expected (%v) got (%v)\n", dc.Records, false, true)
	}
}

func TestRR(t *testing.T) {
	experiment := RecordConfig{
		Type:         "A",
		TTL:          0,
		MxPreference: 0,
	}
	experiment.SetLabel("foo", "example.com")
	experiment.SetTarget("1.2.3.4")
	expected := "foo.example.com.\t300\tIN\tA\t1.2.3.4"
	found := experiment.ToRR().String()
	if found != expected {
		t.Errorf("RR expected (%#v) got (%#v)\n", expected, found)
	}

	experiment = RecordConfig{
		Type:    "CAA",
		TTL:     300,
		CaaTag:  "iodef",
		CaaFlag: 1,
	}
	experiment.SetLabel("@", "example.com")
	experiment.SetTarget("mailto:test@example.com")
	expected = "example.com.\t300\tIN\tCAA\t1 iodef \"mailto:test@example.com\""
	found = experiment.ToRR().String()
	if found != expected {
		t.Errorf("RR expected (%#v) got (%#v)\n", expected, found)
	}

	experiment = RecordConfig{
		Type:             "TLSA",
		TTL:              300,
		TlsaUsage:        0,
		TlsaSelector:     0,
		TlsaMatchingType: 1,
	}
	experiment.SetLabel("@", "_443._tcp.example.com")
	experiment.SetTarget("abcdef0123456789")
	expected = "_443._tcp.example.com.\t300\tIN\tTLSA\t0 0 1 abcdef0123456789"
	found = experiment.ToRR().String()
	if found != expected {
		t.Errorf("RR expected (%#v) got (%#v)\n", expected, found)
	}
}

func TestDowncase(t *testing.T) {
	dc := DomainConfig{Records: Records{
		&RecordConfig{Type: "MX"},
		&RecordConfig{Type: "MX"},
	}}
	dc.Records[0].SetLabel("lower", "example.tld")
	dc.Records[0].SetTarget("targetmx")
	dc.Records[1].SetLabel("UPPER", "example.tld")
	dc.Records[1].SetTarget("TARGETMX")
	downcase(dc.Records)
	if !dc.HasRecordTypeName("MX", "lower") {
		t.Errorf("%v: expected (%v) got (%v)\n", dc.Records, false, true)
	}
	if !dc.HasRecordTypeName("MX", "upper") {
		t.Errorf("%v: expected (%v) got (%v)\n", dc.Records, false, true)
	}
	if dc.Records[0].GetTargetField() != "targetmx" {
		t.Errorf("%v: target0 expected (%v) got (%v)\n", dc.Records, "targetmx", dc.Records[0].GetTargetField())
	}
	if dc.Records[1].GetTargetField() != "targetmx" {
		t.Errorf("%v: target1 expected (%v) got (%v)\n", dc.Records, "targetmx", dc.Records[1].GetTargetField())
	}
}
