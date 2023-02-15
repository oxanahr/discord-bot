package models

import (
	"testing"
)

func TestTask(t *testing.T) {
	t.Parallel()

	taskModel := &Task{
		ID:          1999,
		Priority:    10,
		Name:        "Task name",
		Description: "Task description",
		State:       "completed",
	}

	if taskModel.ID != 1999 {
		t.Errorf("task.ID == %v", taskModel.ID)
	}
	if taskModel.Priority != 10 {
		t.Errorf("task.Priority == %v", taskModel.Priority)
	}

	if taskModel.Name != "Task name" {
		t.Errorf("task.Name== %v", taskModel.Name)
	}
	if taskModel.Description != "Task description" {
		t.Errorf("task.Description == %v", taskModel.Description)
	}
	if taskModel.State != "completed" {
		t.Errorf("task.State == %v", taskModel.State)
	}
}
