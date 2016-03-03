package main

//# * * * * *  command to execute
// # │ │ │ │ │
// # │ │ │ │ │
// # │ │ │ │ └───── day of week (0 - 6) (0 to 6 are Sunday to Saturday, or use names; 7 is Sunday, the same as 0)
// # │ │ │ └────────── month (1 - 12)
// # │ │ └─────────────── day of month (1 - 31)
// # │ └──────────────────── hour (0 - 23)
// # └───────────────────────── min (0 - 59)

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pborman/uuid"
	"github.com/robfig/cron"

	"github.com/briandowns/raceway/config"
)

// ScheduleChan sends and receives Schedule instances
var ScheduleChan = make(chan Schedule, 0)
var c *cron.Cron

type status int

const (
	SCHEDULE_DISABLED status = iota
	SCHEDULE_ENABLED
	SCHEDULE_RUNNING
)

// Schedule represents a scheduled entry
type Schedule struct {
	ID       string
	Schedule string
	Status   status
	Task     Task
}

// newSchedule creates a new isntance "Schedule"
func newSchedule() *Schedule {
	return &Schedule{
		ID:     uuid.NewUUID().String(),
		Status: SCHEDULE_DISABLED,
	}
}

// scheduleDB returns a connection instance to Bolt
func scheduleDB() (*bolt.DB, error) {
	db, err := bolt.Open(Conf.Scheduler.SchedulerDBName, 0644,
		&bolt.Options{
			Timeout: time.Duration(Conf.Scheduler.SchedulerDBTimeout) * time.Second,
		})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}

// ScheduleDBPath returns the path to the scheduler db file
func ScheduleDBPath() (string, error) {
	db, err := scheduleDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	return db.Path(), nil
}

// String provides a string representation of a "Schedule" instance
func (s *Schedule) String() string {
	return fmt.Sprintf("ID: %s, Schedule: %s, Status: %d", s.ID, s.Schedule, s.Status)
}

// StartScheduler will start the scheduler process and handle requests
// in and out.
func StartScheduler(conf *config.Config) {
	c := cron.New()
	c.Start()
	db, err := scheduleDB()
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(conf.Scheduler.SchedulerDBBucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	/*err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("iops"))
		b.ForEach(func(k, v []byte) error {
			c.AddFunc(string(k), func() { "some sort of rally execution code..." })
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}*/
}

// ScheduleExists verifies whether or not a task already has a presence
func scheduleExists(scheduleUUID string) (bool, error) {
	db, err := scheduleDB()
	if err != nil {
		return false, err
	}
	defer db.Close()
	var found bool
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(Conf.Scheduler.SchedulerDBBucket))
		b.ForEach(func(k, v []byte) error {
			if string(k) == scheduleUUID {
				found = true
				return nil
			}
			return nil
		})
		return nil
	})
	return found, nil
}

// ScheduleTask creates a new schedule for a given task
func (s *Schedule) ScheduleTask(task Task, sched string) error {
	db, err := scheduleDB()
	if err != nil {
		return err
	}
	defer db.Close()
	ns := newSchedule()
	ns.Task = task
	return nil
}

// UnscheduleTask will remove a scheduled task
func UnscheduleTask(scheduleUUID string) error {
	db, err := scheduleDB()
	if err != nil {
		return err
	}
	defer db.Close()
	exists, err := scheduleExists(scheduleUUID)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Removing task: %s", scheduleUUID)
	}
	return nil
}

// ShowSchedules retunns a map of all entered schedules
func ShowSchedules() ([]Schedule, error) {
	db, err := scheduleDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var found []Schedule
	err = db.View(func(tx *bolt.Tx) error {
		var result Schedule
		b := tx.Bucket([]byte(Conf.Scheduler.SchedulerDBBucket))
		b.ForEach(func(k, v []byte) error {
			err = json.Unmarshal(v, &result)
			if err != nil {
				return err
			}
			if result.Status == SCHEDULE_ENABLED {
				found = append(found, result)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}
