package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/request"
	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"github.com/Temisaputra/warOnk/internal/delivery/repository"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/xuri/excelize/v2"
)

type SupplierUsecase struct {
	SupplierRepo    repository.SupplierRepository
	AddressRepo     repository.AddressRepository
	ContactRepo     repository.ContactRepository
	GroupRepo       repository.GroupRepository
	MaterialRepo    repository.MaterialRepository
	transactionRepo repository.TransactionRepository
}

func NewSupplierUsecase(supplierRepository repository.SupplierRepository, addressRepository repository.AddressRepository, contactRepository repository.ContactRepository, groupRepository repository.GroupRepository, materialRepository repository.MaterialRepository, transactionRepository repository.TransactionRepository) *SupplierUsecase {
	return &SupplierUsecase{
		SupplierRepo:    supplierRepository,
		AddressRepo:     addressRepository,
		ContactRepo:     contactRepository,
		GroupRepo:       groupRepository,
		MaterialRepo:    materialRepository,
		transactionRepo: transactionRepository,
	}
}

func (u *SupplierUsecase) GetAllSupplier(ctx context.Context, pagination *request.Pagination, params *presenter.SupplierRequest) (suppliers []*presenter.SupplierResponse, meta response.Meta, err error) {
	if pagination.Page == 0 {
		pagination.Page = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}

	res, meta, err := u.SupplierRepo.GetAllSupplier(ctx, pagination, params)
	if err != nil {
		errMsg := fmt.Errorf("[SupplierUsecase-GetAllSupplier] Error when getting all supplier: %w", err)
		return nil, meta, errMsg
	}

	for i := range res {
		address, err := u.AddressRepo.GetAddressByIdSupplier(ctx, res[i].ID)
		if err != nil {
			errMsg := fmt.Errorf("[SupplierUsecase-GetAllSupplier] Error when getting address by supplier id: %w", err)
			return nil, meta, errMsg
		}

		contact, err := u.ContactRepo.GetContactByIdSupplier(ctx, res[i].ID)
		if err != nil {
			errMsg := fmt.Errorf("[SupplierUsecase-GetAllSupplier] Error when getting contact by supplier id: %w", err)
			return nil, meta, errMsg
		}
		res[i].Address = *address.ToPresenter()
		res[i].Contact = *contact.ToPresenter()
	}

	return res, meta, nil
}

func (u *SupplierUsecase) GetSupplierByID(ctx context.Context, id int) (res *presenter.SupplierResponse, err error) {
	res, err = u.SupplierRepo.GetSupplierByID(ctx, id)
	if err != nil {
		errMsg := fmt.Errorf("[SupplierUsecase-GetSupplierByID] Error when getting supplier by id: %w", err)
		return nil, errMsg
	}

	address, err := u.AddressRepo.GetAddressByIdSupplier(ctx, res.ID)
	if err != nil {
		errMsg := fmt.Errorf("[SupplierUsecase-GetSupplierByID] Error when getting address by supplier id: %w", err)
		return nil, errMsg
	}

	contact, err := u.ContactRepo.GetContactByIdSupplier(ctx, res.ID)
	if err != nil {
		errMsg := fmt.Errorf("[SupplierUsecase-GetSupplierByID] Error when getting contact by supplier id: %w", err)
		return nil, errMsg
	}

	group, err := u.GroupRepo.GetGroupByIdSupplier(ctx, res.ID)
	if err != nil {
		errMsg := fmt.Errorf("[SupplierUsecase-GetSupplierByID] Error when getting group by supplier id: %w", err)
		return nil, errMsg
	}
	material, err := u.MaterialRepo.GetMaterialByIdSupplier(ctx, res.ID)
	if err != nil {
		errMsg := fmt.Errorf("[SupplierUsecase-GetSupplierByID] Error when getting material by supplier id: %w", err)
		return nil, errMsg
	}

	res.Address = *address.ToPresenter()
	res.Contact = *contact.ToPresenter()
	res.Group = *group.ToPresenter()
	res.Material = make([]presenter.MaterialResponse, len(*material))
	for i, m := range *material {
		res.Material[i] = *m.ToPresenter()
	}

	return res, nil
}

