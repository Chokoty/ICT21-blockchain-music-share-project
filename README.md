<p align="center"><img src="/imgs/museShare-logo.jpg"></p>

</p>
<p align="center">
		<a href="https://github.com/Chokoty/ICT21-blockchain-music-share-project/blob/main/LICENSE"><img alt="GitHub license" src="https://img.shields.io/github/license/Chokoty/ICT21-blockchain-music-share-project"></a>
</p>

<br><br><br>
<h2 align="center">프로젝트 설명 영상</h2>

<table>
	<tr>
		<td>
			<a href=""><img src="/image/INTRO_THUMBNAIL.png"></a>
		</td>
		<td>
			<a href=""><img src="/image/GUIDE_THUMBNAIL.png"></a>
		</td>
	</tr>
	<tr>
		<td align="center">
			<b>소개 영상</b>
		</td>
		<td align="center">
			<b>가이드 영상</b>
		</td>
	</tr>
</table>

## 프로젝트 설명  
<p><b>museShare</b>는 .</p>
<br>



## 기능 설계
 

## 개발 프로젝트 사용법 (Getting Started)

<p>
	1. 블록체인 네트워크를 구성하고<br>
	2. 체인코드 설치 배포하여 작성된 데이터를 couchdb에서 확인합니다. <br>
	3. 클라이언트 페키지를 설치하고 사용자 인장 발급을 합니다. <br>
	4. 웹에서 요청한 데이터를 하이퍼레져 패브릭 인터페이스를 통해 처리합니다.

</p>
<p>step1. 뮤즈쉐어 저장소를 클론하고 클론한 폴더로 이동합니다.</p>

```bash
$ git clone https://github.com/Chokoty/ICT21-blockchain-music-share-project.git
$ cd ICT21-blockchain-music-share-project
```

<br>
<p>step2. basic-network에서 도커를 이용해 블록체인 네트워크 환경을 구성합니다.</p>

```bash
$ cd basic-network/
// 도커 실행전 준비물 세팅, 새로운 CA생성, 제네시스 블록...(최초 1회)
$ ./generate.sh
// 도커 cli 스크립트 실행 - org peer couchdb ca ... 생성
$ ./start.sh
```

<br>
<p>step3. 체인코드를 설치 배포, 테스트합니다. </p>

```bash
// 체인코드 설치, 배포, 테스트 쉘 스크립트 실행
$ cd ICT21-blockchain-music-share-project/chaincode/musicshare
// 최초 설치 시 다음 명령어 실행
$ ./cc_ms_v1.sh instantiate 1.0
// 이후 재설치 시 버전업으로 실행
$ ./cc_ms_v1.sh upgrade 1.1~
```
<p>couchdb1에서 이력을 확인합니다. http://localhost:5984/_utils/ </p>

<br>
<p>step4. 클라이언트 페키지를 설치하고 사용자 인장 발급을 합니다.</p>


```bash
// 클라이언트 패키지 설치
$ cd ICT21-blockchain-music-share-project/simple-web
$ npm install
// 블록체인네트워크에 접근이 가능하도록 ca와 유저 wallet 등록
$ node enrollAdmin.js
$ node node registerUser.js
```

<br>
<p>step5. 웹서버를 실행합니다.  http://localhost:3000/</p>

```bash
$ npm start
```
 -----

## 팀 정보 (Team Information)
- Kim tae yoon (choko0816@ajou.ac.kr), Github Id: Chokoty
- Gong myung gyu (gmk0904@ajou.ac.kr), Github Id: MyeongQ
- Choi won bin (peactor@gmail.com), Github Id: peactor
 -----
 
## 프로젝트 후원 기관 (Sponsoring Organization)
- 서울ICT이노베이션스퀘어 : https://ict.eksa.or.kr/interact/ict.user
- 아주대학교 LINK+ 사업단 : https://lincplus.ajou.ac.kr/
- NIPA 정보통신산업진흥원 : https://www.nipa.kr/
- KSA 한국표준협회 : https://www.ksa.or.kr/ksa_kr/index.do
 -----

## 저작권 및 사용권 정보 (Copyleft / End User License)
 * [MIT](https://github.com/osam2020-WEB/Sample-ProjectName-TeamName/blob/master/license.md)
 

