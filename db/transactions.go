package db

import(
	"types"
	"upper.io/db.v3"
)


type Transactions string


const (
	transactionsTable = "transactions"
)


func (*Transactions) AddTransaction(transaction types.Transaction) error{
	_, err := currentInstance.instance.
		InsertInto(transactionsTable).Values(transaction).Exec()
	return err
}


func (*Transactions) GetTransactionById(transactionId int64)(*types.Transaction, error){
	var transaction *types.Transaction
	err := currentInstance.instance.Collection(transactionsTable).Find("id",transactionId).One(&transaction)

	if err != nil {
		return nil, err
	}
	return transaction, nil
}


func (*Transactions) GetOldestTransactions(userId, transactionId int64, numberOfTransactions int,isSeller bool)([]*types.Transaction,error){
	var(
		transactions []*types.Transaction
		err          error
	)
	if isSeller{
		err = currentInstance.instance.Collection(transactionsTable).Find(db.Cond{
			"id < ":     transactionId,
			"shop_id =": userId,
		}).OrderBy("-id").Limit(numberOfTransactions).All(&transactions)
	} else {
		err = currentInstance.instance.Collection(transactionsTable).Find(db.Cond{
			"id < ":     transactionId,
			"user_id =": userId,
		}).OrderBy("-id").Limit(numberOfTransactions).All(&transactions)
	}

	if err != nil {
		return nil, err
	}

	return transactions, nil
}


func (*Transactions) GetNewestTransactions(lastTransactionId, userId int64,isSeller bool) ([]*types.Transaction, error) {
	var (
		transactions []*types.Transaction
		err          error
	)
	if isSeller{
		err = currentInstance.instance.Collection(transactionsTable).Find(db.Cond{
		"id > ":     lastTransactionId,
		"shop_id =": userId,
		}).OrderBy("-id").All(&transactions)
	} else {
		err = currentInstance.instance.Collection(transactionsTable).Find(db.Cond{
		"id > ":     lastTransactionId,
		"user_id =": userId,
		}).OrderBy("-id").All(&transactions)
	}

	if err != nil {
		return nil, err
	}

	return transactions, nil
}


func (*Transactions) GetBiggestId() (int64, error) {
	var transaction *types.Transaction
	err := currentInstance.instance.SelectFrom(transactionsTable).Limit(1).OrderBy("-id").One(&transaction)
	if err != nil {
		return -1, err
	} else {
		return transaction.Id, nil
	}
}