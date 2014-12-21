#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall,sleep
import os


def ss_listen(s):
    
    while True:
        sleep(1)
        s.send("hi "+os.urandom(15).encode('hex')+'\n')
        s.recv(1024)

jobs = []

print "Connectting...",
for x in xrange(50):

    ss = socket.socket()
    ss.connect(('localhost', 12345))

    jobs.append(spawn(ss_listen, ss))

print "Done"
joinall(jobs)
