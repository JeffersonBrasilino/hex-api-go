package domain_test

import (
	"strings"
	"testing"

	"github.com/hex-api-go/pkg/core/domain"
)

func TestValidatorInstance(t *testing.T) {
	t.Run("returns non-nil singleton", func(t *testing.T) {
		t.Parallel()
		v := domain.ValidatorInstance()
		if v == nil {
			t.Error("ValidatorInstance() should not return nil")
		}
	})

	t.Run("returns same instance on multiple calls", func(t *testing.T) {
		t.Parallel()
		v1 := domain.ValidatorInstance()
		v2 := domain.ValidatorInstance()
		if v1 != v2 {
			t.Error("ValidatorInstance() should return the same instance")
		}
	})
}

func TestResetDomainValidatorForTest(t *testing.T) {
	t.Run("resets cache so schema is rebuilt", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withRequired struct {
			Name string `domainValidator:"required"`
		}
		dto := &withRequired{Name: "a"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected no failures for valid dto, got %v", result)
		}
		domain.ResetDomainValidatorForTest()
		result2, err2 := v.Validate(dto)
		if err2 != nil {
			t.Fatalf("Validate after reset should not return error: %v", err2)
		}
		if len(result2) != 0 {
			t.Errorf("expected no failures after reset, got %v", result2)
		}
	})
}

