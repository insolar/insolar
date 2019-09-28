.. _architecture:

============
Architecture
============

Below is the platform architecture diagram aimed to address the :ref:`interconnected layers <big_pic>`. The architecture has multiple components and consensuses to address the complexity and variety of requirements.

Components in the diagram are *clickable*, the links will lead you to respective definitions.

.. uml:: architecture.uml

All components communicate via messaging to achieve respective :ref:`consensuses <consensuses>` and use :term:`pulses <pulse>` to stay in sync. Let's decompose the architecture to learn the key design concepts.

.. _fed_of_clouds:

Clouds and Their Federations
----------------------------

:term:`Clouds <cloud>` organize and unify software capabilities, hardware capacities, and the financial and legal liability of :term:`nodes <node>` to ensure seamless operation of business services. The Insolar Platform transparently connects multiple clouds and each cloud is governed independently, e.g., by a community, company, industry consortia, or national agency. Thus, multiple clouds can unite into a federation on the Insolar network.

The cloud itself establishes governance of both network operations and business logic. Therefore, it is a dual entity that controls:

* The :term:`network` and components deployed during :term:`node` setup, such as: 

  * bootstrap configuration; 
  * globula discovery and split-protection protocols;
  * node activation and deactivation protocols with the list of currently active nodes and blacklisted ones;
  * real-time detection protocols of execution fraud.

* A special :term:`domain` that is stored by the cloud itself and carries rigid configuration and rules such as:

  * procedures for registering and deregistering nodes;
  * postexecution fraud detection procedures;
  * compensation and penalization procedures;
  * marketplace rules for processing capacity.

.. _domains:

Domains
-------

Domains establish governance of contracts and nodes, thus, acting as *super contracts* that can contain :term:`objects <object>` and their history (:term:`lifelines <lifeline>`) and can apply varying policies to the lifelines contained within. Policies can differ with regards to particular rules:

* Changing the domain itself.
* Access to/from other domains for lifelines.
* Logic validation, e.g., consensus, number of voters.
* Code mutability -- possibility of changing the code and change procedures.
* Mutability of object history contained in the lifeline. These rules allow to implement GDPR or legal action via authorization requirements defined by the domain.
* Applicability of custom cryptography schemes requested from the cloud that deploys them.

.. _globulas:

Globulas
--------

Globula is a network of up to 1,000 :term:`nodes <node>`. It can run as a truly decentralized network with consistency established by a leaderless, pure BFT-based consensus mechanism, a :ref:`globula network protocol <network_consensus>`.

Insolar also supports larger node networks of up to 100 globulas (a total of 100,000 nodes) that behave transparently across such networks in accordance with whichever contract logic is in place. Such networks rely on the :ref:`inter-globula network protocol <network_consensus>` with leader-based consensus.

.. _multi_role_nodes:

Multi-Role Nodes
----------------

Insolar utilizes a multi-role model for :term:`nodes <node>`: each node has a single :ref:`static role <static_roles>` that defines its primary purpose and a set of :ref:`dynamically assigned roles <dynamic_roles>`. Dynamic role allocation functions enable the :ref:`omni-scaling <omni_scaling>` feature of the Insolar Platform.

.. _static_roles:

Static Roles
~~~~~~~~~~~~

The node’s static role defines what kind of resource and functionality are delivered by that node to the network, and how the network uses such nodes. The network recognizes four static role categories:

* :ref:`virtual <virtual>` -- performs calculations;
* :ref:`light material <light_material>` -- performs short-term data storage and network trafficking;
* :ref:`heavy material <heavy_material>` -- performs long-term data storage;
* :ref:`neutral <neutral>` -- participates in the network consensus (not in the workload distribution) and has at least one utility role.

Static role correlates with the type of resource the node can provide to the cloud, and is a part of the :ref:`omni-scaling <omni_scaling>` feature of the Insolar Platform. All static role categories are detailed below.

.. _neutral:

Neutral nodes
^^^^^^^^^^^^^

Neutral nodes participate in the :ref:`network consensus <network_consensus>` but do not receive any workload automatically distributed by the Insolar network. Neutral nodes serve particular functions:

* API exposure,
* block explorer support,
* discovery support,
* key management.

.. _virtual:

Virtual nodes
^^^^^^^^^^^^^

Virtual nodes are stateless, fast, easy to join and leave, and do not need data recovery. On the Insolar network, virtual nodes do the following:

