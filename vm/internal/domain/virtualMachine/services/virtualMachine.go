package services

import (
	"vm/internal/domain/virtualMachine/repository"
)

type IVMService interface {
	CreateVM() error
	DestroyVM() error
	StartVM() error
	StopVM() error
	RenewVM() error
	GetVMInfo() error
}

type VMService struct {
	VMRepo repository.IVirtualMachineRepository
}

func NewVMService(vmRepo repository.IVirtualMachineRepository) *VMService {
	return &VMService{VMRepo: vmRepo}
}

func (V VMService) CreateVM() (err error) {
	err = V.VMRepo.CreateVM(nil)
	if err != nil {
		return err
	}
	return
}

func (V VMService) DestroyVM() (err error) {
	//TODO implement me
	panic("implement me")
}

func (V VMService) StartVM() (err error) {
	//TODO implement me
	panic("implement me")
}

func (V VMService) StopVM() (err error) {
	//TODO implement me
	panic("implement me")
}

func (V VMService) RenewVM() (err error) {
	//TODO implement me
	panic("implement me")
}

func (V VMService) GetVMInfo() (err error) {
	//TODO implement me
	panic("implement me")
}