func TestDomainValidator_Validate(t *testing.T) {
	t.Run("returns nil map and nil error when data is nil", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		result, err := v.Validate(nil)
		if result != nil {
			t.Errorf("Validate(nil) should return nil map, got %v", result)
		}
		if err != nil {
			t.Errorf("Validate(nil) should return nil error, got %v", err)
		}
	})

	t.Run("returns empty map when struct has no domainValidator tags", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type noTags struct {
			Name string
		}
		dto := &noTags{Name: "x"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty result for struct without tags, got %v", result)
		}
	})

	t.Run("returns failure when required field is empty", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withRequired struct {
			Name string `domainValidator:"required"`
		}
		dto := &withRequired{Name: ""}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("expected one failure, got %d: %v", len(result), result)
		}
		if r, ok := result["Name"]; !ok {
			t.Errorf("expected failure for field Name, got keys: %v", result)
		} else if r.IsValid {
			t.Error("expected IsValid false for empty required field")
		} else if len(r.FailedValidators) != 1 || r.FailedValidators[0] != "required" {
			t.Errorf("expected FailedValidators [required], got %v", r.FailedValidators)
		}
	})

	t.Run("returns empty map when required field is set", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withRequired struct {
			Name string `domainValidator:"required"`
		}
		dto := &withRequired{Name: "ok"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected no failures when required is set, got %v", result)
		}
	})

	t.Run("returns failure when gte validator fails", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withGte struct {
			Value string `domainValidator:"gte=3"`
		}
		dto := &withGte{Value: "ab"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("expected one failure for gte=3 with len 2, got %v", result)
		}
		if r, ok := result["Value"]; !ok {
			t.Errorf("expected failure for Value, got %v", result)
		} else if len(r.FailedValidators) != 1 || !strings.Contains(r.FailedValidators[0], "gte") {
			t.Errorf("expected gte in FailedValidators, got %v", r.FailedValidators)
		}
	})

	t.Run("returns empty map when gte validator passes", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withGte struct {
			Value string `domainValidator:"gte=2"`
		}
		dto := &withGte{Value: "ab"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected no failures when gte passes, got %v", result)
		}
	})

	t.Run("returns failure when lte validator fails", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withLte struct {
			Value string `domainValidator:"lte=2"`
		}
		dto := &withLte{Value: "abc"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("expected one failure for lte=2 with len 3, got %v", result)
		}
		if r, ok := result["Value"]; !ok {
			t.Errorf("expected failure for Value, got %v", result)
		} else if len(r.FailedValidators) != 1 || !strings.Contains(r.FailedValidators[0], "lte") {
			t.Errorf("expected lte in FailedValidators, got %v", r.FailedValidators)
		}
	})

	t.Run("returns failure when len validator fails", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withLen struct {
			Value string `domainValidator:"len=3"`
		}
		dto := &withLen{Value: "abc"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("expected one failure for len=3 with value len 3, got %v", result)
		}
		if r, ok := result["Value"]; !ok {
			t.Errorf("expected failure for Value, got %v", result)
		} else if len(r.FailedValidators) != 1 || !strings.Contains(r.FailedValidators[0], "len") {
			t.Errorf("expected len in FailedValidators, got %v", r.FailedValidators)
		}
	})

	t.Run("returns error when validator name does not exist", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type invalidValidator struct {
			Name string `domainValidator:"nonexistent"`
		}
		dto := &invalidValidator{Name: "x"}
		result, err := v.Validate(dto)
		if err == nil {
			t.Error("Validate should return error for nonexistent validator")
		}
		if result != nil {
			t.Errorf("Validate should return nil result on error, got %v", result)
		}
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("error should mention validator does not exist: %v", err)
		}
	})

	t.Run("validates nested struct fields", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type nested struct {
			Street string `domainValidator:"required"`
		}
		type withNested struct {
			Name    string `domainValidator:"required"`
			Address nested
		}
		dto := &withNested{Name: "ok", Address: nested{Street: ""}}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("expected one failure for nested Street, got %v", result)
		}
		if _, ok := result["Address.Street"]; !ok {
			t.Errorf("expected failure at Address.Street, got keys: %v", result)
		}
	})

	t.Run("accepts pointer to struct", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withRequired struct {
			Name string `domainValidator:"required"`
		}
		dto := &withRequired{Name: "a"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected no failures for valid pointer, got %v", result)
		}
	})

	t.Run("uses cached schema on second validate of same type", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type withRequired struct {
			Name string `domainValidator:"required"`
		}
		dto1 := &withRequired{Name: "first"}
		dto2 := &withRequired{Name: "second"}
		_, _ = v.Validate(dto1)
		result2, err := v.Validate(dto2)
		if err != nil {
			t.Fatalf("second Validate should not return error: %v", err)
		}
		if len(result2) != 0 {
			t.Errorf("second validate should succeed: %v", result2)
		}
	})

	t.Run("multiple validators on same field all run", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type multi struct {
			Value string `domainValidator:"required,gte=2"`
		}
		dtoEmpty := &multi{Value: ""}
		resultEmpty, _ := v.Validate(dtoEmpty)
		if len(resultEmpty) != 1 {
			t.Fatalf("expected one failure for empty Value, got %v", resultEmpty)
		}
		r := resultEmpty["Value"]
		if len(r.FailedValidators) < 1 {
			t.Errorf("expected at least one failed validator, got %v", r.FailedValidators)
		}
		dtoShort := &multi{Value: "a"}
		resultShort, _ := v.Validate(dtoShort)
		if len(resultShort) != 1 {
			t.Fatalf("expected one failure for short Value, got %v", resultShort)
		}
		dtoOk := &multi{Value: "ab"}
		resultOk, _ := v.Validate(dtoOk)
		if len(resultOk) != 0 {
			t.Errorf("expected no failures for valid Value, got %v", resultOk)
		}
	})

	t.Run("field with empty domainValidator tag is skipped", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type emptyTag struct {
			Name string `domainValidator:""`
		}
		dto := &emptyTag{Name: "x"}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected no validation for empty tag, got %v", result)
		}
	})

	t.Run("validates nested struct via pointer field", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type inner struct {
			Street string `domainValidator:"required"`
		}
		type withPtrNested struct {
			Address *inner
		}
		dto := &withPtrNested{Address: &inner{Street: ""}}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("expected one failure for nested pointer struct, got %v", result)
		}
		if _, ok := result["Address.Street"]; !ok {
			t.Errorf("expected failure at Address.Street, got keys: %v", result)
		}
	})

	t.Run("validates nested struct via pointer field when valid", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type inner struct {
			Street string `domainValidator:"required"`
		}
		type withPtrNested struct {
			Address *inner
		}
		dto := &withPtrNested{Address: &inner{Street: "valid"}}
		result, err := v.Validate(dto)
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected no failures when nested pointer struct is valid, got %v", result)
		}
	})
}
