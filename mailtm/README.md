# MailTM Wrapper
A convenience-oriented [mail.tm](https://mail.tm) API wrapper written in Golang

Feel free to report bugs and suggest improvements!

Copy by https://github.com/msuny-c/mailtm/tree/0c39880925d6ce6a0651720dca77fd72cb1e831b

## Installation
```
go get github.com/msuny-c/mailtm
```
## Getting started
### Register
You can create a new account with random credentials
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
}
```
Or provide data directly
```go
import "github.com/msuny-c/mailtm"

func main() {
    opts := mailtm.Options {
        Domain: mailtm.AvailableDomains()[0].Domain,
        Username: "someusername",
        Password: "mypassword",
    }
    account, err := mailtm.NewAccountWithOptions(opts)
    if err != nil {
        panic(err)
    }
}
```
### Login
You can login to your existing account using your address and password
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.Login("username@mail.tm", "mypassword")
    if err != nil {
        panic(err)
    }
}
```
Or using Bearer token
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.LoginWithToken("bearertoken")
    if err != nil {
        panic(err)
    }
}
```
### Working with messages
To get a message you can use the `MessagesAt(id)` method, which returns a slice of messages with their contents on a specific page
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
    msgs, err := account.MessagesAt(1)
    if err != nil {
        print("failed to get messages")
    }
}
```
You can get a message channel that will receive new messages from current moment
```go
import (
    "github.com/msuny-c/mailtm"
    "context"
)

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
    ctx, cancel := context.WithCancel(context.Background())
    ch := account.MessagesChan(ctx)
    for {
        select {
        case msg, ok := <- ch:
            if ok {
                print(msg.HTML)
                cancel()
            }
        }
    }
}
```
Also you can get the last message or by it's id
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
    msg1, err := account.MessageById("someid")
    msg2, err := account.LastMessage()
    if err != nil {
        print("failed to get messages")
    }
}
```
And of course you can delete message
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
    msg, err := account.LastMessage()
    if err != nil {
        print("failed to get message")
    }
    account.DeleteMessage(msg.ID)
}
```
### Account
You can get account's properties (those that are returned in the response by [api.mail.tm](https://api.mail.tm))
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
    print(account.Property("createdAt"))
}
```
Also get address, password and token fields
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
    println(account.Address())
    println(account.Password())
    println(account.Bearer())
}
```
If you wish you can delete your account
```go
import "github.com/msuny-c/mailtm"

func main() {
    account, err := mailtm.NewAccount()
    if err != nil {
        panic(err)
    }
    account.Delete()
}
```
