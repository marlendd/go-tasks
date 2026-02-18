package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"slices"
	"regexp"
)

// начало решения

// Task описывает задачу, выполненную в определенный день
type Task struct {
	Date  time.Time
	Dur   time.Duration
	Title string
}

// ParsePage разбирает страницу журнала
// и возвращает задачи, выполненные за день
func ParsePage(src string) ([]Task, error) {
	lines := strings.Split(src, "\n")
	date, err := parseDate(lines[0])
	tasks, err := parseTasks(date, lines[1:])
	sortTasks(tasks)
	return tasks, err
}

// parseDate разбирает дату в формате дд.мм.гггг
func parseDate(src string) (time.Time, error) {
	return time.Parse("02.01.2006", src)
}

var re = regexp.MustCompile(`(\d+:\d+) - (\d+:\d+) (.+)`)
// parseTasks разбирает задачи из записей журнала
func parseTasks(date time.Time, lines []string) ([]Task, error) {
	taskIdx := make(map[string]int, len(lines))
	var tasks []Task

	for _, task := range lines {
		groups := re.FindStringSubmatch(task)
		if len(groups) < 4 {
			return []Task{}, errors.New("error while parsing regexp")
		}
		start, err := time.Parse("15:04", groups[1]) 
		if err != nil {
			return []Task{}, errors.New("error while parsing")
		}
		end, err := time.Parse("15:04", groups[2])
		if err != nil {
			return []Task{}, errors.New("error while parsing")
		}
		title := groups[3]
		dur := end.Sub(start)
        if dur <= 0 {
            return []Task{}, errors.New("negative dur")
        }
		if idx, exists := taskIdx[title]; !exists {
			taskIdx[title] = len(tasks)
			tasks = append(tasks, Task{
                Date:  date,
                Dur:   dur,
                Title: title,
            })
		} else {
			tasks[idx].Dur += dur
		}
	}
	return tasks, nil
}

// sortTasks упорядочивает задачи по убыванию длительности
func sortTasks(tasks []Task) {
	sortFunc := func(a, b Task) int {
		if a.Dur > b.Dur {
			return -1
		}
		if a.Dur < b.Dur {
			return 1
		}
		return 0
	}
	slices.SortFunc(tasks, sortFunc)
}

// конец решения

func main() {
	page := `15.04.2022
8:00 - 8:30 Завтрак
8:30 - 9:30 Оглаживание кота
9:30 - 10:00 Интернеты
10:00 - 14:00 Напряженная работа
14:00 - 14:45 Обед
14:45 - 15:00 Оглаживание кота
15:00 - 19:00 Напряженная работа
19:00 - 19:30 Интернеты
19:30 - 22:30 Безудержное веселье
22:30 - 23:00 Оглаживание кота`

	entries, err := ParsePage(page)
	if err != nil {
		panic(err)
	}
	fmt.Println("Мои достижения за", entries[0].Date.Format("2006-01-02"))
	for _, entry := range entries {
		fmt.Printf("- %v: %v\n", entry.Title, entry.Dur)
	}

	// ожидаемый результат
	/*
		Мои достижения за 2022-04-15
		- Напряженная работа: 8h0m0s
		- Безудержное веселье: 3h0m0s
		- Оглаживание кота: 1h45m0s
		- Интернеты: 1h0m0s
		- Обед: 45m0s
		- Завтрак: 30m0s
	*/
}
