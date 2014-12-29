#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall, sleep
import sys
import struct
import time

count = 0

def ss_recv(s):
    global count
    while True:
        data = s.recv(1024)
        if len(data) == 0:
            sys.exit(1)
        count += 1

def counter():
    global count
    previous_rps = 0
    while True:
        start_at = time.time()
        sleep(1)
        rps = count/(time.time()-start_at)
        if previous_rps > rps:
            flag = '-'
        else:
            flag = '+'
        sys.stdout.write('\rMSG {0}: {1:.2f}'.format(flag,rps))
        sys.stdout.flush()
        count = 0
        previous_rps = rps

jobs = []


ss = socket.socket()
ss.connect(('localhost', 12345))
ss.send(struct.pack("<i", 5)+"WORLD a\n")
jobs.append(spawn(ss_recv, ss))
jobs.append(spawn(counter))

joinall(jobs)
