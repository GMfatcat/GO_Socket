import tornado.ioloop
import tornado.web
import tornado.websocket
import signal
import asyncio

class WebSocketHandler(tornado.websocket.WebSocketHandler):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.keep_alive_task = None

    async def open(self):
        print("WebSocket 連接已打開")
        self.keep_alive()

    async def on_message(self, message):
        try:
            print(f"收到消息: {message}")
            # 回覆消息給客戶端
            self.write_message(f"回覆: {message}")
        except Exception as e:
            print(f"處理消息時發生錯誤: {e}")

    def on_close(self):
        print("WebSocket 連接已關閉")
        if self.keep_alive_task:
            self.keep_alive_task.cancel()  # 取消心跳任務

    def check_origin(self, origin):
        return True

    def keep_alive(self):
        # 心跳機制，每 30 秒發送一次 ping 來保持連接
        self.ping("heartbeat")
        self.keep_alive_task = tornado.ioloop.IOLoop.current().call_later(30, self.keep_alive)

def signal_handler(signum, frame):
    print("接收到關閉信號，正在關閉伺服器...")
    tornado.ioloop.IOLoop.current().stop()

if __name__ == "__main__":
    app = tornado.web.Application([
        (r"/ws", WebSocketHandler),
    ])

    app.listen(9487)
    print("WebSocket 伺服器正在 port 9487 運行...")

    # 註冊 SIGINT 信號處理器來捕獲 Ctrl+C
    signal.signal(signal.SIGINT, signal_handler)

    try:
        tornado.ioloop.IOLoop.current().start()
    except Exception as e:
        print(f"伺服器關閉時發生錯誤: {e}")
    finally:
        print("伺服器已經關閉")
