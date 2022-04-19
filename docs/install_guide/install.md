# 1. 安装指导

## 1.1. 安装环境

ubuntu 16或18
centOS 7
golang版本1.14及以上，同时需要C编译器

## 1.3. 生成节点初始配置信息

运行代码目录中 src/consensus/hotstuff/tool/encoder_test.go 的TestGenerateAndEncode(), 修改节点数，生成所需配置，并记录拷贝

```
cd src/consensus/hotstuff/tool
go test -v -run TestGenerateAndEncode
```

得到打印数据，保存到一个备份文件中，例如保存到`data.inf`文件中

```
 encoder_test.go:156: addr: 0x4994c8b1e38df4366708129556F74B9378D90bA6, pubKey: 0x02d9e015df9c0879c0a76b840b639f19d711870fd35476499183308a02807d845e, nodeKey: 0xd275abc4ef28f37d7c7ca986b08b702c96925c5b8d32d99f522d16a0ff8ae5cc, static-node-info:d9e015df9c0879c0a76b840b639f19d711870fd35476499183308a02807d845e44222d35fd5449336df7c70aaba9b7d48ad830497c7859e136551e8f9aa77b06
    encoder_test.go:156: addr: 0x793005337B82e2B10888275fe89819704C990514, pubKey: 0x03e3f944c6de95b50f8a050612353b3e9216dd4473e756eee0874b487866879116, nodeKey: 0x0224592e69eda26d752770ffc0c4eca6a25c6cf25d20b2bdd630b6b798f287bc, static-node-info:e3f944c6de95b50f8a050612353b3e9216dd4473e756eee0874b4878668791166e07b6f4643ad50f1ec0d75fd38a9d244e124dcf76fa2f364e7d6a1ce159eb71
    encoder_test.go:156: addr: 0x95749ABe9fD1Db548ad1ae6430f67B503903A620, pubKey: 0x02dd2063081e184165796c767dfa8beb1a40c1c0fafeba0b612a519af2de200cdf, nodeKey: 0x10e8ce765fe72c4f126305c785795272ebc6224f2a9e8c3d7304983546b40e66, static-node-info:dd2063081e184165796c767dfa8beb1a40c1c0fafeba0b612a519af2de200cdf730e7aeedeb763639f4ec892e846cf1d40a10eca061b67a8289cdd7195886046
    encoder_test.go:156: addr: 0xB2FE1032cCcDCA900Cb5C9172035CBAc3EE1D3b1, pubKey: 0x0263b58bc4033eecbe92264c7639da0c076e22cf0b0d0e52dbf429c8af5488f18f, nodeKey: 0x81c8c46dc777339a0b014c2dadd48cb49a04a16940971d278bde1693e859c043, static-node-info:63b58bc4033eecbe92264c7639da0c076e22cf0b0d0e52dbf429c8af5488f18f486c3c65f049d128a48f3d5b1a147172e038c98c2a7aef2a0513be0227616676
    encoder_test.go:156: addr: 0xE392b8ADB88ed03Cb8AFC4F96Ae3A216bEaa4A78, pubKey: 0x03b886df0dc4aef11a15cb714afc48bf50aa7150af327389d7cf117cabc246f8d4, nodeKey: 0x804c8603b7212a567ffb7e80ecbab74a8251e57f0d04f41f4c94871d4276470a, static-node-info:b886df0dc4aef11a15cb714afc48bf50aa7150af327389d7cf117cabc246f8d4c2fa0b4cbe53fe8f52abba3927a98e267cd2611366c5598d798e1dc64cea1d9f
    encoder_test.go:156: addr: 0xa0671B1Eb57A99f517c8a2Fbc6f32809f4aA1199, pubKey: 0x021e29df009fb1ee4c97121d6dfd884a8533dd7a17d069e437e6ce4160ab2bb7aa, nodeKey: 0xf769ea72d4d0322edb69eccf66173846d511614bd12b2ce35ad32b9162db3e88, static-node-info:1e29df009fb1ee4c97121d6dfd884a8533dd7a17d069e437e6ce4160ab2bb7aab33b4a0867ef2a9bbd859390040237bbc3f4dbe34e1033a2fe60e50b045d3fd8
    encoder_test.go:156: addr: 0xb91fD222A217f334b6D478F72296F46e4A486dB7, pubKey: 0x02ca68c41c27bbbbe5bc9d6d458062ea3c910cf27116445ef8833dac5b16fd0e1e, nodeKey: 0xe29f41fdda43413c0caebfb0003c7c9d018799b0bd27e412839faad381ff5358, static-node-info:ca68c41c27bbbbe5bc9d6d458062ea3c910cf27116445ef8833dac5b16fd0e1ef98b487881710ef3b12597dcd25ee8688520d52b834a2a94f5c182dc1cbdc9f0
    encoder_test.go:159: ==================================================================
    encoder_test.go:165: genesis extra 0x0000000000000000000000000000000000000000000000000000000000000000f8daf893944994c8b1e38df4366708129556f74b9378d90ba694793005337b82e2b10888275fe89819704c9905149495749abe9fd1db548ad1ae6430f67b503903a62094b2fe1032cccdca900cb5c9172035cbac3ee1d3b194e392b8adb88ed03cb8afc4f96ae3a216beaa4a7894a0671b1eb57a99f517c8a2fbc6f32809f4aa119994b91fd222a217f334b6d478f72296f46e4a486db7b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080
    encoder_test.go:167: ==================================================================
    encoder_test.go:170: [
        	"enode://d9e015df9c0879c0a76b840b639f19d711870fd35476499183308a02807d845e44222d35fd5449336df7c70aaba9b7d48ad830497c7859e136551e8f9aa77b06@127.0.0.1:30300?discport=0",
        	"enode://e3f944c6de95b50f8a050612353b3e9216dd4473e756eee0874b4878668791166e07b6f4643ad50f1ec0d75fd38a9d244e124dcf76fa2f364e7d6a1ce159eb71@127.0.0.1:30300?discport=0",
        	"enode://dd2063081e184165796c767dfa8beb1a40c1c0fafeba0b612a519af2de200cdf730e7aeedeb763639f4ec892e846cf1d40a10eca061b67a8289cdd7195886046@127.0.0.1:30300?discport=0",
        	"enode://63b58bc4033eecbe92264c7639da0c076e22cf0b0d0e52dbf429c8af5488f18f486c3c65f049d128a48f3d5b1a147172e038c98c2a7aef2a0513be0227616676@127.0.0.1:30300?discport=0",
        	"enode://b886df0dc4aef11a15cb714afc48bf50aa7150af327389d7cf117cabc246f8d4c2fa0b4cbe53fe8f52abba3927a98e267cd2611366c5598d798e1dc64cea1d9f@127.0.0.1:30300?discport=0",
        	"enode://1e29df009fb1ee4c97121d6dfd884a8533dd7a17d069e437e6ce4160ab2bb7aab33b4a0867ef2a9bbd859390040237bbc3f4dbe34e1033a2fe60e50b045d3fd8@127.0.0.1:30300?discport=0",
        	"enode://ca68c41c27bbbbe5bc9d6d458062ea3c910cf27116445ef8833dac5b16fd0e1ef98b487881710ef3b12597dcd25ee8688520d52b834a2a94f5c182dc1cbdc9f0@127.0.0.1:30300?discport=0"
        ]
```

