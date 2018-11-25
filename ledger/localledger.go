package ledger

import (
	"bytes"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/protocol"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/transactionstore"
	"math"
)

type LocalLedger struct {
	ts        *transactionstore.TransactionStore
}

func NewLocalLedger(txStore *transactionstore.TransactionStore) (*LocalLedger) {
	return &LocalLedger{ts:txStore}
}

func (ld *LocalLedger) Initialize(genesisTx *transaction.Transaction) (error) {
	if !ld.ts.IsEmpty() {
		return ErrLedgerAlreadyInitialized
	}

	return ld.saveTransaction(genesisTx)
}

func (ld *LocalLedger) Register(sendTx *transaction.Transaction, receiveTx *transaction.Transaction) (error) {
	if err := ld.VerifyTransaction(sendTx, true); err != nil {
		return err
	}

	if err := ld.VerifyTransaction(receiveTx, true); err != nil {
		return err
	}

	if err := ld.VerifyTransactions(sendTx, receiveTx); err != nil {
		return err
	}

	return ld.saveTransactions(sendTx, receiveTx)
}

func (ld *LocalLedger) GetLastTransaction(address string) (*transaction.Transaction, error) {
	fromTipTx, err := ld.ts.Retrieve(address)
	if err != nil {
		return nil, err
	}
	return fromTipTx, nil
}

