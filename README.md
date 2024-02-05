[![Golang CI](https://github.com/appcrash/media/actions/workflows/golang.yml/badge.svg)](https://github.com/appcrash/media/actions/workflows/golang.yml)
[![codecov](https://codecov.io/gh/appcrash/media/branch/master/graph/badge.svg?token=L76CE5EAKC)](https://codecov.io/gh/appcrash/media)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ec011376a0444882aa74b4e3ae04083e)](https://www.codacy.com/gh/appcrash/media/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=appcrash/media&amp;utm_campaign=Badge_Grade)
[![Coverity Scan](https://scan.coverity.com/projects/23415/badge.svg)](https://scan.coverity.com/projects/appcrash-media)

# What's This
A configurable **Media Server** that can handle, dispatch and manipulate media streams in realtime.

It is mainly developed by golang, other than c/c++ media engine such as gstreamer.
The project focus on low-latency, interactive, dev-efficiency as a generic media engine.

# How it works
The server accepts all kind of different jobs which are described by **Node Language**. 
That means the server can operate on different businesses simultaneously.
The client just describe how media stream should be filtered, routed or any other manipulation when creating session.
After session established, client can send messages to nodes in session to change media stream shape or behaviour at any time.
For example:

````{vabatim}
[rtp_src] -> [pubsub] -> {[audio_file_sink],[rtp_sink]}
````

The node topology is described by DSL called nmd as above.
Intuitively, the text of the example instructs the media server to accept RTP stream and clone two copies. One is 
redirected to another RTP endpoint, the other one is saved to file.
You even add now node and ask pubsub node to connect to it(by sending event) so the original stream is copied into three copies.
Every node of this project provides an atomic capability. Client composes them by nmd languages to finish a specific job.

The project is developed with speed, efficiency and scalability in mind.
It is easy to develop a new node to extend the functionality of media server.
Just inherit ***SessionNode***, write message&event handle functions, that's all.
We provide a cmd tool called **gen_traint**. It will inspect the syntax of newly-added node's source code and generate
necessary stub codes for you.






