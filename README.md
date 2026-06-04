# ChronoFlow

A dynamic time-blocking planner built with Wails and Vue.js.

ChronoFlow helps users create, manage, and optimize their daily routines using a visual time-grid interface. Instead of relying on traditional calendars, the application focuses on structured time blocks, making it easier to build productive schedules, study plans, work routines, fitness programs, and personal projects.

---

## Features

- 📅 Visual weekly planner
- ⏰ Customizable time slots
- 🎨 Color-coded activities
- 🔄 Dynamic schedule management
- 📝 Task and routine organization
- 💻 Cross-platform desktop application
- ⚡ Fast and lightweight interface

---

## Tech Stack

### Backend

- Wails
- Go

### Frontend

- Vue.js
- JavaScript
- HTML5
- CSS3

---

## Project Goals

The main goal of ChronoFlow is to provide a simple and efficient way to organize routines through visual planning.

Potential future features include:

- Drag-and-drop scheduling
- Multiple schedule templates
- Habit tracking
- Productivity statistics
- Recurring activities
- Cloud synchronization
- Dark mode
- Notifications and reminders
- Data export/import

---

## Installation

### Prerequisites

- Go 1.24+
- Node.js 20+
- Wails CLI

Install Wails:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verify installation:

```bash
wails doctor
```

### Clone the Repository

```bash
git clone https://github.com/harcyldowinkelmann/chronoflow.git
cd chronoflow
```

### Install Frontend Dependencies

```bash
cd frontend
npm install
```

### Run in Development Mode

```bash
wails dev
```

### Build for Production

```bash
wails build
```
