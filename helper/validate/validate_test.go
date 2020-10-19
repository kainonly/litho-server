package validate

import "testing"

type LoginBody struct {
	Username string `validate:"required,min=4,max=20"`
	Password string `validate:"required,min=12,max=20"`
}

func TestMake(t *testing.T) {
	message := Message{
		"Username": map[string]string{
			"required": "Submit missing [username] field",
		},
	}
	errs := Make(LoginBody{
		Username: "kain",
		Password: "pass@VAN1234",
	}, message)
	t.Log(errs)
	errs = Make(LoginBody{
		Username: "",
		Password: "pass@VAN1234",
	}, message)
	t.Log(errs)
}
