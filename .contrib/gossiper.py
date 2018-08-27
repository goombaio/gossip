#!/usr/bin/env python

"""

This simple Python program implements a gossip network protocol. It is
intended only for fun and learning.

Written and directed by Stephane Bortzmeyer
<bortz@users.sourceforge.net> Licence as liberal as you want but no
warranty.

Gossip protocols are network protocols where each machine, each
*peer*, does not have a complete list of all peers. Instead, it knows
only a subset of them. In order to spread a message to all the peers,
every "gossiper" transmits the message to all the peers it knows, in
turn, each transmits it to all the peers it knows and so on. As long
as the set of peers is connected, the message will eventually reach
everyone.

The best-known examples of gossip protocols are the Network News (RFC
1036) and BGP (RFC 4271).

An important part of a gossip protocol is the history: peers must
remember which messages they sent, to avoid wasting time (or, worse,
creating endless loops) with peers which already know the message.

Note that a successful gossip protocol does not require every pair of
peers to communicate by the same means (Network News is a good
example: not everyone uses NNTP). But, in this simple example, the
protocol between two peers is fixed. Every peer has an ID, set at
startup. The "server" (the peer which replied to the connection) sends
its ID followed by a comma. The "client" (the peer which initiated the
connection) sends its ID followed by a comma and by the message (one
line only). In this implementation, for each peer, only a tuple (IP
address, port) is used to connect but it is not imposed by the
protocol (machines are identified by the ID, not by the IP address).

Peers remember the messages they have seen (in the global history) and
the messages they sent to each peer (in a per-peer history).

To use it, see the output of the -h option.
 
"""

import socket
import SocketServer
import time
import sys
import os
import signal
import threading
import Queue
import re
import getopt
import random

DEFAULT_PORT = 30480
DEFAULT_TIMESTAMP_DELAY = 300
DEFAULT_SIMULATION_DELAY = 60
DEFAULT_RETRY_DELAY = 120
DEFAULT_RETRY_ATTEMPTS = 10
DEFAULT_MAX_DISPLAY = 40

def current_time():
    return time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(time.time()))
  
def log(msg, peer=None, size=None):
    if peer is None:
        peer_str = ""
    else:
        peer_str = "%s -" % peer[0]
    if size is None:
        size_str = ""
    else:
        size_str = " %i bytes - " % size
    sys.stdout.write("%s %s%s%s\n" % (current_time(), peer_str, size_str,
                                               msg))

def usage(msg=None):
    sys.stderr.write("Usage: %s -i N [-p N] peer...\n" % sys.argv[0])
    sys.stderr.write("Use -h to get detailed help\n")
    if msg is not None:
        sys.stderr.write("%s\n" % msg)

def help(short=True):
    sys.stdout.write("""To use this program, there is one mandatory option, 
-i (or --id) to set the ID of this instance. It must be unique among the swarm.

There is one mandatory argument (at least one must be given), the peer information.
Each peer is to be entered as n,X[:Y] where n is the ID of the peer, X, its IP 
address and Y the facultative port. If X is an IPv6 address, it must be entered 
between brackets. If the port is ommitted, it defaults to %i.

Three examples of legal peer information are 33,[2001:db8:1::dead:babe] (peer 
ID 33, IPv6 address, default port), 629,192.0.2.1:8080 (peer ID 629, IPv4 address, 
port 8080) and 9231,[2001:DB8:99::bad:dcaf]:40000 (peer ID 9231, IPv6 address, port
40000).

The principal facultative option is probably -p which allows you to listen on
a different port than the default one (%i).

Complete example of usage, with three peers:

python gossiper.py -i 2 -p 4242 1,\[::1\] 3,1.2.3.4 5,\[::1\]:6666
\n""" % \
                         (DEFAULT_PORT, DEFAULT_PORT))
    if short:
        sys.stdout.write("Use --long-help to obtain help on all the options.\n")

def long_help():
    help(short=False)
    sys.stdout.write("""The complete list of possible options is:
   * -h or --help: get help
   * -p N or --port=N: sets the listening port
   * -i N or --id=N: sets the ID of this instance
   * -d N or --delay=N: sets the artifical delay we use before connecting 
                        to a peer (to makes things more realistic). The default 
                        is %i seconds.
   * -r N or -- retry=N: sets the delay we wait when a connection fails, before
                         retrying. the default is %i seconds.
   * -t N or --timestamp=N: sets the delay between two timestamps log messages.
                            The default is %i seconds.
\n""" % (DEFAULT_SIMULATION_DELAY, DEFAULT_RETRY_DELAY, DEFAULT_TIMESTAMP_DELAY))

