/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package cross_chain_manager

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/cross_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/side_chain_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/community"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/info_sync"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

var (
	sdb     *state.StateDB
	signers []common.Address
	keys    []*ecdsa.PrivateKey
)

func init() {
	db := rawdb.NewMemoryDatabase()
	sdb, _ = state.New(common.Hash{}, state.NewDatabase(db), nil)
	signers, keys = native.GenerateTestPeers(2)

	node_manager.InitNodeManager()
	side_chain_manager.InitSideChainManager()
	InitCrossChainManager()
	info_sync.InitInfoSync()

	community.StoreCommunityInfo(sdb, big.NewInt(2000), common.EmptyAddress)
	node_manager.StoreGenesisEpoch(sdb, signers, signers)
	node_manager.StoreGenesisGlobalConfig(sdb)

	param := new(side_chain_manager.RegisterSideChainParam)
	param.ChainID = 8
	param.Name = "mychain"
	param.Router = 1

	param1 := new(side_chain_manager.RegisterSideChainParam)
	param1.ChainID = 79
	param1.Name = strings.Repeat("1", 100)
	param1.ExtraInfo = make([]byte, 1000000)
	param1.CCMCAddress = make([]byte, 1000)
	param1.Router = 1

	ccd := common.HexToAddress("0xdedace1809079e241234d546e44517f31b57ab8f")
	param2 := new(side_chain_manager.RegisterSideChainParam)
	param2.ChainID = 10
	param2.Router = 2
	param2.Name = "chain10"
	param2.CCMCAddress = ccd.Bytes()

	param3 := new(side_chain_manager.RegisterSideChainParam)
	param3.ChainID = 11
	param3.Router = 2
	param3.Name = strings.Repeat("1", 100)
	param3.ExtraInfo = make([]byte, 1000000)
	param3.CCMCAddress = ccd.Bytes()

	param4 := *param3
	param4.ChainID = 2

	for _, param := range []*side_chain_manager.RegisterSideChainParam{param, param1, param2, param3, &param4} {
		input, err := utils.PackMethodWithStruct(side_chain_manager.ABI, side_chain_manager_abi.MethodRegisterSideChain, param)
		if err != nil {
			panic(err)
		}
		caller := signers[0]
		contractRef := native.NewContractRef(sdb, caller, caller, big.NewInt(1), common.Hash{}, 10000000, nil)
		_, _, err = contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)
		if err != nil {
			panic(err)
		}
		p := new(side_chain_manager.ChainIDParam)
		p.ChainID = param.ChainID

		input, err = utils.PackMethodWithStruct(side_chain_manager.ABI, side_chain_manager_abi.MethodApproveRegisterSideChain, p)
		if err != nil {
			panic(err)
		}
		contractRef = native.NewContractRef(sdb, caller, caller, big.NewInt(1), common.Hash{}, 10000000, nil)
		_, _, err = contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)
		if err != nil {
			panic(err)
		}
		caller = signers[1]
		contractRef = native.NewContractRef(sdb, caller, caller, big.NewInt(1), common.Hash{}, 10000000, nil)
		_, _, err = contractRef.NativeCall(caller, utils.SideChainManagerContractAddress, input)
		if err != nil {
			panic(err)
		}

		contract := native.NewNativeContract(sdb, contractRef)
		sideChain, err := side_chain_manager.GetSideChainObject(contract, param.ChainID)
		if err != nil {
			panic(err)
		}
		if sideChain == nil {
			panic("side chain not ready yet")
		}
	}
}

