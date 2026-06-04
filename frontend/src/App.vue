<template>
  <div class="app-container" :data-theme="currentTheme">
    <header class="header">
      <div class="header-titles">
        <h1>Dynamic Schedule</h1>
        <p class="subtitle">Organize your routine efficiently and with style</p>
      </div>
      <div class="header-actions">
        <select v-model="currentTheme" @change="changeTheme" class="theme-select">
          <option v-for="t in themes" :key="t.value" :value="t.value">{{ t.label }}</option>
        </select>
        <button class="btn-secondary" @click="openTagsModal">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
          Manage Tags
        </button>
        <button class="btn-primary" @click="addRow">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
          Add Row
        </button>
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
                <div class="task-actions">
                  <input 
                    type="checkbox" 
                    class="task-checkbox" 
                    :checked="getTask(r - 1, d - 1).isCompleted"
                    @click.stop="toggleCompletion(getTask(r - 1, d - 1))"
                  />
                  <button class="btn-icon btn-delete-task" @click.stop="deleteTask(getTask(r - 1, d - 1).id)" title="Delete Task">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                  </button>
                </div>
              </div>
              
              <p class="task-desc">{{ getTask(r - 1, d - 1).description }}</p>
              
              <div class="task-tags" v-if="getTask(r - 1, d - 1).tags && getTask(r - 1, d - 1).tags.length > 0">
                <span v-for="tag in getTask(r - 1, d - 1).tags" :key="tag.id" class="tag-pill" :style="{ backgroundColor: tag.color, color: getContrastYIQ(tag.color) }">
                  {{ tag.name }}
                </span>
              </div>
            </div>
            
            <!-- Empty Cell Hover Effect -->
            <div v-else class="empty-cell-overlay">
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
              :style="{ backgroundColor: isTagSelected(tag) ? tag.color : 'transparent', color: isTagSelected(tag) ? getContrastYIQ(tag.color) : 'inherit', borderColor: tag.color }"
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
            <div class="tag-preview" :style="{ backgroundColor: tag.color, color: getContrastYIQ(tag.color) }">
              {{ tag.name }}
            </div>
            <button class="btn-icon btn-delete-tag" @click="deleteTag(tag.id)" title="Delete Tag">
               <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
            </button>
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
import { GetTasks, SaveTask, DeleteTask, GetTags, SaveTag, DeleteTag, GetTheme, SaveTheme } from '../wailsjs/go/main/App'

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
      newTagColor: '#0052cc',
      
      currentTheme: 'ambitious',
      themes: [
        { value: 'ambitious', label: 'Ambitious' },
        { value: 'sharp-mind', label: 'Sharp Mind' },
        { value: 'courage', label: 'Courage' },
        { value: 'hard-work', label: 'Hard Work' }
      ]
    }
  },
  methods: {
    async loadData() {
      try {
        const theme = await GetTheme();
        this.currentTheme = theme || 'ambitious';
        
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
        this.editingTask = JSON.parse(JSON.stringify(existing));
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
    },
    // Utility to get white/black text depending on background color brightness
    getContrastYIQ(hexcolor){
        hexcolor = hexcolor.replace("#", "");
        if (hexcolor.length === 3) {
            hexcolor = hexcolor.split('').map(c => c+c).join('');
        }
        var r = parseInt(hexcolor.substr(0,2),16);
        var g = parseInt(hexcolor.substr(2,2),16);
        var b = parseInt(hexcolor.substr(4,2),16);
        var yiq = ((r*299)+(g*587)+(b*114))/1000;
        return (yiq >= 128) ? '#172b4d' : '#ffffff';
    },
    async changeTheme() {
      try {
        await SaveTheme(this.currentTheme);
      } catch (err) {
        console.error("Error saving theme:", err);
      }
    }
  },
  mounted() {
    this.loadData();
  }
}
</script>

