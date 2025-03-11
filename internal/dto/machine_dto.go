package dto

import "insist-backend-golang/internal/model"

type CreateMachinePayload struct {
	Machine       model.MstMachine       `json:"machine"`
	MachineDetail model.MstMachineDetail `json:"machine_detail"`
	MachineStatus model.MstMachineStatus `json:"machine_status"`
}
