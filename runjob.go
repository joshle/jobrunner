// A job runner for executing scheduled or ad-hoc tasks asynchronously from HTTP requests.
//
// It adds a couple of features on top of the Robfig cron package:
// 1. Protection against job panics.  (They print to ERROR instead of take down the process)
// 2. (Optional) Limit on the number of jobs that may run simulatenously, to
//    limit resource consumption.
// 3. (Optional) Protection against multiple instances of a single job running
//    concurrently.  If one execution runs into the next, the next will be queued.
// 4. Cron expressions may be defined in app.conf and are reusable across jobs.
// 5. Job status reporting. [WIP]
package jobrunner

import (
	"time"

	"gopkg.in/robfig/cron.v2"
)

// Callers can use jobs.Func to wrap a raw func.
// (Copying the type to this package makes it more visible)
//
// For example:
//    jobrunner.Schedule("cron.frequent", jobs.Func(myFunc))
type Func func()

func (r Func) Run() { r() }

func Schedule(spec string, job cron.Job) (cronId cron.EntryID, err error) {
	sched, err := cron.Parse(spec)
	if err != nil {
		return cronId, err
	}

	cronId = MainCron.Schedule(sched, New(job))
	return cronId, nil
}

// Run the given job at a fixed interval.
// The interval provided is the time between the job ending and the job being run again.
// The time that the job takes to run is not included in the interval.
func Every(duration time.Duration, job cron.Job) {

	MainCron.Schedule(cron.Every(duration), New(job))
}

// Run the given job right now.
func Now(job cron.Job) {
	go New(job).Run()
}

// Run the given job once, after the given delay.
func In(duration time.Duration, job cron.Job) {
	go func() {
		time.Sleep(duration)
		New(job).Run()
	}()
}