func (u *SupplierUsecase) CreateSupplier(ctx context.Context, params *presenter.SupplierCreateRequest) (err error) {
	// Mulai transaksi opsional
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		storedSupplier := &presenter.SupplierRequest{
			Name:     params.Name,
			NickName: params.NickName,
			Status:   params.Status,
		}
		// Semua repo harus pakai txCtx
		supIDInt, err := u.SupplierRepo.CreateSupplier(txCtx, storedSupplier)
		if err != nil {
			return err // rollback otomatis
		}

		supplierID := uint(supIDInt)

		if len(params.Address) > 0 {
			for _, addr := range params.Address {
				err := u.AddressRepo.CreateAddress(txCtx, &entity.Address{
					Name:       addr.Name,
					Address:    addr.Address,
					Main:       addr.Main,
					SupplierID: supplierID,
				})
				if err != nil {
					return err // rollback otomatis
				}

			}

		}

		if len(params.Contact) > 0 {
			for _, cont := range params.Contact {
				err := u.ContactRepo.CreateContact(txCtx, &entity.Contact{
					Name:       cont.Name,
					Phone:      cont.Phone,
					Main:       cont.Main,
					SupplierID: supplierID,
				})
				if err != nil {
					return err // rollback otomatis
				}
			}
		}

		fmt.Println("params.Groups:", params.Groups)

		if len(params.Groups) > 0 {

			// Insert group baru
			var groups []entity.SupplierGroup
			for _, gid := range params.Groups {
				groups = append(groups, entity.SupplierGroup{
					SupplierID: uint(supIDInt),
					GroupID:    uint(gid),
				})
			}
			if err := u.GroupRepo.CreateGroup(txCtx, &groups); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when inserting new groups: %w", err)
			}
		}

		if len(params.Materials) > 0 {
			// Insert material list baru
			var materials []entity.MaterialList
			for _, mt := range params.Materials {
				materials = append(materials, entity.MaterialList{
					SupplierID:    uint(supIDInt),
					MaterialGroup: mt.MaterialGroup,
					MaterialID:    mt.MaterialID,
				})
			}
			if err := u.MaterialRepo.CreateMaterial(txCtx, &materials); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when inserting new material lists: %w", err)
			}
		}

		return nil // commit otomatis
	})
}

