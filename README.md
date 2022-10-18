# iso 8583 golang library && simulator

protocol which may used for passing credit card, debit card and/or 
check information to and from ECHO


### How to use library

```
    go get github.com/mar-tina/iso8583
```

```go
    import (
        "github.com/mar-tina/iso8583/lib/spec"
        "github.com/mar-tina/iso8583/lib/msg"
    )
```

### Define your spec

This is a basic go struct with some tags to provide info about packing
the message

**TO NOTE*** 

The message is will be packed in the order in which the struct fields are defined

```go
type NetMgmtREQ struct {
	Mti                  string `field:"0" ln:"4" json:"0"`
	TxnDateTime          string `field:"7" ln:"10" json:"7"`
	Stan                 string `field:"11" ln:"19" lvar:"1" json:"11"`
	LocalTransactionTime string `field:"12" ln:"6" json:"12"`
	LocalTransactionDate string `field:"13" ln:"6" json:"13"`
	SettlementDate       string `field:"15" ln:"4" json:"15"`
	NetMgmntInfoCode     string `field:"70" ln:"10" json:"70"`
}
```

The fields **field** | **ln** are required fields

```
field represents fieldId i.e 0 to indicate mti
```

### register your spec
library stores an internal signature of your spec and matches it
when you pack the msg.

```go
   err := spec.Register(NetMgmtREQ{}, NetMgmtREQMod{})
   if err != nil {
       //handle error
   }
```

### pack your msg

```go
    ...
    msg1, err := msg.PackMsg(netmgtREQ)
	if err != nil {
		//handle error
	}
    ...

    //0800823a0000000000004000000000000000042009061390000109061304200420001
```

TODO: Add a description and encoding tag

### Project status

Project is under active development

- [x] Generate iso messages
- [ ] Unpack iso messages
- [ ] Network read
- [ ] Testing 
- [ ] Robust Fields