## 1.4. 拷贝setup文件夹, init.sh, start.sh, stop.sh, geth二进制文件到链的安装目录根目录下

```
cp setup /to/your/installpath
cp init.sh /to/your/installpath
cp start.sh /to/your/installpath
cp stop.sh /to/your/installpath
```

## 1.5. 将setup/node i的目录中的nodekey和pubkey改成data.inf对应的key

```
cd setup
```
顺序修改node0-node6中的nodekey和pubkey，将data.inf中的nodeKey和pubKey按顺序填入七节点
其中nodeKey的填入要删除0x前缀

## 1.6. 配置genesis.json和static-nodes.json,替换chainId，alloc，extraData，和ecode url

```
vim setup/genesis.json
```

将alloc的每行的key替换成data.inf的address，publicKey的内容替换成data.inf的pubKey
将extraData的内容替换成data.inf的genesis extra
如果需要修改chainId，则可以将自定义chainId填入

```
vim setup/static-nodes.json
```
将data.inf的encode url信息替换到文件中，注意要根据机器节点情况修改对应ip和port

## 1.8. 执行init.sh

按顺序执行7遍init.sh, 在console的交互中输入0-6的节点号

## 1.9. 修改start.sh

修改start.sh的coinbase，填入七个节点的私钥地址
如果修改了genesis.json文件的chainId的话，要修改启动参数的chainid

## 1.10. 启动节点

执行7遍start.sh, 在console的交互中顺序输入0-6的节点号

## 1.11. 停止节点

执行7遍stop.sh，在console的交互中顺序输入0-6的节点号