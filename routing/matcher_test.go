package routing

import (
    "testing"
    "github.com/zalando/skipper/filters"
)

func TestFailToParse(t *testing.T) {
    r := &Routing{}
    _, err := r.processData("invalid eskip document")
    if err == nil {
        t.Error("failed to fail")
    }
}

func TestCreateShuntBackend(t *testing.T) {
    r := &Routing{}
    m, err := r.processData(`Any() -> <shunt>`)
    if err != nil {
        t.Error(err)
    }

    testMatcher(t, castMatcher(m), &Route{"", true, nil})
}

func TestFailToParseBackend(t *testing.T) {
    r := &Routing{}
    m, err := r.processData(`Any() -> "invalid backend"`)
    if err != nil {
        t.Error(err)
    }

    testMatcher(t, castMatcher(m), nil)
}

func TestParseBackend(t *testing.T) {
    r := &Routing{}
    m, err := r.processData(`Any() -> "https://www.example.org"`)
    if err != nil {
        t.Error(err)
    }

    testMatcher(t, castMatcher(m), &Route{"https://www.example.org", false, nil})
}

func TestFilterNotFound(t *testing.T) {
    spec1 := &testFilter{name: "testFilter1"}
    spec2 := &testFilter{name: "testFilter2"}
    fr := make(filters.Registry)
    fr[spec1.Name()] = spec1
    fr[spec2.Name()] = spec2

    r := &Routing{filterRegistry: fr}
    m, err := r.processData(`Any() -> testFilter3() -> "https://www.example.org"`)
    if err != nil {
        t.Error(err)
    }

    testMatcher(t, castMatcher(m), nil)
}

func TestCreateFilters(t *testing.T) {
    spec1 := &testFilter{name: "testFilter1"}
    spec2 := &testFilter{name: "testFilter2"}
    fr := make(filters.Registry)
    fr[spec1.Name()] = spec1
    fr[spec2.Name()] = spec2

    r := &Routing{filterRegistry: fr}
    m, err := r.processData(`Any() -> testFilter1(1, "one") -> testFilter2(2, "two") -> "https://www.example.org"`)
    if err != nil {
        t.Error(err)
    }

    testMatcher(t, castMatcher(m), &Route{"https://www.example.org", false, []filters.Filter{
        &testFilter{name: "testFilter1", args: []interface{}{float64(1), "one"}},
        &testFilter{name: "testFilter2", args: []interface{}{float64(2), "two"}}}})
}
