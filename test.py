#!/usr/bin/env python
# encoding: utf-8

from gevent import socket, spawn, joinall,sleep
import sys
import random
import struct
import os
import time

msg_pack_num = 0
spread = 10


def ss_send(s):
    global msg_pack_num
    sleep(30)
    while True:
        sleep(random.random()*spread)
        msg = os.urandom(15).encode('hex')
        s.send(struct.pack('<i', 1) +"WORLD " +msg+'\n')
        msg_pack_num += 1

def ss_recv(s):
    while True:
        data = s.recv(1024)
        if len(data) == 0:
            return

def counter():
    global msg_pack_num
    while True:
        start_at = time.time()
        sleep(1)
        sys.stdout.write('\rMSG: {0} /s'.format(msg_pack_num/(time.time()-start_at)))
        sys.stdout.flush()
        msg_pack_num = 0

jobs = []

print "Init WORLD"
ss = socket.socket()
ss.connect(('127.0.0.1', 12345))
ss.send(struct.pack("<i", 3)+ "WORLD a\n")
sleep(1)
CON = 4000
print "Concurrent:%d -> %d msg/s" % (CON, CON/spread)
for x in xrange(CON):

    ss = socket.socket()
    ss.connect(('127.0.0.1', 12345))
    ss.send(struct.pack("<i", 5)+"WORLD a\n")
    jobs.append(spawn(ss_send, ss))
    jobs.append(spawn(ss_recv, ss))

jobs.append(spawn(counter))

print "Done"
joinall(jobs)