func (ld *LocalLedger) GetTransaction(hash string) (*transaction.Transaction, error) {
	tx, err := ld.ts.Retrieve(hash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (ld *LocalLedger) GetAddressStatement(address string) ([]*transaction.Transaction, error) {
	txChain, err := ld.ts.GetTransactionChain(address, false)
	if err != nil {
		return nil, err
	}
	return txChain, nil
}

func (ld *LocalLedger) VerifyTransaction(tx *transaction.Transaction, isNew bool) (error) {
	if ok, err := ld.verifyAddress(tx); !ok {
		return err
	}
	if ok, err := ld.verifyLinkAddress(tx); !ok {
		return err
	}
	if !ld.verifyPow(tx) {
		return ErrInvalidTransactionHash
	}
	if !tx.VerifySignature() {
		return ErrInvalidTransactionSignature
	}

	localTx, err := ld.ts.Retrieve(string(tx.Hash))
	if err != nil {
		return err
	}
	if localTx != nil && isNew {
		return ErrTransactionAlreadyInLedger
	} else if localTx == nil && !isNew {
		return ErrTransactionNotFound
	}

	if tx.Type != protocol.Transaction_OPEN {
		previous, err := ld.ts.Retrieve(string(tx.Previous))
		if err != nil {
			return err
		}
		if previous == nil {
			return ErrPreviousTransactionNotFound
		}
		if isNew {
			head, err := ld.ts.Retrieve(string(previous.Address))
			if err != nil {
				return err
			}
			if head == nil {
				return ErrHeadTransactionNotFound
			}
			if bytes.Compare(head.Hash, previous.Hash) != 0 {
				return ErrPreviousTransactionIsNotHead
			}
		}
	}

	if tx.Type == protocol.Transaction_SEND {
		amount, err := ld.findBalanceDiffWithPrevious(tx)
		if err != nil {
			return err
		}
		if amount < 0.0 || tx.Balance < 0.0 {
			return ErrNotEnoughFunds
		}
	}

	open, _ := ld.getOpenTransaction(tx)
	if open == nil {
		return ErrOpenTransactionNotFound
	}

	return nil
}

func (ld *LocalLedger) VerifyTransactions(sendTx *transaction.Transaction, receiveTx *transaction.Transaction) (error) {
	if sendTx.Type != protocol.Transaction_SEND {
		return ErrInvalidSendTransaction
	}

	if receiveTx.Type != protocol.Transaction_OPEN && receiveTx.Type != protocol.Transaction_RECEIVE {
		return ErrInvalidReceiveTransaction
	}

	if string(receiveTx.Link) != string(sendTx.Hash) {
		return ErrSendReceiveTransactionsNotLinked
	}

	if string(receiveTx.Address) == string(sendTx.Address) {
		return ErrSendReceiveTransactionsCantBeSameAddress
	}

	pending, err := ld.isPendingTransaction(sendTx)
	if err != nil {
		return err
	}
	if !pending {
		return ErrSendTransactionIsNotPending
	}

	sentAmount, err := ld.findAbsoluteBalanceDiffWithPrevious(sendTx)
	if err != nil {
		return err
	}

	receivedAmount, err := ld.findAbsoluteBalanceDiffWithPrevious(receiveTx)
	if err != nil {
		return err
	}

	if sentAmount != receivedAmount {
		return ErrSentAmountDiffersFromReceivedAmount
	}

	return nil
}

func (ld *LocalLedger) saveTransaction(tx *transaction.Transaction) (error) {
	_, err := ld.ts.Store(tx)
	if err != nil {
		return err
	}
	return nil
}

func (ld *LocalLedger) saveTransactions(sendTx *transaction.Transaction, receiveTx *transaction.Transaction) (error) {
	err := ld.VerifyTransactions(sendTx, receiveTx)
	if err != nil {
		return err
	}

	err = ld.saveTransaction(sendTx)
	if err != nil {
		return err
	}

	err = ld.saveTransaction(receiveTx)
	if err != nil {
		return err
	}

	return nil
}

func (ld *LocalLedger) verifyPow(tx *transaction.Transaction) bool {
	ok, _ := tx.VerifyPow()
	return ok
}

func (ld *LocalLedger) isPendingTransaction(tx *transaction.Transaction) (bool, error) {
	if tx.Type != protocol.Transaction_SEND {
		return false, nil
	}
	target, err := ld.GetLastTransaction(string(tx.Link))
	if err != nil {
		return false, err
	}
	if target == nil {
		return true, nil
	}
	chain, err := ld.ts.GetTransactionChain(string(target.Hash), false)
	if err != nil {
		return false, err
	}
	for _, v := range chain {
		if bytes.Equal(tx.Hash, v.Link) {
			return false, nil
		}
	}
	return true, nil
}

func (ld *LocalLedger) findAbsoluteBalanceDiffWithPrevious(tx *transaction.Transaction) (float64, error) {
	amount, err := ld.findBalanceDiffWithPrevious(tx)
	return math.Abs(amount), err
}

func (ld *LocalLedger) findBalanceDiffWithPrevious(tx *transaction.Transaction) (float64, error) {
	if tx.Type == protocol.Transaction_OPEN {
		return tx.Balance, nil
	}
	previous, err := ld.findPrevious(tx)
	if err != nil {
		return 0, err
	}
	if previous == nil {
		return 0, ErrPreviousTransactionNotFound
	}
	return previous.Balance - tx.Balance, nil
}

func (ld *LocalLedger) findPrevious(tx *transaction.Transaction) (*transaction.Transaction, error) {
	if tx.Type != protocol.Transaction_OPEN {
		return ld.ts.Retrieve(string(tx.Previous))
	}
	return ld.ts.Retrieve(string(tx.Link))
}

func (ld *LocalLedger) verifyAddress(tx *transaction.Transaction) (bool, error) {
	if ok, err := address.IsValid(string(tx.Address)); !ok {
		return ok, err
	}
	if !address.MatchesPubKey(tx.Address, tx.PubKey) {
		return false, ErrAddressDoesNotMatchPubKey
	}
	return true, nil
}

func (ld *LocalLedger) verifyLinkAddress(tx *transaction.Transaction) (bool, error) {
	if tx.Type == protocol.Transaction_SEND {
		return address.IsValid(string(tx.Link))
	}
	return true, nil
}

func (ld *LocalLedger) getOpenTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	var ret = tx
	var err error
	for ret != nil && ret.Type != protocol.Transaction_OPEN {
		ret, err = ld.getPreviousTransaction(ret)
		if err != nil {
			return nil, err
		}
	}
	if ret != nil && ret.Type == protocol.Transaction_OPEN {
		return ret, nil
	}
	return nil, nil
}

func (ld *LocalLedger) getPreviousTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	previous, err := ld.ts.Retrieve(string(tx.Previous))
	if err != nil {
		return nil, err
	}
	if previous == nil {
		return nil, ErrPreviousTransactionNotFound
	}
	return previous, nil
}
