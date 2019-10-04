.. _migration_test:

Testing Token Migration
=======================

This test scenario walks you through the process of acquiring test INS tokens and swapping them to test XNS coins.

During the testing period of the Insolar MainNet, all operations will be performed with test XNS coins (native for MainNet).

.. warning:: Do not send real INS tokens from the main Ethereum network to Insolar MainNet during the testing period.

.. _needs_for_migration_test:

What You Will Need
------------------

You need to create two wallets:

* **Ethereum wallet** that supports Ethereum test networks (Ropsten in particular) and custom tokens. This test scenario uses `MetaMask <https://metamask.io/>`_ as an example.
* **Insolar Wallet** in which you can store XNS coins that are native for the MainNet.

Additionally, you need to understand the following:

* **What coins and tokens are**. Coin is a native cryptocurrency of any blockchain and token is a unit of value that may be built on top of it. Coins operate independently, while tokens have specific uses. In the Ethereum network, the coin is the ether (ETH) and the token standard (upon which the INS token is built) is called ERC-20. Test tokens and coins are those that operate in test networks. For testing purposes, we are going to use the Ropsten test network.
* **What a crypto faucet is.** The faucet is a quick and easy way of creating and distributing test coins and tokens in test networks.
* **What a token contract is.** Essentially, a token contract is a smart contract that contains a map of account addresses and their balances. The unit of the balance is the token. When tokens are transferred from one account to another, the token contract updates the balances of the two accounts.

.. _test_overview:

Test Scenario Overview
----------------------

To grasp what you will be doing, take a look at the steps of the test scenario:

#. Create a MetaMask wallet and, inside, switch to Ropsten test network and add the INS custom token.
#. Get Ropsten test ETH coins.
#. Swap the coins to test INS tokens using the INS token contract.
#. Create an Insolar Wallet and, automatically, receive a migration address.
#. Send the test INS tokens from the MetaMask wallet to the migration address. They will automatically swap to XNS native coins and appear in the Insolar Wallet.

All the above steps are described in detail in subsequent sections.

.. _creating_metamask:

Creating and Setting Up Ethereum Wallet
---------------------------------------

To create and set up a MetaMask Ethereum wallet:

#. Go to `metamask.io <https://metamask.io>`_ and add its extension to your browser.
#. Open the extension, click :guilabel:`Get Started`, :guilabel:`Create a Wallet`, and agree or disagree to help MetaMask improve itself.
#. Create a new password, confirm it, and agree to the "Terms of Use".
#. Reveal the backup phrase, copy it, and click :guilabel:`Next`.
#. Click the words of the backup phrase in correct order and then click :guilabel:`Confirm`.
#. Receive congratulations from MetaMask and click :guilabel:`All Done`. You will see the wallet's user interface.
#. In the top right corner, select :guilabel:`Ropsten Test Network` from the drop-down list:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/selecting-ropsten.png
      :width: 700px

   |

#. In the bottom left corner, click :guilabel:`Add Token`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/add-token.png
      :width: 700px

   |

#. On the **Add Tokens** screen, open the :guilabel:`Custom Token` tab:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/custom-token.png
      :width: 300px

   |

#. Copy the INS token contract address -- click the copy icon |copy-icon| in the right corner of the following code block:

   .. |copy-icon| image:: https://github.com/insolar/doc-pics/raw/master/mig-test/copy-icon.png
      :width: 20px

   .. code-block::

      0x7e94f2be613c6846c40325b0f2712269a0d61d10

#. In the :guilabel:`Token Contract Address` field, paste the copied INS token contract address:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ins-token.png
      :width: 300px

   MetaMask will find the INS token symbol and decimals of precision for you. Click :guilabel:`Next`.

#. On the next screen, click :guilabel:`Add Tokens`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/add-ins.png
      :width: 300px

   With that, the MetaMask wallet is set up to operate the test ETH coins and INS tokens:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/wallet-setup.png
      :width: 700px

   |

.. _acquire_test_tokens_and_swap:

Acquiring Test ETH Coins and Swapping Them to Test INS Tokens
-------------------------------------------------------------

To acquire, first, test ETH, then swap them to test INS tokens:

#. In the MetaMask wallet, first, click the :guilabel:`ETH` tab, then :guilabel:`Deposit`.

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/eth-deposit.png
      :width: 700px

   |

