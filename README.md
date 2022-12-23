# disys-exam-peer-to-peer
Made for the 3rd semester subject DISYS

## How to run
To run go-ass4 first start client 1:
```bash
go run peer.go
```
Then if it is the first client input 1, if it is the second input 2 and so on.

This will start three peers. They will start to discuss internally which peer should get control of the plane.

## Good to know
### User Id / Port name
The numbers used as args when starting "main" are interpreted as port numbers. Thus, in the system logs, you will not see the id that has been given but the port number instead. The port number is calculated by adding the id and 5000.

### System log
When a peer request control of plane the reply it will get back will get in the form:
```bash
<time (HH:mm:ss.SSSSSS)> Got reply from id <id> : <request amount> : <is pilot>
```
When a peer is the pilot the log will be in the form:
```bash
<time (HH:mm:ss.SSSSSS)> <id> is now pilot 	-----------------------
```
When they stop being the pilot the log will be in the form:
```bash
<time (HH:mm:ss.SSSSSS)> <id> is not pilot 	-----------
```
When they are starting to try to take control of the plane, the log will be in the form:
```bash
locked
```
When they have either succeeded or failed the log will be in the form:
```bash
unlock
```
