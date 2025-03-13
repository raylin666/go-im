package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mt/errors"
	"mt/pkg/logger"
	"runtime/debug"
	"time"
)

const (
	// 客户端连接超时时间
	heartbeatExpirationTime = 60 * time.Second
)

// Client 客户端连接
type Client struct {
	Ctx           context.Context
	Manager       WSClientManager
	Logger        *logger.Logger
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 连接实例对象
	Send          chan []byte     // 待发送的数据
	ConnectTime   time.Time       // 客户端连接时间
	HeartbeatTime time.Time       // 上次心跳时间
	Account       *Account        // 账号信息
}

func NewClient(
	ctx context.Context,
	manager WSClientManager,
	account *Account,
	conn *websocket.Conn,
	logger *logger.Logger) (client *Client) {
	var currentTime = time.Now()
	client = &Client{
		Ctx:           ctx,
		Manager:       manager,
		Logger:        logger,
		Addr:          conn.RemoteAddr().String(),
		Conn:          conn,
		Send:          make(chan []byte, 100), // 默认预创建容量为100的消息数据包
		ConnectTime:   currentTime,
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

// Read 读取客户端消息
func (c *Client) Read() {
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			c.Logger.Error("读取客户端消息异常", loggerFields...)
		}
	}()

	defer func() {
		// 关闭接收及待发送消息
		close(c.Send)

		c.Logger.Info("读取客户端消息结束, 已关闭数据接收", loggerFields...)
	}()

	for {
		// c.Conn.ReadMessage 该方法会阻塞等待, 直到收到消息才能继续往下执行
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				loggerFields = append(loggerFields, zap.String("close_desc", "socket 连接已被关闭"))
			}

			loggerFields = append(loggerFields, zap.Error(err))
			c.Logger.Error("读取客户端消息失败", loggerFields...)

			return
		}

		loggerFields = append(loggerFields, zap.String("message", string(message)))
		c.Logger.Info("读取客户端消息成功", loggerFields...)

		// 消息解析处理
		c.ParseMessageHandler(message)
	}
}

// Write 写入客户端消息
func (c *Client) Write() {
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			c.Logger.Error("写入客户端消息异常", loggerFields...)
		}
	}()

	defer func() {
		c.Logger.Info("写入客户端消息结束, 已关闭客户端连接", loggerFields...)

		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 写入待发送客户端消息错误并关闭连接
				c.Logger.Error("写入待发送客户端消息错误, 客户端连接将关闭", loggerFields...)

				return
			}

			// 将消息推送至客户端
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// ParseMessageHandler 消息解析处理 (对数据包进行合法校验、解析、事件消息分发及消息发送处理)
func (c *Client) ParseMessageHandler(message []byte) {
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
		zap.String("message", string(message)),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.Any("recover", r))
			c.Logger.Error("消息解析处理异常", loggerFields...)
		}
	}()

	// TODO 数据包合法性校验/解析消息数据包
	request := &MessageRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		loggerFields = append(loggerFields, zap.Error(err))
		c.Logger.Error("消息解析数据包合法性校验失败 json.Unmarshal", loggerFields...)

		// 返回错误给客户端
		c.WriteMessage([]byte("发送消息数据包协议格式错误"))

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		loggerFields = append(loggerFields, zap.Error(err))
		c.Logger.Error("消息解析数据包错误 json.Marshal", loggerFields...)

		// 返回错误给客户端
		c.WriteMessage([]byte("发送消息具体协议数据包格式错误"))

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

	loggerFields = append(loggerFields, zap.String("message_seq", seq), zap.String("message_events", event), zap.String("message_data", string(requestData)))
	c.Logger.Info("消息解析数据包完成", loggerFields...)

	// 采用 MAP 处理事件
	disposeFunc, ok := c.Manager.MessageEvent().GetDisposeFunc(event)
	if !ok {
		errMessage := fmt.Sprintf("事件消息处理: `%s` 事件不存在!", event)
		notEventErr := errors.New().CommandInvalidNotFound()
		responseCode = uint32(notEventErr.GetCode())
		responseMessage = notEventErr.GetMessage()
		c.Logger.Warn(errMessage, loggerFields...)
		// 返回错误给客户端
		c.WriteMessage([]byte(errMessage))

		return
	}

	responseCode, responseMessage, responseData = disposeFunc(c.Ctx, c, seq, requestData)
	// 判断该消息事件是否需要推送给客户端
	if responseSend == false {

		return
	}

	ok, err = c.WriteAgreementEventMessage(event, seq, responseCode, responseMessage, responseData)
	if err != nil {
		loggerFields = append(loggerFields,
			zap.Uint32("response_code", responseCode),
			zap.String("response_message", responseMessage),
			zap.Any("response_data", responseData),
			zap.Bool("response_send", responseSend),
			zap.Error(err))
		c.Logger.Error("事件消息处理响应数据错误 json.Marshal", loggerFields...)

		return
	}

	responseDataJson, _ := json.Marshal(responseData)
	loggerFields = append(loggerFields, zap.String("message", string(responseDataJson)))

	if !ok {
		c.Logger.Error("事件消息发送失败", loggerFields...)
		return
	}

	c.Logger.Info("事件消息发送成功", loggerFields...)
}

// WriteAgreementEventMessage 写入待发送协议消息到通道 (和客户端协商好的协议格式封装)
func (c *Client) WriteAgreementEventMessage(event string, seq string, code uint32, message string, data interface{}) (ok bool, err error) {
	// 事件不存在的消息不推送给客户端
	if hasEvent := c.Manager.MessageEvent().HasClientSupport(event); !hasEvent {
		return false, errors.New().SendMessageTypeNotFound()
	}

	response := NewMessageResponse(seq, event, code, message, data)
	headByte, err := json.Marshal(response)
	if err != nil {
		return false, err
	}

	return c.WriteMessage(headByte), nil
}

// WriteMessage 写入待发送消息到通道
func (c *Client) WriteMessage(message []byte) bool {
	if c == nil {
		return false
	}

	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
		zap.String("message", string(message)),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			c.Logger.Error("消息发送异常", loggerFields...)
		}
	}()

	c.Send <- message

	return true
}
