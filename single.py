#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall
import sys
import struct

def ss_recv(s):
    print "Recv"
    while True:
        data = s.recv(32)
        if len(data) == 0:
            sys.exit(1)
        print data

jobs = []

for x in xrange(1):

    ss = socket.socket()
    ss.connect(('localhost', 12345))
    ss.send(struct.pack("<i", 5)+"WORLD a\n")
    jobs.append(spawn(ss_recv, ss))

joinall(jobs)
