package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mt/pkg/logger"
	"net/http"
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

		// 更新客户端连接心跳
		c.Heartbeat(time.Now())

		// 客户端发送消息解析处理
		c.ParseMessageHandler(message, true)
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

// ParseMessageHandler 消息解析处理 (对数据包进行合法校验、解析、消息事件分发及消息发送处理)
func (c *Client) ParseMessageHandler(message []byte, isClient bool) {
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
		zap.String("message", string(message)),
	}

	msgTitle := "发送消息事件"
	if isClient {
		msgTitle = "客户端" + msgTitle
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.Any("recover", r))
			c.Logger.Error(fmt.Sprintf("%s解析处理异常", msgTitle), loggerFields...)
		}
	}()

	// TODO 数据包合法性校验/解析消息数据包
	request := &MessageRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		loggerFields = append(loggerFields, zap.Error(err))
		c.Logger.Error(fmt.Sprintf("%s解析数据包合法性校验失败 json.Unmarshal", msgTitle), loggerFields...)

		// 返回错误给客户端
		code := http.StatusUnprocessableEntity
		c.WriteAgreementEventMessage("", "", uint32(code), http.StatusText(code), fmt.Sprintf("%s数据包协议格式错误", msgTitle))

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		loggerFields = append(loggerFields, zap.Error(err))
		c.Logger.Error(fmt.Sprintf("%s解析数据包错误 json.Marshal", msgTitle), loggerFields...)

		// 返回错误给客户端
		code := http.StatusUnprocessableEntity
		c.WriteAgreementEventMessage(request.Event, request.Seq, uint32(code), http.StatusText(code), fmt.Sprintf("%s具体协议数据包格式错误", msgTitle))

		return
	}

	// 判断是否客户端请求所支持的消息事件, 不在指定的消息事件客户端无法调用
	if hasEvent := c.Manager.MessageEvent().HasClientSupport(request.Event); !hasEvent {
		// 返回错误给客户端
		code := http.StatusNotFound
		c.WriteAgreementEventMessage(request.Event, request.Seq, uint32(code), http.StatusText(code), fmt.Sprintf("%s协议不存在", msgTitle))

		return
	}

	// 判断消息事件是否存在
	disposeFunc, ok := c.Manager.MessageEvent().GetDisposeFunc(request.Event)
	if !ok {
		// 返回错误给客户端
		code := http.StatusInternalServerError
		c.WriteAgreementEventMessage(request.Event, request.Seq, uint32(code), http.StatusText(code), fmt.Sprintf("服务端消息事件处理: `%s` 事件不存在!", request.Event))

		return
	}

	loggerFields = append(loggerFields, zap.String("message_seq", request.Seq), zap.String("message_event", request.Event), zap.String("message_data", string(requestData)))

	// TODO 将处理完成的数据包返回给客户端
	responseMessages := disposeFunc(c.Ctx, c, request.Seq, requestData)
	// responseMessages 为空时不需要回包, 程序逻辑处理后即完成
	if len(responseMessages) > 0 {
		lenLoggerFields := len(loggerFields)
		for _, responseMessage := range responseMessages {
			// 设置默认响应状态码及响应描述
			if responseMessage.Code == 0 {
				responseMessage.Code = http.StatusOK
				responseMessage.Msg = http.StatusText(int(responseMessage.Code))
			}

			tmpLoggerFields := make([]zap.Field, lenLoggerFields)
			copy(tmpLoggerFields, loggerFields)
			tmpLoggerFields = append(tmpLoggerFields,
				zap.String("response_event", responseMessage.Event),
				zap.Uint32("response_code", responseMessage.Code),
				zap.String("response_message", responseMessage.Msg),
				zap.Any("response_data", responseMessage.Data),
			)

			ok, err = c.WriteAgreementEventMessage(responseMessage.Event, request.Seq, responseMessage.Code, responseMessage.Msg, responseMessage.Data)
			if err != nil {
				tmpLoggerFields = append(tmpLoggerFields, zap.Error(err))
				c.Logger.Error("服务端消息事件发送失败: 消息回包处理数据错误 json.Marshal", tmpLoggerFields...)
				continue
			}

			c.Logger.Info("服务端消息事件发送成功: 消息回包处理完成", tmpLoggerFields...)
		}
	}

	c.Logger.Info(fmt.Sprintf("%s处理成功, 服务端消息回包已完成", msgTitle), loggerFields...)
}

// WriteAgreementEventMessage 写入待发送协议消息到通道 (和客户端协商好的协议格式封装)
func (c *Client) WriteAgreementEventMessage(event string, seq string, code uint32, message string, data interface{}) (ok bool, err error) {
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
			c.Logger.Error("服务端消息事件发送异常", loggerFields...)
		}
	}()

	c.Send <- message

	return true
}
