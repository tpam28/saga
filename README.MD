##About it##
<b>Saga</b> is a codegenerator for a saga orchestrator pattern. This use an yaml format for a config. Look at the example in example/config.

###Config###
A config describe a finite state machine of your transaction.
The config musts contain a milestone which contains states of the transaction.
Any state musts contain state statement.
which of that may be one of three types:
1. Compensatable 
2. Pivot
3. Retriable

Compesatable and Pivot step must contain three states – pending, approval and rejected.
Retriable step contains two states – pending and approval. 

###Code-generation###
See makefile.
it have two parameters intput and output.


<h3>Example</h3>
Use make test-run task=1
```mermaid
sequenceDiagram
  Verify consumer ->> Orchestrator:Accepted
  Orchestrator ->> Create ticket:Pending
  Create ticket ->> Orchestrator:Accepted
  Orchestrator ->> Verify card:Pending
  Verify card ->> Orchestrator:Accepted
  Orchestrator ->> Confirm ticket:Pending
  Confirm ticket ->> Orchestrator:Accepted
  Orchestrator ->> Confirm order:Pending
  Confirm order ->> Orchestrator:Accepted
```
#####What's happened when service rejected tx which was rejected?#####

#####What's happened when retriable tx is rejected?#####