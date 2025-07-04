package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName    = "GOLANG_CLI_REMINDER"
	markValue   = "1"
	notifyIcon  = "assets/information.png"
	logFileName = "reminder.log"
	appTitle    = "Reminder"
)

const (
	exitUsageError   = 1
	exitParseError   = 2
	exitTimePast     = 3
	exitNotifyFailed = 4
	exitExecFailed   = 5
)

func main() {
	// Flag definitions for better CLI UX
	timeInput := flag.String("time", "", "Reminder time (e.g. 'in 5 minutes', 'tomorrow at 10am')")
	message := flag.String("msg", "", "Reminder message")
	flag.Parse()

	// Fallback for legacy args if flags not used
	if *timeInput == "" && *message == "" && len(os.Args) >= 3 {
		*timeInput = os.Args[1]
		*message = strings.Join(os.Args[2:], " ")
	}

	// Validation
	if *timeInput == "" || *message == "" {
		fmt.Printf("Usage: %s -time \"<time>\" -msg \"<message>\"\n", os.Args[0])
		fmt.Println("Example: go run main.go -time \"in 1 minute\" -msg \"Stretch break!\"")
		os.Exit(exitUsageError)
	}

	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(*timeInput, now)
	if err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
		os.Exit(exitParseError)
	}
	if t == nil {
		fmt.Println("Unable to parse time. Try formats like 'in 5 minutes', 'at 4pm', or 'tomorrow at 10am'.")
		os.Exit(exitParseError)
	}
	if now.After(t.Time) {
		fmt.Println("Set a future time!")
		os.Exit(exitTimePast)
	}

	diff := t.Time.Sub(now)
	fmt.Printf("Current time: %s\n", now.Format(time.Kitchen))
	fmt.Printf("Reminder scheduled at: %s (%s from now)\n", t.Time.Format(time.Kitchen), diff.Round(time.Second))

	// Schedule the reminder
	if os.Getenv(markName) == markValue {
		time.Sleep(diff)

		// Check if icon exists
		iconPath := ""
		if _, err := os.Stat(notifyIcon); err == nil {
			iconPath = notifyIcon
		}

		if err := beeep.Alert(appTitle, *message, iconPath); err != nil {
			fmt.Printf("Error displaying notification: %v\n", err)
			os.Exit(exitNotifyFailed)
		}

		// Optional beep
		_ = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)

		logReminder(t.Time, *message)
	} else {
		// Re-exec with env to detach reminder
		cmd := exec.Command(os.Args[0], "-time", *timeInput, "-msg", *message)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting reminder process:", err)
			os.Exit(exitExecFailed)
		}
		fmt.Println("Reminder is scheduled in background.")
		os.Exit(0)
	}
}

func logReminder(reminderTime time.Time, message string) {
	logMsg := fmt.Sprintf("%s | Reminder: %s\n", reminderTime.Format(time.RFC1123), message)
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		_, _ = f.WriteString(logMsg)
	}
}
