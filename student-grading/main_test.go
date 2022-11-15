package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseCSV(t *testing.T) {
	assert := assert.New(t)
	students := parseCSV("grades.csv")

	assert.Equal(30, len(students), "Size should be equal")

	fs := student{"Kaylen", "Johnson", "Duke University", 52, 47, 35, 38}
	assert.Equal(fs, students[0], "First student should be equal")

	ls := student{"Solomon", "Hunter", "Boston University", 45, 62, 32, 58}
	assert.Equal(ls, students[29], "Last student should be equal")
}

func TestCalculateGrade(t *testing.T) {
	assert := assert.New(t)

	students := parseCSV("grades.csv")
	gradedStudents := calculateGrade(students)

	expfs := []float32{43, 59.25, 53, 58.25, 52.25, 50.75, 54.75, 49.25, 64.75, 43.25, 68.5, 57.75, 68.25, 66.75, 45.5, 45.75, 45.5, 58, 56, 60.25, 61, 62.5, 80.5, 53, 30.75, 57.5, 70.75, 48.5, 60.25, 49.25}
	expgrds := []Grade{C, B, B, B, B, B, B, C, B, C, B, B, B, B, C, C, C, B, B, B, B, B, A, B, F, B, A, C, B, C}

	for i, ss := range gradedStudents {
		assert.Equal(expfs[i], ss.finalScore, "Final Score should be equal")
		assert.Equal(expgrds[i], ss.grade, "Grade should be equal")
	}
}

func TestFindOverallTopper(t *testing.T) {
	grds := calculateGrade(parseCSV("grades.csv"))

	got := findOverallTopper(grds).student
	want := student{"Bernard", "Wilson", "Boston University", 90, 85, 76, 71}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestFindTopperPerUniversity(t *testing.T) {
	assert := assert.New(t)

	grds := calculateGrade(parseCSV("grades.csv"))
	tpu := findTopperPerUniversity(grds)

	assert.Equal(student{"Bernard", "Wilson", "Boston University", 90, 85, 76, 71}, tpu["Boston University"].student, "Boston University topper should be 'Bernard Wilson'")
	assert.Equal(student{"Tamara", "Webb", "Duke University", 73, 62, 90, 58}, tpu["Duke University"].student, "Duke University topper should be 'Tamara Webb'")
	assert.Equal(student{"Izayah", "Hunt", "Union College", 29, 78, 41, 85}, tpu["Union College"].student, "Union College topper should be 'Izayah Hunt'")
	assert.Equal(student{"Karina", "Shaw", "University of California", 69, 78, 56, 70}, tpu["University of California"].student, "University of California topper should be 'Karina Shaw'")
	assert.Equal(student{"Nathan", "Gordon", "University of Florida", 53, 79, 84, 51}, tpu["University of Florida"].student, "University of Florida topper should be 'Nathan Gordon'")
}