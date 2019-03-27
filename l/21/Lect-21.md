<style>
.pagebreak { page-break-before: always; }
</style>

Lecture 21 - Solidity - Contract for Hw 7, 8, 9
=======================================================================

## News

1. Yield Curve has Inverted
[https://seekingalpha.com/article/4251080-inverted-yield-curve-important](https://seekingalpha.com/article/4251080-inverted-yield-curve-important)
2. Demand for Insurance and Blockchain
[https://www.ccn.com/bitcoin-expertise-exploding-among-insurance-professionals-in-2019-study](https://www.ccn.com/bitcoin-expertise-exploding-among-insurance-professionals-in-2019-study)
"A few vendors, like Black Insurance, an Estonian firm, are offering more
radical blockchain solutions, with a goal to upend the entire insurance
system. In this model, a broker, not an insurance carrier, creates a new
insurance product, then lists that product through Black Insurance to gauge
interest in capital markets. Risk is transferred directly to the capital
markets, disintermediating (eliminating) traditional insurance carriers
entirely."
3. Whiskey tracked.
[https://cointelegraph.com/news/over-130-year-old-liquor-company-william-grant-sons-to-track-whiskey-on-blockchain](https://cointelegraph.com/news/over-130-year-old-liquor-company-william-grant-sons-to-track-whiskey-on-blockchain)


Solidity
===

## Version of the compiler 

```
pragma solidity >=0.4.21 <0.6.0;
```

## Versions of Tools

```
$ truffle --version
Truffle v5.0.9 - a development framework for Ethereum
    ...
$ node --version
v10.15.3
$ npm --version
6.4.1
$ solc --version
solc, the solidity compiler commandline interface
Version: 0.5.7+commit.6da8b019.Darwin.appleclang

```

## Config File

```
{
    ...
    "Eth_URL_ws"                : "ws://192.168.0.199:8546",
    "Eth_URL_rpc"                 : "http://192.168.0.199:8545",
    "Eth_Account_Address"        : "0x023e291a99d21c944a871adcc44561a58f99bdbc",
    "Eth_Account_KeyFile"        : "./keystore/UTC--2018-07-17T00-00-00.000000000Z--023e291a99d21c944a871adcc44561a58f99bdbc",
    "Eth_Contract_SignedData"    : "0x0134291399821c944a871adcc44561a58f99bdbc"
    ...
}

```

<div class="pagebreak"></div>
## Example for Homework

```
  1 pragma solidity >=0.4.21 <0.6.0;
  2 
  3 import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
  4 
  5 contract SignedDataVersion01 is Ownable {
  6 
  7     address payable owner_address;
  8     uint256 private minPayment;
  9 
 10     mapping(uint256 => mapping(uint256 => bytes32)) dData;
 11     mapping(uint256 => mapping(uint256 => address)) dOowner;
 12     event DataChange(uint256 App, uint256 Name, bytes32 Value, address By);
 13 
 14     event ReceivedFunds(address sender, uint256 value,
            uint256 application, uint256 payFor);
 15     event Withdrawn(address to, uint256 amount);
 16 
 17     constructor() public {
 18         owner_address = msg.sender;
 19         minPayment = 1000;
 20     }
 21 
 22     modifier needMinPayment {
 23         require(msg.value >= minPayment,
                "Insufficient payment.  Must send more than minPayment.");
 24         _;
 25     }
 26 
 27     function init() public {
 28         minPayment = 1000;
 29     }
 30 
 31     function setMinPayment( uint256 _minPayment ) public onlyOwner {
 32         minPayment = _minPayment;
 33     }
 34 
 35     function getMinPayment() public onlyOwner view returns ( uint256 ) {
 36         return ( minPayment );
 37     }
 38 
```
<div class="pagebreak"></div>
```
 41     /**
 42      * @dev TODO if the data is empty, or if the msg.sender is the original
                createor of the data:
 43      *      then : save the msg.sender into dOwner, save the data into dData
 44      *             create a DataChange event.
 45      *      else : revert an error.
 46      */
 47     function setData ( uint256 _app, uint256 _name, bytes32 _data )
        public needMinPayment payable {
 48         // TODO-start - code for students to implement -- about 7 or 8 lines of code
 49         // Create a temorary variable with dOwner[_app][_name]
 50         // If the msg.sender is the owner and or if tmp is address(0) "not assigned" then
 51         // Save the msg.sender to dOnwer
 52         // Save the data.
 53         // Emit DataChange event
 54         // Else
 55         // Revert an error
 57         // TODO-end
 58     }
 59 
 60     /**
 61      * @dev TODO return the data by looking up _app and _name in dData.
 62      */
 63     function getData ( uint256 _app, uint256 _name )
        public view returns ( bytes32 ) {
 64         // TODO-start - code for students to implement -- really just one l.o.c. -
 65         // Lookup the data and return it.
 66         // TODO-end
 67     }
 68 
 71     /**
 72      * @dev payable fallback
 73      */
 74     function () external payable {
 75         emit ReceivedFunds(msg.sender, msg.value, 0, 1);
 76     }
 77 
 78     /**
 79      * @dev genReceiveFunds - generate a receive funds event.
 80      */
 81     function genReceivedFunds ( uint256 application, uint256 payFor )
        public payable {
 82         emit ReceivedFunds(msg.sender, msg.value, application, payFor);
 83     }
 84 
 85     /**
 86      * @dev Withdraw contract value amount.
 87      */
 88     function withdraw( uint256 amount ) public onlyOwner returns(bool) {
 89         address(owner_address).transfer(amount);
 90         // owner_address.send(amount);
 91         emit Withdrawn(owner_address, amount);
 92         return true;
 93     }
 94 
 95     /**
 96      * @dev How much do I got?
 97      */
 98     function getBalanceContract() public view onlyOwner returns(uint256){
 99         return address(this).balance;
100     }
101 
102     /**
103      * @dev For futute to end the contract, take the value.
104      */
105     function kill() public onlyOwner {
106         emit Withdrawn(owner_address, address(this).balance);
107         selfdestruct(owner_address);
108     }
109 }
```

