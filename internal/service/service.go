package service

type Service struct {
	UserService           *UserService
	BalanceHistoryService *BalanceHistoryService
	RefreshTokenService   *RefreshTokenService
	OrderService          *OrderService
	StatusService         *StatusService
}

type Storage struct {
	BalanceHistoryStorage BalanceHistoryStorage
	UserStorage           UserStorage
	OrderStorage          OrderStorage
	TokenStorage          TokenStorage
	StatusStorage         StatusStorage
}

func New(db Storage) *Service {
	return &Service{
		UserService:           NewUserService(db.UserStorage),
		BalanceHistoryService: NewBalanceHistoryService(db.BalanceHistoryStorage),
		RefreshTokenService:   NewRefreshTokenService(db.TokenStorage),
		StatusService:         NewStatusService(db.StatusStorage),
	}
}
