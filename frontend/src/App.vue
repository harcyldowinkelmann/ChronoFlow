<template>
  <div class="app-container">
    <header class="header">
      <div>
        <h1>Dynamic Schedule</h1>
        <p class="subtitle">Organize your routine efficiently and with style</p>
      </div>
      <div class="header-actions">
        <button class="btn-secondary" @click="openTagsModal">Manage Tags</button>
        <button class="btn-primary" @click="addRow">+ Add Row</button>
      </div>
    </header>

    <main class="grid-container">
      <!-- Days of Week Header -->
      <div class="grid-header">
        <div v-for="(day, index) in daysOfWeek" :key="index" class="header-cell">
          {{ day }}
        </div>
      </div>

      <!-- Grid Rows -->
      <div class="grid-body">
        <div v-for="r in rowCount" :key="r" class="grid-row">
          <div 
            v-for="d in 7" 
            :key="`${r}-${d}`" 
            class="grid-cell"
            @click="openTaskModal(r - 1, d - 1)"
          >
            <!-- Task Card -->
            <div v-if="getTask(r - 1, d - 1)" :class="['task-card', { 'task-completed': getTask(r - 1, d - 1).isCompleted }]">
              <div class="task-header">
                <div class="task-time">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                  <span class="time-text">
                    {{ getTask(r - 1, d - 1).startTime }}<template v-if="getTask(r - 1, d - 1).endTime"> - {{ getTask(r - 1, d - 1).endTime }}</template>
                  </span>
                </div>
                <input 
                  type="checkbox" 
                  class="task-checkbox" 
                  :checked="getTask(r - 1, d - 1).isCompleted"
                  @click.stop="toggleCompletion(getTask(r - 1, d - 1))"
                />
              </div>
              
              <p class="task-desc">{{ getTask(r - 1, d - 1).description }}</p>
              
              <div class="task-tags" v-if="getTask(r - 1, d - 1).tags && getTask(r - 1, d - 1).tags.length > 0">
                <span v-for="tag in getTask(r - 1, d - 1).tags" :key="tag.id" class="tag-pill" :style="{ backgroundColor: tag.color + '20', color: tag.color, borderColor: tag.color + '40' }">
                  {{ tag.name }}
                </span>
              </div>

              <button class="btn-delete-task" @click.stop="deleteTask(getTask(r - 1, d - 1).id)" title="Delete Task">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
              </button>
            </div>
            
            <!-- Empty Cell Placeholder -->
            <div v-else class="empty-cell">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Task Modal -->
    <div v-if="isTaskModalOpen" class="modal-overlay" @click="closeTaskModal">
      <div class="modal-content" @click.stop>
        <h2>{{ editingTask.id ? 'Edit Task' : 'New Task' }}</h2>
        <p class="modal-subtitle">{{ daysOfWeek[editingTask.dayOfWeek] }}</p>
        
        <div class="form-row">
          <div class="form-group half">
            <label>Start</label>
            <input type="time" v-model="editingTask.startTime" />
          </div>
          <div class="form-group half">
            <label>End (Optional)</label>
            <input type="time" v-model="editingTask.endTime" />
          </div>
        </div>

        <div class="form-group">
          <label>Description</label>
          <textarea v-model="editingTask.description" rows="3" placeholder="What needs to be done?"></textarea>
        </div>

        <div class="form-group">
          <label>Tags</label>
          <div class="tags-selector">
            <div 
              v-for="tag in allTags" 
              :key="tag.id" 
              :class="['tag-selectable', { 'selected': isTagSelected(tag) }]"
              :style="{ borderLeftColor: tag.color }"
              @click="toggleTagSelection(tag)"
            >
              {{ tag.name }}
            </div>
            <div v-if="allTags.length === 0" class="no-tags-msg">No tags created yet.</div>
          </div>
        </div>

        <div class="modal-actions">
          <button class="btn-cancel" @click="closeTaskModal">Cancel</button>
          <button class="btn-primary" @click="saveTask">Save Task</button>
        </div>
      </div>
    </div>

    <!-- Tags Modal -->
    <div v-if="isTagsModalOpen" class="modal-overlay" @click="closeTagsModal">
      <div class="modal-content" @click.stop>
        <h2>Manage Tags</h2>
        <p class="modal-subtitle">Create tags to categorize your tasks</p>

        <div class="tag-creation-area">
          <input type="text" v-model="newTagName" placeholder="Tag Name" class="input-tag-name" />
          <input type="color" v-model="newTagColor" class="input-tag-color" />
          <button class="btn-primary" @click="createNewTag" :disabled="!newTagName">Create</button>
        </div>

        <div class="tags-list">
          <div v-for="tag in allTags" :key="tag.id" class="tag-list-item">
            <div class="tag-preview" :style="{ backgroundColor: tag.color + '20', color: tag.color, borderColor: tag.color + '40' }">
              <span class="color-dot" :style="{ backgroundColor: tag.color }"></span>
              {{ tag.name }}
            </div>
            <button class="btn-delete-tag" @click="deleteTag(tag.id)">Delete</button>
          </div>
        </div>

        <div class="modal-actions">
          <button class="btn-primary" @click="closeTagsModal">Close</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import { GetTasks, SaveTask, DeleteTask, GetTags, SaveTag, DeleteTag } from '../wailsjs/go/main/App'