#. In the **Deposit Ether** window, click :guilabel:`Get Ether` next to **Test Faucet**:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/get-eth-from-faucet.png
      :width: 700px

   This opens the `MetaMask Ether Faucet page <https://faucet.metamask.io/>`_.

#. On the opened page, click :guilabel:`request 1 ether from faucet`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/request-one-eth.png
      :width: 400px

   MetaMask will ask you to connect the request in the newly opened window. Click :guilabel:`Connect`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/connect-request.png
      :width: 300px

   Once connected, you can click :guilabel:`request 1 ether from faucet` several times more (maximum 6). The corresponding transaction entries will appear below:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/test-eth-txes.png
      :width: 450px

   Wait several seconds to let the transactions be processed by the test network and return to the MetaMask wallet.

#. In the MetaMask wallet's **History**, the confirmed transactions will appear and your balance will be updated. Click :guilabel:`Send`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/meta-balance.png
      :width: 700px

   |

#. Again, copy the INS token contract address -- click the copy icon |copy-icon| in the right corner of the following code block:

   .. code-block::

      0x7e94f2be613c6846c40325b0f2712269a0d61d10

#. On the **Add Recipient** screen, paste the copied address to the search field:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/send-search-field.png
      :width: 300px

   The MetaMask wallet will recognize the INS token contract and display the transfer details:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/meta-transfer-details.png
      :width: 300px

   |

#. On the **Send ETH** screen, you do not need to send actual ETH to the token contract but a small amount of ETH will be automatically subtracted to pay for transaction processing.
   
   However, setting the correct gas value is required. To set it, click :guilabel:`Advanced Options`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/advanced-options.png
      :width: 300px

   |

#. On the **Customize Gas** screen, set the :guilabel:`Gas Limit` to ``80000`` (eighty thousand) and click :guilabel:`Save`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/gas-limit.png
      :width: 300px

   |

   .. caution:: If the gas limit value is lower than 80,000, the token contract will fail.

#. Back on the **Send ETH** screen, click :guilabel:`Next`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/finally-send-eth.png
      :width: 300px

   And, on the next screen, click :guilabel:`Confirm`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/finally-confirm.png
      :width: 300px

   |

#. Repeat the procedure of sending ETH to INS token contract several more times to acquire enough test INS tokens.

   Once the corresponding transactions are confirmed, the MetaMask wallet is set up to operate test INS tokens:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/meta-wallet-setup.png
      :width: 700px

Next, migrate test INS token to the Insolar network. The migration will automatically swap the test INS tokens to test XNS coins.

.. _migrate_test_tokens:

Migrating Test INS Tokens and Swapping Them to Test XNS Coins
-------------------------------------------------------------

To migrate the test INS tokens and swap them to XNS coins:

#. Create your Insolar Wallet. On the Insolar Wallet main page, click :guilabel:`CREATE A NEW WALLET`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/create-ins-wallet.png
      :width: 600px

   The **Wallet creation tutorial** will open. Read through it attentively.

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

#. In the Insolar Wallet, click the avatar icon |avatar-icon| in the upper right corner to open the menu:

   .. |avatar-icon| image:: https://github.com/insolar/doc-pics/raw/master/mig-test/avatar-icon.png
      :width: 30px

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/right-menu.png
      :width: 200px

   In the menu, click :guilabel:`Copy migration address`.

   After that, return to the MetaMask wallet.

#. In the MetaMask wallet, open the :guilabel:`INS` tab and click :guilabel:`Send`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/meta-send-ins.png
      :width: 700px

   |

#. On the **Add Recipient** screen, paste the copied migration address to the search field:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/send-search-field.png
      :width: 300px

   |

#. On the **Send Tokens** screen, first, click :guilabel:`Max`, then :guilabel:`Next`:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/send-ins-to-mig-addr.png
      :width: 300px

   And :guilabel:`Confirm` the transaction:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/confirm-send-to-mig-addr.png
      :width: 300px

   The migration process may take some time.

#. Once the transaction is processed by the Ropsten test network, your test XNS coins will appear in the Insolar Wallet:

   .. image:: https://github.com/insolar/doc-pics/raw/master/mig-test/ins-tokens-hold.png
      :width: 300px

This concludes the migration test.
