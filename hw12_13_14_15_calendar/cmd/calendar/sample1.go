package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/model"
	"github.com/randomurban/hw-test/hw12_13_14_15_calendar/internal/storage"
)

func sample1(ctx context.Context, store storage.EventStorage, logg *logger.Logger) {
	loc, _ := time.LoadLocation("Asia/Omsk")
	t1 := time.Date(2020, 10, 20, 15, 10, 0, 0, time.UTC)
	t2 := time.Date(2020, 10, 20, 23, 10, 0, 0, loc)
	nt := time.Minute * 15
	ev1 := &model.Event{
		ID:          0,
		Title:       "event1",
		Start:       t1,
		End:         t2,
		Owner:       1,
		Description: "desc1",
		NoticeTime:  nt,
	}
	id, err := store.Create(ctx, *ev1)
	if err != nil {
		logg.Error("failed to create event: " + err.Error())
	} else {
		logg.Info("created " + strconv.Itoa(int(id)))
	}

	ev1.Title = "updated"
	var ok bool
	ok, err = store.Update(ctx, id, *ev1)
	if err != nil {
		logg.Error("failed to update event: " + err.Error())
	}
	if !ok {
		logg.Error("failed to update event")
	} else {
		logg.Info("updated " + strconv.Itoa(int(id)))
	}

	// ok, err = store.Delete(ctx, id)
	// if err != nil {
	// 	logg.Error("failed to delete event: " + err.Error())
	// }
	// if !ok {
	// 	logg.Error("failed to delete event")
	// } else {
	// 	logg.Info("deleted " + strconv.Itoa(int(id)))
	// }

	var ev2 *model.Event
	ev2, err = store.GetByID(ctx, 2)
	if err != nil {
		logg.Error("failed to get event by id: " + err.Error())
	} else {
		logg.Info("id=2: " + ev2.Title)
	}

	var ev *[]model.Event
	ev, err = store.GetDayFromTo(ctx, t1, t2)
	if err != nil {
		logg.Error("failed to get event by id: " + err.Error())
	} else {
		fmt.Printf("%v", ev)
	}
}
