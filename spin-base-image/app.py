import traceback
import json
import os
import time
import uuid
import socket
import array

import tornado
import tornado.ioloop
import tornado.web
import tornado.httpserver
import tornado.netutil

if os.environ.get('IMPORT_EXTRA_LIBS') is not None:
    import numpy

file_sock_path = 'fork.sock'
file_sock = None

def start_app_server():
    # print("daemon.py: start app server on fd: %d" % file_sock.fileno())

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
                tornado.ioloop.IOLoop.instance().stop() # Stop the server immediately after receiving the first request
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

def main():
    global file_sock
    file_sock = tornado.netutil.bind_unix_socket(file_sock_path)
    start_app_server()

if __name__ == '__main__':
    main()