export default {
  data() {
    return {
      daysOfWeek: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],
      rowCount: 5,
      tasks: [],
      allTags: [],
      
      isTaskModalOpen: false,
      isTagsModalOpen: false,
      
      editingTask: {
        id: 0,
        scheduleId: 1,
        rowIndex: 0,
        dayOfWeek: 0,
        startTime: '',
        endTime: '',
        description: '',
        isCompleted: false,
        tags: []
      },

      newTagName: '',
      newTagColor: '#3b82f6'
    }
  },
  methods: {
    async loadData() {
      try {
        const tasks = await GetTasks();
        this.tasks = tasks || [];
        
        const tags = await GetTags();
        this.allTags = tags || [];
        
        let maxRow = 0;
        this.tasks.forEach(t => {
          if (t.rowIndex >= maxRow) maxRow = t.rowIndex;
        });
        if (maxRow >= this.rowCount) {
          this.rowCount = maxRow + 2;
        }
      } catch (err) {
        console.error("Error loading data:", err);
      }
    },
    getTask(rowIndex, dayOfWeek) {
      return this.tasks.find(t => t.rowIndex === rowIndex && t.dayOfWeek === dayOfWeek);
    },
    addRow() {
      this.rowCount++;
    },
    
    // --- Task Methods ---
    openTaskModal(rowIndex, dayOfWeek) {
      const existing = this.getTask(rowIndex, dayOfWeek);
      if (existing) {
        this.editingTask = JSON.parse(JSON.stringify(existing)); // Deep copy to avoid reactive mess before save
        if (!this.editingTask.tags) this.editingTask.tags = [];
      } else {
        this.editingTask = {
          id: 0,
          scheduleId: 1,
          rowIndex,
          dayOfWeek,
          startTime: '09:00',
          endTime: '',
          description: '',
          isCompleted: false,
          tags: []
        };
      }
      this.isTaskModalOpen = true;
    },
    closeTaskModal() {
      this.isTaskModalOpen = false;
    },
    async saveTask() {
      if (!this.editingTask.startTime || !this.editingTask.description) {
        alert("Please provide start time and description!");
        return;
      }
      try {
        await SaveTask(this.editingTask);
        this.closeTaskModal();
        await this.loadData();
      } catch (err) {
        console.error("Error saving task:", err);
      }
    },
    async deleteTask(id) {
      if(confirm("Are you sure you want to delete this task?")) {
        try {
          await DeleteTask(id);
          await this.loadData();
        } catch (err) {
          console.error("Error deleting task:", err);
        }
      }
    },
    async toggleCompletion(task) {
      task.isCompleted = !task.isCompleted;
      try {
        await SaveTask(task);
        await this.loadData();
      } catch (err) {
        console.error("Error toggling completion:", err);
      }
    },

    // --- Tag Methods ---
    openTagsModal() {
      this.isTagsModalOpen = true;
    },
    closeTagsModal() {
      this.isTagsModalOpen = false;
    },
    async createNewTag() {
      try {
        await SaveTag(this.newTagName, this.newTagColor);
        this.newTagName = '';
        await this.loadData();
      } catch (err) {
        console.error("Error creating tag:", err);
      }
    },
    async deleteTag(id) {
      if(confirm("Are you sure you want to delete this tag? It will be removed from all tasks.")) {
        try {
          await DeleteTag(id);
          await this.loadData();
        } catch (err) {
          console.error("Error deleting tag:", err);
        }
      }
    },
    isTagSelected(tag) {
      return this.editingTask.tags.some(t => t.id === tag.id);
    },
    toggleTagSelection(tag) {
      if (this.isTagSelected(tag)) {
        this.editingTask.tags = this.editingTask.tags.filter(t => t.id !== tag.id);
      } else {
        this.editingTask.tags.push(tag);
      }
    }
  },
  mounted() {
    this.loadData();
  }
}
</script>

