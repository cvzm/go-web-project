package usecase

import "github.com/cvzm/go-web-project/doamin"

type eventUsecase struct {
	eventRepo doamin.EventRepository
}

func NewEventUsecase(repo doamin.EventRepository) doamin.EventUsecase {
	return &eventUsecase{eventRepo: repo}
}

func (u *eventUsecase) Create(cloudEvent doamin.CloudEvent) error {
	event, err := cloudEvent.Parse()
	if err != nil {
		return err
	}

	// More business logic
	// e.g slack notification

	return u.eventRepo.Save(&event)
}
