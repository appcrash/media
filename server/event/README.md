# What is event graph
It is a customized version of event system for rtp server.

```

+-------------------+           +--------------------+
|                   |           |                    |
|    session A      |           |     session B      |
|                   |           |                    |
|                   |           |                    |
|  +-------------+  |           |   +------------+   |
|  |             |  |           |   |            |   |
|  |   NODE      |  |  message  |   |   NODE     |   |
|  |             +--+-----------+--->            |   |
|  | audio sink  |  |           |   |audio mixing|   |
|  |             |  |           |   |            |   |
|  |             |  |           |   |            |   |
|  +-------------+  |           |   +-------+----+   |
|                   |           |           |        |
|                   |           |           |        |
|                   |           |           |        |
|  +-------------+  |           |           |        |
|  |             |  |           |           |        |
|  |    NODE     |  |           |           |        |
|  |             |  |  message  |           |        |
|  |audio source <--+-----------+-----------+        |
|  |             |  |           |                    |
|  |             |  |           |                    |
|  +-------------+  |           |                    |
|                   |           |                    |
|                   |           |                    |
|                   |           |                    |
+-------------------+           +--------------------+

```

Media server runs with lots of ongoing rtp session. it is often required to communicate to each other between 
sessions and other components. The goal of event graph is to simplify the communication pattern and also tailored 
for rtp server in which realtime is first priority.

# Concepts
There are three roles:
* node: basic communication unit, can receive events from other node or send events by link
* link: directed link, when two nodes get connected with a link, one can only send and the other one can only receive
* graph: consists of nodes and links, responsible for adding or removing nodes, build up or tear down links between 
  nodes

It is really very simple, here comes more details.

## Node && Event
Node has two important attributes: scope and name. It is straight mapped to the rtp session style, scope is the 
session name, and many nodes resides in the session with different names. Nodes can have the same name as long as 
they are in different scope. By scope and name every node can be located. Node can receive events from any other 
nodes provided links are ready. However, node doesn't know which node the event comes from. This is deliberate 
design. The node writer should try his best to keep node stateless no matter where the event comes from. What if a 
node requires response after sending the event. It can put the return-address in the event, when receiver finish its 
job and send the response to that address. In contrast to request/response model such as http, event delivering and 
receiving are nonblock. Once a event is delivered sender can keep delivering the other events and won't wait for 
response. The response comes into sender's receiving event queue asynchronously. 

## Link
Link is a must if two nodes needs communicating. As the event flows in one direction, two links are required if the 
two nodes need to talk to each other. Node can send link-up or link-down request to graph. Then node can use the link id to deliver events. Node only be avail of 
output links, however, new input link up is transparent to receiver node. Receiver doesn't know how many input links 
connected to it, but only pull events from his own event channel. Every two nodes can't establish more than one link 
in each direction, in other words,(sender,receiver,direction) tuple must be unique across whole graph.

## Graph
Graph add nodes into the event network. A node can be located only after added by graph, then comes to event 
receiving. Graph is bookkeeping of nodes and links info, such as how many input/output links a node owns, whether a 
link is available etc. When a node behaviour abnormally, graph can notify all senders that their receiver crashed 
and tear down the links from all senders. Graph would remove the bad node out after notified all senders.


# Implementation
User code implements *Node* interface then add it to *Graph* by **AddNode** method. If node struct defined some fields 
that graph interested such as **maxLink**, initialize them before **AddNode** is called, this can change the options 
of the added node instance. Refer to *Node* definition for details.

In *Node* struct, all state change events are notified by **OnXXX** callback methods. These methods are invoked one by 
one, which protected by a mutex, but not in the same goroutine. After node is added to graph, **OnEnter** will be 
invoked. As everything is async, other node can connect to the newly added node before **OnEnter** is called because 
graph has already record node info and the node is seen by everybody even though its initialization isn't done. So 
don't rely on **OnEnter**, it is just a notification not for initialization. Node should be ready for events after 
**AddNode** returns. **OnEnter** just tells the node it is able to deliver events now. 

In fact, user struct implements *Node* never talks to graph directly. The graph would create *NodeDelegate* for 
every node entering event graph, config the delegate then pass it to node by **OnEnter**. The delegate is the only 
way to communicate with others from the node's point of view: building up or tearing down links, request exiting from 
graph etc. Supposing node asks to talk to another one, it first requests link up through the delegate, the delegate 
passes the request to graph, after graph's sanity checking done, delegate get the new link id. Delegate return the 
created linkId to caller that invoke delegate's **RequestLinkUp**, and this link is buffered by delegate. Everytime 
user 
ask to 
**Deliver** event providing a link id, delegate will locate the receiver by inspecting link info, then call 
target node's receive method. The important thing is every node's max output link number is fixed once added to graph. 
User can change the max link number by define a private field **maxLink** as described before if the default max number 
can not satisfy him. Fixing the output link number is for performance while being reasonable in most real-world cases.  
The benefit is most frequently called method **Deliver** of delegate is lock-free because link buffer itself is 
fixed(read only) although elements of array may change. An Atomic load is enough to complete the receiver locating, 
so the cost is acceptable. As for most business, a node can foresee its link usage. Hold the link until node exits 
graph, while cases requiring dynamically add/remove links should be rare.

Event graph has only a control channel along with a receiving loop. Every node delegate has a control channel and a 
data channel along with receiving loops each. Delegate expose API to send control event to graph and graph also sends 
results of the control events back though delegate's control channel. So control events handling is transparent to 
user, under the framework control, so it should be bug-free. But inter-node data events never gives a guarantee of 
good behaviour. Sometimes bugs can panic the node. Here we limit the bug to the root cause of events passing, such 
as an event with null object that crash the poor receiver without any sanity check. If the panic happened in 
**OnEvent** handling routine, the delegate would catch the error, ask graph to break all links to crashing node then 
remove it out. 

User code should never put heavy jobs in **OnXXX** callbacks which would decrease the throughput of that node. Node 
delegate initialize buffered data channel for its node, the size of which can be overridden by **dataChannelSize** 
of node instance. If the node event handling can not catch up with senders, buffered channel is full and sender 
would fail to deliver. The contract is that, sender should not send events anymore if it found a delivery failure. 
The failure can be:
 * receiver event channel is full
 * link becomes invalid before node gets a notification, such as receiver exited graph(crash or requested), sender 
   requested the link down but still keep that link id to deliver event(user bugs)

Because everything is async, node encountering failure should retreat, then stop or try again later. Node should 
always keep function simple and focused, robust and testable. Compose simple nodes to a powerful monster!