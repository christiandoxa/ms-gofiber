package structbase

import (
	"reflect"
	"testing"

	v10 "github.com/go-playground/validator/v10"

	"ms-gofiber/internal/dto"
)

type stubStructLevel struct {
	current  reflect.Value
	reported int
}

func (s *stubStructLevel) Validator() *v10.Validate { return v10.New() }
func (s *stubStructLevel) Top() reflect.Value       { return s.current }
func (s *stubStructLevel) Parent() reflect.Value    { return s.current }
func (s *stubStructLevel) Current() reflect.Value   { return s.current }
func (s *stubStructLevel) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return field, field.Kind(), false
}
func (s *stubStructLevel) ReportError(interface{}, string, string, string, string) { s.reported++ }
func (s *stubStructLevel) ReportValidationErrors(string, string, v10.ValidationErrors) {
	s.reported++
}

func TestTodoUpsertStructRule(t *testing.T) {
	wrong := &stubStructLevel{current: reflect.ValueOf(struct{}{})}
	TodoUpsertStructRule(wrong)
	if wrong.reported != 0 {
		t.Fatalf("expected 0 report for wrong type, got %d", wrong.reported)
	}

	trim := &stubStructLevel{current: reflect.ValueOf(dto.TodoUpsertRequest{Title: " abc "})}
	TodoUpsertStructRule(trim)
	if trim.reported != 1 {
		t.Fatalf("expected 1 trim report, got %d", trim.reported)
	}

	blank := &stubStructLevel{current: reflect.ValueOf(dto.TodoUpsertRequest{Title: "   "})}
	TodoUpsertStructRule(blank)
	if blank.reported != 2 {
		t.Fatalf("expected 2 reports for blank+trim, got %d", blank.reported)
	}

	ok := &stubStructLevel{current: reflect.ValueOf(dto.TodoUpsertRequest{Title: "abc"})}
	TodoUpsertStructRule(ok)
	if ok.reported != 0 {
		t.Fatalf("expected 0 report for valid title, got %d", ok.reported)
	}
}

func TestPrepareExampleStructRuleAndCheckOsType(t *testing.T) {
	wrong := &stubStructLevel{current: reflect.ValueOf(struct{}{})}
	PrepareExampleStructRule(wrong)
	if wrong.reported != 0 {
		t.Fatalf("expected 0 report for wrong type, got %d", wrong.reported)
	}

	invalidApp := &stubStructLevel{current: reflect.ValueOf(dto.PrepareExampleRequest{TerminalType: "APP", OsType: "BAD"})}
	PrepareExampleStructRule(invalidApp)
	if invalidApp.reported != 1 {
		t.Fatalf("expected invalid app osType report, got %d", invalidApp.reported)
	}

	nonAppWithOS := &stubStructLevel{current: reflect.ValueOf(dto.PrepareExampleRequest{TerminalType: "WEB", OsType: "IOS", OsVersion: "1"})}
	PrepareExampleStructRule(nonAppWithOS)
	if nonAppWithOS.reported != 2 {
		t.Fatalf("expected non-app report count 2, got %d", nonAppWithOS.reported)
	}

	validApp := &stubStructLevel{current: reflect.ValueOf(dto.PrepareExampleRequest{TerminalType: "APP", OsType: "ANDROID", OsVersion: "14"})}
	PrepareExampleStructRule(validApp)
	if validApp.reported != 0 {
		t.Fatalf("expected 0 report for valid app os fields, got %d", validApp.reported)
	}

	// Cover direct helper path as well.
	h := &stubStructLevel{current: reflect.ValueOf(struct{}{})}
	checkOsType(h, "WEB", "", "")
	if h.reported != 0 {
		t.Fatalf("expected no helper reports, got %d", h.reported)
	}
}