* receive and handle requests to execute contracts;
* :ref:`execute and validate contracts <execution_validation>`;
* read the latest :term:`contract <object>` state and generate updates (i.e., new :term:`records <record>`) for material nodes;
* enable CPU scalability;
* handle contract-related data encryption when provided with access to relevant key storages.

.. _light_material:

Light material nodes
^^^^^^^^^^^^^^^^^^^^

Light material nodes are stateful and they automatically collect hot data and indices upon restart. On the Insolar network, light material nodes do the following:

* build blocks;
* manage data access and do audit;
* provide caching for recent data;
* enable scalability of network throughput;
* perform data retrieval and storage operations for :ref:`virtual nodes <virtual>`;
* redirect requests to relevant material nodes when the required data is not available;
* maintain indices of the most recent records, attribute indices, and other functions;
* deduplicate and recover requests in case of virtual node failures;
* assist :ref:`heavy material nodes <heavy_material>` by serving as temporary backup and cache for individual blocks;
* serve as integrity validators, recovery sources, proof-of-storage approvers, and handover voters;
* collect and register :term:`dust` (e.g., service inconsistency reports, long operations, logs).

Although light nodes can add dust, in case of :term:`lifelines <lifeline>`, they can only add records on behalf of relevant :ref:`virtual nodes <virtual>`. This is enforced by signatures and their checks during new :ref:`block validations <material_execution_validation>`.

.. _heavy_material:

Heavy material nodes
^^^^^^^^^^^^^^^^^^^^

Heavy material nodes are stateful and require recovery and content revalidation (proof-of-storage), both periodically and upon rejoining the network. On the Insolar network, heavy material nodes do the following:

* provide long-term data storage and scalability of storage capacity;
* store all data received from :ref:`light material nodes <light_material>` (and, in turn, from :ref:`virtual nodes <virtual>`);
* check data integrity but are unable to introduce or change data or form a block;
* ensure the required level of block replication and the maximum data density (scattering) to reduce the impact of data leakage from a single material node (heavy or light).

Heavy material nodes differ significantly from other nodes -- they store lots of data and must take additional measures to mitigate the following risks:

* losing (or corrupting) data but not having enough copies, or
* data leakage caused by the accumulation of too much data on a single node.

Heavy material node's implementation is simplified for the TestNet 1.1 and will gradually extend during the development of Insolar's enterprise version.

Moreover, additional network protocol is implemented to maintain backups and archival storage nodes without burdening the main Insolar network consensus.

.. _dynamic_roles:

Dynamic Roles
~~~~~~~~~~~~~

In addition to the node's static role, it can be equipped with dynamic ones -- roles able to change.

:ref:`Virtual nodes <virtual>` can have the following roles and respective responsibilities:

* **Virtual executor** handles operations on a :term:`lifeline` and builds new :term:`object <object>` states.
* **Virtual validator** verifies virtual executor's actions from previous :term:`pulses <pulse>`.

:ref:`Light material nodes <light_material>` can have the following roles and respective responsibilities:

* **Material executor** forms new :term:`blocks <jet drop>` and grants access to previous blocks.
* **Material validator** checks the block's validity and consistency.
* **Material stash** caches hot data and relevant indices (current states of all :term:`objects <object>`) and syncs the indices among other stash nodes.

In essence, all the nodes take part in two kinds of :ref:`execution and validation <execution_validation>` procedures, depending on their dynamic roles: **virtual** and **material**. :ref:`Heavy material nodes <heavy_material>` rely on validation performed by light material ones.

A node can have multiple dynamic roles, e.g., a virtual node can be selected via the :term:`entropy <pulse>` to be an executor for one :term:`lifeline` and a validator of another.

Dynamic roles are designed to:

* enable dynamic and straightforward scaling of the network;
* require minimal preparation to become operational;
* get new workload allocations while dynamic roles of all the nodes change with every :term:`pulse`.

.. _utulity_roles:

Delegated and Utility Roles
~~~~~~~~~~~~~~~~~~~~~~~~~~~

In addition to static and dynamic roles, nodes can take on delegated and utility roles that serve additional functions: caching, inter-globula coordination, and node joining.

.. _contracts:

Contracts
---------

The Insolar's main principle is that everything is a :term:`contract <object>` on the Insolar Platform. Contracts are stored as :term:`lifelines <lifeline>` in the :ref:`ledger <ledger>` and are based on general-purpose programming languages such as Golang or Java. They allow existing practices, libraries, and development environments to be used straightforwardly.

