package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStats struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {

	students := make([]student, 0)

	lines, err := ReadLines(filePath)
	if err != nil {
		log.SetPrefix("Error: ")
		log.Print(err.Error())
		return students
	}

	//empty check
	if len(lines) <= 1 {
		return students
	}
	//converting each line in file to student type
	for _, line := range lines[1:] {
		sData := strings.Split(line, ",")
		students = append(students,
			student{sData[0], sData[1], sData[2],
				parseToInt(sData[3]), parseToInt(sData[4]), parseToInt(sData[5]), parseToInt(sData[6])})
	}
	return students
}

func calculateGrade(students []student) []studentStats {
	studentStatSlice := make([]studentStats, 0)

	for _, student := range students {
		grade, finalScore := getGradeAndFinalScore(student)
		studentStatSlice = append(studentStatSlice, studentStats{student: student, finalScore: finalScore, grade: grade})
	}
	return studentStatSlice
}

func findOverallTopper(gradedStudents []studentStats) studentStats {

	if len(gradedStudents) < 1 {
		return studentStats{}
	}

	max := gradedStudents[0]

	for i := 0; i < len(gradedStudents); i++ {
		if max.finalScore < gradedStudents[i].finalScore {
			max = gradedStudents[i]
		}
	}
	return max
}

func findTopperPerUniversity(gradedStudents []studentStats) []studentStats {
	gradedStudentsCopy := make([]studentStats, len(gradedStudents))
	copy(gradedStudentsCopy, gradedStudents)
	sort.SliceStable(gradedStudentsCopy, func(i, j int) bool {
		return gradedStudentsCopy[i].student.university < gradedStudentsCopy[j].student.university
	})
	var topperPerUniversitySlice []studentStats
	j := 0
	i := 0
	for i = 0; i < len(gradedStudentsCopy); i++ {
		if gradedStudentsCopy[j].student.university != gradedStudentsCopy[i].student.university {
			topperPerUniversitySlice = append(topperPerUniversitySlice, findOverallTopper(gradedStudentsCopy[j:i]))
			j = i
		}
	}
	topperPerUniversitySlice = append(topperPerUniversitySlice, findOverallTopper(gradedStudentsCopy[j:]))

	return topperPerUniversitySlice
}

func getGradeAndFinalScore(s student) (grade Grade, finalScore float32) {
	fs := getFinalScore(s)
	var grd Grade

	if fs < 35 {
		grd = F
	} else if fs >= 35 && fs < 50 {
		grd = C
	} else if fs >= 50 && fs < 70 {
		grd = B
	} else {
		grd = A
	}

	return grd, fs
}

func getFinalScore(s student) float32 {
	return float32(s.test1Score+s.test2Score+s.test3Score+s.test4Score) / 4.0
}

func parseToInt(input string) int {
	number, err := strconv.ParseInt(input, 10, 32)
	//check for error
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return int(number)
}

func main() {
	gradedStudentsSlice := calculateGrade(parseCSV("grades.csv"))

	overallTopper := findOverallTopper(gradedStudentsSlice)

	fmt.Println("Overall topper: ", overallTopper)

	topperPerUniversitySlice := findTopperPerUniversity(gradedStudentsSlice)
	fmt.Println("------Topper per university------")
	for i, s := range topperPerUniversitySlice {
		fmt.Println(i, ":", s)
	}
}
