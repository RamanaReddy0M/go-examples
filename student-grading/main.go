package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
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

	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return make([]student, 0)
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	//ignoring headers
	reader.Read()
	var students []student
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Println(error)
		}

		students = append(students, student{
			firstName:  line[0],
			lastName:   line[1],
			university: line[2],
			test1Score: parseToInt(line[3]),
			test2Score: parseToInt(line[4]),
			test3Score: parseToInt(line[5]),
			test4Score: parseToInt(line[6]),
		})
	}
	return students
}

func calculateGrade(students []student) []studentStats {
	studentStatSlice := make([]studentStats, 0)

	for _, student := range students {
		studentStat := studentStats{student: student, finalScore: 0, grade: F}
		setGradeAndFinalScore(&studentStat)
		studentStatSlice = append(studentStatSlice, studentStat)
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

func findTopperPerUniversity(gradedStudents []studentStats) map[string]studentStats {
	gradedStudentsCopy := make([]studentStats, len(gradedStudents))
	copy(gradedStudentsCopy, gradedStudents)
	sort.SliceStable(gradedStudentsCopy, func(i, j int) bool {
		return gradedStudentsCopy[i].student.university < gradedStudentsCopy[j].student.university
	})
	var topperPerUniversityMap map[string]studentStats = map[string]studentStats{}
	j := 0
	i := 0
	for i = 0; i < len(gradedStudentsCopy); i++ {
		if gradedStudentsCopy[j].student.university != gradedStudentsCopy[i].student.university {
			topperPerUniversityMap[gradedStudentsCopy[i].university] = findOverallTopper(gradedStudentsCopy[j:i])
			j = i
		}
	}
	topperPerUniversityMap[gradedStudentsCopy[j].university] = findOverallTopper(gradedStudentsCopy[j:])

	return topperPerUniversityMap
}

func setGradeAndFinalScore(ss *studentStats) {
	fs := getFinalScore(ss.student)
	ss.finalScore = fs

	if fs < 35 {
		ss.grade = F
	} else if fs >= 35 && fs < 50 {
		ss.grade = C
	} else if fs >= 50 && fs < 70 {
		ss.grade = B
	} else {
		ss.grade = A
	}
}

func getFinalScore(s student) float32 {
	return float32(s.test1Score+s.test2Score+s.test3Score+s.test4Score) / 4.0
}

func parseToInt(input string) int {
	number, err := strconv.ParseInt(input, 10, 32)
	//check for error
	if err != nil {
		log.Println(err.Error())
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
