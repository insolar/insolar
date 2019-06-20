.. _integration:

========================
Integrating with Insolar
========================

.. note::

   Upon the public network release, this chapter will be expanded to include the appropriate integration instructions.

Join or set up the Insolar network:

* :ref:`Connect to TestNet 1.1 <connecting_to_testnet>`. Participation in this network is permissioned, with participants invited by the Insolar Core Development Team based on their ability to fulfill the respective SLA.
* :ref:`Set up a network <setting_up_devnet>` locally for development and test purposes. The local setup is done on one computer with no particular system requirements, and the 'network nodes' are simply services listening on different ports.

.. _connecting_to_testnet:

Connecting to Test Network
--------------------------

To connect to Insolar TestNet 1.1:

#. Skim through :ref:`known issues and limitations <issues_and_limitations>` and :ref:`system requirements <sys_requirements>`.
#. :ref:`Set up and connect <connect_to_testnet>` a node.

.. _issues_and_limitations:

Known Issues and Limitations
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. note:: Issues below will be addressed in future releases.

On TestNet 1.1:

**Node Maintenance**

* Only computational (virtual) nodes are available to external participants. Data storage is provided by Insolar nodes.
* All discovery nodes are hosted by Insolar. Other nodes use them to reconnect to the network.
* In the unlikely event of short-term storage (light material) nodes having to reconnect, multiple errors may occur for a few pulses.
* Only one long-term (heavy material) node and one pulsar are deployed. If either of the nodes is missing, the network will go down but, upon the node's restart, will recover.
* Storage node crash may lead to data loss.
* Under certain conditions, the node’s process (``insolard`` daemon) may exit and its Docker container will restart it automatically. Insolar may also ask node holders for assistance with a manual restart.
* Nodes joining the network produce errors when other nodes are leaving the network.

**Security & data consistency**

* Smart contract validation is disabled. Therefore, any execution result returned by a virtual node is treated as verified.
* Distributed transactions are not yet implemented. This can lead to decorrelated object changes. For example, an interrupted 'money' transfer from one wallet to another decreases the source wallet's balance but leaves the target wallet's balance unchanged.
* Operations executed during pulse changes will be declined with the ``Incorrect message pulse`` error.
* Currently, all the data is stored and transferred unencrypted.
* All network messages are signed but signature checks are disabled.

**Performance**

* A simplistic rate limiter is implemented for light material nodes, so they reject incoming requests when the number of pending requests reaches a certain limit. This results in an exponential backoff on our benchmark tool -- the retry interval increases exponentially.
* The rate limiter does not consider the request's origin. So, when a user puts an excessive load on the network, other users may suffer.

**Application level**

* Only pre-built smart contracts are available. Custom contracts will be available on TestNet 2.0.
* All user wallets are created with a starting balance of 1,000,000 coins.
* Only one contract can be called via the node's API. All other methods (e.g., coin transfer) are called via the ``Call`` method on a ``member`` object.

.. _sys_requirements:

Test Network System Requirements
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Preferable system requirements for virtual nodes on TestNet 1.1 are as follows:

+-----------+-------+------------+
| Processor | RAM   | Storage    |
+===========+=======+============+
| 40 cores  | 32 GB | 256 GB SSD |
+-----------+-------+------------+

Insolar is meant to run on Linux, e.g., **CentOS**.

All servers wishing to join the Insolar test network must have **public IP addresses**.

.. _connect_to_testnet:

Connecting to TestNet 1.1
~~~~~~~~~~~~~~~~~~~~~~~~~

To connect to the Insolar test network, do the following:

#. Install `Docker and Docker Compose <https://docs.docker.com/v17.12/install/>`_ and run the Docker daemon.

#. Download the Insolar's ``insolar-node-<version>.tar.gz`` archive from the `latest release <https://github.com/insolar/insolar/releases>`_. You can find it under the :guilabel:`Assets` drop-down list.

#. Unpack the archive on your server. A good place is under the ``/opt/insolar`` directory.

#. Go to the unpacked directory, open the ``docker-compose.yml`` file in a text editor, and insert your server's public IP address to the ``INSOLARD_TRANSPORT_FIXED_ADDRESS`` field.

#. Acquire ``cert.json`` and ``keys.json`` files from Insolar. You can ask for them in our `Telegram developer's chat <https://t.me/InsolarTech>`_.

   Put the files to the ``configs`` directory.

#. Run ``docker-compose up -d``.

Enjoy being a part of the Insolar Network!

.. note:: The Insolar's API is under development and not yet finalized. Please, await its first release.

In addition to the Insolar node, the Docker Compose starts Kibana and Grafana services to take care of :ref:`logging and monitoring <logs_and_monitor>`.

.. _ports_used:

Ports Used
~~~~~~~~~~

Insolar uses the following ports:

+--------------+----------+-----------------------------------------------------+
| Port         | Protocol | Description                                         |
+--------------+----------+-----------------------------------------------------+
| 7900, 7901   | TCP, UDP | Nodes intercommunication.                           |
|              |          | The node must be publicly available on these ports. |
+--------------+----------+-----------------------------------------------------+
| 8090         | TCP      | Node-pulsar communication.                          |
|              |          | The node must be publicly available on this port.   |
+--------------+----------+-----------------------------------------------------+
| 18181, 18182 | TCP      | Communication between the main node daemon and the  |
|              |          | smart contract executor daemon.                     |
+--------------+----------+-----------------------------------------------------+
| 19191        | TCP      | Node's JSON-RPC API.                                |
+--------------+----------+-----------------------------------------------------+
| 8080         | TCP      | Prometheus metrics endpoint.                        |
+--------------+----------+-----------------------------------------------------+

.. _setting_up_devnet:

Setting Up Network Locally
--------------------------

To set up the network locally, do the following:

#. Since Insolar is written in Go, install its `programming tools <https://golang.org/doc/install#install>`_.

   .. note:: Make sure the ``$GOPATH`` environment variable is set. 

#. Download the Insolar package:

   .. code-block:: bash

      go get github.com/insolar/insolar

#. Go to the package directory:

   .. code-block:: bash

      cd $GOPATH/src/github.com/insolar/insolar

#. Install dependencies and build binaries: simply run ``make``.

#. Take a look at the ``scripts/insolard/bootstrap_template.yaml`` file. Here, you can find a list of nodes to be launched. In local setup, the 'nodes' are simply services listening on different ports.

   To add more nodes to the 'network', uncomment some.

#. Run the launcher:

   .. code-block:: bash

      scripts/insolard/launchnet.sh -g

   The launcher generates bootstrap data, starts the nodes and a pulse watcher, and logs events to ``.artifacts/launchnet/logs``.

When the pulse watcher says ``INSOLAR STATE: READY``, the network is up and has achieved consensus. You can start running test scripts and `benchmarks <https://github.com/insolar/insolar/blob/master/cmd/benchmark/README.md>`_.

Also, you can manually bring up :ref:`logging and monitoring <logs_and_monitor>` by running ``scripts/monitor.sh``.

.. _logs_and_monitor:

Logging and Monitoring
----------------------

To see the node’s logs, open Kibana in a web browser (``http://<your_server_IP>:5601/``) and click :guilabel:`Discover` in the menu.

To see the monitoring dashboard, open ``http://<your_server_IP>:3000/``, log in to Grafana (login: ``admin``, password: ``pass``), click :guilabel:`Home`, and open the :guilabel:`Insolar Dashboard`.
