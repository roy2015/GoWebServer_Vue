package test

type BaseTest struct {
	curUserid   int
	curUsername string
}

func (this *BaseTest) GetUserName() string {
	return ""
}

func (this *BaseTest) SetUserName(name string) {

}
