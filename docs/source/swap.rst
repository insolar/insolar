.. _swap:

Swap INS for XNS
================

To swap your INS tokens to XNS coins and simultaneously migrate them from the Ethereum network to Insolar MainNet, go through the following steps:

#. Open the `Insolar Wallet <https://wallet.insolar.io>`_ website and make sure to select :guilabel:`MAINNET` from the drop-down list.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/select-mainnet.png
      :width: 600px

#. Click :guilabel:`CREATE A NEW WALLET`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/create-ins-wallet.png
      :width: 600px

   This opens a **Wallet creation tutorial**. Read through it attentively.

   Upon creation, the Wallet takes care of security for you:

   #. Generates a backup phrase and private key using randomization. They are synonymous in function.
   #. Encrypts the key with your password and puts it in a keystore file. You can use this file to access your wallet and authorize operations.
   #. Ensures that you make a record of the backup phrase. Using this phrase, you can restore the Wallet in case you lose the private key or the keystore file and your password.

   .. caution:: You are solely responsible for keeping your funds as no one else can recover your Wallet. Insolar does not store your credentials, encrypted or otherwise.

#. On the **Create a new Wallet** page:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ins-wallet-password.png
      :width: 370px

   #. Enter a new password. It should be at least 8 characters long and contain a mix of numbers, uppercase, and lowercase letters.
   #. Re-enter the password to confirm it.
   #. Agree to the "Term of Use".
   #. Allow anonymous data collection to improve the service.
   #. Click :guilabel:`NEXT`.

#. On the next screen, click :guilabel:`REVEAL TEXT` to see the backup phrase:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ins-reveal-phrase.png
      :width: 450px

   The secret backup phrase is a series of words that store all the information needed to recover Insolar Wallet. The backup phrase and private key are synonymous in function.

   .. warning:: Never disclose your backup phrase (or private key).

   .. tip::

      Security tips:

      * Store the backup phrase in a password manager.
      * Write the phrase down on several pieces of paper and store them in different locations.
      * Memorize the phrase.

   Once you have secured the backup phrase, click :guilabel:`NEXT`.

#. On the next screen, enter the requested words in the correct order and click :guilabel:`OPEN MY WALLET`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ins-word-order.png
      :width: 350px

#. Wait for the Wallet validation to complete and all features to become available:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/one-more-thing.png
      :width: 400px

#. Once the Wallet is created, receive congratulations from Insolar:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ins-congrats.png
      :width: 400px

   And save the keystore file in one of the following ways:

   * Click :guilabel:`SAVE LOCALLY` to save it to your browser’s local storage. Keeping the file locally allows easier access from the browser on the device you are using.
   * Click :guilabel:`DOWNLOAD` to save it to your computer. In this case, you can move it to another device via, for example, a USB drive.

   Later, you can log in using one of the following:

   * (Recommended) Your password and the keystore file.
   * Unencrypted private key.

   Either way, the Wallet does not store the private key. Instead, it uses the private key provided every time to authorize login and operations. While logged in, you can copy your unencrypted private key, but keep in mind, this is its most vulnerable form.

#. In the Insolar Wallet, open the :guilabel:`SWAP` tab and copy your migration address.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/wlt-open-swap-tab.png
      :width: 600px

   This is a special address in the Ethereum network. Insolar monitors INS tokens sent to it and automatically migrates and swaps them to XNS coins in the Insolar network.

#. Open your ERC-20 Ethereum wallet where you hold your XNS, for example, MetaMask:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/open-eth-wallet.png
      :width: 300px

   Make sure to select :guilabel:`Main Ethereum Network` and that you have some ETH for the transaction fee.

#. In the Ethereum wallet, select INS tokens and click :guilabel:`SEND`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/eth-wlt-send-ins.png?get
      :width: 300px

#. Paste the migration address to the :guilabel:`Add Recipient` field, enter the INS amount, select the transaction fee (in ETH), and click :guilabel:`NEXT`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ins-transfer-details.png
      :width: 300px

#. Confirm the transaction details:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/confirm-eth-tx.png
      :width: 300px

#. Wait for the transaction to go through in the Ethereum network. Optionally, check the transaction status at `Etherscan <https://etherscan.io>`_ — click the arrow button to view the transaction:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/view-on-ethscan.png
      :width: 300px

   It usually takes 20 processed blocks to confirm the transaction.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/eth-scan-tx.png
      :width: 600px

#. Go back to the :guilabel:`SWAP` tab in your Insolar Wallet.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/swap-and-release.png
      :width: 600px

Congratulations! You swapped your INS tokens to XNS coins and they are now stored in your Insolar Wallet.
