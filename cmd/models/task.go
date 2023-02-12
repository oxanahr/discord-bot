package models

import (
	"github.com/oxanahr/discord-bot/cmd/database"
	"log"
	"time"
)

type Task struct {
	ID             uint64     `json:"id" gorm:"primaryKey"`
	AssignedUserID *string    `json:"assignedUserID"`
	Priority       int        `json:"priority" gorm:"not null"`
	Name           string     `json:"name" gorm:"not null"`
	Description    string     `json:"description" gorm:"not null"`
	State          string     `json:"state" gorm:"not null"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"not null"`
	Deadline       *time.Time `json:"deadline"`
	Comments       []Comment  `json:"comments"`
}

func (t *Task) Create() error {
	t.State = "not_started"
	t.CreatedAt = time.Now()
	return database.DB.Create(t).Error
}

func StartTask(id uint64) error {
	return database.DB.Model(&Task{}).Where("id = ?", id).Update("state", "in_progress").Error
}

func CompleteTask(id uint64) error {
	return database.DB.Model(&Task{}).Where("id = ?", id).Update("state", "completed").Error
}

func AssignTask(id uint64, userID string) error {
	return database.DB.Model(&Task{}).Where("id = ?", id).Update("assigned_user_id", userID).Error
}

func GetTasks(assignedUserID *string, sort string, soon bool, unassigned bool) ([]Task, error) {
	q := database.DB.Model(&Task{}).Preload("Comments").Where("state != ?", "completed")
	if assignedUserID != nil {
		q.Where("assigned_user_id = ?", *assignedUserID)
	}
	if soon {
		now := time.Now()
		weekday := time.Duration(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		//TODO is this working ok?
		monday := now.Add(-1 * weekday * 24 * time.Hour)
		sunday := now.Add((6 - weekday) * 24 * time.Hour)
		log.Println(monday, sunday)
		q.Where("CAST(deadline AS DATE) BETWEEN CAST(? AS DATE) AND CAST(? AS DATE)", monday, sunday)
	}
	if unassigned {
		q.Where("assigned_user_id is null")
	}
	if sort == "deadline" {
		q.Order("deadline is null, deadline")
	} else if sort == "priority" {
		q.Order("priority desc")
	}
	var tasks []Task
	q.Find(&tasks)
	return tasks, nil
}

func GetTasksEndingTomorrow() ([]Task, error) {
	q := database.DB.Model(&Task{}).Preload("Comments").
		Where("state != ?", "completed").
		Where("assigned_user_id is not null").
		Where("CAST(deadline AS DATE) = CAST(? AS DATE)", time.Now().Add(24*time.Hour))

	var tasks []Task
	q.Find(&tasks)
	return tasks, nil
}

func GetInProgressTasks() ([]Task, error) {
	q := database.DB.Model(&Task{}).Preload("Comments").
		Where("state = ?", "in_progress").
		Where("assigned_user_id is not null")

	var tasks []Task
	q.Find(&tasks)
	return tasks, nil
}