func TestImportOuterTransfer(t *testing.T) {

	syncRoot := func(chainID uint64, rootHash common.Hash) {
		data, err := json.Marshal(struct {
			Root common.Hash `json:"stateRoot" gencodec:"required"`
		}{rootHash})
		assert.Nil(t, err)
		data, err = rlp.EncodeToBytes(&info_sync.RootInfo{Height: 12641624, Info: data})
		assert.Nil(t, err)
		param := &info_sync.SyncRootInfoParam{
			ChainID:   chainID,
			RootInfos: [][]byte{data},
		}
		for i := 0; i < 2; i++ {
			digest, err := param.Digest()
			assert.Nil(t, err)
			param.Signature, err = crypto.Sign(digest, keys[i])
			assert.Nil(t, err)

			input, err := utils.PackMethodWithStruct(info_sync.ABI, info_sync.MethodSyncRootInfo, param)
			assert.Nil(t, err)

			blockNumber := big.NewInt(1)
			caller := common.Address{}
			contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, 10000000, nil)
			_, _, err = contractRef.NativeCall(caller, utils.InfoSyncContractAddress, input)
			assert.Nil(t, err)
		}
	}

	event, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000004f00000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000001e00000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000290000000000000000000000000000000000000000000000000000000000000020103f9e71e982c9d4e250e988ce7ed99c220e80c8184a94d1f28f8b23c8b8fe300000000000000000000000000000000000000000000000000000000000000014d4c894eb6829301f23bc2777a532209c7c11f4f50000000000000000000000000000000000000000000000000000000000000000000000000000000000000014f0a8515244b2dc9c7885cc3d83b04d976803c1980000000000000000000000000000000000000000000000000000000000000000000000000000000000000006756e6c6f636b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004a147466366d6d303f41a1876c45b1acfbc2b17123e4140c888cca1190940ebc156d4cf13cbf880a83e4a3010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.Nil(t, err)
	proof, err := hex.DecodeString("7b2261646472657373223a22307864656461636531383039303739653234313233346435343665343435313766333162353761623866222c2262616c616e6365223a22307830222c22636f646548617368223a22307835373230313031626533363130656362643638303639383834366462373463356162366264356534363335363963353232633536373639613331316439353434222c226e6f6e6365223a22307831222c2273746f7261676548617368223a22307861333061666636623433313763306630396465366163393736633238343261376633373762323163633831373731316366363037643232333162383735663731222c226163636f756e7450726f6f66223a5b2230786639303231316130353438366461373535616134396235663233643534393363646632396566323437316332303432623133306339616534393131396661333133366235373237356130663935373262613233663663343739373463393838326366356630343764643634373466663330613630376562643537643237656530356431646432386532306130326532303536656564343838383739346565343938383632353137326466633731373434653232303666343962613738363337386266633139616230616630336130613939316466626336316363356539323563346133333162633838343137323831333461646163323931376166356439666663663333616333376462666463366130323535356331363735653635386164386232353561343731353931393335316236663630333331656339313566323436633961336565313665393062326333626130343863336261386366316537613164333532363465303931643665323139666165303438383037363032623034613635326630313736313031393933373636616130343338336432376638356631613263373731343938353135313331303936623265646337633537366435303963333339373537386664613362326336343232316130373032333839363236333439613962636133303066316332336630383830643830623865343431656566623832353234646233613366393536613566613639356130626263323630316333656462613630353832363962656465376333323330623762616132373761626132663862323362613036356563316238393166333065666130616161303662356666333562653837366330643863363465383232666434336639356263663062393138313164666330373132353466623730373331656431386130363534643835316434613032346163653531643333363666396231363666613366383136306230383264323937333733366462373264326365383163323963626130323939393365316533626538393837373932633237656264376636383438323034383066636430623734396538303532373962373032363937646530633464376130376365353161613437376330346633313863363832346532346334616661363566616530613036396537343432643265363064343463303561613062626366636130613533386634393863386332643436306663396663626463623531343966383864613734333765613363663561653466343838383132363764336537323637646130336232376261653936623631623239656636353661313764623736613837316162653334356661643639656636616436663732643466313531323730383135386130366266633035373666623738363263356363373631646132316464363765373030333435663030376362316638376333326435633365303562653032323063363830222c2230786639303231316130653663383364376634663665356134373562666538386333623239653165336534383761643938333536653332313437386465633166393534353939643831306130363930306365626664616239656532613933303664373733653439346562353639353061353530303838353039303630346234643536353630396563636262306130353632386566343135663061323661376239323930633266366235353262353965643565356132353165373534396362646465393538646234323963323935646130393137363834393235386464383333383332363836633062323464333830386161616231646632393037623135643735333938666565333366333465313464336130353562633537626463356666356236393464333665333831353563366665653063623738633739323866346332663631363130313335653539373935313533366130323131353264363163313737333433336637326232396361373638636136346634363539383563653363656264393733383537646661333061383432346433306130333162386330393264383035613462363735303934643330313036666166336232333262303135326164363932333430373039306235616437643064633734646130376636333261613130646530623238333639353339373463616233663538316136303639656630636236313465323839353164396432613039636531363361336130316630373766363462633038373733333639336538626333623362396261626662653035326261383739336537633837326566363633323564633665316331356130336431346233383464333166396239326161393961333330346232396135313865653962363336626136366639303562366165633433366638316136653835376130636533313866313465613630373639373535663065666532366433393463383666303161623164386463383961633939616562383132303733326434623032636130336533353636623732376135383838333836386665326631646331383262366431383835316137613832303765636466366239666366666432623937383939386130363564396463343435623262353230363934393231663736333363656634623836653563643566646566336261646538646266363538623765306661613365336130363933343032396232336239346261303732663462393835323633386236643935636661306539376530316136366634383837353839306131326461333535306130653032333637323731643631393939656363396433313762633032656130303563353266393561303137376136316530316432303733313933363133333738616130323463633962303964346636393536373138363134383930623465353131623037313730373237633833613566326634633166363038613962383766356362323830222c2230786639303231316130616465336637646431393136393833303537626139363339303833356639653165303235333930393437383537363634616130623864343335393732383738616130616533376637643135306138613661616564376161666165633038656534386630313635393035383033383530393663313961393333623563353335623164626130326139653136653466353835353438323462303061646430646637346266393237666232616262333961623335626630656231346262633939646131346263316130393562613539666238333738373862666163386234336432333566643233326632313536646533363738343931323239353736636434646330313531316234636130323037383639353834626238316237356639623331303136643063363262616332303566393937623464316333363439623236396435393363613637346330326130626535616130306465666130363933376534373566306633633036303163303833303961373731323664646332306331393836316538303730613165363538376130383733393837306530656538383439343230373461363064343739613935363133346666333566323863393064626138656263333566353733346535623839636130613532383161386436636235306338643733356537323165613636343562336637666661366235643864353563636635323130656530646539313532613361666130633437663936613064366465613135313362363362653666396630343964303563326234393731633737306638616138653465336364393535626365303932336130366161353439393836383938356238313965356436313931346137646331623433336539316336313630663662356537623162343336323465623364626562376130333662653132353338643931356339613835336130373265363438636163643166383336363638303338633230373030623461343361303832633731663432656130616530613734326364363762373834626532633339663565383066626130336335643664356666346165356434626130633564656535616538373765643861646130306134623734393036363131643836303838376238373937643939323164616561386564383232663134363866373938343965653531353963636638313966666130636365633266363061343836313266363666636436623361343266393462326663333966393339336666366631376130326662366232643865353964383337346130386239353330333231643138613664323161653032653238323364333935323839356364623731346132393131343933653865383130636130386363343565656130376665373031656462626234343235326533303534383734386462626262613739356139333435666462356438306461653866383533326230663066383235613830222c2230786639303231316130323732653865393764653230346437326539643465336130613564663336316265336232356138653131303263313036326665343537393962643434393135366130646339646563353232346433366664616436633136313630323330313664653134313538353364303462666239343432366338623133363362373837373162646130646631653739303064633132613164363336303336313262633935633961333062326334663339303661613731383064646462653665326365303437353636376130333761623536663462396135353864633561346131656337326536633663376565343361653762623163616264366661626234333666333030656630303639656130383234613339626266353961393731383636613863393163643534346239303761623963303033656162626133646362393633346331333333666637316462346130613763306139363833633564346333643563666630313137373232303464633535653538313639326464396362643537663735346339336561613464643463646130393131386232353962363735373763353933616132633732616565623764633730313333636666636430633963613437383963333338623766353462363166326130633763616630373637333131626564326536343434666161373863656462333733343439623538646561333830353534373863666533366535373535613833616130633866623863323966323736663462323032316333346635653135323865376331366633333132633561353066653632663661373365363461326463663961656130363736303135383337333265386439383036653462663935613931653938323931383263356663643565383936353539633331623330353536376661376366656130333839393330656632356136353039393561633330326562623962623966333666383236313438656439356538656166343739366536333436656266393935376130643231393539616639663736646563623630383866643662613737383438313630623736393332633161323461373030656161393136323962373337376232646130313332336132393732643937646439316263633066623931663130636234613932313833316537353965393066333033343761623637316166373533633566356130306230373231373262323730396562336331653631656138653066363832333164623261316265306632376265323638626634353766306138326131343935656130323332386339663432356435613261343134396435323536383836323634613766393166626462363639313932313836356237346131393964363164343833376130393133393932303465306365646136323439386164383939643637386661373436656662313731666634646461653937656434366635633234643262373964303830222c2230786639303231316130393463653237643039333931303463616534633338336438663237373361633137663966363064633665633739623432306263313434313065643430323339356130313039616335353333663935376234393465643430623166366134356230373731373235656133306264376365656537663736623838306664346430653734326130303137356230353663393839353662333933383531323232306335353536623266376431316164643337353930323234333532656562376635306662626533626130646361363330633035643033653633373364636133333430653761396333303335623263616536316235613465336136366665613231626439346438356666336130613630326663646261646430653630373566356666323638353530353832623564323435643064623961643437323661623734336165356538646565373935316130376537656239316335303866343264633466666130333930333735643431353037396333343138363238613765323632396536393664303261343333353737386130616465616639343836656365653736653135636239643731313730623337646664356636373737303261363661333034616537633166316431613535356630356130636262316133643264333232323536643065613630366535613265323366326132343433343130613033393563633063333039363931336536336636363266326130343938303234656366353939323262643631393135353066383661616166363534623937306534663732643731623066343339323231326361656431626362626130616630316235373765383964333462643464303063396335346139353833653762386366616438653137623966643366333137316663623862313332626536366130666466376331343366643965383335336139383830346239376433393733656131613437616637643132393065616237653833666231383934653139653761656130343563396138643535396537363361356238396666383563663266356163393963356532393062393634653531613239613735396637653861363762393532646130663931303664646334316434373464393965336465666335303832636434323963346166336566346433343761303530646165643134363438656130323464366130613263366132383631643139396530643164646165396266366262366666643733663239613861333036353634313662613133326131616531386333303165346130616533386433316634383064643435666531333761363136623733323437623335656664396631396664376365663935323133383739643837633066373161636130396434383830366461656336623734373065316264366330313664353132303831666638633532303261303361333066323737613032633534383263366631653830222c22307866393031663161306265376535386338363236616633626632353330663936393930666263366261376233656535303935383361353663343861363239366231663835353832383561303038396464396636363062646530363436303333663761373262663339376165366134623066363330323833313233656534616463396234613966393535643961303535363436386264386633656431383431623537306562323235363332386435643461323962396633356537313539636164303566303964323234663832626161306139326365303939623364326238396335303036663566393531663964313036306164326561616561303738643934613237613737643264366562363935393161303661373936303161633132653430646138396661626563646234356361643762653364393837333035366231623566653563343566313530303935616439353661306664616334343464643032343737333734323239643063643938333936663065363036356432646639356132366362393932356466353935373834666663393161303231363830346161646530383835663435373061316336646161383733323438376232666262353663663733336262373039383930333931353539336534653761303332373331663063363834646466356234386532353865313433343661366361616165666332613935613037333861393035346263613430633331613666303261303335326632336635616265363463386135633861396435356461396264383734393735636135366631393339623961393965323363343961373866623930393161303565613435646365613765643931323263313439373261663839333364373534633132633739373238396364353934666539613965636466663032333561343061306661623865656232353531393535663465653261356432323337636165336438393565363665643037373432623763623639396638356338346239383061366161306634356263346262386434643736383161326636656661643230336134313263373765313639373637653832636535343565373063656235633630613634353661303233326461343161613334656663393730653435633065373634363433656333656464643431333737393030636634643137646162616535343163353065656461306666336166653566653431366430666232653433343561376437316661663763636137666339313131393337336335356661346439353134663131393066656461303030306538306432653832313264616534316230646532663030383061316431666238366265323331373430613932626366323230643265643766333765323438303830222c2230786638373138303830383038303830383038303830613066313936633134396232633262656238393765666531363237656363396430343032356563316430613133623866623730373461663761633830356462383138613066393336396133633236313863623566633061653766346539393834373161373938656465613933363235646365353966663833636134613639346131646361383038303830383038306130633032393630656234646334356666366566316533313431656464376634386631326630333931323766616433373664623437346566653665653038663462643830222c22307866383636396433326139376164343930666566386539313336323537666232336561613833363538333934313437653437613862663830303732323233336131623834366638343430313830613061333061666636623433313763306630396465366163393736633238343261376633373762323163633831373731316366363037643232333162383735663731613035373230313031626533363130656362643638303639383834366462373463356162366264356534363335363963353232633536373639613331316439353434225d2c2273746f7261676550726f6f66223a5b7b226b6579223a22307834643038636233666661333335393266336235656562383737363437393834663430373064316432393937643537633961656134633365333464646436653865222c2276616c7565223a22307832646335646135643261333662623263613135613437353336636133383033316337316162616537363332316536316365333263333163663261663333303136222c2270726f6f66223a5b2230786639303231316130363431633763303038616161373162376136616161346236626639313162303266623238386633313735383939376132303762313837366361363138386137636130643530353732346633643037303866613464333131366361353966633364343339346332663635383463396634303731316131373964616433363664373665336130313663373533323763353165356232663732613430313466623539633238653662383632326666356630333838346264336538376263333630363137376663346130663362316638306533616365366237363437616261316233316436373461383933316238336162346539656164633137633162663835306131343634666366306130623938373432363939363033373965636661336639336138646632313330613930333263313733653838643365363061383434646432386430653433643334336130393462386633663666363365313062366164316537633537306635666330656661343134303934636331656466366637333265616132353165643034653166316130386230393739346438666539386564393435666237366461383336633564626562356636366530383834646138656338363739363964636532643562393532326130326639633839383430613665393365383332623736643062303434626235373238396435343432326564356433366533623838633730393439356533366639666130663737656134343464646239373732666239343131623233356131363932653935346135633462643437383732646565613334666331396163653739643838316130643735393135383535393362333733363136303563643035313930353231333664613234663563613631653530626566373732613734643161303965616364396130333865363265633737646361346261663661616638646131393063613233363639663532663834373565353964306330656438316531636665383938393838646130373465393961626335373333373434353062386135373761343463323633633838316663656232386263386232316630636564666633373862393363356265376130316166663262646430313264306431346432343534336338613061633162376337636431663036303362646139313835393136363630383431333733633733616130386363393737383939353764626532343334306562326136383938363464333537616331326639643234336236323937316139366139333133643839316437336130386335356132383261373964346265366431306264303233376338613730393766366565323064613435663739626130363534643534356437656332616261306130373963633465623037363861646334343461303663393538666665613166666236636139663136636563363930383766336530313065636661643835383132353830222c2230786638373161306133613338616262326639356339626262306232366532623739323863343130326564303965653237313434306431333935623266633265343865303133346538306130376264306133626662623439643262333933353365333230393037373437623931373939373963653335383764666662633231633966636133343064363839363830383038303830383038303830383038303830383061303135383730346665616232626438333337663936376133626133613165633430323965633831333966656464363761613236666238313738366663376333343238303830222c223078663834336130323036643235646636663135353135386639353934313062646539646138633563643639373766373461393635386464346466386137666161326537333339386131613032646335646135643261333662623263613135613437353336636133383033316337316162616537363332316536316365333263333163663261663333303136225d7d5d7d")
	assert.Nil(t, err)

	param := new(scom.EntranceParam)
	param.SourceChainID = 8
	param.Extra = event

	param1 := new(scom.EntranceParam)
	param1.SourceChainID = 79
	param1.Extra = event

	param2 := new(scom.EntranceParam)
	param2.SourceChainID = 10
	param2.Height = 12641624
	syncRoot(10, common.HexToHash("0bdb33c1f2e4a23a8429c61c1bd31aaccc38795655e3ef9e2baf10f6567bbe3c"))
	param2.Extra = event
	param2.Proof = proof

	param3 := new(scom.EntranceParam)
	param3.SourceChainID = 11
	param3.Height = 12641624
	param3.Extra = event
	assert.Nil(t, err)
	// param3.Proof, err = hex.DecodeString("7b2261646472657373223a22307831316532613731386434366562653937363435623837663233363361666531626632386332363732222c2262616c616e6365223a22307830222c22636f646548617368223a22307830636561363334383038323338373837343332343564373739373966336338636639346338623162396433316133343764643533376466313330613966346265222c226e6f6e6365223a22307831222c2273746f7261676548617368223a22307830303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030222c226163636f756e7450726f6f66223a5b5d2c2273746f7261676550726f6f66223a5b7b226b6579223a22307831656366633263343264383065613636306330383834633739373065363764616138356265353035363238656236623363616464616563323036316332636163222c2276616c7565223a22307834323037353035336163306137393939323939663835663438653364616233363830323262316361396563313461613433663233656362666664373536666535222c2270726f6f66223a5b5d7d5d7d")
	param3.Proof = proof

	assert.Nil(t, err)
	syncRoot(11, common.HexToHash("0bdb33c1f2e4a23a8429c61c1bd31aaccc38795655e3ef9e2baf10f6567bbe3c"))

	tr := native.NewTimer(scom.MethodImportOuterTransfer)
	for _, param := range []*scom.EntranceParam{param, param1} {
		for i := 0; i < 2; i++ {
			digest, err := param.Digest()
			assert.Nil(t, err)
			param.Signature, err = crypto.Sign(digest, keys[i])
			assert.Nil(t, err)

			input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodImportOuterTransfer, param)
			assert.Nil(t, err)

			blockNumber := big.NewInt(1)
			extra := uint64(21000)
			caller := common.Address{}
			contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodImportOuterTransfer]+extra, nil)
			tr.Start()
			ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
			tr.Stop()
			assert.Nil(t, err)
			result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodImportOuterTransfer, true)
			assert.Nil(t, err)
			assert.Equal(t, ret, result)
			assert.Equal(t, leftOverGas, uint64(0))
		}
	}

	for _, param := range []*scom.EntranceParam{param2, param3} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodImportOuterTransfer, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(21000)
		caller := common.Address{}
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodImportOuterTransfer]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodImportOuterTransfer, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, uint64(0))
	}
	tr.Dump()
}

