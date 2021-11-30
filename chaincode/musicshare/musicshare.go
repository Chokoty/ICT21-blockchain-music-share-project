package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

type Music struct {
	MusicID     string     `json:"musichash"`
	Title        string     `json:"title"`
	Artist      string     `json:"artist"`
	Length      string     `json:"length"`
	EntireStake int        `json:"entire_stake"`
	Remains     int        `json : "remains"`
	Contracts   []Contract `json:"contracts"`
	Donors      []Donor    `json: "donor"`
}
type Contract struct {
	MusicID string `json:"musicID"`
	Owner   string `json:"owner"`   // condition?
	Buyer   string `json:"buyer"`
	Date    string `json:"date"`    // condition
	Expired string `json:"expired"` // condition
}
type Donor struct {
	Buyer   string `json: "buyer"`
	Stake   int    `json: "stake"`
	Date    string `json: "date"`
	Expired string `json: "expired"`
}

// 객체 = 구조체 + 메서드
type MusicAsset struct {
}

// 객체를 shim에 전달할때 필요한 함수 호출
// t *MusicAsset : MusicAsset에 속하는 메서드
// 입력값 형식 : APIstub
// 출력값 형식 : peer.Resoponse
func (t *MusicAsset) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *MusicAsset) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	// CLI나 Application에서 전달되는 Argument 가져오기
	fn, arg := APIstub.GetFunctionAndParameters()

	// Invoke : sCoet & get > 원장에 데이터를 저장 / 조회 (key:value 형식)
	var result string
	var err error

	// 기능 함수들
	if fn == "register" {
		result, err = t.registerMusic(APIstub, arg)
	} else if fn == "initmusic" {
		result, err = t.initMusic(APIstub)
	} else if fn == "add" {
		result, err = t.addCondition(APIstub, arg)
	} else if fn == "make" {
		result, err = t.makeContract(APIstub, arg)
	} else if fn == "donate" {
		result, err = t.donate(APIstub, arg)
	} else if fn == "mquery" {
		result, err = t.showMusicList(APIstub)
	} else if fn == "cquery" {
		result, err = t.queryContract(APIstub, arg)
	} else if fn == "dquery" {
		result, err = t.queryDonor(APIstub, arg)
	} else if fn == "expire" {
		result, err = t.expire(APIstub, arg)
	} else if fn == "share" {
		result, err = t.shareProfit(APIstub, arg)
	} else {
		return shim.Error("Not supporte chaincode function!!")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

// 초기 음원데이터 생성
func (t *MusicAsset) initMusic(APIstub shim.ChaincodeStubInterface) (string, error) {
	musics := []Music{
		Music{
			MusicID: "0001", Title: "Savage", Artist: "Aespa", Length: "3:59", EntireStake: 200, Remains: 200, Contracts: nil},
		Music{
			MusicID: "0002", Title: "작은것들을위한시", Artist: "BTS", Length: "3:12", EntireStake: 300, Remains: 300, Contracts: nil},
		Music{
			MusicID: "0003", Title: "술한잔해요", Artist: "지아", Length: "4:13", EntireStake: 50, Remains: 50, Contracts: nil},
		Music{
			MusicID: "0004", Title: "Tempo", Artist: "EXO", Length: "3:46", EntireStake: 150, Remains: 150, Contracts: nil},
		Music{
			MusicID: "0005", Title: "Enemy", Artist: "Imagine Dragons", Length: "2:53", EntireStake: 200, Remains: 200, Contracts: nil}}

	for i := 0; i < len(musics); i++ {
		fmt.Println("i is ", i)
		musicAsBytes, _ := json.Marshal(musics[i])
		APIstub.PutState(musics[i].MusicID, musicAsBytes)
		fmt.Println("Added", musics[i])
	}
	return "", nil
}

// 음원 등록
func (t *MusicAsset) registerMusic(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	
	// musicID, title, artist, length, 지분 생성 개수
	if len(args) != 5 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}

	// 원작자는 전체 지분을 설정
	stake, _ := strconv.Atoi(args[4])

	var music = Music{MusicID: args[0], Title: args[1], Artist: args[2], Length: args[3], EntireStake: stake, Remains: stake, Contracts: nil, Donors: nil}
	musicAsBytes, _ := json.Marshal(music)
	err := APIstub.PutState(args[0], musicAsBytes)

	if err != nil {
		return "", fmt.Errorf("failed to set asset: %s", err)
	}

	return string(musicAsBytes), nil
}

