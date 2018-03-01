


eth.syncinINFO [09-26|15:37:16] Imported new state entries               count=1170 elapsed=12.782ms  processed=8073 pending=12844 retry=0 duplicate=0 unexpected=0
> eth.syncing
{
  currentBlock: 1237724,
  highestBlock: 1748738,
  knownStates: 20917,
  pulledStates: 8073,
  startingBlock: 1232543
}
> 
> eth.blockNumber
0
> 

tokens on ropsten

0x6b6414efb7c3775666e74077afa549d83bdeda68   digipulse
0x95a48dca999c89e4e284930d9b9af973a7481287   bet
0xc8927c83d088dd913cbb6a29cc358718260f1bea   bcd
0x1d9c8cbc24b2eb42b3028d055bdc86d272fa730d   bcpt

-----

Transaction Receipts

PARITY

{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "transactionHash": "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238",
    "transactionIndex": "0x1", // 1
    "blockNumber": "0xb", // 11
    "blockHash": "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
    "cumulativeGasUsed": "0x33bc", // 13244
    "gasUsed": "0x4dc", // 1244
    "contractAddress": "0xb60e8dd61c5d32be8058bb8eb970870f07233155", // or null, if none was created
    "logs": [{ ... }, { ... }, ...]] // logs as returned by eth_getFilterLogs, etc.
  }
}

GETH

{
"id":1,
"jsonrpc":"2.0",
"result": {
     transactionHash: '0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238',
     transactionIndex:  '0x1', // 1
     blockNumber: '0xb', // 11
     blockHash: '0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b',
     cumulativeGasUsed: '0x33bc', // 13244
     gasUsed: '0x4dc', // 1244
     contractAddress: '0xb60e8dd61c5d32be8058bb8eb970870f07233155' // or null, if none was created
     logs: [{
         // logs as returned by getFilterLogs, etc.
     }, ...],
     status: '0x1'
  }
}