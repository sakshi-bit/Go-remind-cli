# Go Remind CLI

A cross-platform command-line tool written in Go that lets you schedule desktop reminders using natural language time inputs.

---

Installation

1. Clone the repository:
   git clone https://github.com/sakshi-bit/go-remind-cli.git

2. Navigate to the project folder:
   cd go-remind-cli

3. Build the executable:
   go build -o reminder main.go

---

Usage

./reminder -time "in 10 minutes" -msg "Check the oven"

Supported time formats:
- in 5 minutes
- at 4pm
- tomorrow at 10am
- next Monday at 9am

---

Features

- Natural language time parsing
- Cross-platform desktop notifications
- Optional sound/beep alert
- Background scheduling
- Reminder logging in reminder.log

---

Author
Sakshi Srivastava  


