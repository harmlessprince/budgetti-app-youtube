package requests

type IDParamRequest struct {
	ID uint `param:"id" binding:"required"`
}
