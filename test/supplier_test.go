package module

import (
	"context"
	"errors"
	"testing"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/Temisaputra/warOnk/internal/usecase"
	mockRepo "github.com/Temisaputra/warOnk/shared/mock/repository"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestSupplierUsecase(t *testing.T) {
	convey.Convey("SupplierUsecase", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSupplierRepo := mockRepo.NewMockSupplierRepository(ctrl)
		mockAddressRepo := mockRepo.NewMockAddressRepository(ctrl)
		mockContactRepo := mockRepo.NewMockContactRepository(ctrl)
		mockGroupRepo := mockRepo.NewMockGroupRepository(ctrl)
		mockMaterialRepo := mockRepo.NewMockMaterialRepository(ctrl)
		mockTransactionRepo := mockRepo.NewMockTransactionRepository(ctrl)

		supplierUC := usecase.NewSupplierUsecase(
			mockSupplierRepo,
			mockAddressRepo,
			mockContactRepo,
			mockGroupRepo,
			mockMaterialRepo,
			mockTransactionRepo,
		)

		ctx := context.Background()

		supplierID := 1
		supplier := &presenter.SupplierResponse{
			ID:     supplierID,
			Name:   "Test Supplier",
			Status: "Active",
		}
		address := &entity.Address{
			ID:         1,
			Name:       "Main Address",
			SupplierID: uint(supplierID),
		}
		contact := &entity.Contact{
			ID:         1,
			Name:       "Main Contact",
			SupplierID: uint(supplierID),
		}
		group := &entity.Group{
			ID:   1,
			Name: "Supplier Group",
		}
		materials := &[]entity.MaterialList{
			{
				ID:            1,
				MaterialGroup: "Steel",
				MaterialID:    1,
				SupplierID:    uint(supplierID),
				IsActive:      true,
			},
		}

		convey.Convey("GetSupplierByID Success", func() {
			mockSupplierRepo.EXPECT().GetSupplierByID(ctx, supplierID).Return(supplier, nil)
			mockAddressRepo.EXPECT().GetAddressByIdSupplier(ctx, supplierID).Return(address, nil)
			mockContactRepo.EXPECT().GetContactByIdSupplier(ctx, supplierID).Return(contact, nil)
			mockGroupRepo.EXPECT().GetGroupByIdSupplier(ctx, supplierID).Return(group, nil)
			mockMaterialRepo.EXPECT().GetMaterialByIdSupplier(ctx, supplierID).Return(materials, nil)

			result, err := supplierUC.GetSupplierByID(ctx, supplierID)

			convey.So(err, convey.ShouldBeNil)
			convey.So(result, convey.ShouldNotBeNil)
			convey.So(result.Name, convey.ShouldEqual, supplier.Name)
			convey.So(result.Address.Name, convey.ShouldEqual, address.Name)
			convey.So(result.Contact.Name, convey.ShouldEqual, contact.Name)
			convey.So(result.Group.GroupName, convey.ShouldEqual, group.Name)
			convey.So(result.Material[0].MaterialGroup, convey.ShouldEqual, (*materials)[0].MaterialGroup)
		})

		convey.Convey("CreateSupplier Success", func() {
			input := &presenter.SupplierCreateRequest{
				Name:     "Test Supplier",
				NickName: "TS",
				Status:   "Active",
				Address: []presenter.AddressRequest{
					{Name: "Main Address", Address: "Street 123", Main: true},
				},
				Contact: []presenter.ContactRequest{
					{Name: "Main Contact", Phone: "08123456789", Main: true},
				},
				Groups: []uint{1},
				Materials: []presenter.MaterialRequest{
					{MaterialID: 1, MaterialGroup: "Steel"},
				},
			}

			mockTransactionRepo.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
				func(ctx context.Context, f func(context.Context) error) error {
					txCtx := ctx

					mockSupplierRepo.EXPECT().CreateSupplier(txCtx, gomock.Any()).Return(1, nil)
					mockAddressRepo.EXPECT().CreateAddress(txCtx, gomock.Any()).Return(nil)
					mockContactRepo.EXPECT().CreateContact(txCtx, gomock.Any()).Return(nil)
					mockGroupRepo.EXPECT().CreateGroup(txCtx, gomock.Any()).Return(nil)
					mockMaterialRepo.EXPECT().CreateMaterial(txCtx, gomock.Any()).Return(nil)

					return f(txCtx)
				})

			err := supplierUC.CreateSupplier(ctx, input)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("UpdateSupplier Success", func() {
			input := &presenter.SupplierCreateRequest{
				Name:     "Updated Supplier",
				NickName: "US",
				Status:   "Active",
				Address: []presenter.AddressRequest{
					{Name: "New Address", Address: "New Street", Main: true},
				},
				Contact: []presenter.ContactRequest{
					{Name: "New Contact", Phone: "0888888888", Main: true},
				},
				Groups: []uint{2},
				Materials: []presenter.MaterialRequest{
					{MaterialID: 2, MaterialGroup: "Aluminum"},
				},
			}

			mockTransactionRepo.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
				func(ctx context.Context, f func(context.Context) error) error {
					txCtx := ctx

					mockSupplierRepo.EXPECT().GetSupplierByID(txCtx, supplierID).Return(supplier, nil)
					mockSupplierRepo.EXPECT().UpdateSupplier(txCtx, gomock.Any(), supplierID).Return(nil)
					mockAddressRepo.EXPECT().DeleteAddressByIdSupplier(txCtx, supplierID).Return(nil)
					mockAddressRepo.EXPECT().CreateAddress(txCtx, gomock.Any()).Return(nil)
					mockContactRepo.EXPECT().DeleteContactByIdSupplier(txCtx, supplierID).Return(nil)
					mockContactRepo.EXPECT().CreateContact(txCtx, gomock.Any()).Return(nil)
					mockGroupRepo.EXPECT().DeleteGroupByIdSupplier(txCtx, supplierID).Return(nil)
					mockGroupRepo.EXPECT().CreateGroup(txCtx, gomock.Any()).Return(nil)
					mockMaterialRepo.EXPECT().DeleteMaterialByIdSupplier(txCtx, supplierID).Return(nil)
					mockMaterialRepo.EXPECT().CreateMaterial(txCtx, gomock.Any()).Return(nil)

					return f(txCtx)
				})

			err := supplierUC.UpdateSupplier(ctx, input, supplierID)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("DeleteSupplier Success", func() {
			mockTransactionRepo.EXPECT().WithTransaction(ctx, gomock.Any()).DoAndReturn(
				func(ctx context.Context, f func(context.Context) error) error {
					txCtx := ctx

					mockSupplierRepo.EXPECT().GetSupplierByID(txCtx, supplierID).Return(supplier, nil)
					mockAddressRepo.EXPECT().DeleteAddressByIdSupplier(txCtx, supplierID).Return(nil)
					mockContactRepo.EXPECT().DeleteContactByIdSupplier(txCtx, supplierID).Return(nil)
					mockGroupRepo.EXPECT().DeleteGroupByIdSupplier(txCtx, supplierID).Return(nil)
					mockMaterialRepo.EXPECT().DeleteMaterialByIdSupplier(txCtx, supplierID).Return(nil)
					mockSupplierRepo.EXPECT().DeleteSupplier(txCtx, supplierID).Return(nil)

					return f(txCtx)
				})

			err := supplierUC.DeleteSupplier(ctx, supplierID)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("GetSupplierByID Failed - Not Found", func() {
			mockSupplierRepo.EXPECT().GetSupplierByID(ctx, supplierID).Return(nil, errors.New("not found"))

			result, err := supplierUC.GetSupplierByID(ctx, supplierID)

			convey.So(err, convey.ShouldNotBeNil)
			convey.So(result, convey.ShouldBeNil)
		})
	})
}
