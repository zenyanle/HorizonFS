package types

// StreamEventType 定义 StreamEvent 的类型
type StreamEventType int

const (
	StreamEventError           StreamEventType = iota // 通用错误
	StreamEventSendError                              // 发送数据时发生错误
	StreamEventRecvError                              // 接收数据时发生错误
	StreamEventStateChange                            // 流状态改变事件 (例如: StreamActive, StreamClosed, StreamPaused)
	StreamEventFileProcessDone                        // 文件处理完成事件
	// ... 你可以根据需要添加更多事件类型，例如：
	// StreamEventChunkProcessed               // 数据块处理完成事件
	// StreamEventConnectionLost                 // 连接丢失事件
)

// StreamEvent 定义流事件的结构体
type StreamEvent struct {
	Type     StreamEventType `json:"type"`            // 事件类型
	StreamID string          `json:"stream_id"`       // 发生事件的 Stream ID
	Error    error           `json:"error,omitempty"` // 具体的错误信息 (如果 Type 是 Error 类型，否则为 nil)
	//State    StreamState       `json:"state,omitempty"`    // 如果是状态改变事件，包含新的状态
	FileName string `json:"file_name,omitempty"` // 如果事件与文件相关，可以包含文件名
	// ... 可以根据需要添加其他字段，例如时间戳、数据块 ID 等
}
