## Phase 1

### Core Test Metrics
- DNS
- TCP (round trip, handshake, tls done)
- ICMP

### Sensor Code
1. Implement the structure to run desired tests per domain
    - Should be concurrent per domain
    - Take instructions from the outside
    - Should be able to stop on CNC call
    - Keep track of traffic volume
2. Authentication with the command server (in web3 phase?)
3. Describe common information to store about each sensor
    - Geo location
    - ISP info
    - Default NS info
    - Local
    - Hardware id?
4. Protect against:
    - Fake results
    - VPN
5. Ensure limits
    - User defined network traffic

### Server (CNC)
1. Implement monitoring & communication with all sensors
2. IF the communications is through ws (needed for real time)
    - Implement lower rank server in order to scale and protect the CNC

## Phase 2

### WEB3
1. TODO

### Sensor Infra
1. A sensor described as an NFT?
2. Run as a docker; Wrap for any OS with user-friendly install

## Phase 3

### UI
1. Can start with just prometheus?
    - For a sensor owner
    - For internal usage
    - For client usage

## Phase X
#### Run customer supplied script code for personalised metrics
#### ISP or and DNS Server rating