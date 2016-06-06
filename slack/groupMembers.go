package slack

import(
  "strings"
)

func (a *Api) groupId(groupName string) (string, error){
  var ugId string
  g, err := a.api.GetUserGroups()
  if err != nil {
    return ugId, err
  }

  for _,g := range g{
    if g.Handle != groupName{
      continue
    }
    ugId = g.ID
  }

  return ugId, nil
}

func (a *Api) GroupMembers(groupName string) ([]string, error){
  var usr []string
  ugId, err := a.groupId(groupName)
  if err != nil {
    return usr, err
  }

  m, err := a.api.GetUserGroupMembers(ugId)
  if err != nil {
    return usr, err
  }

  for _,id := range m{
    usr = append(usr, id)
  }

  return usr, nil
}

func (a *Api) UpdateMembers(groupName string, users []string) error{
  var userList string = strings.Join(users, ",")
  ugId, err := a.groupId(groupName)
  if err != nil {
    return err
  }

  _, err = a.api.UpdateUserGroupMembers(ugId, userList)
  return err
}
