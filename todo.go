package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ساختار مربوط به هر تسک
type Task struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Note      string     `json:"note,omitempty"`
	Done      bool       `json:"done"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Due       *time.Time `json:"due,omitempty"`
	Priority  int        `json:"priority,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
}

// ساختار مدیریت‌کنندهٔ تسک‌ها
type TaskManager struct {
	tasks    []Task
	filePath string
}

// سازندهٔ TaskManager
func NewTaskManager(filepath string) *TaskManager {
	tm := &TaskManager{
		tasks:    make([]Task, 0),
		filePath: filepath,
	}
	tm.Load()
	return tm
}

// بارگذاری تسک‌ها از فایل JSON
func (tm *TaskManager) Load() error {
	if _, err := os.Stat(tm.filePath); os.IsNotExist(err) {
		tm.tasks = make([]Task, 0)
		return nil
	}
	data, err := os.ReadFile(tm.filePath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		tm.tasks = make([]Task, 0)
		return nil
	}
	return json.Unmarshal(data, &tm.tasks)
}

// ذخیره‌سازی تسک‌ها در فایل JSON
func (tm *TaskManager) Save() error {
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tm.filePath, data, 0644)
}

// افزودن تسک جدید
func (tm *TaskManager) AddTask(title, note string, priority int, tags []string) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	task := Task{
		ID:        uuid.New().String(),
		Title:     title,
		Note:      note,
		Priority:  priority,
		Tags:      tags,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	tm.tasks = append(tm.tasks, task)
	if err := tm.Save(); err != nil {
		return nil, err
	}
	return &task, nil
}

// نمایش تمام تسک‌ها
func (tm *TaskManager) ListTasks(showDone bool) {
	if len(tm.tasks) == 0 {
		fmt.Println("\nNo tasks found.")
		return
	}

	fmt.Println("\nYour Tasks:")
	fmt.Println(strings.Repeat("-", 50))

	for index, task := range tm.tasks {
		if !showDone && task.Done {
			continue
		}

		status := "[ ]"
		if task.Done {
			status = "[X]"
		}

		priorityStr := ""
		switch task.Priority {
		case 1:
			priorityStr = "Low"
		case 2:
			priorityStr = "Medium"
		case 3:
			priorityStr = "High"
		}

		fmt.Printf("%s [%d] %s\n", status, index+1, task.Title)
		fmt.Printf("   Priority: %s | ID: %s\n", priorityStr, task.ID[:8])

		if task.Note != "" {
			fmt.Printf("   Note: %s\n", task.Note)
		}
		if len(task.Tags) > 0 {
			fmt.Printf("   Tags: %s\n", strings.Join(task.Tags, ", "))
		}
		fmt.Println()
	}
}

// تغییر وضعیت تسک به انجام‌شده
func (tm *TaskManager) MarkDone(index int) error {
	if index < 0 || index >= len(tm.tasks) {
		return fmt.Errorf("invalid task number")
	}
	tm.tasks[index].Done = true
	tm.tasks[index].UpdatedAt = time.Now()
	return tm.Save()
}

// حذف تسک
func (tm *TaskManager) DeleteTask(index int) error {
	if index < 0 || index >= len(tm.tasks) {
		return fmt.Errorf("invalid task number")
	}
	taskTitle := tm.tasks[index].Title
	tm.tasks = append(tm.tasks[:index], tm.tasks[index+1:]...)
	err := tm.Save()
	if err != nil {
		return err
	}
	fmt.Printf("Task deleted: %s\n", taskTitle)
	return nil
}

// خواندن ورودی از کاربر
func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// نمایش منوی اصلی
func showMenu() {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("TODO LIST MANAGER")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println("1. Add New Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Mark Task as Done")
	fmt.Println("4. Delete Task")
	fmt.Println("5. Exit")
	fmt.Println(strings.Repeat("=", 40))
}

