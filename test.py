#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall,sleep
import os
import random

def ss_listen(s):
    
    while True:
        sleep(random.random())
        s.send("NOR WORLD "+os.urandom(15).encode('hex')+'\n')
        data = s.recv(1024)

jobs = []

print "Connectting...",
for x in xrange(223):

    ss = socket.socket()
    ss.connect(('localhost', 12345))
    ss.send("JOI WORLD %s:%d\n" % ss.getsockname())
    jobs.append(spawn(ss_listen, ss))

print "Done"
joinall(jobs)