<style>
/* CSS Resets & Variables */
:root {
  --primary: #4f46e5;
  --primary-hover: #4338ca;
  --secondary: #f1f5f9;
  --secondary-hover: #e2e8f0;
  
  --bg-color: #f8fafc;
  --surface: #ffffff;
  
  --text-main: #0f172a;
  --text-muted: #64748b;
  --border: #e2e8f0;
  --danger: #ef4444;
}

body {
  margin: 0;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  background-color: var(--bg-color);
  color: var(--text-main);
}

.app-container {
  max-width: 1500px;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  height: 100vh;
  box-sizing: border-box;
}

/* Header */
.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 2rem;
}

.header h1 {
  font-size: 2rem;
  font-weight: 800;
  margin: 0 0 0.2rem 0;
  color: var(--text-main);
  letter-spacing: -0.02em;
}

.subtitle {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.95rem;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

button {
  font-family: inherit;
  transition: all 0.2s ease;
}

.btn-primary {
  background-color: var(--primary);
  color: white;
  border: none;
  padding: 0.7rem 1.4rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  box-shadow: 0 4px 6px -1px rgba(79, 70, 229, 0.2);
}
.btn-primary:hover:not(:disabled) { 
  background-color: var(--primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 6px 8px -1px rgba(79, 70, 229, 0.3);
}
.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  box-shadow: none;
}

.btn-secondary {
  background-color: var(--secondary);
  color: var(--text-main);
  border: 1px solid var(--border);
  padding: 0.7rem 1.4rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
}
.btn-secondary:hover { background-color: var(--secondary-hover); }

/* Grid */
.grid-container {
  background: var(--surface);
  border-radius: 16px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.05), 0 4px 6px -4px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(226, 232, 240, 0.8);
  overflow: auto;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.grid-header {
  display: grid;
  grid-template-columns: repeat(7, minmax(180px, 1fr));
  background-color: rgba(248, 250, 252, 0.95);
  backdrop-filter: blur(8px);
  border-bottom: 2px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-cell {
  padding: 1.2rem 1rem;
  text-align: center;
  font-weight: 700;
  color: var(--text-main);
  text-transform: uppercase;
  font-size: 0.8rem;
  letter-spacing: 0.1em;
  border-right: 1px solid var(--border);
  overflow: hidden;
  text-overflow: ellipsis;
}
.header-cell:last-child { border-right: none; }

.grid-body {
  display: flex;
  flex-direction: column;
  background-color: #fafbfc;
}

.grid-row {
  display: grid;
  grid-template-columns: repeat(7, minmax(180px, 1fr));
  border-bottom: 1px solid var(--border);
  min-height: 140px;
}

.grid-cell {
  border-right: 1px solid var(--border);
  padding: 0.6rem;
  position: relative;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  min-width: 0; /* CRITICAL for grid children to prevent overflow */
  overflow: hidden; /* Ensures content doesn't break boundaries */
}
.grid-cell:last-child { border-right: none; }
.grid-cell:hover { background-color: rgba(241, 245, 249, 0.5); }
.grid-cell:hover .empty-cell { opacity: 1; transform: scale(1); }

.empty-cell {
  margin: auto;
  color: #cbd5e1;
  opacity: 0;
  transform: scale(0.8);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Task Card */
.task-card {
  background-color: var(--surface);
  border: 1px solid var(--border);
  padding: 0.8rem;
  border-radius: 10px;
  height: 100%;
  box-sizing: border-box;
  position: relative;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  width: 100%;
  max-width: 100%; /* Prevent overflow */
}
.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 12px -2px rgba(0,0,0,0.08);
  border-color: #cbd5e1;
}

.task-completed {
  opacity: 0.6;
  background-color: #f8fafc;
}
.task-completed .task-desc {
  text-decoration: line-through;
  color: var(--text-muted);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
}

.task-time {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--primary);
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: nowrap;
  min-width: 0;
  white-space: nowrap;
}
.time-text {
  white-space: nowrap;
}

.task-checkbox {
  cursor: pointer;
  width: 16px;
  height: 16px;
  min-width: 16px;
  accent-color: var(--primary);
  margin-top: 2px;
}

.task-desc {
  font-size: 0.85rem;
  margin: 0;
  color: var(--text-main);
  line-height: 1.5;
  word-break: normal;
  overflow-wrap: break-word;
  flex-grow: 1;
}

.task-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: auto;
}