class Entity:
    """ Models a peer on the network """    
    def __init__(self, family, address, port, id):
        self.family = family
        # TODO: allows a list of (address, port), not just one
        self.address = address
        self.port = port
        self.id = id

class InvalidPeerID(Exception):
    pass

class InvalidAddress(Exception):
    pass

class InvalidPort(Exception):
    pass

def parse(str):
    str = str.strip()
    try:
        (id, loc) = str.split(",")
        id = int(id)
    except ValueError:
        raise InvalidPeerID
    # TODO: raw IPv6 addresses without brackets
    if loc[0] == '[':
        match = re.search("^\[([0-9A-Za-z:]+)\](:(.*))?$", loc)
        if not match:
            raise InvalidAddress
        address = match.group(1)
        port = match.group(3)
        if port is None:
            port = DEFAULT_PORT
        try:
            binary_address = socket.inet_pton(socket.AF_INET6, address)
            port = int(port)
        except socket.error:
            raise InvalidAddress
        except ValueError:
            raise InvalidPort
        return Entity(socket.AF_INET6, 
                      address, 
                      port, 
                      id)
    else:
        if loc.find(':') >= 0:
            (address, port) = loc.split(':')
        else:
            address = loc
            port = DEFAULT_PORT
        try:
            binary_address = socket.inet_pton(socket.AF_INET, address)
            port = int(port)
        except socket.error:
            raise InvalidAddress
        except ValueError:
            raise InvalidPort
        return Entity(socket.AF_INET, 
                      address, 
                      port, 
                      id)

def pretty(server):
    if server.family == socket.AF_INET6:
        return ("[%s]:%i" % (server.address, 
                             server.port))
    elif server.family == socket.AF_INET:
        return ("%s:%i" % (server.address, 
                           server.port))
    else:
        return ("Unknown address family %i" % server.family)

class Timestamper(threading.Thread):

    def __init__(self, delay=DEFAULT_TIMESTAMP_DELAY):
        self.delay = delay
        threading.Thread.__init__(self)

    def run(self):
        while True:
            time.sleep(self.delay)
            log("Time stamp")

class Sender(threading.Thread):

    def __init__(self, peer, channel, simulation_delay=DEFAULT_SIMULATION_DELAY,
                 retry_delay=DEFAULT_RETRY_DELAY):
        self.peer = peer
        self.channel = channel
        self.history = {}
        self.simulation_delay = simulation_delay
        self.retry_delay = retry_delay
        threading.Thread.__init__(self)

    def run(self):
        while True:
            (myid, itsid, msg) = self.channel.get()
            if msg in self.history:
                continue
            # This is to simulate network delays
            time.sleep(generator.randint(1, self.simulation_delay))
            log("Sender task %s received \"%s\" from %i, connecting to %i (%s)" % \
                    (self.getName(), msg, itsid, self.peer.id, pretty(self.peer)))
            done = False
            attempts = 0
            while not done:
                try:
                    attempts += 1
                    self.s = socket.socket(self.peer.family, socket.SOCK_STREAM)
                    self.s.connect((self.peer.address, self.peer.port))
                    outf = self.s.makefile('w')
                    inf = self.s.makefile('r')
                    outf.write("%i,%s\n" % (myid, msg))
                    outf.close()
                    # TODO: read the peer ID and check it to be sure it is
                    # the one we wanted to talk with
                    inf.close()
                    self.s.shutdown(socket.SHUT_RDWR)
                    self.history[msg] = True
                    done = True
                except socket.error, error_msg:
                    log("Cannot connect to %s: %s" % (pretty(self.peer), error_msg))
                    if attempts > DEFAULT_RETRY_ATTEMPTS:
                        break
                self.s.close()
                time.sleep(generator.randint(self.retry_delay/2,self.retry_delay))
       
class RequestHandler(SocketServer.StreamRequestHandler):

    def handle(self):
        result = ""
        data = "DUMMY"
        size = 0
        self.wfile.write("%i,\n" % self.server.id)
        while data != "":
            data = self.rfile.read(1)
            if data == "\n" or data == "":
                break
            size = size + len(data)
            result = result + data
        try:
            (peer_id, message) = result.split(',')
            peer_id = int(peer_id)
            if message not in self.server.history:
                self.server.history[message] = True
                log("NEW message received from peer %i: \"%s...\"" % \
                        (peer_id, message[:DEFAULT_MAX_DISPLAY]), 
                    self.client_address, size)
                for mysender in mysenders.keys():
                    if mysender != peer_id:
                        mysenders[mysender].channel.put((self.server.id, 
                                                         peer_id, message))
            else:
                log("Ignoring known message from peer %i: \"%s...\"" % \
                        (peer_id, message[:DEFAULT_MAX_DISPLAY]), 
                    self.client_address, size)
        except ValueError: # Not a well-formatted message
            log("Ignoring badly formatted message \"%s...\"" % \
                    message[:DEFAULT_MAX_DISPLAY], self.client_address, size)
