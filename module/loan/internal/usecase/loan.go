package usecase

import (
	"context"
	"fmt"
	"time"

	"example.com/loan/internal/util"
	"example.com/loan/module/loan/entity"
	"example.com/loan/module/loan/internal/repository"
	"example.com/loan/module/payment/client"
)

var (
	RepaymentReferencePrefix = "REPAY"
)

type LoanUsecase interface {
	GetOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error)
	PayBilling(ctx context.Context, req entity.PayBillingRequest) (entity.LoanBilling, error)
	Create(ctx context.Context, loan *entity.Loan) error
	ForceDisburse(ctx context.Context, loanID string) error
}

type loanUsecase struct {
	paymentService  client.PaymentService
	loanBillingRepo repository.LoanBillingRepository
	repaymentRepo   repository.RepaymentRepository
	loanRepo        repository.LoanRepository
	clock           util.Clock
}

func (u *loanUsecase) Create(ctx context.Context, loan *entity.Loan) error {
	return u.loanRepo.Create(ctx, loan)
}

func (u *loanUsecase) GetOutstandingLoans(ctx context.Context, loanID string) (entity.GetOutstandingLoansResponse, error) {
	response, err := u.loanBillingRepo.SumOutstandingLoans(ctx, loanID)
	if err != nil {
		return entity.GetOutstandingLoansResponse{}, err
	}

	return response, nil
}

func (u *loanUsecase) PayBilling(ctx context.Context, req entity.PayBillingRequest) (entity.LoanBilling, error) {
	billing, err := u.loanBillingRepo.GetByID(ctx, req.BillingID)
	if err != nil {
		return entity.LoanBilling{}, err
	}

	if billing.Status == entity.LoanBillingStatusPaid {
		return billing, nil
	}

	repayment := entity.Repayment{
		LoanBillingID: billing.ID,
		Amount:        req.Amount,
		Status:        entity.RepaymentStatusCreated,
	}
	if err := u.repaymentRepo.Create(ctx, &repayment); err != nil {
		return entity.LoanBilling{}, err
	}

	_, err = u.paymentService.CreateAndPayMutation(ctx, client.CreateAndPayMutationRequest{
		UserID:    req.UserID,
		Amount:    req.Amount,
		Reference: fmt.Sprintf("%s-%s", RepaymentReferencePrefix, repayment.ID),
	})
	if err != nil {
		return entity.LoanBilling{}, err
	}

	billing.Status = entity.LoanBillingStatusPaid
	if err := u.loanBillingRepo.Update(ctx, billing); err != nil {
		return entity.LoanBilling{}, err
	}

	repayment.Reference = fmt.Sprintf("%s-%s", RepaymentReferencePrefix, repayment.ID)
	if err := u.repaymentRepo.Update(ctx, repayment); err != nil {
		return entity.LoanBilling{}, err
	}

	return billing, nil
}

func (u *loanUsecase) ForceDisburse(ctx context.Context, loanID string) error {
	loan, err := u.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return err
	}

	loan.Status = entity.LoanStatusDisbursed
	if err := u.loanRepo.Update(ctx, loan); err != nil {
		return err
	}

	gmt7 := time.FixedZone("GMT+7", 7*60*60)
	now := u.clock.Now().In(gmt7)
	var billings []entity.LoanBilling
	for i := 0; i < loan.Term; i++ {
		billings = append(billings, entity.LoanBilling{
			LoanID:  loan.ID,
			Amount:  loan.WeeklyInstallment,
			Status:  entity.LoanBillingStatusCreated,
			DueDate: now.AddDate(0, 0, 7*(i+1)),
		})
	}

	return u.loanBillingRepo.BulkCreate(ctx, billings)
}

func NewLoanUsecase(
	paymentService client.PaymentService,
	loanBillingRepo repository.LoanBillingRepository,
	repaymentRepo repository.RepaymentRepository,
	loanRepo repository.LoanRepository,
	clock util.Clock,
) LoanUsecase {
	if clock == nil {
		clock = util.RealClock{}
	}

	return &loanUsecase{
		paymentService:  paymentService,
		loanBillingRepo: loanBillingRepo,
		repaymentRepo:   repaymentRepo,
		loanRepo:        loanRepo,
		clock:           clock,
	}
}
