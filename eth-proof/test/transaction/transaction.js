var Web3 = require('web3')
var EP   = require('./../../index')
var eP   = new EP(new Web3.providers.HttpProvider("https://gmainnet.infura.io"))


describe('getTransactionProof', function () {
  it('should be able to request a proof from web3 and verify it', function (done) {
    eP.getTransactionProof('0xb53f752216120e8cbe18783f41c6d960254ad59fac16229d4eaec5f7591319de').then((result)=>{
      console.log(result)
      EP.transaction(result.path, result.value, result.parentNodes, result.header, result.blockHash).should.be.true(result.path)
      // console.log('0x'+result.blockHash.toString('hex'))
      done()
    }).catch((e)=>{console.log(e)})
  });

});
