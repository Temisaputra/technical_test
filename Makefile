generate-mocks:
	mockgen -source=delivery/repository/supplier_repository.go -destination=./shared/mock/repository/supplier_repository_mock.go -package=repository
	mockgen -source=delivery/repository/transaction_repository.go -destination=./shared/mock/repository/transaction_repository_mock.go -package=repository