A contract developer may focus solely on the contract logic and calls of other contracts, while such details as location & implementation of other contracts are managed transparently by the platform. Every contract has :ref:`domain-level <domains>` managed rules that define the contracts handling:

* policies for code updates,
* validation requirements,
* inbound or outbound call permissions.

In addition to :ref:`governance <domains>` with logical rules, domains can also be deployed in separate :ref:`clouds <fed_of_clouds>` for stronger network security and data inspection on network edges, while contract/business logic can dynamically tune validation performed by the Insolar Platform to balance **costs**, **risks**, and **performance** by adjusting *quantity* and *quality* (stake or liability levels) of :ref:`validators <dynamic_roles>` involved.

Contracts also have individual time tracking and resources which can be subsequently connected to custom billing procedures and prepaid (or on-spot) allocation of :ref:`hardware capacities <multi_role_nodes>`. Moreover, the :ref:`ledger <ledger>` that stores contract data applies strict controls on the following:

* Data access by requiring signatures from :ref:`nodes <multi_role_nodes>` that need the access;
* Scattering of versioned data across multiple :ref:`storage nodes <heavy_material>` to significantly reduce risks of fraud, intrusions, or data leaks.

Furthermore, Insolar guarantees to execute any contract and ensures duplicate calls will not emerge in case of hardware, system, or network failure.

For practical enterprise use, Insolar contracts can store and transfer large data :term:`objects <object>` with the following benefits:

* on-chain, without the need for additional systems integrations;
* with algorithms to provide :ref:`network traffic <globulas>`, :ref:`CPU <virtual>`, and :ref:`storage <heavy_material>` scalabilities.

.. _contract_determinism:

Contract Determinism
~~~~~~~~~~~~~~~~~~~~

As the platform already reduces determinism via network messaging, Insolar applies relatively relaxed requirements regarding the determinism of :ref:`contracts <contracts>`. As such, a method invocation:

* on the same :term:`object <object>` state,
* with the same parameters,
* and on the same :term:`pulse`;

Should:

* produce exactly the same results,
* consume roughly the same amount of :ref:`CPU resources <virtual>`.

Contract execution methods that run longer than one full pulse must be explicitly declared with an *execution duration* policy.

A contract that does not produce the same results under given conditions will not pass :ref:`validation <execution_validation>`. In this case, all expended efforts will be at the cost of the party that deploys the contract (as opposed to the caller). Insolar records information on spent efforts in :term:`sidelines <sideline>` and can track assigned limits, however, the actual billing and payment execution must be handled by :ref:`governance logic <domains>` (i.e., by other contracts).

Although :ref:`virtual nodes <virtual>` are used to isolate contracts incompatible with security or governance rules, the new contract's code can only be introduced to Insolar as source code, with compilation and static inspection performed by :ref:`nodes <multi_role_nodes>` in accordance with an applicable :ref:`governance model <fed_of_clouds>`.

To provide contract execution determinism, Insolar utilizes its :ref:`network consistency <network_consistency>`.

.. _network_consistency:

Network Consistency
~~~~~~~~~~~~~~~~~~~

Insolar uses the :ref:`network layer <network_consensus>` to ensure view consistency across the whole network. The next step is to facilitate the efficient and secure execution of contracts across all :ref:`virtual nodes <virtual>`.

To this end, Insolar:

* :ref:`sets apart the functionality <multi_role_nodes>` requiring different resources and permissions,
* distributes workloads across all available/active nodes of the Insolar network using entropy.

As a result, all nodes have:

* the same :ref:`entropy <pulsars>` value,
* a list of active :ref:`nodes <multi_role_nodes>`.

Insolar does not use node workload statistics to provide network consistency, instead, it implements pseudo-random workload distribution.

The reason is simple: a trustful workload factor in distributed systems requires full visibility and operations aggregation but they still do not guarantee smooth workload distribution when workloads fluctuate faster than the average duration of a workload control cycle (aggregate statistics – balance – execute). 

Pseudo-random workload distribution can cause distribution anomalies within a workload control cycle but it provides a relatively smooth distribution on longer timescales, without the need for full visibility and operations aggregation.

Such a workload distribution and the entorpy-based allocation functions for :ref:`dynamic roles <dynamic_roles>` are the core instruments that enable the :ref:`omni-scaling <omni_scaling>` feature of the Insolar Platform. This feature provides a balance in accordance with client's needs.