func TestReplenish(t *testing.T) {
	param := new(scom.ReplenishParam)
	param.ChainID = 8
	param.TxHashes = []string{"0x74676ce6389bbb479ffc9afe720749ad28b9500ff09c7ae8f19bd1e543f8845f"}

	param1 := new(scom.ReplenishParam)
	param1.ChainID = 9
	for i := 0; i < 200; i++ {
		param1.TxHashes = append(param1.TxHashes, "0x74676ce6389bbb479ffc9afe720749ad28b9500ff09c7ae8f19bd1e543f8845f")
	}

	tr := native.NewTimer(scom.MethodReplenish)
	for _, param := range []*scom.ReplenishParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodReplenish, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(21000)
		caller := common.Address{}
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodReplenish]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodReplenish, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, uint64(0))
	}
	tr.Dump()
}

func TestCheckDone(t *testing.T) {
	param := new(scom.CheckDoneParam)
	param.CrossChainID = make([]byte, 32)

	param1 := new(scom.CheckDoneParam)
	param1.CrossChainID = make([]byte, 2000)

	tr := native.NewTimer(scom.MethodCheckDone)
	for _, param := range []*scom.CheckDoneParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodCheckDone, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(21000)
		caller := common.Address{}
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodCheckDone]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodCheckDone, false)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, uint64(0))
	}
	tr.Dump()
}

