const ethers = require('ethers')
const abi = require('./abi.json');


const provider = getProvider();
const signer = getSigner(provider);
const contract = new ethers.Contract('0x33a226fdA02004c2e6e3b4Ad13DfF79187DD1B80', abi, signer);


// These constants must match the ones used in the smart contract.
const SIGNING_DOMAIN_NAME = "GameItem"
const SIGNING_DOMAIN_VERSION = "1"


async function signData() {
  const voucher = { 
      tokenId: 123, 
      itemType: 1,
      strength: 5,
      level: 1,
      expireTime: 123456789//getExpireTime()
  }
  const domain = await getSigningDomain()
  const types = {
    ItemInfo: [
      {name: "tokenId", type: "uint256"},
      {name: "itemType", type: "uint256"},
      {name: "strength", type: "uint256"},
      {name: "level", type: "uint256"},
      {name: "expireTime", type: "uint256"}
    ]
  }
  const signature = await signer._signTypedData(domain, types, voucher)
  return {
    ...voucher,
    signature,
  };
}

signData().then(async result => {
  console.log(result);
  try {
    //var contractWithSigner = contract.connect(signer);
    console.log(signer.address);
    var tx = await contract.mintTokenWithSignedMessage(signer.address, result, {
      gasLimit: 250000
    });
    console.log(tx.hash);
    var r = await tx.wait();
    console.log(r);
  } catch (e) {
    console.log(e);
  }
});



async function getSigningDomain() {
    //const chainId = await contract.getChainID()
    const { chainId } = await provider.getNetwork()
    console.log('chainId', chainId);
    const domain = {
      name: SIGNING_DOMAIN_NAME,
      version: SIGNING_DOMAIN_VERSION,
      verifyingContract: contract.address,
      chainId,
    }
    return domain
}

function getExpireTime() {
    var minutesToAdd=15;
    var currentDate = new Date();
    var futureDate = new Date(currentDate.getTime() + minutesToAdd*60000);
    return Math.floor(futureDate.getTime() / 1000);
}

function getProvider() {
  const NODE_URL = "https://data-seed-prebsc-1-s1.binance.org:8545/";
  const provider = new ethers.providers.JsonRpcProvider(NODE_URL);

  // provider is read-only get a signer for on-chain transactions
  //const signer = provider.getSigner();

  return provider;
}

function getSigner(provider) {

  const privateKey = 'ae33f4fe38e8f26fe52095e2fdbfd22aff37ff177d1f6233e65d046291384278'
  //const privateKey = '503f38a9c967ed597e47fe25643985f032b072db8075426a92110f82df48dfcb';
  var wallet = new ethers.Wallet(privateKey, provider);
  
  return wallet;
}