func (u *SupplierUsecase) UpdateSupplier(ctx context.Context, params *presenter.SupplierCreateRequest, id int) (err error) {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Ambil supplier dulu, pakai txCtx
		supplier, err := u.SupplierRepo.GetSupplierByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("[SupplierUsecase-UpdateSupplier] Error when getting supplier by id: %w", err)
		}

		if supplier == nil {
			return fmt.Errorf("[SupplierUsecase-UpdateSupplier] Supplier not found")
		}

		updatedData := &presenter.SupplierRequest{
			Name:     params.Name,
			NickName: params.NickName,
			Status:   params.Status,
		}

		// Update supplier pakai txCtx
		if err := u.SupplierRepo.UpdateSupplier(txCtx, updatedData, id); err != nil {
			return fmt.Errorf("[SupplierUsecase-UpdateSupplier] Error when updating supplier: %w", err)
		}

		if len(params.Address) > 0 {
			// Hapus address lama
			if err := u.AddressRepo.DeleteAddressByIdSupplier(txCtx, id); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when deleting old addresses: %w", err)
			}

			// Insert address baru
			for _, addr := range params.Address {
				err := u.AddressRepo.CreateAddress(txCtx, &entity.Address{
					Name:       addr.Name,
					Address:    addr.Address,
					Main:       addr.Main,
					SupplierID: uint(id),
				})
				if err != nil {
					return err // rollback otomatis
				}

			}
		}

		if len(params.Contact) > 0 {
			// Hapus contact lama
			if err := u.ContactRepo.DeleteContactByIdSupplier(txCtx, id); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when deleting old contacts: %w", err)
			}

			// Insert contact baru
			for _, cont := range params.Contact {
				err := u.ContactRepo.CreateContact(txCtx, &entity.Contact{
					Name:       cont.Name,
					Phone:      cont.Phone,
					Main:       cont.Main,
					SupplierID: uint(id),
				})
				if err != nil {
					return err // rollback otomatis
				}
			}
		}

		fmt.Println("params.Groups:", params.Groups)

		if len(params.Groups) > 0 {
			// Hapus group lama
			if err := u.GroupRepo.DeleteGroupByIdSupplier(txCtx, id); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when deleting old groups: %w", err)
			}

			// Insert group baru
			var groups []entity.SupplierGroup
			for _, gid := range params.Groups {
				groups = append(groups, entity.SupplierGroup{
					SupplierID: uint(id),
					GroupID:    uint(gid),
				})
			}
			if err := u.GroupRepo.CreateGroup(txCtx, &groups); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when inserting new groups: %w", err)
			}
		}

		if len(params.Materials) > 0 {
			// Hapus material list lama
			if err := u.MaterialRepo.DeleteMaterialByIdSupplier(txCtx, id); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when deleting old material lists: %w", err)
			}

			// Insert material list baru
			var materials []entity.MaterialList
			for _, mt := range params.Materials {
				materials = append(materials, entity.MaterialList{
					SupplierID:    uint(id),
					MaterialGroup: mt.MaterialGroup,
					MaterialID:    mt.MaterialID,
				})
			}
			if err := u.MaterialRepo.CreateMaterial(txCtx, &materials); err != nil {
				return fmt.Errorf("[SupplierUsecase-UpdateSupplier] error when inserting new material lists: %w", err)
			}
		}

		return nil // commit otomatis kalau tidak ada error
	})
}

func (u *SupplierUsecase) DeleteSupplier(ctx context.Context, id int) (err error) {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Ambil supplier dulu, pakai txCtx
		supplier, err := u.SupplierRepo.GetSupplierByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("[SupplierUsecase-DeleteSupplier] Error when getting supplier by id: %w", err)
		}

		if supplier == nil {
			return fmt.Errorf("[SupplierUsecase-DeleteSupplier] Supplier not found")
		}

		// Hapus address terkait
		if err := u.AddressRepo.DeleteAddressByIdSupplier(txCtx, id); err != nil {
			return fmt.Errorf("[SupplierUsecase-DeleteSupplier] error when deleting related addresses: %w", err)
		}

		// Hapus contact terkait
		if err := u.ContactRepo.DeleteContactByIdSupplier(txCtx, id); err != nil {
			return fmt.Errorf("[SupplierUsecase-DeleteSupplier] error when deleting related contacts: %w", err)
		}

		// Hapus group terkait
		if err := u.GroupRepo.DeleteGroupByIdSupplier(txCtx, id); err != nil {
			return fmt.Errorf("[SupplierUsecase-DeleteSupplier] error when deleting related groups: %w", err)
		}

		// Hapus material list terkait
		if err := u.MaterialRepo.DeleteMaterialByIdSupplier(txCtx, id); err != nil {
			return fmt.Errorf("[SupplierUsecase-DeleteSupplier] error when deleting related material lists: %w", err)
		}

		// Delete supplier pakai txCtx
		if err := u.SupplierRepo.DeleteSupplier(txCtx, id); err != nil {
			return fmt.Errorf("[SupplierUsecase-DeleteSupplier] Error when deleting supplier: %w", err)
		}

		return nil // commit otomatis kalau tidak ada error
	})
}

