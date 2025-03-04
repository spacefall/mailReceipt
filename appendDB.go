package main

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func dbAppend(id string, event string, data TrackData) {
	// Check for valid UUID
	if err := uuid.Validate(id); err != nil {
		//log.Println("Invalid UUID: " + id)
		log.Errorf("Couldn't update tracking table, as an invalid id (%s) was passed", id)
		return
	}

	// Insert into DB
	// I know that stitching event like this isn't the best way
	//goland:noinspection ALL
	_, err := dbpool.Exec(context.Background(), "UPDATE mail_receipts SET "+event+" = array_append("+event+", $1) WHERE id = $2", data, id)
	if err != nil {
		log.Errorf("Couldn't update tracking for id %s: %v", id, err)
		return
	}
}
