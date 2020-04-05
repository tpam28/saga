## About it ##
<b>Saga</b> is a codegenerator for a saga orchestrator pattern. This use an yaml format for a config. Look at the example in example/config.

### Config ###
A config describe a finite state machine of your transaction.
The config musts contain a milestone which contains states of the transaction.
Any state musts contain state statement.
which of that may be one of three types:
1. Compensatable 
2. Pivot
3. Retriable

Compesatable and Pivot step must contain three states – pending, approval and rejected.
Retriable step contains two states – pending and approval. 

### Code-generation ###
See makefile.
it have two parameters intput and output.


<h3>Example</h3>
Use make test-run task=1
![First task]
(https://mermaid.ink/img/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gIGF1dG9udW1iZXJcbiAgVGl0bGU6IEZpcnN0IHRhc2tcbiAgVmVyaWZ5IGNvbnN1bWVyIC0-PiBPcmNoZXN0cmF0b3I6QWNjZXB0ZWRcbiAgT3JjaGVzdHJhdG9yIC0-PiBDcmVhdGUgdGlja2V0OlBlbmRpbmdcbiAgQ3JlYXRlIHRpY2tldCAtLT4-IE9yY2hlc3RyYXRvcjpBY2NlcHRlZFxuICBPcmNoZXN0cmF0b3IgLT4-IFZlcmlmeSBjYXJkOlBlbmRpbmdcbiAgVmVyaWZ5IGNhcmQgLS0-PiBPcmNoZXN0cmF0b3I6QWNjZXB0ZWRcbiAgT3JjaGVzdHJhdG9yIC0-PiBDb25maXJtIHRpY2tldDpQZW5kaW5nXG4gIENvbmZpcm0gdGlja2V0IC0tPj4gT3JjaGVzdHJhdG9yOkFjY2VwdGVkXG4gIE9yY2hlc3RyYXRvciAtPj4gQ29uZmlybSBvcmRlcjpQZW5kaW5nXG4gIENvbmZpcm0gb3JkZXIgLS0-PiBPcmNoZXN0cmF0b3I6QWNjZXB0ZWQiLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9LCJ1cGRhdGVFZGl0b3IiOmZhbHNlfQ)
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
##### What's happened when service rejected tx which was rejected? #####

##### What's happened when retriable tx is rejected? #####