package jobrunner

import (
	"time"

	"gopkg.in/robfig/cron.v2"
)

type StatusData struct {
	Id        cron.EntryID
	JobRunner *Job
	Next      time.Time
	Prev      time.Time
}

// Return detailed list of currently running recurring jobs
// to remove an entry, first retrieve the ID of entry
func Entries() []cron.Entry {
	return MainCron.Entries()
}

// Return detailed a specific job from running
func Entry(id cron.EntryID) cron.Entry {
	return MainCron.Entry(id)
}

// Return detailed a specific job from running
func EntryStatus(id cron.EntryID) StatusData {
	v := Entry(id)

	return StatusData{
		Id:        v.ID,
		JobRunner: AddJob(v.Job),
		Next:      v.Next,
		Prev:      v.Prev,
	}
}

func StatusPage() []StatusData {

	ents := MainCron.Entries()

	Statuses := make([]StatusData, len(ents))
	for k, v := range ents {
		Statuses[k].Id = v.ID
		Statuses[k].JobRunner = AddJob(v.Job)
		Statuses[k].Next = v.Next
		Statuses[k].Prev = v.Prev

	}

	// t := template.New("status_page")

	// var data bytes.Buffer
	// t, _ = t.ParseFiles("views/Status.html")

	// t.ExecuteTemplate(&data, "status_page", Statuses())
	return Statuses
}

func StatusJson() map[string]interface{} {

	return map[string]interface{}{
		"jobrunner": StatusPage(),
	}

}

func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
