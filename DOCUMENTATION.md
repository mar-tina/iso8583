# iso 8583

#Authorization only
- no money movement
- can be for pre-authorization for future money movement

```
   FIELD 0 ll:4]
     - MTI --> 0100
   FIELD 1 [l:16]
     - BITMAP --> 
   FIELD 2 [l:19]
     - PAN --> 400555000000****
   FIELD 3 [l:2]
     - PROCESSING CODE --> 00 [purchase of goods with card] | 01 [cash]
   FIELD 4 [l:12]
     - TRANSACTION AMOUNT --> 000000000100 
   FIELD 7 [l: 10]
     - TRANSMISSION DATE AND TIME --> 0302211500
   FIELD 11 [l:6]
     - SYSTEMS TRACE AUDIT NUMBER --> 000001 - 999999
   FIELD 42 [l:10]
     - CARD ACCEPTOR ID CODE --> 1231234567
```
