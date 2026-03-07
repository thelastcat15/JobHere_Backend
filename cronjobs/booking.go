package cronjobs

import (
	"sync"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var once sync.Once

func StartBookingCron(db *gorm.DB) {
	once.Do(func() {

		c := cron.New()

		c.AddFunc("@every 1m", func() {

			db.Exec(`
			WITH expired_bookings AS (
				UPDATE bookings
				SET status = 'EXPIRED',
				    updated_at = NOW()
				WHERE
				    status = 'CONFIRMED'
				    AND checkin_time IS NULL
				    AND NOW() > booked_time_start + (grace_minutes * INTERVAL '1 minute')
				RETURNING slot_id
			)
			UPDATE parking_slots
			SET status = 'available'
			WHERE id IN (SELECT slot_id FROM expired_bookings)
			`)

		})

		c.Start()
	})
}
