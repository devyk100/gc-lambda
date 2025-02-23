package types

import "time"

type Course struct {
	Description   string `json:"description"`
	Id            int    `json:"id"`
	Language      string `json:"language"`
	Name          string `json:"name"`
	Instructor    string `json:"instructor"`
	ImageUrl      string `json:"image_url"`
	InstructorUrl string `json:"instructor_img_url"`
	IsPublic      bool   `json:"is_public"`
}

type Lesson struct {
	Id        int       `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	IsPublic  bool      `json:"is_public"`
}
