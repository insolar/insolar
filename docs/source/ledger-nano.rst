.. _ledger-nano-user-guide:

Using the Insolar Application on Ledger Nano S
===============================================

Ledger Nano S is a hardware wallet that lets you securely store your crypto assets.

Insolar application stores your Insolar private key on the Ledger Nano S and lets you securely manage XNS coins via the Insolar Web Wallet.

The application is developed by Insolar team. To request technical support, send an email to support@insolar.io.

Prerequisites
-------------

To use the Insolar application on your Ledger Nano S hardware wallet, make sure to install and set up the following:

#. `Install Google Chrome <https://www.google.com/chrome/>`_.
#. `Install the Ledger Live application <https://support.ledger.com/hc/en-us/articles/360006395553/>`_.
#. `Set up your Ledger Nano S <https://support.ledger.com/hc/en-us/articles/360000613793>`_.

   .. important:: During the setup, you are required to choose a PIN and write down the recovery phrase. Without the PIN you will not be able to unlock the device and without the recovery phrase—restore access to your wallet. The device recovery phrase differs from the backup phrase of the Insolar Web Wallet, so you cannot use the recovery phrase to log in to the Web Wallet—only the device itself.
   
   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ledger-s-pin.png
      :width: 400px

#. `Install the latest firmware on your Ledger Nano S <https://support.ledger.com/hc/en-us/articles/360002731113-Update-Ledger-Nano-S-firmware>`_.

Install the Insolar Application on Ledger Nano S
------------------------------------------------

To install the Insolar application on Ledger Nano S:

#. Open the :guilabel:`Manager` tab in Ledger Live.
#. Connect and unlock you Ledger Nano S.
#. If prompted, press both buttons simultaneously on the device to allow the manager connection.
#. Find :guilabel:`Insolar` in the application catalog and click :guilabel:`Install` next to it.

   This displays the installation window with a progress bar. Wait for the installation to complete.
#. In the dashboard of the Ledger Nano S device, press :guilabel:`left` or :guilabel:`right` buttons to find the Insolar application.
#. Once found, press both :guilabel:`left` and :guilabel:`right` buttons simultaneously to launch the application.

Once the Insolar application is launched, proceed to creating an Insolar Wallet on Ledger Nano S.

Create a Connected Insolar Wallet
-----------------------------------

To create an Insolar Wallet using the Insolar application on Ledger Nano S, go through the following steps:

#. In Google Chrome, open the `Insolar Web Wallet <https://wallet.insolar.io>`_ and click :guilabel:`CREATE A NEW WALLET`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/create-ins-wlt.png
      :width: 400px

#. On the **Create a new Wallet** screen, click :guilabel:`USE LEDGER NANO S`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/use-ledger-n.png
      :width: 400px

#. Make sure your Ledger Nano S is connected, unlocked, and the Insolar application is launched on it.

   .. _enter_key_number:

#. If required, enter the key number. Ledger Nano S can store multiple private keys—each to an individual Insolar MainNet Wallet. Every key stored in the device has a number. By default, the number of the first key is ``0``.

   .. important:: Remember the number of this private key. You are required to specify it upon every login to use a particular Insolar MainNet Wallet.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/key-number.png
      :width: 500px

#. Check the boxes to allow anonymous data collection and agree to the terms of use and click :guilabel:`CONNECT TO LEDGER NANO S`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/connect-n.png
      :width: 500px

#. In the browser's prompt window, select the :guilabel:`Nano S` device and click :guilabel:`Connect`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/select-n.png
      :width: 400px

#. In the dashboard of the Ledger Nano S device, the Insolar application prompts you to confirm the :guilabel:`Create Account` command.
   
   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ledger-s-create-account.png
      :width: 400px

   Press both :guilabel:`left` and :guilabel:`right` to open the signing options and press both :guilabel:`left` and :guilabel:`right` buttons again to sign the command.
      
   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ledger-s-create-account-sign.png
      :width: 400px

   This securely stores the private key on the device.

#. Once signed, the Insolar Web Wallet displays a wallet validation window.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/one-more-thing.png
      :width: 400px

#. Wait for the validation to complete and see the congratulations message.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ledger-n-congrats.png
      :width: 400px

Once the Wallet is created, you can manage your XNS with it. Every login and XNS transfer operation requires the associated private key stored on the Ledger Nano S, so the device must be connected to confirm these actions.

Log In the Connected Wallet and View Balance
--------------------------------------------

To log in the Insolar Wallet connected to your Ledger Nano S, go through the following steps:

#. In Google Chrome, open the `Insolar Web Wallet <https://wallet.insolar.io>`_ and click :guilabel:`LOG IN`.
#. In the **Log in** panel, click the :guilabel:`Hardware` tab.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/login-hw.png
      :width: 400px

