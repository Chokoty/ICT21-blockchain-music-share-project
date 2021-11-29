// ExpressJS Setup
const express = require('express');
const app = express();
var bodyParser = require('body-parser');

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname,'connection-org1.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 3000;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views')));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get('/', (req, res)=>{
    res.sendFile(__dirname + '/index-musicshare.html');
});

// dashboard page routing
app.get('/dashboard', (req, res) => {
    res.sendFile(__dirname + '/dashboard.html')
})

// music post 라우팅 - 음원등록
app.post('/music', async(req, res)=>{

    //nmid, ntitle, nartist, nlength, ntotal - args 5
    const nmid = req.body.nmid;
    const ntitle = req.body.ntitle;
    const nartist = req.body.nartist;
    const nlength = req.body.nlength;
    const ntotal = req.body.ntotal;

    try {
        console.log(`music post routing - ${nmid}`);
        await cc_call('register', [nmid,ntitle,nartist,nlength,ntotal])
        console.log("succeed")
        
        const myobj = {result: "success"}
        res.status(200).json(myobj)
    }
    catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        const myobj = {result: "failed"}
        res.status(200).json(myobj);
    }
});

// music post 라우팅 - 음원후원 ? put
app.post('/music', async(req, res)=>{

    //nmid, nbuyid, ndue - args 4
    const nmid = req.body.nmid;
    const nbuyid = req.body.nbuyid;
    const nstake = req.body.nstake;
    const ndue = req.body.ndue;

    console.log("got!")

    try {
        console.log(`make donate post routing - ${nmid}`);
        await cc_call('donate', [nmid,nbuyid,nstake,ndue])
        console.log("succeed")

        const myobj = {result: "success"}
        res.status(200).json(myobj)
    }
    catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        const myobj = {result: "failed"}
        res.status(200).json(myobj);
    }
});

// music post 라우팅 - 수익공유???
app.post('/music', async(req, res)=>{

    //nmid, nbuyid, ndue - args 4
    const nmid = req.body.nmid;
    const nbuyid = req.body.nbuyid;
    const nstake = req.body.nstake;
    const ndue = req.body.ndue;

    console.log("got!")

    try {
        console.log(`make share post routing - ${nmid}`);
        await cc_call('share', [nmid,nbuyid,nstake,ndue], res)
        
        console.log("succeed")

        const myobj = {result: "success"}
        res.status(200).json(myobj)
    }
    catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        const myobj = {result: "failed"}
        res.status(200).json(myobj);
    }
});


// // music get 라우팅 - 음악조회
// app.get('/music', async(req, res)=>{
//     // mquery
    
// });

// // music get 라우팅 - 후원자조회
// app.get('/music', async(req, res)=>{
//     // mquery
    
// });

// music get 라우팅 - 계약조회
app.get('/music', async(req, res)=>{
    // mquery
    try {
        const result = await cc_call("mquery");
        console.log(result.toString());
        const json = JSON.parse(result.toString());
        const data = json.map((row) => row.Record);
        res.status(200).send({ result: "succeed", data })
    } catch(err) {
        console.log(err);
        res.status(200).send({ result: "failed"})
    }
});

// conatract post 라우팅 - 계약조건추가
app.post('/contract', async(req, res)=>{
  
    // nmid, nbuyer, nstake, ncondition - args 4
    const nmid = req.body.nmid;
    const nbuyer = req.body.nbuyer;
    const nstake = req.body.nstake;
    const ncondition = req.body.ncondition;
    
    console.log("got!")

    try {
        console.log(`contract post routing - ${nmid}`);
        await cc_call('add', [nmid,nbuyer,nstake,ncondition])
        
        console.log("succeed")

        const myobj = {result: "success"}
        res.status(200).json(myobj)
    }
    catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        const myobj = {result: "failed"}
        res.status(200).json(myobj);
    }

});

// conatract post 라우팅 - 계약서작성
app.post('/contract', async(req, res)=>{
   
    //nmid, nbuyid, ndue - args 3
    const nmid = req.body.nmid;
    const nbuyid = req.body.nbuyid;
    const ndue = req.body.ndue;

    console.log("got!")

    try {
        console.log(`contract post routing - ${nmid}`);
        await cc_call('make', [nmid,nbuyid,ndue])

        console.log("succeed")

        const myobj = {result: "success"}
        res.status(200).json(myobj)
    }
    catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        const myobj = {result: "failed"}
        res.status(200).json(myobj);
});

// contract put 라우팅 - 만기계약리셋
app.put('/contract', async(req,res)=>{
    //nmid
    const nmid = req.body.nmid;

    try {
        console.log(`expire post routing - ${nmid}`);
        await cc_call('expire', [nmid])
        console.log("succeed")
        
        const myobj = {result: "success"}
        res.status(200).json(myobj)
    }
    catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        const myobj = {result: "failed"}
        res.status(200).json(myobj);
    }
});

async function cc_call(fn_name, args){

    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log(`cc_call`);
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    const gateway = new Gateway();
    await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: false } });
    const network = await gateway.getNetwork('msharenet');
    const contract = network.getContract('musicshare');

    var result;

    if(fn_name == 'register'){
        console.log(`register ${args.toString()}`);
        result = await contract.submitTransaction('register', args[0],args[1],args[2],args[3],args[4]);
    }else if(fn_name == 'initmusic'){
        result = await contract.submitTransaction('initmusic');
    }else if(fn_name == 'add'){
        result = await contract.submitTransaction('add', args[0],args[1],args[2],args[3]);
    }else if(fn_name == 'make'){
        result = await contract.submitTransaction('make', args[0],args[1],args[2]);
    }else if(fn_name == 'donate'){
        result = await contract.submitTransaction('donate', args[0],args[1],args[2],args[3]);
    }else if(fn_name == 'mquery'){
        result = await contract.evaluateTransaction('mquery');
    }else if(fn_name == 'cquery'){
        console.log(`cquery ${args[0]}`);
        result = await contract.evaluateTransaction('cquery', args[0]);
        
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        const myobj = JSON.parse(result);
        res.status(200).json(myobj);
    }else if(fn_name == 'dquery'){
        result = await contract.evaluateTransaction('dquery', args[0]);
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        const myobj = JSON.parse(result);
        res.status(200).json(myobj);
    }else if(fn_name == 'expire'){
        result = await contract.submitTransaction('expire', args[0]);
    }else if(fn_name == 'share'){
        result = await contract.submitTransaction('share', args[0],args[1]);
    }else{
    result = 'not supported function'
    }
    // sample
    // if(fn_name == 'history'){
    //     result = await contract.evaluateTransaction('history', args[0]);
    //     console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
    //     // 1
    //     const myobj = JSON.parse(result);
    //     res.status(200).json(myobj);
    //     // 2
    //     res.status(200).json(JSON.parse(result))
    // }
    gateway.disconnect();

    return result;
}

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);