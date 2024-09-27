package usecase

import "github.com/cvzm/go-web-project/domain"

type eventUsecase struct {
	eventRepo domain.EventRepository
}

func NewEventUsecase(repo domain.EventRepository) domain.EventUsecase {
	return &eventUsecase{eventRepo: repo}
}

func (u *eventUsecase) Save(cloudEvent domain.CloudEvent) error {
	event, err := cloudEvent.Parse()
	if err != nil {
		return err
	}

	// TODO: Check idempotence

	// More business logic
	// e.g slack notification

	return u.eventRepo.Save(&event)
}
