package shared

import "testing"

func TestGetFunctionName(t *testing.T) {
	if got := helperFunctionName(); got != "helperFunctionName" {
		t.Fatalf("unexpected function name: %s", got)
	}
}

func helperFunctionName() string {
	return GetFunctionName()
}

type mapper struct{}

func (mapper) FromModel(value string) int { return len(value) }
func (mapper) ToModel() string            { return "value" }

func TestMapperInterface(t *testing.T) {
	var m Mapper[string, int] = mapper{}
	if m.FromModel("abc") != 3 || m.ToModel() != "value" {
		t.Fatalf("unexpected mapper behavior")
	}
}
