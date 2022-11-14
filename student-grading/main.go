package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
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

func (s student) string() string {
	return fmt.Sprintf("%v %v %v %v %v %v %v", s.firstName, s.lastName, s.university, s.test1Score, s.test2Score, s.test3Score, s.test4Score)
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func (ss studentStat) string() string {
	return fmt.Sprintf("%v %v %v",  ss.student.string(), ss.finalScore, ss.grade)
}

func main() {
	gradedStudents := calculateGrade(parseCSV("grades.csv"))

	overallTopper := findOverallTopper(gradedStudents)
	fmt.Println("Overall topper: ", overallTopper.string())

	topperPerUniversity := findTopperPerUniversity(gradedStudents)
	fmt.Println("------ Topper per university ------")
	for k, v := range topperPerUniversity {
		fmt.Println(k, " : ", v.string())
	}

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
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
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

func calculateGrade(students []student) []studentStat {
	studentStats := make([]studentStat, 0)

	for _, s := range students {
		ss := studentStat{student: s, finalScore: 0, grade: F}
		setGradeAndFinalScore(&ss)
		studentStats = append(studentStats, ss)
	}
	return studentStats
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	var max studentStat

	if len(gradedStudents) < 1 {
		return max
	}
	max = gradedStudents[0]

	for i := 0; i < len(gradedStudents); i++ {
		if max.finalScore < gradedStudents[i].finalScore {
			max = gradedStudents[i]
		}
	}
	return max
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	gsm := make(map[string][]studentStat, 0)

	//group student by university using map
	for _, ss := range gs {
		val, ok := gsm[ss.student.university]
		if !ok {
			gsm[ss.student.university] = []studentStat{{student: ss.student, finalScore: ss.finalScore, grade: ss.grade}}
			continue
		}
		val = append(val, studentStat{student: ss.student, finalScore: ss.finalScore, grade: ss.grade})
		gsm[ss.student.university] = val
	}

	tpm := make(map[string]studentStat)
	fmt.Println("len gsm: ", len(gsm))
	for k, v := range gsm {
		tpm[k] = findOverallTopper(v)
	}
	return tpm
}

func setGradeAndFinalScore(ss *studentStat) {
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
	num, err := strconv.ParseInt(input, 10, 32)
	//check for error
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return int(num)
}