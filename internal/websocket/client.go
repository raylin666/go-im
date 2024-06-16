package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/websocket/types"
	"runtime/debug"
	"time"
)

const (
	// 客户端连接超时时间
	heartbeatExpirationTime = 60 * time.Second
)

// Client 客户端连接
type Client struct {
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 连接实例对象
	Send          chan []byte     // 待发送的数据
	FirstTime     time.Time       // 首次连接时间
	HeartbeatTime time.Time       // 上次心跳时间
	Account       Account         // 账号信息
}

func NewClient(account Account, conn *websocket.Conn) (client *Client) {
	var currentTime = time.Now()
	client = &Client{
		Addr:          conn.RemoteAddr().String(),
		Conn:          conn,
		Send:          make(chan []byte, 100), // 默认预创建容量为100的消息数据包
		FirstTime:     currentTime,
		HeartbeatTime: currentTime,
		Account:       account,
	}

	return
}

// Heartbeat 更新连接心跳时间
func (c *Client) Heartbeat(currentTime time.Time) {
	c.HeartbeatTime = currentTime

	return
}

// IsHeartbeatTimeout 判断连接心跳是否超时
func (c *Client) IsHeartbeatTimeout(currentTime time.Time) (timeout bool) {
	if c.HeartbeatTime.Add(heartbeatExpirationTime).Before(currentTime) {
		timeout = true
	}

	return
}

// loggerConnFields 获取连接日志字段信息
func (c *Client) loggerConnFields() []zap.Field {
	var addr = zap.String("address", c.Addr)
	var accountId = zap.String("account_id", c.Account.ID)
	var firstTime = zap.Time("first_time", c.FirstTime)
	var heartbeatTime = zap.Time("heartbeat_time", c.HeartbeatTime)

	return []zap.Field{addr, accountId, firstTime, heartbeatTime}
}

// Read 读取客户端消息
func (c *Client) Read(ctx context.Context) {
	var loggerConnFields = c.loggerConnFields()

	defer func() {
		if r := recover(); r != nil {
			loggerConnFields = append(loggerConnFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			Logger(ctx).Error("读取客户端消息异常", loggerConnFields...)
		}
	}()

	defer func() {
		// 关闭接收及待发送消息
		close(c.Send)

		Logger(ctx).Debug("读取客户端消息结束, 已关闭数据接收", loggerConnFields...)
	}()

	for {
		// c.Conn.ReadMessage 该方法会阻塞等待, 直到收到消息才能继续往下执行
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			loggerConnFields = append(loggerConnFields, zap.Error(err))
			Logger(ctx).Error("读取客户端消息失败", loggerConnFields...)

			return
		}

		loggerConnFields = append(loggerConnFields, zap.String("message", string(message)))
		Logger(ctx).Info("读取客户端消息成功", loggerConnFields...)

		// 事件消息处理
		c.EventMessageHandler(ctx, message)
	}
}

