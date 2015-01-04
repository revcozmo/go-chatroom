#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall, sleep
import sys
import struct
import time

count = 0
rps = 0

def ss_recv(s):
    global count
    while True:
        data = s.recv(1024)
        if len(data) == 0:
            sys.exit(1)
        count += 1
        sys.stdout.write('\rMSG: {0:.2f} {1}\n'.format(rps, data))
        sys.stdout.flush()


def counter():
    global count, rps
    while True:
        start_at = time.time()
        sleep(1)
        rps = count/(time.time()-start_at)
        count = 0

jobs = []


ss = socket.socket()
ss.connect(('localhost', 12345))
ss.send(struct.pack("<i", 5)+"WORLD a\n")
jobs.append(spawn(ss_recv, ss))
jobs.append(spawn(counter))

joinall(jobs)
