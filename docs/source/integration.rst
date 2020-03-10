.. _integration:

======================
Set Up Insolar Network
======================

To set up an Insolar network:

#. Check the :ref:`hardware requirements <sys_requirements>`.
#. :ref:`Deploy the network locally <deploying_devnet>` for development and test purposes. The local setup is done on one computer, and the "network nodes" are simply services listening to different ports.

.. _sys_requirements:

Hardware Requirements
---------------------

The recommended setup for a proof-of-concept private Insolar network is to consist of **at least 5 nodes** that may be deployed both on virtual or physical servers in a data center.

The minimal hardware requirements for all servers are as follows:

+-------------------------+-------+---------+-------------------+
| Processor               | RAM   | Storage | Network bandwidth |
+=========================+=======+=========+===================+
| 4 cores (8 recommended) | 16 GB | 50 GB   | 1 Gbps            |
+-------------------------+-------+---------+-------------------+

.. note:: The storage capacity may need to be expanded depending on the size of the data to be stored.

Insolar runs on Linux, e.g., **CentOS**.

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

.. _deploying_devnet:

Deploying Network Locally
-------------------------

To set up the network locally, do the following:

#. Since Insolar is written in Go, install the latest 1.12 version of the `Golang programming tools <https://golang.org/doc/install#install>`_.

   .. note:: Make sure the ``$GOPATH`` environment variable is set. 

#. Download the Insolar package:

   .. code-block:: bash

      go get github.com/insolar/insolar

#. Go to the package directory:

   .. code-block:: bash

      cd $GOPATH/src/github.com/insolar/insolar

#. Install dependencies and build binaries: simply run ``make``.

#. Take a look at the ``scripts/insolard/bootstrap_template.yaml`` file. Here, you can find a list of nodes to be launched. In local setup, the "nodes" are simply services listening to different ports.

   To add more nodes to the "network", uncomment some.

#. Run the launcher:

   .. code-block:: bash

      insolar-scripts/insolard/launchnet.sh -g

   The launcher generates bootstrap data, starts the nodes and a pulse watcher, and logs events to ``.artifacts/launchnet/logs``.

When the pulse watcher says ``INSOLAR STATE: READY``, the network is up and has achieved consensus. You can start running test scripts and `benchmarks <https://github.com/insolar/insolar/blob/master/cmd/benchmark/README.md>`_.

.. _logs_and_monitor:

Logging and Monitoring
----------------------

To manually bring up :ref:`logging and monitoring <logs_and_monitor>`, run ``insolar-scripts/monitor.sh``.

To see the node’s logs, open Kibana in a web browser (``http://<your_server_IP>:5601/``) and click :guilabel:`Discover` in the menu.

To see the monitoring dashboard, open ``http://<your_server_IP>:3000/``, log in to Grafana (login: ``admin``, password: ``pass``), click :guilabel:`Home`, and open the :guilabel:`Insolar Dashboard`.
