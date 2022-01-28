using System;
using System.Text;
using System.Collections.Generic;
using System.Numerics;
using System.Threading.Tasks;

using Nethereum.Web3;
using Nethereum.Hex.HexConvertors.Extensions;
using Nethereum.Signer;
using Nethereum.Signer.EIP712;
using Nethereum.Util;
using Nethereum.ABI.FunctionEncoding.Attributes;


public class Program
{

    static async Task Main(string[] args)
    {
      var web3 = new Web3("https://data-seed-prebsc-1-s1.binance.org:8545/");
      var key = new EthECKey("ae33f4fe38e8f26fe52095e2fdbfd22aff37ff177d1f6233e65d046291384278");
      var address = key.GetPublicAddress();
      Console.WriteLine("address = " + address);
      

      Eip712TypedDataSigner signer = new Eip712TypedDataSigner();

      var typedData = new TypedData
      {
          Domain = new Domain
          {
              Name = "GameItem",
              Version = "1",
              ChainId = 97, // bsc-testnet
              VerifyingContract = "0x33a226fdA02004c2e6e3b4Ad13DfF79187DD1B80"
          },
          Types = new Dictionary<string, MemberDescription[]>
          {
              ["EIP712Domain"] = new[]
              {
                  new MemberDescription {Name = "name", Type = "string"},
                  new MemberDescription {Name = "version", Type = "string"},
                  new MemberDescription {Name = "chainId", Type = "uint256"},
                  new MemberDescription {Name = "verifyingContract", Type = "address"},
              },
              ["ItemInfo"] = new[]
              {
                  new MemberDescription {Name = "tokenId", Type = "uint256"},
                  new MemberDescription {Name = "itemType", Type = "uint256"},
                  new MemberDescription {Name = "strength", Type = "uint256"},
                  new MemberDescription {Name = "level", Type = "uint256"},
                  new MemberDescription {Name = "expireTime", Type = "uint256"},
              }
          },
          PrimaryType = "ItemInfo",
          Message = new[]
          {
              new MemberValue { TypeName = "uint256", Value = 123 },
              new MemberValue { TypeName = "uint256", Value = 1 },
              new MemberValue { TypeName = "uint256", Value = 5 },
              new MemberValue { TypeName = "uint256", Value = 1 },
              new MemberValue { TypeName = "uint256", Value = 123456789 /*getExpireTime()*/ },
          }
      };

      var encodedData = signer.EncodeTypedData(typedData);
      var hashedData = Sha3Keccack.Current.CalculateHash(encodedData);
      var signature = key.SignAndCalculateV(hashedData);
      var signatureStr =  EthECDSASignature.CreateStringSignature(signature);
      Console.WriteLine("signature = " + signatureStr);
      
      var addressRecovered = new EthereumMessageSigner().EcRecover(hashedData, signature);
      Console.WriteLine("addressRecovered = " + addressRecovered);

    }

    static int getExpireTime() 
    {
      TimeSpan t = DateTime.UtcNow - new DateTime(1970, 1, 1);
      int secondsSinceEpoch = (int)t.TotalSeconds;
      Console.WriteLine("Expire = " + (secondsSinceEpoch + (15 * 60)));
      return secondsSinceEpoch + (15 * 60);
    }

}