// 음악 목록 출력
func (t *MusicAsset) showMusicList(APIstub shim.ChaincodeStubInterface) (string, error) {
	startKey := "0001"
	endKey := "1000"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return "Iterator error !!", err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "Iterator error !!", err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(", ")
		}
		buffer.WriteString("{\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\":")

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- 음원 리스트 :\n%s\n", buffer.String())

	return buffer.String(), nil
	// 안될경우 바이트코드로 변환
}

// 원작자 계약조건 작성
func (t *MusicAsset) addCondition(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	
	// musicID, Buyer, Stake, 계약조건
	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}

	// musicID로 조회
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	stake2int, _ := strconv.Atoi(args[2])
	entire2int := music.EntireStake

	err := json.Unmarshal(musicAsBytes, &music)

	// list := music.Contracts
	// preset := list[0]

	// 후원자는 entireStake 줄이기
	if stake2int != 0 {
		// 남은 stake 체크
		if music.EntireStake-entire2int < 0 {
			return "", fmt.Errorf("No more stake. Out of EntireStake.")
		}

		music.EntireStake -= entire2int
	}

	// 계약내용 채우기
	date := time.Now().Format("2006-01-02")

	contract := Contract{MusicID: args[0], Owner: args[1], Buyer: args[2], Date: date, Expired: args[4]}

	music.Contracts = append(music.Contracts, contract)

	err = APIstub.PutState(args[0], musicAsBytes)
	return "", err
}

// 계약서 작성 : musicId, myId, duration
// date 추가하기 add condition
func (t *MusicAsset) makeContract(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect number of arguments !!")
	}
	// 음원ID에 해당하는 Music 갖고오기
	musicAsBytes, err := APIstub.GetState(args[0])
	music := Music{}
	json.Unmarshal(musicAsBytes, &music)
	// stake, _ := strconv.Atoi(args[2])
	owner := music.Artist

	// 판매가능한 지분이 없을 경우 종료
	// if music.Remains-stake < 0 {
	// 	return "no more stakes !! 남은 지분 수 : ", nil
	// }

	// 만료기한 설정
	now := time.Now()
	convDays, _ := strconv.Atoi(args[2])
	expired := now.AddDate(0, 0, convDays).Format("2006-01-02")
	// 계약서 작성
	contract := Contract{MusicID: args[0], Owner: owner, Buyer: args[1], Date: now.Format("2006-01-02"), Expired: expired}
	// 지분 차감 & 저장
	// music.EntireStake -= stake
	music.Contracts = append(music.Contracts, contract)
	musicAsBytes, _ = json.Marshal(music)

	err = APIstub.PutState(args[0], musicAsBytes)
	return "the contract has successfully set", err
}

// 후원하기
// musicID, myID, stake, duration
func (t *MusicAsset) donate(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}

	// musicID로 조회
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	stake, _ := strconv.Atoi(args[2])
	duration, _ := strconv.Atoi(args[3])
	err := json.Unmarshal(musicAsBytes, &music)
	if err != nil {
		return "", fmt.Errorf("error on unmarshalling")
	}

	// 후원자는 Remains 줄이기
	// 남은 stake 체크
	fmt.Printf("%s %s", music.Remains, stake)
	if music.Remains-stake < 0 {
		return "", fmt.Errorf("No more stake. Out of Remains.")
	}

	music.Remains -= stake

	// 후원자 명단 생성
	now := time.Now()
	expired := now.AddDate(0, 0, duration).Format("2006-01-02")

	donor := Donor{Buyer: args[1], Stake: stake, Date: now.Format("2006-01-02"), Expired: expired}

	music.Donors = append(music.Donors, donor)
	musicAsBytes, _ = json.Marshal(music)
	err = APIstub.PutState(args[0], musicAsBytes)
	return "", err
}

