package grpc_helpers

import (
	"connectrpc.com/connect"
	"github.com/kavkaco/Kavka-Core/utils/vali"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func VarrorAsGrpcErrDetails(varror *vali.Varror) (*connect.ErrorDetail, error) {
	fieldViolations := []*errdetails.BadRequest_FieldViolation{}

	for _, ve := range varror.ValidationErrors {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       ve.Field(),
			Description: ve.Translate(vali.ValiTranslator),
		})
	}

	return connect.NewErrorDetail(&errdetails.BadRequest{
		FieldViolations: fieldViolations,
	})
}

func GrpcVarror(varror *vali.Varror, code connect.Code) error {
	connectErr := connect.NewError(code, varror.Error)
	varrorDetail, err := VarrorAsGrpcErrDetails(varror)
	if err != nil {
		return err
	}

	if len(varror.ValidationErrors) > 0 {
		connectErr.AddDetail(varrorDetail)
	}

	return connectErr
}
