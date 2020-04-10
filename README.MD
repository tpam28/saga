## About it ##
<b>Saga</b> is a codegenerator for a saga orchestrator pattern.
 I read this pattern in Microservices patterns book by Kris Richardson and decided writing useful tool for generation distributed transaction.   

### Config ###
This use an yaml format for a config. Look at the example in example/config.
A config describe a finite state machine of your transaction.
The config musts contain a milestone which contains states of the transaction.
Any step musts contain state statement.
which of that may be one of three types:
1. Compensatable 
2. Pivot
3. Retriable

Compesatable and Pivot step must contain three states – pending, approval and rejected.
Retriable step contains two states – pending and approval. 

### Code-generation ###
See makefile.
it have two parameters intput and output.


### Example ###
1. Use make gen for regenerate example code.
2. Use make example-build for rebuild example applications.
3. Use  make example-run TASK=1.

<details>
  <summary>make example-run TASK=1</summary>
 
```mermaid
sequenceDiagram
  autonumber
  Title: First task
  Verify consumer ->> Orchestrator:Accepted
  Orchestrator ->> Create ticket:Pending
  Create ticket -->> Orchestrator:Accepted
  Orchestrator ->> Verify card:Pending
  Verify card -->> Orchestrator:Accepted
  Orchestrator ->> Confirm ticket:Pending
  Confirm ticket -->> Orchestrator:Accepted
  Orchestrator ->> Confirm order:Pending
  Confirm order -->> Orchestrator:Accepted
```

</details>

![First task](first_task.png)


<details>
  <summary>make example-run TASK=2</summary>


```mermaid
sequenceDiagram
  autonumber
  Title: Second task
  Verify consumer ->> Orchestrator:Accepted
  Orchestrator ->> Create ticket:Pending
  Create ticket -->> Orchestrator:Rejected
  Orchestrator ->> Verify consumer:Rejected
  Verify consumer -->> Orchestrator:Accepted
```

</details>

![Second task](second_task.png)


<details>
  <summary>make example-run TASK=3</summary>


```mermaid
sequenceDiagram
  autonumber
  Title: Third task
  Verify consumer ->> Orchestrator:Accepted
  Orchestrator ->> Create ticket:Pending
  Create ticket -->> Orchestrator:Accepted
  Orchestrator ->> Verify card:Pending
  Verify card -->> Orchestrator:Accepted
  Orchestrator ->> Confirm ticket:Pending
  loop Retry x3
  Confirm ticket ->> Confirm ticket:Rejected
  end
  Confirm ticket -->> Orchestrator:Accepted
  Orchestrator ->> Confirm order:Pending
  Confirm order -->> Orchestrator:Accepted
```

</details>

![Third task](third_task.png)


<details>
  <summary>make example-run TASK=4</summary>



```mermaid
sequenceDiagram
  autonumber
  Title: Fourth task
  Verify consumer ->> Orchestrator:Accepted
  Orchestrator ->> Create ticket:Pending
  Create ticket -->> Orchestrator:Rejected
  Orchestrator ->> Verify consumer:Rejected
  loop Retry x3
  Verify consumer ->> Verify consumer:Rejected
  end
  Verify consumer -->> Orchestrator:Accepted
```

</details>

![Fourth task](fourth_task.png)


Every message in loop also goes to orchestrator then come back.  
