document.addEventListener('DOMContentLoaded', function() {
    // Initialize Bootstrap modals
    const editModal = new bootstrap.Modal(document.getElementById('editTaskModal'));
    const deleteModal = new bootstrap.Modal(document.getElementById('deleteConfirmModal'));
    
    // Add new task functionality
    document.getElementById('addNewTask').addEventListener('click', function(e) {
        e.preventDefault();
        alert('Add new task functionality will be implemented here.');
        // In a real app, you would show a form and make an API call to create a new task
    });
    
    // Edit task functionality
    document.querySelectorAll('.edit-task').forEach(button => {
        button.addEventListener('click', function() {
            const taskId = this.getAttribute('data-task-id');
            const taskCard = document.querySelector(`.task-card[data-task-id="${taskId}"]`);
            const taskName = taskCard.querySelector('.card-title').textContent;
            const taskStatus = taskCard.getAttribute('data-task-status').toLowerCase();
            
            document.getElementById('editTaskId').value = taskId;
            document.getElementById('editTaskName').value = taskName;
            document.getElementById('editTaskStatus').value = taskStatus === 'completed' ? 'done' : 'pending';
            
            editModal.show();
        });
    });
    
    // Save task changes
    document.getElementById('saveTaskChanges').addEventListener('click', function() {
        const taskId = document.getElementById('editTaskId').value;
        const taskName = document.getElementById('editTaskName').value;
        const taskStatus = document.getElementById('editTaskStatus').value;
        
        // In a real app, you would make an API call here to update the task
        alert(`Task ${taskId} would be updated to "${taskName}" with status "${taskStatus}"`);
        
        // For demo purposes, update the UI directly
        const taskCard = document.querySelector(`.task-card[data-task-id="${taskId}"]`);
        taskCard.querySelector('.card-title').textContent = taskName;
        
        const statusBadge = taskCard.querySelector('.badge');
        if (taskStatus === 'done') {
            statusBadge.textContent = 'Completed';
            statusBadge.className = 'badge rounded-pill badge-success';
            taskCard.setAttribute('data-task-status', 'Completed');
        } else {
            statusBadge.textContent = 'Pending';
            statusBadge.className = 'badge rounded-pill badge-warning';
            taskCard.setAttribute('data-task-status', 'Pending');
        }
        
        editModal.hide();
    });
    
    // Delete task functionality
    document.querySelectorAll('.delete-task').forEach(button => {
        button.addEventListener('click', function() {
            const taskId = this.getAttribute('data-task-id');
            document.getElementById('deleteTaskId').value = taskId;
            deleteModal.show();
        });
    });
    
    // Confirm delete
    document.getElementById('confirmDelete').addEventListener('click', function() {
        const taskId = document.getElementById('deleteTaskId').value;
        
        // In a real app, you would make an API call here to delete the task
        alert(`Task ${taskId} would be deleted`);
        
        // For demo purposes, just hide the task card
        const taskCard = document.querySelector(`.task-card[data-task-id="${taskId}"]`);
        taskCard.remove();
        
        // Update task count
        const taskCount = document.querySelectorAll('.task-card').length;
        document.getElementById('task-count').textContent = taskCount;
        
        deleteModal.hide();
    });
    
    // Toggle status functionality
    document.querySelectorAll('.toggle-status').forEach(button => {
        button.addEventListener('click', function() {
            const taskId = this.getAttribute('data-task-id');
            const currentStatus = this.getAttribute('data-task-status');
            const newStatus = currentStatus === 'Completed' ? 'Pending' : 'Completed';
            
            // In a real app, you would make an API call here to update the status
            alert(`Task ${taskId} status would be changed from "${currentStatus}" to "${newStatus}"`);
            
            // For demo purposes, update the UI directly
            const taskCard = document.querySelector(`.task-card[data-task-id="${taskId}"]`);
            const statusBadge = taskCard.querySelector('.badge');
            
            if (newStatus === 'Completed') {
                statusBadge.textContent = 'Completed';
                statusBadge.className = 'badge rounded-pill badge-success';
            } else {
                statusBadge.textContent = 'Pending';
                statusBadge.className = 'badge rounded-pill badge-warning';
            }
            
            this.setAttribute('data-task-status', newStatus);
            taskCard.setAttribute('data-task-status', newStatus);
        });
    });
    
    // Filter functionality
    const searchInput = document.getElementById('searchInput');
    const statusFilter = document.getElementById('statusFilter');
    const clearFiltersBtn = document.getElementById('clearFilters');
    
    function applyFilters() {
        const searchTerm = searchInput.value.toLowerCase();
        const statusValue = statusFilter.value;
        let visibleCount = 0;
        
        document.querySelectorAll('.task-card').forEach(card => {
            const taskName = card.querySelector('.card-title').textContent.toLowerCase();
            const taskStatus = card.getAttribute('data-task-status').toLowerCase();
            
            // Check if task matches both filters
            const matchesSearch = taskName.includes(searchTerm);
            const matchesStatus = statusValue === 'all' || 
                                (statusValue === 'completed' && taskStatus === 'completed') ||
                                (statusValue === 'pending' && taskStatus === 'pending');
            
            if (matchesSearch && matchesStatus) {
                card.classList.remove('hidden');
                visibleCount++;
            } else {
                card.classList.add('hidden');
            }
        });
        
        // Update visible task count
        document.getElementById('task-count').textContent = visibleCount;
    }
    
    searchInput.addEventListener('input', applyFilters);
    statusFilter.addEventListener('change', applyFilters);
    
    clearFiltersBtn.addEventListener('click', function() {
        searchInput.value = '';
        statusFilter.value = 'all';
        
        document.querySelectorAll('.task-card').forEach(card => {
            card.classList.remove('hidden');
        });
        
        // Reset task count
        document.getElementById('task-count').textContent = document.querySelectorAll('.task-card').length;
    });
}); 