func (u *SupplierUsecase) BlockUnblockSupplier(ctx context.Context, id int, params presenter.SupplierBlockUnblockRequest) (err error) {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		// Ambil supplier dulu, pakai txCtx
		supplier, err := u.SupplierRepo.GetSupplierByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("[SupplierUsecase-BlockUnblockSupplier] Error when getting supplier by id: %w", err)
		}
		if supplier == nil {
			return fmt.Errorf("[SupplierUsecase-BlockUnblockSupplier] Supplier not found")
		}
		// Block supplier pakai txCtx
		if err := u.SupplierRepo.BlockUnblockSupplier(txCtx, id, supplier.Status); err != nil {
			return fmt.Errorf("[SupplierUsecase-BlockUnblockSupplier] Error when blocking supplier: %w", err)
		}

		var reason string
		var status string
		if supplier.Status == "Active" {
			status = "Blocked"
		} else {
			status = "Active"
		}
		if supplier.Status == "Active" {
			reason = "Supplier blocked: " + params.Reason
		} else {
			reason = "Supplier unblocked: " + params.Reason
		}
		// insert log block/unblock supplier pakai txCtx
		if err := u.SupplierRepo.InsertHistorySupplier(txCtx, &entity.StatusHistory{
			SupplierID: uint(id),
			Status:     status,
			Reason:     reason,
		}); err != nil {
			return fmt.Errorf("[SupplierUsecase-BlockUnblockSupplier] Error when inserting log block/unblock supplier: %w", err)
		}
		return nil // commit otomatis kalau tidak ada error
	})
}

func (u *SupplierUsecase) ExportSupplier(ctx context.Context, params *presenter.SupplierRequest) (fileName string, err error) {

	exportData, err := u.SupplierRepo.ExportSupplier(ctx, params)
	if err != nil {
		return "", fmt.Errorf("[SupplierUsecase-ExportSupplier] Error when getting supplier data: %w", err)
	}

	if len(exportData) < 1 {
		return "", fmt.Errorf("no data to export")
	}

	for i := range exportData {
		address, err := u.AddressRepo.GetAddressByIdSupplier(ctx, exportData[i].ID)
		if err != nil {
			errMsg := fmt.Errorf("[SupplierUsecase-GetAllSupplier] Error when getting address by supplier id: %w", err)
			return "", errMsg
		}

		contact, err := u.ContactRepo.GetContactByIdSupplier(ctx, exportData[i].ID)
		if err != nil {
			errMsg := fmt.Errorf("[SupplierUsecase-GetAllSupplier] Error when getting contact by supplier id: %w", err)
			return "", errMsg
		}
		exportData[i].Address = *address.ToPresenter()
		exportData[i].Contact = *contact.ToPresenter()
	}

	excelFile := excelize.NewFile()

	sheetName := excelFile.GetSheetName(0)

	title := "Export Supplier"
	excelFile.SetCellValue(sheetName, "A1", title)

	excelFile.MergeCell(sheetName, "A1", "H1")

	style := &excelize.Style{}
	err = json.Unmarshal([]byte(`{
		"font": {
			"bold": false,
			"size": 10,
			"style": "arial",
		},
		"alignment": {
			"horizontal": "left",
			"vertical": "left"
		}
	}`), style)

	styleIndex, err := excelFile.NewStyle(style)
	if err != nil {
		return "", fmt.Errorf("failed to create excel style: %w", err)
	}

	excelFile.SetCellStyle(sheetName, "A1", "H1", styleIndex)
	headers := []string{"ID", "NAMA SUPPLIER", "ADDRESS", "CONTACT", "STATUS"}
	for col, header := range headers {
		cell := fmt.Sprintf("%c2", 'A'+col)
		excelFile.SetCellValue(sheetName, cell, header)
	}
	sort.Slice(exportData, func(i, j int) bool {
		return exportData[i].Name < exportData[j].Name
	})

	for i, exportSupplier := range exportData {
		slice := []any{
			exportSupplier.ID,
			exportSupplier.Name,
			exportSupplier.Address.Name,
			exportSupplier.Contact.Name,
			exportSupplier.Status,
		}
		excelFile.SetSheetRow(excelFile.GetSheetName(0), fmt.Sprintf("A%d", i+3), &slice)
	}

	currentTime := time.Now().Format("2006-01-02 15_04_05")

	fileName = fmt.Sprintf("Export_Supplier_%s.xlsx", currentTime)

	fileReader, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create excel file: %w", err)
	}
	defer fileReader.Close()

	err = excelFile.SaveAs(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create excel file: %w", err)
	}

	return fileName, nil
}
