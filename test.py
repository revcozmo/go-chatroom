#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall,sleep
import os
import sys
import random
import struct

def ss_listen(s):
    
    while True:
        sleep(random.random()*120)
        msg = os.urandom(15).encode('hex')
        s.send(struct.pack('<ii', 1, 0)+msg+'\n')

def ss_recv(s):
    while True:
        data = s.recv(1024)
        if len(data) == 0:
            sys.exit(1)

jobs = []

print "Init WORLD"
ss = socket.socket()
ss.connect(('localhost', 12345))
ss.send(struct.pack("<ii", 3, 0)+"a\n")
sleep(1)
print "Connectting...",
for x in xrange(3000):

    ss = socket.socket()
    ss.connect(('localhost', 12345))
    ss.send(struct.pack("<ii", 5, 0)+"a\n")
    jobs.append(spawn(ss_listen, ss))
    jobs.append(spawn(ss_recv, ss))

print "Done"
joinall(jobs)