// Write 写入客户端消息
func (c *Client) Write(ctx context.Context) {
	var loggerConnFields = c.loggerConnFields()

	defer func() {
		if r := recover(); r != nil {
			loggerConnFields = append(loggerConnFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			Logger(ctx).Error("写入客户端消息异常", loggerConnFields...)
		}
	}()

	defer func() {
		Logger(ctx).Debug("写入客户端消息结束, 已关闭客户端连接", loggerConnFields...)
		ManagerInstance().ClientManager().UnRegister <- c
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 写入待发送客户端消息错误并关闭连接
				Logger(ctx).Error("写入待发送客户端消息错误, 客户端连接将关闭", loggerConnFields...)

				return
			}

			// 将消息推送至客户端
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// WriteMessage 写入待发送消息到通道
func (c *Client) WriteMessage(ctx context.Context, message []byte) bool {
	if c == nil {
		return false
	}

	defer func() {
		if r := recover(); r != nil {
			loggerConnFields := c.loggerConnFields()
			loggerConnFields = append(loggerConnFields, zap.String("message", string(message)), zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			Logger(ctx).Error("发送消息异常", loggerConnFields...)
		}
	}()

	c.Send <- message

	return true
}

// WriteEventMessage 写入待发送事件消息到通道
func (c *Client) WriteEventMessage(ctx context.Context, event string, seq string, code uint32, msg string, data interface{}) (ok bool, err error) {
	// 事件不存在的消息不推送给客户端
	if _, hasEvent := ManagerInstance().GetEventHandler(event); !hasEvent {
		return false, nil
	}

	response := types.NewResponse(seq, event, code, msg, data)
	headByte, err := json.Marshal(response)
	if err != nil {
		return false, err
	}

	return c.WriteMessage(ctx, headByte), nil
}

// EventMessageHandler 事件消息处理, 对数据包进行合法校验、解析、事件消息分发及消息发送处理
func (c *Client) EventMessageHandler(ctx context.Context, message []byte) {
	var loggerConnFields = c.loggerConnFields()

	Logger(ctx).Info("进入事件消息处理", loggerConnFields...)

	defer func() {
		if r := recover(); r != nil {
			loggerConnFields = append(loggerConnFields, zap.Any("recover", r))
			Logger(ctx).Error("事件消息处理异常", loggerConnFields...)
		}
	}()

	// TODO 数据包合法性校验/解析消息数据包
	request := &types.Request{}
	err := json.Unmarshal(message, request)
	if err != nil {
		loggerConnFields = append(loggerConnFields, zap.Error(err))
		Logger(ctx).Error("事件消息数据包合法性校验失败 json.Unmarshal", loggerConnFields...)

		// 返回错误给客户端
		c.WriteMessage(ctx, []byte("事件消息发送数据包协议格式错误"))

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		loggerConnFields = append(loggerConnFields, zap.Error(err))
		Logger(ctx).Error("事件消息解析数据包错误 json.Marshal", loggerConnFields...)

		// 返回错误给客户端
		c.WriteMessage(ctx, []byte("事件消息协议格式错误"))

		return
	}

	// TODO 将处理完成的数据包返回给客户端
	var (
		seq   = request.Seq
		event = request.Event

		responseCode    uint32
		responseMessage string
		responseData    interface{}
		responseSend    bool
	)

	loggerConnFields = append(loggerConnFields, zap.String("seq", seq), zap.String("event", event), zap.String("data", string(requestData)))
	Logger(ctx).Info("事件消息解析数据包完成", loggerConnFields...)

	// 采用 MAP 处理事件
	if value, ok := ManagerInstance().GetEventHandler(event); ok {
		responseCode, responseMessage, responseData, responseSend = value(ctx, c, seq, requestData)
		// 判断该消息事件是否需要推送给客户端
		if responseSend == false {

			return
		}
	} else {
		e := defined.ErrorCommandInvalidNotFound
		responseCode = uint32(e.Code)
		responseMessage = e.Message
		Logger(ctx).Warn(fmt.Sprintf("事件消息处理: `%s` 事件不存在!", event), loggerConnFields...)
	}

	ok, err := c.WriteEventMessage(ctx, event, seq, responseCode, responseMessage, responseData)
	if err != nil {
		loggerConnFields = append(loggerConnFields,
			zap.Uint32("response_code", responseCode),
			zap.String("response_message", responseMessage),
			zap.Any("response_data", responseData),
			zap.Bool("response_send", responseSend),
			zap.Error(err))
		Logger(ctx).Error("事件消息处理响应数据错误 json.Marshal", loggerConnFields...)

		return
	}

	responseDataJson, _ := json.Marshal(responseData)
	loggerConnFields = append(loggerConnFields, zap.String("message", string(responseDataJson)))

	if !ok {
		Logger(ctx).Error("事件消息发送失败", loggerConnFields...)
		return
	}

	Logger(ctx).Info("事件消息发送成功", loggerConnFields...)
}