<style>
/* Themes & Palettes */
[data-theme="ambitious"] {
  --primary: #1a472a;        
  --primary-hover: #2a623d;  
  --secondary: #dcdcdc;      
  --secondary-hover: #c0c0c0;
  --bg-gradient: linear-gradient(135deg, #0b1c11 0%, #1a472a 100%);
  --surface-table: #f2f5f4;  
  --surface-header: #ffffff; 
  --surface-card: #ffffff;
  --text-main: #0d1a12;      
  --text-muted: #5e6c64;     
  --border: #c4cdc8;         
  --danger: #d32f2f;
  --danger-hover: #b71c1c;
}

[data-theme="sharp-mind"] {
  --primary: #0e1a40;
  --primary-hover: #222f5b;
  --secondary: #946b2d;
  --secondary-hover: #7a5825;
  --bg-gradient: linear-gradient(135deg, #050a1f 0%, #0e1a40 100%);
  --surface-table: #f2f4f7;
  --surface-header: #ffffff;
  --surface-card: #ffffff;
  --text-main: #070d20;
  --text-muted: #5e6c84;
  --border: #b0bac7;
  --danger: #d32f2f;
  --danger-hover: #b71c1c;
}

[data-theme="courage"] {
  --primary: #740001;
  --primary-hover: #ae0001;
  --secondary: #d3a625;
  --secondary-hover: #eeba30;
  --bg-gradient: linear-gradient(135deg, #3a0000 0%, #740001 100%);
  --surface-table: #faf4f4;
  --surface-header: #ffffff;
  --surface-card: #ffffff;
  --text-main: #3a0000;
  --text-muted: #740001;
  --border: #e6c8a8;
  --danger: #d32f2f;
  --danger-hover: #b71c1c;
}

[data-theme="hard-work"] {
  --primary: #eeb939;
  --primary-hover: #f0c75e;
  --secondary: #111111;
  --secondary-hover: #222222;
  --bg-gradient: linear-gradient(135deg, #000000 0%, #1a1a1a 100%);
  --surface-table: #fffef5;
  --surface-header: #ffffff;
  --surface-card: #ffffff;
  --text-main: #111111;
  --text-muted: #555555;
  --border: #eeb939;
  --danger: #d32f2f;
  --danger-hover: #b71c1c;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Noto Sans', 'Ubuntu', 'Droid Sans', 'Helvetica Neue', sans-serif;
}

.app-container {
  background: var(--bg-gradient);
  background-attachment: fixed;
  color: var(--text-main);
  width: 100%;
  min-height: 100vh;
  padding: 3rem 5rem;
  display: flex;
  flex-direction: column;
  height: 100vh;
  box-sizing: border-box;
}

/* Header */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  background: rgba(255, 255, 255, 0.85);
  padding: 1rem 1.5rem;
  border-radius: 12px;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.header-titles h1 {
  font-size: 1.8rem;
  font-weight: 800;
  margin: 0 0 0.2rem 0;
  color: var(--text-main);
  letter-spacing: -0.01em;
}

.subtitle {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.95rem;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

button {
  font-family: inherit;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn-primary {
  background-color: var(--primary);
  color: white;
  border: none;
  padding: 0.6rem 1.2rem;
  border-radius: 3px; 
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
}
.btn-primary:hover:not(:disabled) { 
  background-color: var(--primary-hover);
}
.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  box-shadow: none;
}

.btn-secondary {
  background-color: var(--secondary);
  color: var(--text-main);
  border: none;
  padding: 0.6rem 1.2rem;
  border-radius: 3px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
}
.btn-secondary:hover { background-color: var(--secondary-hover); }

.btn-icon {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 4px;
  border-radius: 3px;
  color: var(--text-muted);
}
.btn-icon:hover {
  background: var(--secondary);
  color: var(--text-main);
}

.theme-select {
  padding: 0.6rem 1.2rem;
  border-radius: 3px;
  border: 1px solid var(--border);
  background: var(--surface-header);
  color: var(--text-main);
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
}

/* TABLE SYSTEM - Hard lines restored, styled elegantly */
.grid-container {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  overflow: auto;
  background-color: var(--surface-table);
  border-radius: 12px;
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255,255,255,0.2);
}

.grid-header {
  display: grid;
  grid-template-columns: repeat(7, minmax(220px, 1fr));
  background-color: var(--surface-header);
  border-bottom: 2px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-cell {
  padding: 1rem;
  text-align: center;
  font-weight: 700;
  color: var(--text-main);
  font-size: 0.85rem;
  text-transform: uppercase;
  border-right: 1px solid var(--border);
}
.header-cell:last-child {
  border-right: none;
}

.grid-body {
  display: flex;
  flex-direction: column;
  background-color: var(--surface-table);
}

.grid-row {
  display: grid;
  grid-template-columns: repeat(7, minmax(220px, 1fr));
  border-bottom: 1px solid var(--border);
  min-height: 140px;
}
.grid-row:last-child {
  border-bottom: none;
}

.grid-cell {
  border-right: 1px solid var(--border);
  position: relative;
  display: flex;
  flex-direction: column;
  padding: 8px; /* Internal padding so cards don't touch the borders */
  cursor: pointer;
  transition: background-color 0.2s ease;
}
.grid-cell:last-child {
  border-right: none;
}
.grid-cell:hover {
  background-color: rgba(9, 30, 66, 0.03);
}

/* Empty Cell Placeholder */
.empty-cell-overlay {
  margin: auto;
  color: var(--text-muted);
  opacity: 0;
  transform: scale(0.8);
  transition: all 0.2s ease;
}
.grid-cell:hover .empty-cell-overlay {
  opacity: 0.5;
  transform: scale(1);
}

/* Task Card - Trello Style inside the table cell */
.task-card {
  background-color: var(--surface-card);
  border-radius: 6px;
  padding: 10px 12px;
  box-shadow: 0 1px 2px rgba(9, 30, 66, 0.25);
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
  height: 100%;
}
.task-card:hover {
  box-shadow: 0 4px 8px rgba(9, 30, 66, 0.2);
}

.task-completed {
  opacity: 0.65;
  background-color: #fafbfc;
}
.task-completed .task-desc {
  text-decoration: line-through;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.task-time {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--secondary);
  padding: 2px 6px;
  border-radius: 3px;
}
.time-text {
  white-space: nowrap;
}

.task-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.task-checkbox {
  cursor: pointer;
  width: 16px;
  height: 16px;
  accent-color: var(--primary);
  margin: 0;
}

.btn-delete-task {
  opacity: 0;
  color: var(--text-muted);
  padding: 2px;
}
.task-card:hover .btn-delete-task { 
  opacity: 1; 
}
.btn-delete-task:hover {
  color: var(--danger);
  background: rgba(235, 90, 70, 0.1);
}

.task-desc {
  font-size: 0.9rem;
  margin: 0;
  color: var(--text-main);
  line-height: 1.4;
  word-break: normal;
  overflow-wrap: break-word;
  flex-grow: 1;
}

.task-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.tag-pill {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
  white-space: nowrap;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Modals */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(9, 30, 66, 0.54); 
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  animation: fadeIn 0.15s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: var(--surface-card);
  padding: 2rem;
  border-radius: 3px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 8px 16px -4px rgba(9, 30, 66, 0.25);
  animation: slideUp 0.2s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-content h2 { margin: 0 0 0.2rem 0; font-size: 1.25rem; color: var(--text-main); font-weight: 600; }
.modal-subtitle { color: var(--text-muted); margin-bottom: 1.5rem; font-size: 0.9rem; }

.form-group {
  margin-bottom: 1.2rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-group.half {
  flex: 1;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
}

.form-group input[type="time"], .form-group textarea, .input-tag-name {
  padding: 0.6rem 0.8rem;
  border: 2px solid var(--secondary);
  border-radius: 3px;
  font-family: inherit;
  font-size: 0.95rem;
  transition: all 0.2s;
  background-color: #fafbfc;
  color: var(--text-main);
}
.form-group input:focus, .form-group textarea:focus, .input-tag-name:focus {
  outline: none;
  border-color: var(--primary);
  background-color: #fff;
}

.tags-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 0.2rem 0;
}

.tag-selectable {
  padding: 0.4rem 0.8rem;
  border: 2px solid transparent;
  border-radius: 4px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.1s;
}
.tag-selectable:hover { filter: brightness(0.95); }
.tag-selectable.selected {
  box-shadow: 0 0 0 2px var(--surface-card), 0 0 0 4px var(--primary);
}

.no-tags-msg {
  font-size: 0.85rem;
  color: var(--text-muted);
  font-style: italic;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 2rem;
}

.btn-cancel {
  background: transparent;
  color: var(--text-main);
  border: none;
  padding: 0.6rem 1.2rem;
  border-radius: 3px;
  cursor: pointer;
  font-weight: 500;
}
.btn-cancel:hover { background: var(--secondary); }

/* Tag Management Area */
.tag-creation-area {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}
.input-tag-name { flex-grow: 1; }
.input-tag-color {
  width: 42px;
  height: 42px;
  padding: 0;
  border: 2px solid var(--secondary);
  border-radius: 3px;
  cursor: pointer;
  background: #fff;
}
.input-tag-color::-webkit-color-swatch-wrapper { padding: 2px; }
.input-tag-color::-webkit-color-swatch { border: none; border-radius: 2px; }

.tags-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  max-height: 250px;
  overflow-y: auto;
}

.tag-list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0.5rem 0.5rem 1rem;
  border-radius: 4px;
  background: var(--surface-card);
  box-shadow: 0 1px 1px rgba(9,30,66,0.1);
  border: 1px solid var(--secondary);
}

.tag-preview {
  font-size: 0.85rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
}

.btn-delete-tag {
  color: var(--danger);
}
.btn-delete-tag:hover {
  background: rgba(235, 90, 70, 0.1);
}
</style>
