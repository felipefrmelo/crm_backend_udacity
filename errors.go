package crm

type ResponseError struct {
	Message string `json:"message"`
}

func (r *ResponseError) Error() string {
	return r.Message

}

const NotFoundErrorMsg = "Customer not found"