// 계약서 조회 : musicId
func (t *MusicAsset) queryContract(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	// musicID로 조회 & 해당 음악의 지분 리스트 추출
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	err := json.Unmarshal(musicAsBytes, &music)

	list := music.Contracts
	if list == nil {
		return "the list is empty !!", nil
	}
	var buffer bytes.Buffer

	buffer.WriteString(music.MusicID)
	buffer.WriteString(" 저작권 대여 리스트 [\n")
	for _, c := range list {
		buffer.WriteString("소유주: ")
		buffer.WriteString(c.Buyer)
		buffer.WriteString(" 구매일: ")
		buffer.WriteString(c.Date)
		buffer.WriteString(" 만료일: ")
		buffer.WriteString(c.Expired)
		buffer.WriteString("\n")
	}

	buffer.WriteString("]")
	fmt.Printf("- 계약내용조회:\n%s\n", buffer.String())

	return buffer.String(), err
}

// 후원자 조회
func (t *MusicAsset) queryDonor(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	
	// musicID
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}

	// musicID로 조회 & 해당 음악의 지분 리스트 추출
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	err := json.Unmarshal(musicAsBytes, &music)

	list := music.Donors
	if list == nil {
		return "the list is empty !!", nil
	}

	var buffer bytes.Buffer

	buffer.WriteString(music.MusicID)
	buffer.WriteString(" 후원 리스트\n[")
	for _, d := range list {
		buffer.WriteString("소유주: ")
		buffer.WriteString(d.Buyer)
		buffer.WriteString(" 지분: ")
		stake2str := strconv.Itoa(d.Stake)
		buffer.WriteString(stake2str)
		buffer.WriteString(" 구매일: ")
		buffer.WriteString(d.Date)
		buffer.WriteString(" 만료일: ")
		buffer.WriteString(d.Expired)
		buffer.WriteString("\n")
	}

	fmt.Printf("- 후원이력조회:\n%s\n", buffer.String())

	return buffer.String(), err
}

// 만기 지분 수거 : 아직 웹에서 미사용
func (t *MusicAsset) expire(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	
	// musicId
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	// musicID로 조회
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	err := json.Unmarshal(musicAsBytes, &music)

	list := music.Donors
	if list == nil {
		return "the list is empty !!", nil
	}

	// now := time.Now().Format("2006-01-02")
	now := time.Now()

	// 만기된 지분 수거

	for _, l := range list {
		expired, _ := time.Parse("2006-01-02", l.Expired)
		if now.After(expired) {
			music.Remains += l.Stake
		}
	}
	// 원장에 저장
	musicAsBytes, _ = json.Marshal(music)
	err = APIstub.PutState(args[0], musicAsBytes)
	if err != nil {
		return "Can't update the whole stake !!", err
	}

	return "지분 수거 완료", nil
}

// 수익 분배 : musicId, profit(원작자 것만 구현)
// 원작자 수익 -> 지분에 따라 후원자에게 분배
// 2차 창작자 수익 -> 계약 내용에 따라 원작자에게 분배
func (t *MusicAsset) shareProfit(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	// musicID로 조회 & 해당 음악의 지분 리스트 추출
	musicAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return "Invalid music ID !!", err
	}
	music := Music{}
	whole, _ := strconv.Atoi(args[1]) // 총수익
	err = json.Unmarshal(musicAsBytes, &music)

	list := music.Donors
	if list == nil {
		return "the list is empty !!", nil
	}
	list_d := []Donor{} // 유효 후원자 list

	wholeStake := music.EntireStake
	now := time.Now().Format("2006-01-02")
	now2int, _ := strconv.Atoi(now)

	// 만료되지 않은 지분에 대해서 지분 계산
	for _, l := range list {
		expired2int, _ := strconv.Atoi(l.Expired)
		left := expired2int - now2int
		if l.Stake != 0 && (left >= 0) {
			list_d = append(list_d, l)
		}
	}

	var buffer bytes.Buffer

	if list_d != nil {
		buffer.WriteString(music.MusicID)
		buffer.WriteString(" 수익분배표\n")
		for _, d := range list_d {
			share := whole * d.Stake / wholeStake
			buffer.WriteString("소유주: ")
			buffer.WriteString(d.Buyer)
			buffer.WriteString(" 배당금: ")
			buffer.WriteString(strconv.Itoa(share))
			buffer.WriteString("\n")
		}
	}
	fmt.Printf("- 수익분배조회:\n%s\n", buffer.String())

	return buffer.String(), nil
}

func main() {
	err := shim.Start(new(MusicAsset)) // shim에 chaincode 객체를 인자로 전달
	if err != nil {
		fmt.Printf("error creating new Smart Contract: %s", err)
	}
}
