package usecase

import (
	"context"
	"fmt"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/delivery/repository"
)

type GroupUsecase struct {
	groupRepo       repository.GroupRepository
	transactionRepo repository.TransactionRepository
}

func NewGroupUsecase(groupRepository repository.GroupRepository, transactionRepository repository.TransactionRepository) *GroupUsecase {
	return &GroupUsecase{
		groupRepo:       groupRepository,
		transactionRepo: transactionRepository,
	}
}

func (u *GroupUsecase) GetAllGroup(ctx context.Context, pagination *request.Pagination) (groups []*presenter.GroupResponse, meta response.Meta, err error) {
	if pagination.Page == 0 {
		pagination.Page = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	res, meta, err := u.groupRepo.GetAllGroups(ctx, pagination)
	if err != nil {
		errMsg := fmt.Errorf("[GroupUsecase-GetAllGroup] Error when getting all group: %w", err)
		return nil, meta, errMsg
	}

	return res, meta, nil
}

func (u *GroupUsecase) GetGroupByID(ctx context.Context, id int) (res *presenter.GroupResponse, err error) {
	res, err = u.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		errMsg := fmt.Errorf("[GroupUsecase-GetGroupByID] Error when getting group by id: %w", err)
		return nil, errMsg
	}

	return res, nil
}

func (u *GroupUsecase) CreateGroup(ctx context.Context, params *presenter.GroupRequest) (err error) {
	// Mulai transaksi opsional
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Semua repo harus pakai txCtx
		if err := u.groupRepo.CreateGroups(txCtx, *params); err != nil {
			return err // rollback otomatis
		}

		// if err := u.orderRepo.CreateOrder(txCtx, order); err != nil {
		// 	return err // rollback otomatis
		// }

		return nil // commit otomatis
	})
}

func (u *GroupUsecase) UpdateGroup(ctx context.Context, params *presenter.GroupRequest, id int) (err error) {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Ambil group dulu, pakai txCtx
		group, err := u.groupRepo.GetGroupByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("[GroupUsecase-UpdateGroup] Error when getting group by id: %w", err)
		}

		if group == nil {
			return fmt.Errorf("[GroupUsecase-UpdateGroup] Group not found")
		}

		// Update group pakai txCtx
		if err := u.groupRepo.UpdateGroups(txCtx, *params, id); err != nil {
			return fmt.Errorf("[GroupUsecase-UpdateGroup] Error when updating group: %w", err)
		}

		return nil // commit otomatis kalau tidak ada error
	})
}

func (u *GroupUsecase) DeleteGroup(ctx context.Context, id int) (err error) {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Ambil group dulu, pakai txCtx
		group, err := u.groupRepo.GetGroupByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("[GroupUsecase-DeleteGroup] Error when getting group by id: %w", err)
		}

		if group == nil {
			return fmt.Errorf("[GroupUsecase-DeleteGroup] Group not found")
		}

		// Delete group pakai txCtx
		if err := u.groupRepo.DeleteGroup(txCtx, id); err != nil {
			return fmt.Errorf("[GroupUsecase-DeleteGroup] Error when deleting group: %w", err)
		}

		return nil // commit otomatis kalau tidak ada error
	})
}
