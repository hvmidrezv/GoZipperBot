package bot

import "sync"

type UserState struct {
	State   string // "idle", "awaiting_files", "awaiting_zip"
	TempDir string // پوشه موقت برای ذخیره فایل‌ها
}

type StateManager struct {
	users map[int64]*UserState
	mutex sync.RWMutex
}

func NewStateManager() *StateManager {
	return &StateManager{
		users: make(map[int64]*UserState),
	}
}

func (sm *StateManager) Get(chatID int64) *UserState {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if _, ok := sm.users[chatID]; !ok {
		sm.mutex.RUnlock()
		sm.mutex.Lock()
		sm.users[chatID] = &UserState{State: "idle"}
		sm.mutex.Unlock()
		sm.mutex.RLock()
	}

	return sm.users[chatID]
}

// Set وضعیت یک کاربر را تنظیم می‌کند
func (sm *StateManager) Set(chatID int64, state *UserState) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.users[chatID] = state
}
