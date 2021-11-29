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

// music post 라우팅 - 음원등록 음원후원, 수익공유
app.post('/music', async(req, res)=>{
    const mode = req.body.mode;

    if(mode == 'register')
    {
        //nmid, ntitle, nartist, nlength, ntotal - args 5
        const nmid = req.body.nmid;
        const ntitle = req.body.ntitle;
        const nartist = req.body.nartist;
        const nlength = req.body.nlength;
        const ntotal = req.body.ntotal;
        
        console.log("got!")

        try {
            console.log(`music post routing - ${nmid}`);
            result = cc_call('register', [nmid,ntitle,nartist,nlength,ntotal], res)
        }
        catch (error) {
            console.error(`Failed to submit transaction: ${error}`);
        }

    }else if(mode == 'donate'){

        //nmid, nbuyid, ndue - args 4
        const nmid = req.body.nmid;
        const nbuyid = req.body.nbuyid;
        const nstake = req.body.nstake;
        const ndue = req.body.ndue;

        console.log("got!")

        try {
            console.log(`make contract post routing - ${nmid}`);
            result = cc_call('donate', [nmid,nbuyid,nstake,ndue], res)
        }
        catch (error) {
            console.error(`Failed to submit transaction: ${error}`);
        }
     }else if(mode == 'share'){ // ###yet
        //nmid, nbuyid, ndue - args 4
        const nmid = req.body.nmid;
        const nbuyid = req.body.nbuyid;
        const nstake = req.body.nstake;
        const ndue = req.body.ndue;

        console.log("got!")

        try {
            console.log(`make contract post routing - ${nmid}`);
            result = cc_call('share', [nmid,nbuyid,nstake,ndue], res)
        }
        catch (error) {
            console.error(`Failed to submit transaction: ${error}`);
        }
    }

    // const myobj = {result: "success"}
    // res.status(200).json(myobj)
});

// music get 라우팅 - 음악조회, 후원자조회, 계약조회
app.get('/music', async(req, res)=>{
    const mode = req.body.mode;
    
    if(mode == 'mquery') 
    {
    
    }else if(mode == 'dquery'){

    }else if(mode == 'cquery'){

    }

});
// conatract post 라우팅 - 계약조건추가, 계약서작성, 만기리셋
app.post('/contract', async(req, res)=>{
    const mode = req.body.mode;

    if(mode == 'add') // ###yet
    {
        //nmid, ntitle, nartist, nlength, ntotal - args 5
        const nmid = req.body.nmid;
        const ntitle = req.body.ntitle;
        const nartist = req.body.nartist;
        const nlength = req.body.nlength;
        const ntotal = req.body.ntotal;
        
        console.log("got!")

        try {
            console.log(`music post routing - ${nmid}`);
            result = cc_call('register', [nmid,ntitle,nartist,nlength,ntotal], res)
        }
        catch (error) {
            console.error(`Failed to submit transaction: ${error}`);
        }

    }else if(mode == 'make'){

        //nmid, nbuyid, ndue - args 3
        const nmid = req.body.nmid;
        const nbuyid = req.body.nbuyid;
        const ndue = req.body.ndue;

        console.log("got!")

        try {
            console.log(`make contract post routing - ${nmid}`);
            cc_call('make', [nmid,nbuyid,ndue], res)
        }
        catch (error) {
            console.error(`Failed to submit transaction: ${error}`);
        }
     }else if(mode == 'expire'){ // ###yet
        //nmid, nbuyid, ndue - args 4
        const nmid = req.body.nmid;
        const nbuyid = req.body.nbuyid;
        const nstake = req.body.nstake;
        const ndue = req.body.ndue;

        console.log("got!")

        try {
            console.log(`make contract post routing - ${nmid}`);
            result = cc_call('share', [nmid,nbuyid,nstake,ndue], res)
        }
        catch (error) {
            console.error(`Failed to submit transaction: ${error}`);
        }
    }

    // const myobj = {result: "success"}
    // res.status(200).json(myobj)
});


// paper issue // 생성
app.post('/paper', async(req, res)=>{
    const mode = req.body.mode;

    if(mode == 'issue')
    {
        const issuer = req.body.issuer;
        const pid = req.body.pid;
        const idate = req.body.idate;
        const mdate = req.body.mdate;
        const fvalue = req.body.fvalue;
        result = cc_call('issue', [issuer, pid, idate, mdate, fvalue])
    }else if(mode == 'buy'){
        const pid = req.body.pid;
        const from = req.body.from;
        const to = req.body.to;
        const price = req.body.price;
        result = cc_call('buy', [pid, from, to, price])
    }else if(mode == 'redeem'){
        const pid = req.body.pid;
        const from = req.body.from;
        const to = req.body.to;
        result = cc_call('redeem', [pid, from, to])
    }
    const myobj = {result: "success"}
    res.status(200).json(myobj)
});

// history
app.get('/paper', async(req, res)=>{

    try {
        const id = req.query.pid;
        console.log(`${id}`);

        result = cc_call('history', [id]);

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
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('papercontract');
        
        result = await contract.evaluateTransaction('history', id);
        
        gateway.disconnect();
        
        console.log(`${result}`);

        const myobj = JSON.parse(result);
        res.status(200).json(myobj);

    }
    catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        //process.exit(1);
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
    await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });
    const network = await gateway.getNetwork('msharenet');
    const contract = network.getContract('musicshare');

    var result;

<<<<<<< HEAD
    if(fn_name == 'register'){
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
        result = await contract.evaluateTransaction('mquery', args[0]);
=======
    if(fn_name == 'issue'){
        result = await contract.submitTransaction('issue', args[0],args[1],args[2],args[3],args[4]); //invoke
    }else if(fn_name == 'buy'){
        result = await contract.submitTransaction('buy', args[0],args[1],args[2],args[3]);
    }else if(fn_name == 'redeem'){
        result = await contract.submitTransaction('redeem', args[0],args[1],args[2]);
    }else if(fn_name == 'history'){
        result = await contract.evaluateTransaction('history', args[0]); // query
>>>>>>> d170cf8befea26a883b4fd97554ae7c2de461bf1
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        const myobj = JSON.parse(result);
        res.status(200).json(myobj);
    }else if(fn_name == 'cquery'){
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

    return result   ;
}

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);