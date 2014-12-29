#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall,sleep
import sys
import random
import struct
import os
def ss_send(s):
    while True:
        sleep(random.random()*85)
        msg = os.urandom(15).encode('hex')
        s.send(struct.pack('<i', 1) +"WORLD " +msg+'\n')

def ss_recv(s):
    while True:
        data = s.recv(1024)
        if len(data) == 0:
            sys.exit(1)

jobs = []

print "Init WORLD"
ss = socket.socket()
ss.connect(('127.0.0.1', 12345))
ss.send(struct.pack("<i", 3)+ "WORLD a\n")
sleep(1)
print "Connectting...",
for x in xrange(10000):

    ss = socket.socket()
    ss.connect(('127.0.0.1', 12345))
    ss.send(struct.pack("<i", 5)+"WORLD a\n")
    jobs.append(spawn(ss_send, ss))
    jobs.append(spawn(ss_recv, ss))

print "Done"
joinall(jobs)
