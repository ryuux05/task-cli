# Task CLI HTML View Documentation

The Task CLI application provides beautiful and functional HTML views for your tasks. This document explains how to use these views and their features.

## Overview

Task CLI offers two HTML view options:
1. **Single Task View** - View details of a single task in HTML format
2. **All Tasks View** - View and manage all tasks in an interactive HTML interface

Both views are automatically opened in your default web browser when you use the relevant commands.

## Single Task View

### Usage

To view a single task in HTML format:

```
task view <task_id>
```

For example:
```
task view 1
```

### Features

The single task view shows:
- Task name
- Task status (Pending or Completed)
- Creation date
- Options to edit or delete the task

## All Tasks View

### Usage

To view all tasks in HTML format:

```
task view-all
```

You can also specify the format:
```
task view-all --format html
```

### Features

The all tasks view provides a comprehensive task management interface:

#### Task Listing

- All tasks are displayed as cards
- Each card shows:
  - Task name
  - Status badge (Pending or Completed)
  - Creation date
  - Action buttons (Edit, Toggle Status, Delete)

#### Filtering & Searching

- **Search Box**: Filter tasks by name
- **Status Filter**: Filter by status (All, Pending, Completed)
- **Clear Filters** button to reset all filters

#### Task Management

- **Add Task**: Button to add a new task (currently displays an alert for demonstration)
- **Edit Task**: Opens a modal with form fields for task name and status
- **Toggle Status**: Quickly change a task's status between Pending and Completed
- **Delete Task**: Removes a task with confirmation dialog

#### UI Components

1. **Filter Section**:
   - Located at the top of the page
   - Contains search input, status dropdown, and action buttons

2. **Task Cards**:
   - Each task is displayed as a card with consistent styling
   - Status is color-coded (yellow for Pending, green for Completed)

3. **Edit Modal**:
   - Form to edit task details
   - Fields for task name and status
   - Cancel and Save buttons

4. **Delete Confirmation Modal**:
   - Confirmation dialog before deletion
   - Warning about the action being irreversible

5. **Footer**:
   - Displays the total count of tasks (or filtered tasks when filters are applied)

## Technical Details

### File Structure

```
public/
├── assets/
│   ├── css/
│   │   └── tasks_list.css   # Styling for the tasks view
│   └── js/
│       └── tasks_list.js    # Interactive functionality for the tasks view
└── templates/
    ├── task_view.html       # Template for single task view
    └── tasks_list.html      # Template for all tasks view
```

### CSS Styling

The CSS for the HTML views provides:
- Clean, modern styling with Bootstrap integration
- Responsive design that works on mobile and desktop
- Custom styling for task cards, buttons, and status badges
- Proper spacing and layout for optimal readability

### JavaScript Functionality

The JavaScript enables:
- Modal dialogs for editing and confirming deletions
- Task filtering by name and status
- Status toggling with visual feedback
- Task editing with form validation
- Deletion with confirmation
- Dynamic task count updates

### Template Data

The templates receive data in the following format:

For single task view:
```go
TaskViewModel{
    Id:          task.Id,
    Name:        task.Name,
    StatusText:  "Pending" or "Completed",
    StatusClass: "badge-warning" or "badge-success",
    CreatedAt:   task.CreatedAt,
}
```

For all tasks view:
```go
TasksListViewModel{
    Tasks:      []TaskViewModel,  // Array of task view models
    TotalTasks: count,            // Total number of tasks
}
```

## Browser Compatibility

The HTML views are designed to work in all modern browsers including:
- Chrome
- Firefox
- Safari
- Edge

## Customization

The HTML views use Bootstrap CSS and custom styling. You can customize the appearance by modifying:
- `public/assets/css/tasks_list.css` for styling
- `public/templates/*.html` for HTML structure

## Limitations

- The browser views are read-only and don't currently sync changes back to the CLI
- The "Add Task" button in the HTML view shows a placeholder message
- The HTML files are temporary (for single task view) or stored in your home directory (for all tasks view)

## Troubleshooting

If the browser doesn't open automatically:
1. Check the console output for the HTML file path
2. Open the file manually in your browser
3. Look for any error messages in the console output 