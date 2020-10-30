# gofmg

A very minimal wrapper for FMG.

Usage:
```
var FMG fmg.FMG
FMG.Username = target.Username
FMG.Password = target.Apikey
FMG.IP = target.Firewallip

err := FMG.Login()
if err != nil {
  log.Println(err)
  return make([]byte, 0), err
}
res, err := FMG.Call(method, path, nil)
if err != nil {
  fmt.Println(err)
  return make([]byte, 0), err
}
output, err := json.Marshal(res)
if err != nil {
  fmt.Println(err)
  return make([]byte, 0), err
}
```
