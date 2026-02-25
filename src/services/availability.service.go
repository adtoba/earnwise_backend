package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/adtoba/earnwise_backend/src/models"
	"gorm.io/gorm"
)

type Slot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func combineDateAndTime(date time.Time, timeStr string, loc *time.Location) time.Time {
	t, _ := time.Parse("15:04", timeStr)

	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		loc,
	)
}

func GenerateSlotsForDate(
	db *gorm.DB,
	expertID string,
	date time.Time,
	durationMins int,
	timezone string,
) ([]Slot, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	var expert models.ExpertProfile
	result := db.First(&expert, "id = ?", expertID)
	if result.Error != nil {
		return nil, result.Error
	}

	var availability *models.Availability
	weekdayStr := date.Weekday().String()

	for i := range expert.Availability {
		if expert.Availability[i].Day == weekdayStr {
			availability = &expert.Availability[i]
			break
		}
	}

	if availability == nil {
		return nil, errors.New("availability not found")
	}

	if availability.Status != "available" {
		return nil, errors.New("availability is not available")
	}

	fmt.Println(availability.Day)

	start := combineDateAndTime(date, availability.Start, loc)
	end := combineDateAndTime(date, availability.End, loc)

	var slots []Slot

	for {
		slotEnd := start.Add(time.Duration(durationMins) * time.Minute)

		if slotEnd.After(end) {
			break
		}
		slots = append(slots, Slot{
			Start: start.In(loc),
			End:   slotEnd.In(loc),
		})

		start = slotEnd
	}

	return slots, nil
}

func FilterConflictingSlots(
	db *gorm.DB,
	expertID string,
	slots []Slot,
) ([]Slot, error) {
	var calls []models.Call
	result := db.Where("expert_id = ? AND status IN ?", expertID, []string{models.CallStatusPending, models.CallStatusAccepted}).Find(&calls)
	if result.Error != nil {
		return nil, result.Error
	}

	var available []Slot

	for _, slot := range slots {
		conflict := false

		for _, call := range calls {
			callEnd := call.ScheduledAt.Add(
				time.Duration(call.DurationMins) * time.Minute,
			)

			if slot.Start.Before(callEnd) && slot.End.After(call.ScheduledAt) {
				conflict = true
				break
			}
		}

		if !conflict {
			available = append(available, slot)
		}
	}

	return available, nil
}

func GetAvailableSlots(
	db *gorm.DB,
	expertID string,
	date time.Time,
	durationMins int,
	timezone string,
) ([]Slot, error) {
	slots, err := GenerateSlotsForDate(db, expertID, date, durationMins, timezone)
	if err != nil {
		return nil, err
	}

	slots, err = FilterConflictingSlots(db, expertID, slots)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	var futureSlots []Slot
	for _, slot := range slots {
		if slot.Start.After(now) {
			futureSlots = append(futureSlots, slot)
		}
	}

	return futureSlots, nil
}
