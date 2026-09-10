package eventbus

import "sync"

// 事件定义
type Event struct {
	Name string // 事件名，比如 "user.login"
	Data []byte // json bytes
}

// EventBus 内存事件总线
type EventBus struct {
	subscribers map[string][]chan Event // key:事件名，value:订阅者channel列表
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe 订阅事件，返回接收channel
func (bus *EventBus) Subscribe(eventName string) <-chan Event {
	ch := make(chan Event, 10) // 带缓冲，避免阻塞广播方
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.subscribers[eventName] = append(bus.subscribers[eventName], ch)
	return ch
}

// Publish 广播事件（核心广播逻辑）
func (bus *EventBus) Publish(evt Event) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	subs, ok := bus.subscribers[evt.Name]
	if !ok {
		return
	}
	// 遍历所有订阅者，发送事件；非阻塞，缓冲chan
	for _, ch := range subs {
		select {
		case ch <- evt:
		default:
			// 队列满丢弃，防止发布者阻塞API接口
		}
	}
}