Processing costs can be traded off against:

* **Uninsured risks**. Suitable for situations where a cheaper transaction is executed but fewer validators verify said transaction, meaning greater risk of loss.
* **Processing speed**. It can be increased to the detriment of operational risk:

  * frequent transactions could be processed without awaiting validation, or
  * validations may be batched together and processed following some delay, leading to the possibility of resource-consuming rollbacks.

.. _execution_validation:

Execution & Validation
----------------------

The Insolar Platform works on the principle of actions executed by one node, validated by many.

The number of selected validators can be determined in accordance with the :ref:`business process <domains>` at hand and, since validators in shared enterprise networks will have liability and legal guarantees, this works as transaction insurance.

As described in the :ref:`network consistency section <network_consistency>`, validator selections are *not* based on voting; instead, they are part of the :ref:`omni-scaling <omni_scaling>` feature. Insolar uses the active node list and :ref:`entropy <pulsars>` generated by consensus of the :ref:`globula network protocol <network_consensus>`, and then applies deterministic allocation functions for :ref:`node roles <dynamic_roles>`. This avoids wasting efforts on numerous per-transaction and network-wide consensuses.

Since Insolar sets apart functionality using :ref:`node roles <multi_role_nodes>`, it has two sets of execution & validation procedures: **virtual** and **material**.

.. _virtual_execution_validation:

Virtual Execution & Validation
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Nodes with :ref:`virtual static roles <virtual>` carry out **virtual** execution & validation:

#. The network selects (determines based on :term:`entropy <pulse>`) a specific virtual node to become a :ref:`virtual executor <dynamic_roles>`. Upon receiving the request, the executor:

   #. Registers the request within the current :term:`pulse`.

      In case the request arrives to a 'busy' virtual executor, it can delegate the execution of an :term:`object <object>` to other virtual nodes (not necessary to virtual executors). Moreover, multiple requests can be executed within the same pulse when opportunistic execution/validation is allowed by the caller or by the called object.

   #. Executes the request on the :term:`object <object>` (contract).
   #. Collects the results of outbound calls.
   #. Provides :term:`lifeline <lifeline>` and :term:`sideline <sideline>` updates for validation by other nodes.

#. Once the executor’s status expires, the network selects :ref:`virtual validators <dynamic_roles>` from the list of active :ref:`virtual nodes <virtual>` on a new :term:`pulse <pulse>` (new entropy), meaning executors cannot predict which nodes will validate transactions, thereby avoiding a collusion scenario. 

#. Each virtual validator:

   #. Checks that the request is legitimate.
   #. Executes the request on the :term:`object <object>` (contract) a second time.
   #. Checks that the request returns the same response given the :ref:`same arguments <contract_determinism>`.
   #. Checks that the request performs the same outbound calls.

#. Lastly, the outbound calls validation is stacked into a single validation round as validators use signed results collected by previous executors.

A single virtual executor can execute long requests that span several pulses. To do this, the virtual node that started the execution asks current executors in each pulse for tokens that give the execution permission.

.. _material_execution_validation:

Material Execution & Validation
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Nodes with :ref:`light material static roles <virtual>` carry out **material** execution & validation:

#. The network selects (determines based on :term:`entropy <pulse>`) a specific light material node to become a :ref:`light material executor <dynamic_roles>`. Upon receiving data requests from the virtual executor in the current :term:`pulse <pulse>`, the light material executor:

   #. Manages data access for :term:`contracts <object>`.
   #. Performs data retrieval and storage operations for :ref:`virtual executors <dynamic_roles>`.
   #. Builds a new :term:`block <jet drop>` from the :term:`lifeline <lifeline>` & :term:`sideline <sideline>` updates sent by the virtual executor.
   #. Splits (or merges) :term:`jets <jet>` if required.

#. Once the executor’s status expires, the network selects :ref:`material validators <dynamic_roles>` from the list of active :ref:`light material nodes <light_material>` on a new :term:`pulse <pulse>` (new entropy), meaning executors cannot predict which nodes will validate transactions, thereby avoiding a collusion scenario. 

#. Each material validator checks that the light material executor has formed the last :term:`block <jet drop>` correctly. The block must have:

   * Correct hashes.
   * Correct order of new :term:`records <record>` in the affected :term:`filaments <filament>`. 
   * No contradictions between records in the filaments.

   In addition, each validator ensures that the executor made the right decision to split (or merge) the corresponding :term:`jet <jet>`.

