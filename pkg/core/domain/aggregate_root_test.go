package domain_test

import (
	"testing"
	"time"

	"github.com/hex-api-go/pkg/core/domain"
)

type fakeEvent struct{}
func (f *fakeEvent) Payload() any {
	return "ok"
}
func (f *fakeEvent) OcurredOn() time.Time {
	return time.Now()
}
func (f *fakeEvent) Uuid() string {
	return "uuid"
}

func TestNewAggregateRoot(t *testing.T) {
	t.Parallel()
	got := domain.NewAggregateRoot("mock")
	if got == nil {
		t.Error("AggregateRoot should return a non-nil instance")
	}
}

func TestAggregateRoot_AddDomainEvent(t *testing.T) {
	t.Run("Should success when add event", func(t *testing.T) {
		t.Parallel()
		ar := domain.NewAggregateRoot("mock")
		got := ar.AddDomainEvent(&fakeEvent{})
		if got != nil{
			t.Errorf("AddDomainEvent should return a nil error, got: %v", got)
		}
	})

	t.Run("Should error when add nil event", func(t *testing.T) {
		t.Parallel()
		ar := domain.NewAggregateRoot("mock")
		got := ar.AddDomainEvent(nil)
		if got == nil{
			t.Error("AddDomainEvent should return a non-nil error")
		}
	})
}

func TestAggregateRoot_DomainEvents(t *testing.T){
	t.Run("Should return domain event list", func(t *testing.T) {
		ar := domain.NewAggregateRoot("mock")
		ar.AddDomainEvent(&fakeEvent{})
		got := ar.DomainEvents()
		
		if got == nil{
			t.Error("DomainEvents should return a non-nil list")
		}

		if len(got) != 1{
			t.Error("DomainEvents should return a exactly 1 event added")
		}
	})
}

func TestAggregateRoot_RemoveDomainEvent(t *testing.T){
	t.Run("Should return domain event list", func(t *testing.T) {
		t.Parallel()
		ar := domain.NewAggregateRoot("mock")
		ev:= &fakeEvent{}
		ar.AddDomainEvent(ev)
		ar.RemoveDomainEvent(ev.Uuid())
		
		if ar.DomainEvents() == nil{
			t.Error("DomainEvents should return a non-nil list")
		}

		if len(ar.DomainEvents()) != 0{
			t.Error("DomainEvents should return a exactly 0 event added")
		}
	})
}

func TestAggregateRoot_ClearEvents(t *testing.T){
	t.Run("Should return domain event list", func(t *testing.T) {
		t.Parallel()
		ar := domain.NewAggregateRoot("mock")
		ev:= &fakeEvent{}
		ar.AddDomainEvent(ev)
		ar.ClearEvents()
		
		if ar.DomainEvents() == nil{
			t.Error("DomainEvents should return a non-nil list")
		}

		if len(ar.DomainEvents()) != 0{
			t.Error("DomainEvents should return a exactly 0 event added")
		}
	})
}