// افزودن تسک جدید از طریق ورودی تعاملی
func addTaskInteractive(tm *TaskManager) {
	fmt.Println("\nAdd New Task")
	fmt.Println(strings.Repeat("-", 30))

	title := readInput("Task title: ")
	if title == "" {
		fmt.Println("Error: Title cannot be empty!")
		return
	}

	note := readInput("Note (optional): ")

	priorityStr := readInput("Priority (1=Low, 2=Medium, 3=High) [2]: ")
	priority := 2
	if priorityStr != "" {
		if p, err := strconv.Atoi(priorityStr); err == nil && p >= 1 && p <= 3 {
			priority = p
		} else {
			fmt.Println("Error: Invalid priority! Using default (2)")
		}
	}

	tagsInput := readInput("Tags (comma separated, optional): ")
	tags := []string{}
	if tagsInput != "" {
		tags = strings.Split(tagsInput, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	task, err := tm.AddTask(title, note, priority, tags)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Task added successfully! (ID: %s)\n", task.ID[:8])
}

// لیست تسک‌ها به صورت تعاملی
func listTasksInteractive(tm *TaskManager) {
	fmt.Println("\nTask List Options:")
	fmt.Println("1. List incomplete tasks")
	fmt.Println("2. List all tasks")

	choice := readInput("Select option [1]: ")
	showAll := false

	if choice == "2" {
		showAll = true
	}

	tm.ListTasks(showAll)
}

// علامت‌گذاری تسک به عنوان انجام شده
func markDoneInteractive(tm *TaskManager) {
	if len(tm.tasks) == 0 {
		fmt.Println("Error: No tasks available!")
		return
	}

	tm.ListTasks(false) // فقط تسک‌های انجام نشده

	taskNumStr := readInput("\nEnter task number to mark as done: ")
	if taskNumStr == "" {
		return
	}

	taskNum, err := strconv.Atoi(taskNumStr)
	if err != nil || taskNum < 1 || taskNum > len(tm.tasks) {
		fmt.Println("Error: Invalid task number!")
		return
	}

	// پیدا کردن ایندکس واقعی تسک (با توجه به اینکه ممکن است فقط تسک‌های انجام نشده نمایش داده شده باشند)
	actualIndex := -1
	count := 0
	for i, task := range tm.tasks {
		if !task.Done {
			count++
			if count == taskNum {
				actualIndex = i
				break
			}
		}
	}

	if actualIndex == -1 {
		fmt.Println("Error: Task not found!")
		return
	}

	err = tm.MarkDone(actualIndex)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Task marked as done: %s\n", tm.tasks[actualIndex].Title)
}

// حذف تسک به صورت تعاملی
func deleteTaskInteractive(tm *TaskManager) {
	if len(tm.tasks) == 0 {
		fmt.Println("Error: No tasks available!")
		return
	}

	tm.ListTasks(true) // نمایش همه تسک‌ها

	taskNumStr := readInput("\nEnter task number to delete: ")
	if taskNumStr == "" {
		return
	}

	taskNum, err := strconv.Atoi(taskNumStr)
	if err != nil || taskNum < 1 || taskNum > len(tm.tasks) {
		fmt.Println("Error: Invalid task number!")
		return
	}

	// تأیید حذف
	taskTitle := tm.tasks[taskNum-1].Title
	confirm := readInput(fmt.Sprintf("Are you sure you want to delete '%s'? (y/N): ", taskTitle))

	if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
		err = tm.DeleteTask(taskNum - 1)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		fmt.Println("Deletion cancelled.")
	}
}

func main() {
	fmt.Println("Starting Todo List Manager...")
	tm := NewTaskManager("tasks.json")

	// نمایش آمار اولیه
	incompleteCount := 0
	for _, task := range tm.tasks {
		if !task.Done {
			incompleteCount++
		}
	}

	fmt.Printf("\nYou have %d tasks (%d incomplete)\n", len(tm.tasks), incompleteCount)

	// حلقه اصلی برنامه
	for {
		showMenu()
		choice := readInput("Select an option (1-5): ")

		switch choice {
		case "1":
			addTaskInteractive(tm)
		case "2":
			listTasksInteractive(tm)
		case "3":
			markDoneInteractive(tm)
		case "4":
			deleteTaskInteractive(tm)
		case "5":
			fmt.Println("\nThank you for using Todo List Manager! Goodbye!")
			return
		case "":
			fmt.Println("Error: Please select an option")
		default:
			fmt.Println("Error: Invalid option! Please choose 1-5")
		}

		// مکث کوتاه قبل از نمایش مجدد منو
		fmt.Print("\nPress Enter to continue...")
		bufio.NewScanner(os.Stdin).Scan()
	}
}
