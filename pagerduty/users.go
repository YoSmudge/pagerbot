package pagerduty

type Users []User

type User struct{
  Id        string
  Name      string
  Email     string
}

func (a *Api) Users() (Users, error){
  var usr Users

  err := a.requestThing("users", &usr)
  if err != nil {
    return usr, err
  }

  return usr, nil
}
