# go-authorizer
Simple **command line**  application that authorizes a transaction 
for a specific account following a set of predefined rules.

## Commands

### `make setup`
Install dependency modules.

### `make format`
Format all files using `go fmt`.

### `make build`
Build source files into an executable binary called `cmd.bin`.

### `make test`                    
Execute all available tests.

### `make run`
Runs the application and reads input from `stdin`.

### `make docker`
Runs the application on a `Docker` image and reads input from `stdin`.

## Operations
The program handles two kinds of operations, deciding on which one according to the line that is being processed.

### Account creation
Creates the account with `availableLimit` and `activeCard` set.

###### input 
    { "account": { "activeCard": true, "availableLimit": 100 }  }
###### output 
    { "account": { "activeCard": true, "availableLimit": 100 }, "violations": [] }
###### expected violations
    ["account-already-initialized"]

### Transaction authorization
Tries to authorize a transaction for a particular `merchant`, `amount` and `time` given the account's state 
and last authorized transactions.

###### input 
    { "transaction": { "merchant": "Acme Corporation", "amount": 20, "time": "2020-07-12T10:00:00.000Z" } }
###### output 
    { "account": { "activeCard": true, "availableLimit": 80 }, "violations": [] }
###### expected violations
    ["insufficient-limit", "card-not-active", "high-frequency-small-interval", "doubled-transaction"]
