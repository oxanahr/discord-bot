package models

import (
	"testing"
)

func TestComment(t *testing.T) {
	t.Parallel()

	comment := &Comment{
		ID:       19,
		TaskID:   35,
		AuthorID: "1976",
		Text:     "Task description",
	}

	if comment.ID != 19 {
		t.Errorf("comment.ID == %v", comment.ID)
	}
	if comment.TaskID != 35 {
		t.Errorf("comment.TaskID == %v", comment.TaskID)
	}

	if comment.AuthorID != "1976" {
		t.Errorf("comment.AuthorID == %v", comment.AuthorID)
	}
	if comment.Text != "Task description" {
		t.Errorf("comment.Text == %v", comment.Text)
	}
}