Upon each pulse, every light material node sends the data they formed to :ref:`heavy material nodes <heavy_material>`. However, light nodes keep hot data and share hot indices among a number of :ref:`light material stash <dynamic_roles>` nodes.

Light material stash nodes are nodes which have been :ref:`light material executors <dynamic_roles>` for a number of past :term:`pulses <pulse>`. The number is called a *stash history limit* and its default value is 5 but it is configurable within a :ref:`cloud <fed_of_clouds>`. Thus, stash material nodes provide caching for recent data.

.. _consensuses:

Consensuses
-----------

Consensus procedures vary in their degree of control by business logic, with two consensus procedures available:

* **Domain-defined consensus**: procedures that are a set of Raft-like protocols with :ref:`entropy-controlled <pulsars>` voter selection. These protocols are applied to an :term:`object <object>` after a series of changes. Such protocols can be chosen at the :ref:`domain <domains>` level and configured at the transaction level.
* **Utility consensus**: procedures -- a set of protocols -- that cover various platform operations not directly operated or required by business logic, including network consensus, pulsar consensus, and traffic cascade.

Different sets of consensus procedures affect every action applied to :term:`lifelines <lifeline>`: :ref:`logic <logic_consensus>`, :ref:`storage <storage_consensus>`, :ref:`network <network_consensus>`, and :ref:`pulsar <pulsar_consensus>` consensuses.

.. _logic_consensus:

Logic Consensus
~~~~~~~~~~~~~~~

Ensures that actions applied to an :term:`object` were performed correctly considering the object’s state, input parameters, and external dependencies (calls).

For more information on logic consensus, see the :ref:`virtual execution & validation section <virtual_execution_validation>`.

.. _storage_consensus:

Storage Consensus
~~~~~~~~~~~~~~~~~

Ensures that:

#. :term:`Nodes <node>` which participated in logical consensus had allocated roles.
#. :term:`Records <record>` generated by the nodes are structurally and referentially valid.

For more information on storage consensus, see the :ref:`material execution & validation section <material_execution_validation>`.

.. _network_consensus:

Network Consensus
~~~~~~~~~~~~~~~~~

Ensures :term:`node` availability and synchronization of time and state among nodes and provides consistent allocation of :ref:`dynamic roles <dynamic_roles>` to nodes. There are two consensus protocols behind the network consensus:

* **Globula network protocol**: a truly decentralized BFT-like protocol without any consensus leader that establishes the consistency of a globula (a smaller network of up to 1,000 nodes).
* **Inter-globula network protocol**: a leader-based protocol that extends the GNP and establishes consistency among globulas of the Insolar network (up to 100 globulas or 100,000 nodes).

The network layer of Insolar deals with the consistency of network node's view and :term:`pulse` distribution. Pulse is a signal carrying entropy (randomness) that triggers the production of a new :term:`block <jet drop>`.

The entropy's consistency and the set of active nodes on the network are vital for the methodology of executed by one node, validated by many. Nodes are selected from the active node list to perform :ref:`different functions <dynamic_roles>`, while entropy and consistency ensure behavioral consensus across all nodes. :ref:`Validator <dynamic_roles>` nodes are selected only on a new pulse to ensure that :ref:`executor <dynamic_roles>` nodes cannot collude with validators.

In addition to the aforementioned consensuses, :ref:`pulsars <pulsars>` can have their :ref:`own <pulsar_consensus>`.

.. _pulsars:

Pulsars
-------

Pulsars running on a pulsar protocol represent a separate logical layer that is responsible for network synchronization and provides a source of randomness (:term:`pulses <pulse>`). Interoperability of :term:`nodes <node>` within a single :term:`cloud` depends on pulses and all nodes must be on the same pulse to process new requests or operations.

Pulsars can run either on the same network or an entirely separate one. Cases of the former include:

* private networks that can implement a dedicated server;
* cross-enterprise and hybrid networks that can use a shared network of pulsars yet run individual installations of Insolar networks;
* and public networks that can use trusted pulsar nodes or run the pulsar function on other nodes.

In case of multiple pulsars on the network, their consensus generates the :term:`pulses <pulse>`.

.. _pulsar_consensus:

Pulsar Consensus
~~~~~~~~~~~~~~~~