# TODO: two servers, for IPv4 and IPv6?
class Server(SocketServer.ThreadingMixIn, SocketServer.TCPServer): 

    def __init__(self, address, handler, id, num):
        self.id = id
        self.history = {}
        SocketServer.TCPServer.__init__(self, address, handler)
        log("Starting server %i at address %s, %i peers...\n" % \
                             (id, address, num))

class Watcher:
    """
    http://code.activestate.com/recipes/496735/
    
    This class solves two problems with multithreaded
    programs in Python, (1) a signal might be delivered
    to any thread (which is just a malfeature) and (2) if
    the thread that gets the signal is waiting, the signal
    is ignored (which is a bug).

    The watcher is a concurrent process (not thread) that
    waits for a signal and the process that contains the
    threads.  See Appendix A of The Little Book of Semaphores.
    http://greenteapress.com/semaphores/

    I have only tested this on Linux.  I would expect it to
    work on the Macintosh and not work on Windows.
    """
    
    def __init__(self):
        """ Creates a child thread, which returns.  The parent
            thread waits for a KeyboardInterrupt and then kills
            the child thread.
        """
        self.child = os.fork()
        if self.child == 0:
            return
        else:
            self.watch()

    def watch(self):
        try:
            os.wait()
        except KeyboardInterrupt:
            log("KeyBoardInterrupt")
            self.kill()
        sys.exit()

    def kill(self):
        try:
            os.kill(self.child, signal.SIGKILL)
        except OSError: pass

port = DEFAULT_PORT   
myid = None
timestamp_delay = DEFAULT_TIMESTAMP_DELAY
simulation_delay = DEFAULT_SIMULATION_DELAY
retry_delay = DEFAULT_RETRY_DELAY
try:
    optlist, args = getopt.getopt (sys.argv[1:], "p:i:t:d:r:h",
                               ["port=", "id=", "delay=", "retry=",
                                "timestamp=", "help", "long-help"])
    for option, value in optlist:
        if option == "--help" or option == "-h":
            help(short=True)
            sys.exit(0)
        elif option == "--long-help":
            long_help()
            sys.exit(0)
        elif option == "--port" or option == "-p":
            port = int(value) # TODO: handle the possible conversion exception
            # to provide a better error message?
        elif option == "--retry" or option == "-r":
            retry_delay = int(value) # TODO: handle the possible conversion exception
            # to provide a better error message?
        elif option == "--delay" or option == "-d":
            simulation_delay = int(value) # TODO: handle the possible conversion exception
            # to provide a better error message?
        elif option == "--timestamp" or option == "-t":
            timestamp_delay = int(value) # TODO: handle the possible conversion exception
            # to provide a better error message?
        elif option == "--id" or option == "-i":
            myid = int(value) # TODO: handle the possible conversion exception
            # to provide a better error message?
        else:
            # Should never occur, it is trapped by getopt
            usage("Unhandled option %s" % option)
            sys.exit(1)
except getopt.error, reason:
    usage(reason)
    sys.exit(1)
if len(args) == 0:
    usage("Not enough peers indicated")
    sys.exit(1)
if myid is None:
    usage("No ID indicated")
    sys.exit(1)
peers = []
mysenders = {}
for arg in args:
    try:
        peer = parse(arg)
        peers.append(peer)
    except InvalidPeerID:
        sys.stderr.write("No peer ID in %s (must precede the IP address and a comma)\n" % arg)
        sys.exit(1)
    except InvalidAddress:
        sys.stderr.write("No legal IP address in %s\n" % arg)
        sys.exit(1)
    except InvalidPort:
        sys.stderr.write("No legal port in %s\n" % arg)
        sys.exit(1)
generator = random.Random()
Watcher()
for peer in peers:
    channel = Queue.Queue()
    mysenders[peer.id] = Sender(peer, channel, simulation_delay, retry_delay)
    mysenders[peer.id].start()
Server.allow_reuse_address = True
# To have as many IP addresses as we want, we use IPV6. Depending on
# the system used and on options like Linux's sys.net.ipv6.bindv6only,
# the use of AF_INET6 may or may not allow also IPv4 connections.
Server.address_family = socket.AF_INET6
# TODO: an option to bind on a specific address
myserver = Server(("", port), RequestHandler, myid, len(mczysenders))
stamper = Timestamper(timestamp_delay)
stamper.start()
run_server = threading.Thread(target=myserver.serve_forever)
run_server.start()