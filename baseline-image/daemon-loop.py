# This python file is the runtime for directly startup from `runc run` and call code in /code/index.py
import traceback
import json
import os
import time

import importlib.util
import sys
import base64

func = None

def start_faas_server():
    global func
    sys.path.append("/code")
    # load code
    if func is None:
        func = importlib.import_module('index')

    ####### hard code start ######
    # invoke the function
    print("image_resize output: ")
    f = open("/code/test.jpg", 'rb')
    print(func.handler({'img': LoadTestImage(), 'height': 200, 'width': 200}))

def LoadTestImage():
    f = open("/code/test.jpg", 'rb')
    return str(base64.b64encode(f.read()), encoding='ascii')

def main():
    start_faas_server()

    i=0
    while True:
        time.sleep(1)
        #i = i+1



if __name__ == '__main__':
    main()

