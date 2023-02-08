package models

import (
	"fmt"
	"github/oxanahr/discordBot/cmd/database"
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

func GetTasks(assignedUserID *string, sort string, soon bool) ([]Task, error) {
	q := database.DB.Model(&Task{}).Where("state != ?", "completed")
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
		fmt.Println(monday, sunday)
		q.Where("CAST(deadline AS DATE) BETWEEN CAST(? AS DATE) AND CAST(? AS DATE)", monday, sunday)
	}
	if sort == "deadline" {
		q.Order("deadline desc")
	} else if sort == "priority" {
		q.Order("priority desc")
	}
	rows, err := q.Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := []Task{}
	for rows.Next() {
		var t Task
		database.DB.ScanRows(rows, &t)
		result = append(result, t)
	}
	return result, nil
}

func GetTasksEndingTomorrow() ([]Task, error) {
	q := database.DB.Model(&Task{}).
		Where("state != ?", "completed").
		Where("assigned_user_id is not null").
		Where("CAST(deadline AS DATE) = CAST(? AS DATE)", time.Now().Add(24*time.Hour))

	rows, err := q.Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := []Task{}
	for rows.Next() {
		var t Task
		database.DB.ScanRows(rows, &t)
		result = append(result, t)
	}
	return result, nil
}

func GetInProgressTasks() ([]Task, error) {
	q := database.DB.Model(&Task{}).
		Where("state = ?", "in_progress").
		Where("assigned_user_id is not null")

	rows, err := q.Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := []Task{}
	for rows.Next() {
		var t Task
		database.DB.ScanRows(rows, &t)
		result = append(result, t)
	}
	return result, nil
}