.tag-pill {
  font-size: 0.65rem;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 12px;
  border: 1px solid transparent;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  max-width: 100%;
}

.btn-delete-task {
  position: absolute;
  top: -8px;
  right: -8px;
  background: var(--surface);
  border: 1px solid var(--border);
  color: var(--danger);
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0;
  transform: scale(0.8);
  transition: all 0.2s;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  z-index: 5;
}
.task-card:hover .btn-delete-task { 
  opacity: 1; 
  transform: scale(1);
}
.btn-delete-task:hover {
  background: var(--danger);
  color: white;
  border-color: var(--danger);
}

/* Modals */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 23, 42, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  backdrop-filter: blur(4px);
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: var(--surface);
  padding: 2.5rem;
  border-radius: 16px;
  width: 100%;
  max-width: 450px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-content h2 { margin: 0 0 0.5rem 0; font-size: 1.5rem; color: var(--text-main); }
.modal-subtitle { color: var(--text-muted); margin-bottom: 2rem; font-size: 0.95rem; }

.form-group {
  margin-bottom: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-group.half {
  flex: 1;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-main);
}

.form-group input[type="time"], .form-group textarea, .input-tag-name {
  padding: 0.8rem;
  border: 1px solid var(--border);
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.95rem;
  transition: all 0.2s;
  background-color: var(--secondary);
}
.form-group input:focus, .form-group textarea:focus, .input-tag-name:focus {
  outline: none;
  border-color: var(--primary);
  background-color: var(--surface);
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.tags-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 0.5rem 0;
}

.tag-selectable {
  padding: 0.4rem 0.8rem;
  border: 1px solid var(--border);
  border-left-width: 4px;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--surface);
}
.tag-selectable:hover { background: var(--secondary); }
.tag-selectable.selected {
  background: #eff6ff;
  border-color: var(--primary);
  border-left-color: var(--primary) !important;
  color: var(--primary-hover);
}

.no-tags-msg {
  font-size: 0.85rem;
  color: var(--text-muted);
  font-style: italic;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2.5rem;
}

.btn-cancel {
  background: transparent;
  color: var(--text-main);
  border: 1px solid var(--border);
  padding: 0.7rem 1.4rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: background 0.2s;
}
.btn-cancel:hover { background: var(--secondary); }

/* Tag Management Area */
.tag-creation-area {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 2rem;
}
.input-tag-name { flex-grow: 1; }
.input-tag-color {
  width: 46px;
  height: 46px;
  padding: 0;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  background: var(--surface);
}
.input-tag-color::-webkit-color-swatch-wrapper { padding: 4px; }
.input-tag-color::-webkit-color-swatch { border: none; border-radius: 4px; }

.tags-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 200px;
  overflow-y: auto;
}

.tag-list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--surface);
}

.tag-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 12px;
  border: 1px solid;
}
.color-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.btn-delete-tag {
  background: transparent;
  border: none;
  color: var(--danger);
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
}
.btn-delete-tag:hover { text-decoration: underline; }
</style>