:term:`Clouds <cloud>` define the pulsar selection rules and they can vary significantly. On enterprise networks, servers that complete no other operations manage the selection, whereas on public networks, it may be a random subset of 10 to 50 nodes with high uptime. Other configurations are also possible for different network types.

Default :term:`pulse` generation is based on BFT-consensus among pulsars, where *each member contributes* to entropy and *none can predict it*. The pulsar protocol enables entropy generation in a way that prevents individual nodes from being able to predictably manipulate the entropy through vote withdrawals.

This protocol does not include negotiations related to pulsar membership or pulse duration -- such parameters are considered as preconfigured or preagreed. The default pulse duration is 10 seconds.

As a consensus result, pulsars distribute the collaboratively-generated entropy signed by every pulsar to every node on the network.

.. _ledger:

Ledger
------

Ledger is a common term for distributed storage, a network of nodes that store data.

As described in the :ref:`static roles section <static_roles>`, material nodes are responsible for storing data and providing it on requests for :ref:`virtual nodes <virtual>`. Virtual nodes create and sign new information and pass it to material nodes to store. So, material nodes do not create or modify information (:term:`objects <object>`) with the exception of specifically defined meta data.

A typical :term:`object <object>` workflow is as follows:

.. uml::

   skinparam backgroundColor transparent
   skinparam entity { 
     backgroundColor transparent
   }

   entity "Virtual node" as v [[../architecture.html#virtual]]
   entity "Material node" as m [[../architecture.html#light-material]]

   v -> m : Get Object
   m -> v : [[../glossary.html#term-object Object]]
   v -> v : Perform calculations
   v -> m : Add modification [[../glossary.html#term-record record]] to the object

.. _records:

Records
~~~~~~~

Data is stored in the ledger as a series of immutable :term:`records <record>`. All records are created and signed by :ref:`virtual nodes <virtual>`. Each record is addressed by its hash and a :term:`pulse <pulse>` number. Records can contain a reference to another record, thus, creating a chain. An example of a chain is the :term:`object's <object>` :term:`lifeline <lifeline>`. Each :ref:`material node <static_roles>` is responsible for its own lifelines determined by their hashes.

In the Insolar's key-value storage, the key is a fixed structure -- a combination of a pulse number and a value hash. The value can be one of several types:

* :term:`Record <record>` -- immutable structured data unit. Can form chains if each record references a previous one in succession.
* Index -- meta information about record chains, e.g., pointers to the latest record in a chain. Represents an :term:`object <object>`.
* Blob -- immutable payload. Used to store (potentially big) chunks of serialized data, e.g., object's memory. Usually, records refer to blobs to store application data.

.. _requests:

Requests
~~~~~~~~

Each operation performed by :ref:`virtual nodes <virtual>` is registered as a request in the ledger. Request is a single :ref:`record <records>` that contains information necessary to perform an operation. Each request belongs to an :term:`object <object>` and is affined to it.

.. _results:

Results
~~~~~~~

Each operation performed by :ref:`virtual nodes <virtual>` has exactly one result. Although an operation can have many side effects (:term:`records <record>` stored in the ledger), result represents a summary of that operation. So, each finished request has its own result, i.e., result references its request. A request without an associated result stored in the ledger is a *pending* one.

.. _objects:

Objects
~~~~~~~

:term:`Objects <object>` (contracts) are fundamental application building blocks. Borrowing OOP terminology, an object is a class instance. In other words, an object is a series of :ref:`records <records>` that can be accessed via an index.

Each record represents an object's state at a certain point. The state can contain the object's memory at the point. Memory is a binary blob stored in the ledger and a contract can put any data it needs into it.

In a blockchain, objects cannot be modified, only appended by another record. Therefore, object states can be one of the following types:

* **Activated** -- the :term:`object <object>` has been initialized. This is the first state of any object and it contains initial memory.
* **Amended** -- the object's memory has been modified. Contains new memory. 
* **Deactivated** --  the object has been "removed" from the system. Since data cannot be removed from the chain, objects are simply marked as *removed*.

A succession of object records (states) is called a :term:`lifeline <lifeline>`:

.. uml::

   skinparam backgroundColor transparent
   skinparam object { 
     backgroundColor transparent
   }

   package "[[../glossary.html#term-lifeline Lifeline]]" as Lifeline {
      object Request
      object Activate
      object "Amend 1" as Amend1
      object "Amend 2" as Amend2
      object Deactivate
   }
   object Index

   Amend2 <|-- Deactivate
   Amend1 <|-- Amend2
   Activate <|-- Amend1
   Request <|-- Activate

   Request : key = 1
   Activate : key = 2
   Amend1 : key = 3
   Amend2 : key = 4
   Deactivate : key = 5

   Index : key = 1
   Index : stateKey = 5

   Lifeline -[hidden]r- Index

   Index -l- Request
   Index -l-> Deactivate

An object is assembled from a lifeline via its index. As stated above, index is a collection of pointers to object's records (states, requests, etc.). So, to get an object, all we need is its index. The ledger stores multiple versions of the object's index depending on the pulse.

To preserve consistency, each operation is performed on a particular object's version. To get an object to execute on, a :ref:`virtual node <virtual>` sends an operation request based on which the object's version is calculated. This way, two concurrent operations can be performed on different versions of said object.

Object's lifeline is not the only chain, though. The ledger stores any requests that belong to an object in a :term:`sideline <sideline>`. The general term for all the chains (lines) is a :term:`filament <filament>`. So, a more complex object structure including all filaments is as follows:

.. uml::

   skinparam backgroundColor transparent
   skinparam package { 
     backgroundColor transparent
   }
   skinparam object { 
     backgroundColor transparent
   }

   package "[[../glossary.html#term-lifeline Lifeline]]" as Lifeline {
      object Request
      object Activate
      object "Amend 1" as Amend1
      object "Amend 2" as Amend2
      object Deactivate
   }
   object Index

   Amend2 <|-- Deactivate
   Amend1 <|-- Amend2
   Activate <|-- Amend1
   Request <|-- Activate

   package "[[../glossary.html#term-sideline Requests sideline]]" as rsl {
      object "Request 1" as Req1
      object "Request 2" as Req2
      object "Result 1" as Res1
      object "Request 3" as Req3
   }

   Req1 <|-- Req2
   Req2 <|-- Res1
   Res1 <|-- Req3

   Request : key = 11
   Activate : key = 12
   Amend1 : key = 13
   Amend2 : key = 14
   Deactivate : key = 15

   Req1 : key = 31
   Req2 : key = 32
   Res1 : key = 33
   Req3 : key = 34

   Index : key = 11
   Index : stateKey = 15
   Index : requestKey = 34

   Index -- Request
   Index --> Deactivate
   Index --> Req3
   Lifeline -[hidden]r- rsl
   rsl  -[hidden]r- Index

.. _object_address:

Object's Address
^^^^^^^^^^^^^^^^

Object's address is more complicated than that of a simple :ref:`record <records>`. An :term:`object <object>` consists of many :ref:`records <records>` but should have only one address. So, the ledger considers the address to be a pointer to the creation request's record. The object's index can be found via this address.

.. _relations:

Relations
~~~~~~~~~

Objects have relations to other entities and to each other. Most of those relations are references in the object's :ref:`activation record <objects>`.

Key figures in those relations are:

* **Object**. Directly references a prototype. This reference cannot be changed during the object's lifetime, although multiple objects can have the same prototype. Serves as an *instance* of a prototype.
* **Prototype**. Special kind of :term:`object <object>` that acts as a template for building other objects. It contains default memory and directly refers to relevant code.
* **Code**. Single immutable :ref:`record <records>` which contains code for :ref:`virtual nodes <virtual>` to execute. They perform operations on the referenced object. The same code can be referenced by multiple prototypes.

Relations between the entities are as follows:

.. uml::

   skinparam backgroundColor transparent
   skinparam object { 
     backgroundColor transparent
   }

   object "Code 1" as Code1
   object "Prototype 1 (Object)" as Proto1
   object "Instance 1 (Object)" as Inst1

   object "Code 2" as Code2
   object "Prototype 2 (Object)" as Proto2
   object "Instance 2 (Object)" as Inst2

   object "Instance 3 (Object)" as Inst3

   object "Prototype 3 (Object)" as Proto3

   Code1 <|-- Proto1 : Image
   Proto1 <|-- Inst1 : Image

   Code2 <|-- Proto2 : Image
   Proto2 <|-- Inst2 : Image

   Proto2 <|-- Inst3 : Image
   Code2 <|-- Proto3 : Image

Since both prototype and object are technically :term:`objects <object>`, they contain a reference to either:

* prototype in case of an object, or 
* code in case of a prototype.

The general term for this reference is an *image*. In other words, object's image is its prototype, and prototype's image is its code. 