func TestWhiteChain(t *testing.T) {
	param := new(scom.BlackChainParam)
	param.ChainID = 8

	param1 := new(scom.BlackChainParam)
	param1.ChainID = 9

	tr := native.NewTimer(scom.MethodBlackChain)
	for _, param := range []*scom.BlackChainParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodWhiteChain, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(21000)
		caller := signers[0]
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodWhiteChain]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodWhiteChain, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, uint64(0))
	}
	tr.Dump()
}

func TestBlackChain(t *testing.T) {
	param := new(scom.BlackChainParam)
	param.ChainID = 8

	param1 := new(scom.BlackChainParam)
	param1.ChainID = 9

	tr := native.NewTimer(scom.MethodBlackChain)
	for _, param := range []*scom.BlackChainParam{param, param1} {
		input, err := utils.PackMethodWithStruct(scom.ABI, cross_chain_manager_abi.MethodBlackChain, param)
		assert.Nil(t, err)

		blockNumber := big.NewInt(1)
		extra := uint64(21000)
		caller := signers[0]
		contractRef := native.NewContractRef(sdb, caller, caller, blockNumber, common.Hash{}, gasTable[cross_chain_manager_abi.MethodBlackChain]+extra, nil)
		tr.Start()
		ret, leftOverGas, err := contractRef.NativeCall(caller, utils.CrossChainManagerContractAddress, input)
		tr.Stop()
		assert.Nil(t, err)
		result, err := utils.PackOutputs(scom.ABI, cross_chain_manager_abi.MethodBlackChain, true)
		assert.Nil(t, err)
		assert.Equal(t, ret, result)
		assert.Equal(t, leftOverGas, uint64(0))
	}
	tr.Dump()
}
