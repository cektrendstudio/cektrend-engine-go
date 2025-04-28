generate-mocks:
	# repositories
	@mockgen -destination=./service/repository/mocks/mock_user_repository.go -package=mocks github.com/cektrendstudio/cektrend-engine-go/service UserRepository
	@mockgen -destination=./service/repository/mocks/mock_auth_repository.go -package=mocks github.com/cektrendstudio/cektrend-engine-go/service AuthRepository