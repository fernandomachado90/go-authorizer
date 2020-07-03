# go-authorizer
Simple **command line**  application that authorizes a transaction 
for a specific account following a set of predefined rules.

The program expects `json` lines as inputs in the `stdin`, 
and provides a `json` line output for each one.

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

## Design choices

### Architecture

I decided to use `Go` due to simplicity reasons since the problem statement didn't require anything too elaborate. I've 
been writing `Go` recently and wanted to experiment more with it and keep the flow going.

The code had to be organized on a **single package**  due to `Go` limitations. In order to refer to code from other packages,
I would need to import them from a repository, which would go against the **anonymity** requirement of this challenge.

Even on this flat architecture, I was able to make a clear split between  **pure** and **impure** logic, 
making it easy to refactor the code into dedicated **core domain** and **interface adapter** packages in the future.

### Solution

#### Input decoding

After the program starts, every line received on `stdin` is tentatively parsed to a known structure so we can understand
if the input is related to an  **Account creation** or a **Transaction authorization** operation. 

In case the program is unable to identify the input, an empty body is printed on `stdout` as a form of feedback 
but the execution does not stop.

#### Account creation

Initializes the `CurrentAccount` global variable with the `account` informed or returns the `account-already-initialized` 
violation if an account is already set.

#### Transaction authorization

Tries to authorize a `transaction`. Updates the `CurrentAccount` state in case of success. 

The validations access simple properties directly from the `CurrentAccount` state 
to check for `insufficient-limit` and `card-not-active` violations or iterates 
through a last authorized `transactions` array to count matches in order to detect 
`high-frequency-small-interval` and `doubled-transaction` violations.

In the future, the last authorized `transactions` array could be improved to keep track of only the events 
that happened during the last **2 minutes** (interval customizable on the `IntervalMinutes` constant).

#### Output encoding

After either operation is done, a payload containing the `CurrentAccount` state is encoded along with any
`violations` that might have happened during execution and then forwarded to `stdout`.
