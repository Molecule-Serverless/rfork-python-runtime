import traceback
import json
import os
import time

import tornado
import tornado.ioloop
import tornado.web
import tornado.httpserver
import tornado.netutil

file_sock_path = 'fork.sock'
file_sock = None

def start_app_server():
    print("daemon.py: start app server on fd: %d" % file_sock.fileno())

    class SockFileHandler(tornado.web.RequestHandler):
        def post(self):
            try:
                data = self.request.body
                try :
                    event = json.loads(data)
                except:
                    self.set_status(400)
                    self.write('bad POST data: "%s"'%str(data))
                    return
                self.write(event)
            except Exception:
                self.set_status(500) # internal error
                self.write(traceback.format_exc())

    tornado_app = tornado.web.Application([
        (".*", SockFileHandler),
    ])
    server = tornado.httpserver.HTTPServer(tornado_app)
    server.add_socket(file_sock)
    tornado.ioloop.IOLoop.instance().start()
    server.start()

def start_fork_server():
    global file_sock
    global file_sock_path
    print("daemon.py: start fork server on fd: %d" % file_sock.fileno())
    file_sock.setblocking(True)

    while True:
        client, info = file_sock.accept()
        start = time.perf_counter_ns()
        pid = os.fork()
        end = time.perf_counter_ns()

        if pid:
            client.sendall(bytes(str(pid), 'utf8'))
            client.sendall(bytes('\n', 'utf8'))
            client.sendall(bytes(str(end-start), 'utf8'))
            client.sendall(bytes('ns', 'utf8'))
            client.sendall(bytes('\n', 'utf8'))
            client.close()
        else:
            file_sock.close()
            file_sock = None

            file_sock_path = 'fork.sock.' + str(os.getpid())
            file_sock = tornado.netutil.bind_unix_socket(file_sock_path)
            client.close()
            start_fork_server()


def main():
    global file_sock
    print("daemon.py: main")
    file_sock = tornado.netutil.bind_unix_socket(file_sock_path)
    start_fork_server()

if __name__ == '__main__':
    main()
