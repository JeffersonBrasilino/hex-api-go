package domain_test

import (
	"strings"
	"testing"

	"github.com/hex-api-go/pkg/core/domain"
)

func TestRequiredValidator(t *testing.T) {
	t.Run("fails when string field is empty", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Name string `domainValidator:"required"`
		}
		result, err := v.Validate(&s{Name: ""})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 || result["Name"].IsValid {
			t.Errorf("required should fail for empty string, got %v", result)
		}
	})

	t.Run("passes when string field is non-empty", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Name string `domainValidator:"required"`
		}
		result, err := v.Validate(&s{Name: "x"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("required should pass for non-empty string, got %v", result)
		}
	})

	t.Run("fails when slice field is nil", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Tags []string `domainValidator:"required"`
		}
		result, err := v.Validate(&s{Tags: nil})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 || result["Tags"].IsValid {
			t.Errorf("required should fail for nil slice, got %v", result)
		}
	})

	t.Run("passes when slice field is non-nil", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Tags []string `domainValidator:"required"`
		}
		result, err := v.Validate(&s{Tags: []string{}})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("required should pass for non-nil slice, got %v", result)
		}
	})

	t.Run("fails when map field is nil", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Meta map[string]int `domainValidator:"required"`
		}
		result, err := v.Validate(&s{Meta: nil})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 || result["Meta"].IsValid {
			t.Errorf("required should fail for nil map, got %v", result)
		}
	})

	t.Run("fails when interface field is nil", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Data interface{} `domainValidator:"required"`
		}
		result, err := v.Validate(&s{Data: nil})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 || result["Data"].IsValid {
			t.Errorf("required should fail for nil interface, got %v", result)
		}
	})

	t.Run("fails when nested struct field is zero value", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type inner struct {
			X string `domainValidator:"required"`
		}
		type s struct {
			Inner inner
		}
		result, err := v.Validate(&s{Inner: inner{}})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 || result["Inner.X"].IsValid {
			t.Errorf("required should fail for zero nested struct field, got %v", result)
		}
	})
}

func TestGteValidator(t *testing.T) {
	t.Run("fails when value length is less than param", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"gte=3"`
		}
		result, err := v.Validate(&s{Value: "ab"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("gte=3 should fail for len 2, got %v", result)
		}
	})

	t.Run("passes when value length is greater or equal to param", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"gte=2"`
		}
		result, err := v.Validate(&s{Value: "ab"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("gte=2 should pass for len 2, got %v", result)
		}
	})

	t.Run("fails when param is not a valid integer", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"gte=abc"`
		}
		result, err := v.Validate(&s{Value: "anything"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("gte=abc should fail validation, got %v", result)
		}
		if r, ok := result["Value"]; ok && len(r.FailedValidators) != 1 {
			t.Errorf("expected one failed validator, got %v", r.FailedValidators)
		}
	})
}

func TestLteValidator(t *testing.T) {
	t.Run("fails when value length is greater than param", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"lte=2"`
		}
		result, err := v.Validate(&s{Value: "abc"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("lte=2 should fail for len 3, got %v", result)
		}
	})

	t.Run("passes when value length is less or equal to param", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"lte=3"`
		}
		result, err := v.Validate(&s{Value: "ab"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("lte=3 should pass for len 2, got %v", result)
		}
	})

	t.Run("fails when param is not a valid integer", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"lte=xyz"`
		}
		result, err := v.Validate(&s{Value: "a"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("lte=xyz should fail validation, got %v", result)
		}
		if r, ok := result["Value"]; ok && len(r.FailedValidators) != 1 {
			t.Errorf("expected one failed validator, got %v", r.FailedValidators)
		}
	})
}

func TestLenValidator(t *testing.T) {
	t.Run("fails when value length equals param", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"len=3"`
		}
		result, err := v.Validate(&s{Value: "abc"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("len=3 should fail when length is 3, got %v", result)
		}
		if r, ok := result["Value"]; !ok || !strings.Contains(r.FailedValidators[0], "len") {
			t.Errorf("expected len validator to fail, got %v", result)
		}
	})

	t.Run("passes when value length is not equal to param", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"len=3"`
		}
		result, err := v.Validate(&s{Value: "ab"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("len=3 should pass when length is not 3, got %v", result)
		}
	})

	t.Run("fails when param is not a valid integer", func(t *testing.T) {
		domain.ResetDomainValidatorForTest()
		v := domain.ValidatorInstance()
		type s struct {
			Value string `domainValidator:"len=invalid"`
		}
		result, err := v.Validate(&s{Value: "x"})
		if err != nil {
			t.Fatalf("Validate should not return error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("len=invalid should fail validation, got %v", result)
		}
		if r, ok := result["Value"]; ok && len(r.FailedValidators) != 1 {
			t.Errorf("expected one failed validator, got %v", r.FailedValidators)
		}
	})
}
