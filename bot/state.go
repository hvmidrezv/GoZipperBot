package bot

import "sync"

// UserState وضعیت هر کاربر را نگه می‌دارد
type UserState struct {
	State   string // "idle", "awaiting_files", "awaiting_zip"
	TempDir string // پوشه موقت برای ذخیره فایل‌ها
}

// StateManager وضعیت تمام کاربران را به صورت امن مدیریت می‌کند
type StateManager struct {
	users map[int64]*UserState
	mutex sync.RWMutex
}

// NewStateManager یک مدیر وضعیت جدید ایجاد می‌کند
func NewStateManager() *StateManager {
	return &StateManager{
		users: make(map[int64]*UserState),
	}
}

// Get وضعیت یک کاربر را برمی‌گرداند
func (sm *StateManager) Get(chatID int64) *UserState {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// اگر کاربر وجود نداشت، یک وضعیت پیش‌فرض برای او ایجاد می‌کنیم
	if _, ok := sm.users[chatID]; !ok {
		sm.mutex.RUnlock() // قفل خواندن را آزاد می‌کنیم تا بتوانیم بنویسیم
		sm.mutex.Lock()
		sm.users[chatID] = &UserState{State: "idle"}
		sm.mutex.Unlock()
		sm.mutex.RLock() // دوباره قفل خواندن را می‌گیریم
	}

	return sm.users[chatID]
}

// Set وضعیت یک کاربر را تنظیم می‌کند
func (sm *StateManager) Set(chatID int64, state *UserState) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.users[chatID] = state
}