#. Make sure your Ledger Nano S is connected, unlocked, and the Insolar application is launched on it.
#. Specify the key number you chose upon :ref:`wallet creation <enter_key_number>` and click :guilabel:`CONNECT TO LEDGER NANO S`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/enter-key-number.png
      :width: 400px

Insolar Web Wallet recognizes the launched application on the device and automatically logs into the wallet. Once logged in, you can see your balance on the :guilabel:`Dashboard` tab.

Receive XNS
-----------

To receive XNS, do the following:

#. Open the dashboard of the Insolar Web Wallet and click the avatar icon the in upper-right corner.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/click-avatar.png
      :width: 250px

#. In the **Your Wallet** panel, click :guilabel:`Copy XNS address`. This copies the address to the clipboard.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/copy-xns-address.png
      :width: 200px

#. Reveal the address to anyone who wishes to transfer XNS to you and wait for the incoming transaction.
#. View the incoming transactions: in the **Your Wallet** panel, click :guilabel:`Transaction history`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/click-history.png
      :width: 200px

#. On the **Transaction history** screen, open the :guilabel:`RECEIVED` tab.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/click-received.png
      :width: 400px

Once you receive the XNS, the balance on the :guilabel:`Dashboard` tab increases.

Send XNS
--------

To send XNS, do the following:

#. Open the :guilabel:`Dashboard` tab in the Insolar Web Wallet and click :guilabel:`SEND`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/click-send.png
      :width: 150px

#. On the **Send XNS** screen, fill in the recipient address, amount of XNS to send, and click :guilabel:`NEXT`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/send-xns.png
      :width: 400px

#. Make sure your Ledger Nano S is connected, unlocked, and the Insolar application is launched on it.
#. On the **Send XNS** screen, check the following transaction details and click :guilabel:`SEND`:

   * recipient address,
   * amount of XNS to send,
   * transaction fee,
   * total amount — including the fee.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/check-details.png
      :width: 400px

#. In the dashboard of the Ledger Nano S device, the application prompts you to verify the transfer details and sign the :guilabel:`Send XNS` command. Click the :guilabel:`right` button to cycle through the details and check that they are the same as in the web wallet.
  
#. Press both :guilabel:`left` and :guilabel:`right` buttons to sign the :guilabel:`Send XNS` command.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ledger-s-create-account-sign.png
      :width: 400px

#. View the outgoing transactions: in the **Your Wallet** panel, click :guilabel:`Transaction history`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/click-history.png
      :width: 200px

#. On the **Transaction history** screen, open the :guilabel:`SENT` tab.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/click-sent.png
      :width: 400px

Once you send the XNS, the balance in the :guilabel:`Dashboard` tab decreases.

Transfer Swapped XNS from Deposit to Your Main Account
------------------------------------------------------

Once you've `swapped your INS into XNS <./swap.html>`_ your XNS are stored in your Insolar Web Wallet—on a deposit account. Each swap operation creates a separate deposit account that goes from the status :guilabel:`ON HOLD` to :guilabel:`RELEASED` upon a successful swap.

You can transfer your released XNS from deposit to your main account to perform further operations on them. 

#. In the Insolar Wallet, open the SWAP tab, choose the deposit account, and click :guilabel:`TRANSFER`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/transfer-xns-deposit-to-main-account.png
      :width: 600px

#. On the next screen, choose the amount of XNS you want to transfer or click :guilabel:`Use all` to transfer all XNS from this deposit account. Click :guilabel:`TRANSFER` again. 

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/transfer-xns-deposit-to-main-account-use-all.png
      :width: 600px

Insolar Web Wallet asks you to follow instructions on your Ledger Nano S device.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/transfer-xns-deposit-to-main-nano-s.png
      :width: 600px      

#. In the dashboard of the device, the Insolar application prompts you to verify the transfer details and sign the :guilabel:`Transfer` command. Click the right button to cycle through the details.

   Press both :guilabel:`left` and :guilabel:`right` to open the signing options and press both :guilabel:`left` and :guilabel:`right` buttons again to sign the command.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ledger-s-create-account-sign.png
      :width: 400px

#. Insolar Web Wallet shows you a :guilabel:`Transfer initiated` popup message.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/transfer-xns-deposit-to-main-success.png
      :width: 600px

#. View the incoming transactions: in the **Your Wallet** panel, click :guilabel:`Transaction history`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/transfer-xns-deposit-to-main-transaction-history.png
      :width: 800px

Once the transfer operation finishes, the balance in the :guilabel:`Dashboard` tab increases.
